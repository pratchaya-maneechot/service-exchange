package grpc

import (
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// StringValueToPtr unwraps a *wrapperspb.StringValue into a *string.
// Returns nil if the input *wrapperspb.StringValue is nil.
func StringValueToPtr(sv *wrapperspb.StringValue) *string {
	if sv == nil {
		return nil
	}
	val := sv.GetValue()
	return &val
}

// PtrToStringValue wraps a *string into a *wrapperspb.StringValue.
// Returns nil if the input *string is nil.
func PtrToStringValue(s *string) *wrapperspb.StringValue {
	if s == nil {
		return nil
	}
	return wrapperspb.String(*s)
}

// AnyMapToStringMap converts a map[string]any to a map[string]string,
// including only values that are actually strings.
func AnyMapToStringMap(m map[string]any) map[string]string {
	if m == nil {
		return nil
	}
	res := make(map[string]string)
	for k, v := range m {
		if _v, ok := v.(string); ok {
			res[k] = _v
		}
	}
	return res
}

// StringMapToAnyMap converts a map[string]string to a map[string]any.
func StringMapToAnyMap(m map[string]string) map[string]any {
	if m == nil {
		return nil
	}
	res := make(map[string]any)
	for k, v := range m {
		res[k] = v
	}
	return res
}
