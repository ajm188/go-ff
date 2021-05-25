package feature

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"

	"github.com/Knetic/govaluate"
	"github.com/golang/protobuf/jsonpb"

	featurepb "github.com/ajm188/go-ff/proto/feature"
)

var (
	ErrInvalidFeature     = errors.New("invalid feature spec")
	ErrNoFeature          = errors.New("no such feature")
	ErrUnknownFeatureType = errors.New("unknown feature type")
)

// Feature wraps an underlying Feature protobuf message.
type Feature struct {
	*featurepb.Feature
	expr *govaluate.EvaluableExpression
}

// IsEnabled returns whether the given feature is enabled. It returns an error
// either if the feature has an unknown type, or if it is an EXPRESSION feature
// and an error was encountered during expression evaluation.
//
// This function is equivalent to calling f.IsEnabledForParameters(nil).
func (f *Feature) IsEnabled() (bool, error) {
	return f.IsEnabledForParameters(nil)
}

// IsEnabled returns whether the given feature is enabled for the given
// parameters. It returns an error either if the feature has an unknown type,
// or if it is an EXPRESSION feature and an error was encountered during
// expression evaluation.
func (f *Feature) IsEnabledForParameters(parameters map[string]interface{}) (bool, error) {
	switch f.Type {
	case featurepb.Feature_CONSTANT:
		return f.Enabled, nil
	case featurepb.Feature_PERCENTAGE_BASED:
		n := rand.Intn(100)
		return uint32(n) < f.Percentage, nil
	case featurepb.Feature_EXPRESSION:
		if err := f.parseExpression(); err != nil {
			return false, err
		}

		result, err := f.expr.Evaluate(parameters)
		if err != nil {
			return false, err
		}

		v, ok := result.(bool)
		if !ok {
			return false, fmt.Errorf("expression %v did not return a bool: %v", f.expr, result)
		}

		return v, nil
	}

	return false, fmt.Errorf("%w %v for %s", ErrUnknownFeatureType, f.Type, f.Name)
}

func (f *Feature) Validate() (bool, error) {
	switch f.Type {
	case featurepb.Feature_CONSTANT:
		return true, nil
	case featurepb.Feature_PERCENTAGE_BASED:
		if f.Percentage < 0 || f.Percentage > 100 {
			return false, fmt.Errorf("%w: percentage must be in [0, 100] (got %d)", ErrInvalidFeature, f.Percentage)
		}

		return true, nil
	case featurepb.Feature_EXPRESSION:
		if f.Expression == "" {
			return false, fmt.Errorf("%w: expression cannot be empty", ErrInvalidFeature)
		}

		return true, nil
	}

	return false, fmt.Errorf("%w: %v", ErrUnknownFeatureType, f.Type)
}

// MarshalJSON implements json.Marshaler for Feature. It marshals only the
// underlying protobuf message.
func (f *Feature) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	m := jsonpb.Marshaler{
		EnumsAsInts: false,
		Indent:      "    ",
	}

	if err := m.Marshal(buf, f.Feature); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// UnmarshalJSON implements json.Unmarshaler for Feature. It unmarshals the
// underlying protobuf message, and, if the Feature is an EXPRESSION type,
// parses the expression string as well.
func (f *Feature) UnmarshalJSON(data []byte) error {
	if f.Feature == nil {
		f.Feature = &featurepb.Feature{}
	}

	m := jsonpb.Unmarshaler{}
	if err := m.Unmarshal(bytes.NewBuffer(data), f.Feature); err != nil {
		return err
	}

	if f.Type == featurepb.Feature_EXPRESSION {
		return f.parseExpression()
	}

	return nil
}

// parseExpression parses the feature's expression string. The resulting
// expression is reused, so parseExpression will be a no-op on subsequent calls.
//
// Callers must call this function before attempting to use f.expr. Note that
// f.UnmarshalJSON calls this function if the feature is of type EXPRESSION.
func (f *Feature) parseExpression() (err error) {
	if f.expr != nil {
		return nil
	}

	f.expr, err = govaluate.NewEvaluableExpression(f.Expression)
	return err
}

// Get returns whether a feature is enabled or not.
func Get(name string, parameters map[string]interface{}) (bool, error) {
	feat, err := inst.getFeature(name)
	if err != nil {
		return false, err
	}

	return feat.IsEnabledForParameters(parameters)
}
