package literally

import "testing"

func TestTypeName(t *testing.T) {
	type (
		myInt    int
		myStruct struct{}
		mySet    map[myInt]myStruct
	)

	tests := []struct {
		v        interface{}
		expected string
	}{
		{[]int{}, "[]int{}"},
		{[]*int{}, "[]*int{}"},
		{map[*int]struct{}{}, "map[*int]struct{}{}"},
		{[]map[*int]struct{}{}, "[]map[*int]struct{}{}"},
		{[]myInt{}, "[]literally.myInt{}"},
		{myStruct{}, "literally.myStruct{}"},
		{mySet{1: myStruct{}}, "literally.mySet{1: literally.myStruct{}}"},
		{[]interface{}{}, "[]interface{}{}"},
		{[]func(string) error{}, "[]func(string) error{}"},
		{[]func(int){}, "[]func(int){}"},
		{[]func(string, bool) int{}, "[]func(string, bool) int{}"},
		{[]func() (int, string){}, "[]func() (int, string){}"},
	}
	for _, test := range tests {
		actual := Figurative{}.Literally(test.v)
		if test.expected != actual {
			t.Errorf("%v (%T):\nexpected:\n\"%s\"\ngot:\n\"%s\"", test.v, test.v, test.expected, actual)
		}
	}
}
