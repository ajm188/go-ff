package feature

import (
	"embed"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"sort"
	"strings"
	"text/template"

	featurepb "github.com/ajm188/go-ff/proto/feature"
)

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

				return strings.Title(strings.ToLower(strings.ReplaceAll(name, "_", "-")))
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
	sort.SliceStable(resp.Features, func(i, j int) bool {
		left, right := resp.Features[i], resp.Features[j]
		return compare(left, right)
	})
	if err := indexTmpl.Execute(w, resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
}

func compare(left, right *featurepb.Feature) bool {
	leftEnabled := isNaivelyEnabled(left)
	rightEnabled := isNaivelyEnabled(right)

	if leftEnabled {
		if rightEnabled {
			if left.Type == right.Type {
				if left.Type == featurepb.Feature_PERCENTAGE_BASED {
					return (left.Percentage > right.Percentage) || (left.Name < right.Name)
				}

				return left.Name < right.Name
			}

			// Constants go before PBR
			if left.Type == featurepb.Feature_CONSTANT {
				return true
			}

			if right.Type == featurepb.Feature_CONSTANT {
				return false
			}

			// Technically, since isNaivelyEnabled can only return true for
			// CONSTANT and PERCENTAGE_BASED types, this is unreachable.
		}

		return true
	}

	if rightEnabled {
		return false
	}

	if left.Type == right.Type {
		return left.Name < right.Name
	}

	switch left.Type {
	case featurepb.Feature_CONSTANT:
		return true
	case featurepb.Feature_PERCENTAGE_BASED:
		switch right.Type {
		case featurepb.Feature_CONSTANT:
			return false
		default:
			return true
		}
	default:
		return false
	}
}

func isNaivelyEnabled(f *featurepb.Feature) bool {
	switch f.Type {
	case featurepb.Feature_CONSTANT:
		return f.Enabled
	case featurepb.Feature_PERCENTAGE_BASED:
		return f.Percentage > 0
	}

	return false
}
