package feature

import (
	"context"
	"embed"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
	"text/template"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"

	featurepb "github.com/ajm188/go-ff/proto/feature"
)

var (
	// Global singleton.
	inst = &server{
		features: map[string]*Feature{},
	}

	_ featurepb.FeaturesServer = (*server)(nil)
)

type server struct {
	m        sync.RWMutex
	features map[string]*Feature
}

// DeleteFeature is part of the featurepb.FeaturesServer interface.
func (s *server) DeleteFeature(ctx context.Context, req *featurepb.DeleteFeatureRequest) (*featurepb.DeleteFeatureResponse, error) {
	s.m.Lock()
	defer s.m.Unlock()

	if feat, ok := s.features[req.Name]; ok {
		delete(s.features, req.Name)

		return &featurepb.DeleteFeatureResponse{
			Feature: feat.Feature,
		}, nil
	}

	return &featurepb.DeleteFeatureResponse{}, nil
}

// GetFeature is part of the featurepb.FeaturesServer interface.
func (s *server) GetFeature(ctx context.Context, req *featurepb.GetFeatureRequest) (*featurepb.GetFeatureResponse, error) {
	feat, err := s.getFeature(req.Name)
	if err != nil {
		return nil, err
	}

	return &featurepb.GetFeatureResponse{
		Feature: feat.Feature,
	}, nil
}

// getFeature separates the concurrent-safe Feature lookup from the gRPC service
// implementation so this can be reused by the package-level Get function.
func (s *server) getFeature(name string) (*Feature, error) {
	s.m.RLock()
	defer s.m.RUnlock()

	if feat, ok := s.features[name]; ok {
		return feat, nil
	}

	return nil, fmt.Errorf("%w with name %s", ErrNoFeature, name)
}

// GetFeatures is part of the featurepb.FeaturesServer interface.
func (s *server) GetFeatures(ctx context.Context, req *featurepb.GetFeaturesRequest) (*featurepb.GetFeaturesResponse, error) {
	s.m.RLock()
	defer s.m.RUnlock()

	var (
		features []*featurepb.Feature
		names    = make([]string, 0, len(s.features))
	)

	if !req.NamesOnly {
		features = make([]*featurepb.Feature, 0, len(s.features))
	}

	for name, feat := range s.features {
		names = append(names, name)

		if !req.NamesOnly {
			features = append(features, feat.Feature)
		}
	}

	return &featurepb.GetFeaturesResponse{
		Features: features,
		Names:    names,
	}, nil
}

// SetFeature is part of the featurepb.FeaturesServer interface.
func (s *server) SetFeature(ctx context.Context, req *featurepb.SetFeatureRequest) (*featurepb.SetFeatureResponse, error) {
	s.m.Lock()
	defer s.m.Unlock()

	var (
		before *featurepb.Feature
		after  = proto.Clone(req.Feature).(*featurepb.Feature)
	)

	if feat, ok := s.features[req.Feature.Name]; ok {
		before = feat.Feature
	}

	f := &Feature{Feature: after}

	switch f.Type {
	case featurepb.Feature_PERCENTAGE_BASED:
		if f.Percentage < 0 || f.Percentage > 100 {
			return nil, fmt.Errorf("%w percentage must be in [0, 100]; have %d", ErrInvalidFeature, f.Percentage)
		}
	case featurepb.Feature_EXPRESSION:
		if err := f.parseExpression(); err != nil {
			return nil, fmt.Errorf("could not parse expression %s: %w", f.Expression, err)
		}
	}

	s.features[req.Feature.Name] = f
	return &featurepb.SetFeatureResponse{
		Before: before,
		After:  after,
	}, nil
}

// RegisterServer adds the global feature server instance to the given gRPC
// server.
func RegisterServer(s *grpc.Server) {
	featurepb.RegisterFeaturesServer(s, inst)
}

var (
	//go:embed assets
	assets embed.FS

	indexTmpl     *template.Template
	liveReloadDir string
)

func init() {
	if v := os.Getenv("HTTP_LIVE_RELOAD_DIR"); v != "" {
		log.Printf("using HTTP_LIVE_RELOAD_DIR %s", v)
		liveReloadDir = v
	}
}

func Index(w http.ResponseWriter, r *http.Request) {
	if indexTmpl == nil || liveReloadDir != "" {
		var (
			data []byte
			err  error
		)
		if liveReloadDir != "" {
			data, err = ioutil.ReadFile(path.Join(liveReloadDir, "assets/index.html.tmpl"))
			if err != nil {
				err = fmt.Errorf("%w (HTTP_LIVE_RELOAD_DIR=%s)", err, liveReloadDir)
			}
		} else {
			data, err = assets.ReadFile("assets/index.html.tmpl")
		}

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		indexTmpl, err = template.New("index").Funcs(map[string]interface{}{
			"featureTypeToString": func(ftype featurepb.Feature_Type) string {
				name, ok := featurepb.Feature_Type_name[int32(ftype)]
				if !ok {
					return "unknown"
				}

				return strings.ToLower(name)
			},
			"featureSettings": func(feat *featurepb.Feature) string {
				switch feat.Type {
				case featurepb.Feature_CONSTANT:
					if feat.Enabled {
						return "on"
					} else {
						return "off"
					}
				case featurepb.Feature_PERCENTAGE_BASED:
					return fmt.Sprintf("%d%%", feat.Percentage)
				case featurepb.Feature_EXPRESSION:
					return feat.Expression
				}

				return ""
			},
		}).Parse(string(data))
		if err != nil {
			indexTmpl = nil

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	}

	resp, _ := inst.GetFeatures(r.Context(), &featurepb.GetFeaturesRequest{})
	if err := indexTmpl.Execute(w, resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
}
