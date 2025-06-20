package utils

import (
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// GetValuePointer converts google.protobuf.StringValue to *string
func GetValuePointer(sv *wrapperspb.StringValue) *string {
	if sv == nil {
		return nil
	}
	val := sv.GetValue()
	return &val
}

// ToStringInterfaceMap converts map<string, string> from proto to map[string]any
func ToStringInterfaceMap(m map[string]string) map[string]any {
	if m == nil {
		return nil
	}
	res := make(map[string]any)
	for k, v := range m {
		res[k] = v
	}
	return res
}
