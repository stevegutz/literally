package literally

import (
	"sync"
	"testing"
	"time"
)

var literallyTests = []struct {
	v        interface{}
	expected string
}{
	{nil, "nil"},
	{[]int{1, 2, 3, 4}, "[]int{1, 2, 3, 4}"},
	{[4]int{1, 2, 3, 4}, "[4]int{1, 2, 3, 4}"},
	// Maps with multiple entries are not stable for testing
	//
	// {map[int]string{1: "a", 2: "b"}, `map[int]string{1: "a", 2: "b"}`},
	// {map[int][]int64{1: {10, 11, 12}, 2: {20, 30, 40}},
	// 	"map[int][]int64{1: []int64{10, 11, 12}, 2: []int64{20, 30, 40}}"},
	{map[int]string{1: "a"}, `map[int]string{1: "a"}`},
	{map[int][]int64{1: []int64{10, 11, 12}}, "map[int][]int64{1: []int64{10, 11, 12}}"},
	// Channels
	{make(chan int), "make(chan int, 0)"},
	{make(chan []*string, 42), "make(chan []*string, 42)"},
	// Structs
	{struct{}{}, "struct{}{}"},
	{testStruct{Int: 23, String: "str", IntSlice: []int{1, 2, 3}},
		`literally.testStruct{Int: 23, String: "str", IntSlice: []int{1, 2, 3}}`},
	{&testStruct{}, `&literally.testStruct{Int: 0, String: "", IntSlice: []int{}}`},
	{sync.Mutex{}, "sync.Mutex{}"},
	{&sync.Mutex{}, "&sync.Mutex{}"},
	{testEmbeddedStruct{Bool: true},
		`literally.testEmbeddedStruct{Mutex: sync.Mutex{}, testStruct: literally.testStruct{Int: 0, String: "", IntSlice: []int{}}, Bool: true}`},
	// TODO Fix runes
	// {map[rune]string{'a': "a"}, `map[rune]string{'a': "a"}`},
}

type testStruct struct {
	Int      int
	String   string
	IntSlice []int
}

type testEmbeddedStruct struct {
	sync.Mutex
	testStruct
	Bool bool
}

func TestLiterally(t *testing.T) {
	for _, test := range literallyTests {
		actual := Literally(test.v)
		if test.expected != actual {
			t.Errorf("%v (%T):\nexpected:\n\"%s\"\ngot:\n\"%s\"", test.v, test.v, test.expected, actual)
		}
	}
}

func TestConstructorProviders(t *testing.T) {
	f := NewFigurative()
	// Note: do not use time pointers in real code
	f.ConstructorProviders[TypeKey{"time", "Time", true}] = func(f Figurative, v interface{}) string {
		return "&" + f.Qualifier("time") + "Time{}"
	}

	var tests = []struct {
		v        interface{}
		expected string
	}{
		{time.Date(2015, time.April, 14, 14, 56, 3, 55635788, time.UTC),
			"time.Date(2015, time.April, 14, 14, 56, 3, 55635788, time.UTC)"},
		{time.Time{}, "time.Time{}"},
		{&time.Time{}, "&time.Time{}"},
	}
	for _, test := range tests {
		actual := f.Literally(test.v)
		if test.expected != actual {
			t.Errorf("%v (%T):\nexpected:\n\"%s\"\ngot:\n\"%s\"", test.v, test.v, test.expected, actual)
		}
	}
}
