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

func GetStringValue(v *string) *wrapperspb.StringValue {
	if v == nil {
		return nil
	}
	return wrapperspb.String(*v)
}

func GetInterfaceString(m map[string]any) map[string]string {
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

func GetStringInterface(m map[string]string) map[string]any {
	if m == nil {
		return nil
	}
	res := make(map[string]any)
	for k, v := range m {
		res[k] = v
	}
	return res
}
