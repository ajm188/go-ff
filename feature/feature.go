package feature

import (
	"context"
	"errors"

	featurepb "github.com/ajm188/goff/proto/feature"
)

var ErrNoFeature = errors.New("no such feature")

// Get returns whether a feature is enabled or not.
func Get(name string) (bool, error) {
	feat, err := inst.GetFeature(context.Background(), &featurepb.GetFeatureRequest{Name: name})
	if err != nil {
		return false, err
	}

	return feat.Feature.Enabled, nil
}
