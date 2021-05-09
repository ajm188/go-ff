package feature

import (
	"fmt"
	"strings"

	featurepb "github.com/ajm188/go-ff/proto/feature"
)

// ParseType converts a string into a featurepb.Feature_Type enum, returning an
// error if the lowercased input name is not in the enum mapping.
func ParseType(s string) (featurepb.Feature_Type, error) {
	if t, ok := featurepb.Feature_Type_value[strings.ToUpper(s)]; ok {
		return featurepb.Feature_Type(t), nil
	}

	return featurepb.Feature_UNKNOWN, fmt.Errorf("%w %s", ErrUnknownFeatureType, s)
}
