package literally

import (
	"reflect"
	"testing"
)

func TestPtr(t *testing.T) {
	tests := []struct {
		v        interface{}
		expected string
	}{
		{BoolPtr(true), "BoolPtr(true)"},
		{IntPtr(42), "IntPtr(42)"},
		{Int8Ptr(42), "Int8Ptr(42)"},
		{Int16Ptr(42), "Int16Ptr(42)"},
		{Int32Ptr(42), "Int32Ptr(42)"},
		{Int64Ptr(42), "Int64Ptr(42)"},
		{UintPtr(42), "UintPtr(42)"},
		{Uint8Ptr(42), "Uint8Ptr(42)"},
		{Uint16Ptr(42), "Uint16Ptr(42)"},
		{Uint32Ptr(42), "Uint32Ptr(42)"},
		{Uint64Ptr(42), "Uint64Ptr(42)"},
		{Float32Ptr(4.2), "Float32Ptr(4.2)"},
		{Float64Ptr(4.2), "Float64Ptr(4.2)"},
		{Complex64Ptr((-5 + 12i)), "Complex64Ptr((-5+12i))"},
		{Complex128Ptr((-5 + 12i)), "Complex128Ptr((-5+12i))"},
		{StringPtr("some string"), `StringPtr("some string")`},
	}
	f := Figurative{PkgNames: make(map[string]string)}
	f.PkgNames[reflect.TypeOf(f).PkgPath()] = ""
	for _, test := range tests {
		actual := f.Literally(test.v)
		if test.expected != actual {
			t.Errorf("%v (%T):\nexpected:\n\"%s\"\ngot:\n\"%s\"", test.v, test.v, test.expected, actual)
		}
	}
}
