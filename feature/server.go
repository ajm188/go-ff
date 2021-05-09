package feature

import (
	"context"
	"fmt"
	"sync"

	featurepb "github.com/ajm188/goff/proto/feature"
	"google.golang.org/grpc"
)

var (
	// Global singleton.
	inst = &server{
		features: map[string]bool{},
	}

	_ featurepb.FeaturesServer = (*server)(nil)
)

type server struct {
	m        sync.RWMutex
	features map[string]bool
}

// DeleteFeature is part of the featurepb.FeaturesServer interface.
func (s *server) DeleteFeature(ctx context.Context, req *featurepb.DeleteFeatureRequest) (*featurepb.DeleteFeatureResponse, error) {
	s.m.Lock()
	defer s.m.Unlock()

	if feat, ok := s.features[req.Name]; ok {
		delete(s.features, req.Name)

		return &featurepb.DeleteFeatureResponse{
			Feature: &featurepb.Feature{
				Name:    req.Name,
				Enabled: feat,
			},
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
			Feature: &featurepb.Feature{
				Name:    req.Name,
				Enabled: feat,
			},
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
			features = append(features, &featurepb.Feature{
				Name:    name,
				Enabled: feat,
			})
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
		after  = &featurepb.Feature{
			Name:    req.Feature.Name,
			Enabled: req.Feature.Enabled,
		}
	)

	if feat, ok := s.features[req.Feature.Name]; ok {
		before = &featurepb.Feature{
			Name:    req.Feature.Name,
			Enabled: feat,
		}
	}

	s.features[req.Feature.Name] = after.Enabled
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
