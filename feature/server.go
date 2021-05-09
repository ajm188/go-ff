package feature

import (
	"context"
	"fmt"
	"sync"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"

	featurepb "github.com/ajm188/goff/proto/feature"
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
	s.m.RLock()
	defer s.m.RUnlock()

	if feat, ok := s.features[req.Name]; ok {
		return &featurepb.GetFeatureResponse{
			Feature: feat.Feature,
		}, nil
	}

	return nil, fmt.Errorf("%w with name %s", ErrNoFeature, req.Name)
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

	if after.Type == featurepb.Feature_PERCENTAGE_BASED {
		if after.Percentage < 0 || after.Percentage > 100 {
			return nil, fmt.Errorf("%w percentage must be in [0, 100]; have %d", ErrInvalidFeature, after.Percentage)
		}
	}

	s.features[req.Feature.Name] = &Feature{Feature: after}
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
