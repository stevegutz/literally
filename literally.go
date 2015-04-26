/*
Package literally allows for the serialization of values to their literal representation.

You can use literally by calling literally.Literally
	literally.Literally(myValue)
or by using a figurative instance:
	f := literally.NewFigurative()
	// Set-up figurative
	f.Literally(myValue)
The default top-level Figurative is defined by literally.ConFigurative
*/
package literally

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// ConFigurative is the active configuration used by top-level functions
var ConFigurative = NewFigurative()

// DefaultFigurative returns the default Figurative instance
func NewFigurative() Figurative {
	return Figurative{
		PkgNames: map[string]string{},
		ConstructorProviders: map[TypeKey]ConstructorProvider{
			TypeKey{"time", "Time", false}: timeConstructor,
		},
		Panic: false,
	}
}

// Figurative is used to configure literalizations
type Figurative struct {
	// PkgNames allows aliasing packages, with "" meaning that the mapped package
	// should be referenced without qualification.
	PkgNames map[string]string
	// ConstructorProviders gives a mapping from types to special functions for literalizing them
	ConstructorProviders map[TypeKey]ConstructorProvider
	// Panic indicates whether to panic when a type is unsupported.
	// TODO: Add a better error pattern
	Panic bool
}

// TypeKey defines a [pointer to a] type
type TypeKey struct {
	Package string
	Name    string
	// No pointer to pointer for now (-ever?)
	Pointer bool
}

// ConstructorProvider returns a constructor for a given type
type ConstructorProvider func(Figurative, interface{}) string

func timeConstructor(f Figurative, v interface{}) string {
	var (
		t = v.(time.Time).In(time.UTC)
		q = f.Qualifier("time")
	)
	if t.IsZero() {
		return q + "Time{}"
	}
	return fmt.Sprintf("%sDate(%d, %s%s, %d, %d, %d, %d, %d, %sUTC)",
		q, t.Year(), q, t.Month().String(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), q)
}

// Literally returns a string representing the provided interface{} as a literal value
func Literally(v interface{}) string {
	return ConFigurative.Literally(v)
}

// Literally returns a string representing the provided interface{} as a literal value
func (f Figurative) Literally(v interface{}) string {
	if v == nil {
		return "nil"
	}

	// Check for a ConstructorProvider
	var (
		value = reflect.ValueOf(v)
		typ   = value.Type()
		kind  = value.Kind()
		key   TypeKey
	)
	if kind == reflect.Ptr {
		// Won't work on Ptr to Ptr
		elemTyp := value.Elem().Type()
		key = TypeKey{elemTyp.PkgPath(), elemTyp.Name(), true}
	} else {
		key = TypeKey{typ.PkgPath(), typ.Name(), false}
	}
	if cp, ok := f.ConstructorProviders[key]; ok {
		return cp(f, v)
	}

	// TODO: Currently uint8 always handles byte and uint32 always handles rune
	switch kind {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		// TODO: It could be nice to special case complex (to format better and handle Inf, NaN, et al.)
		reflect.Complex64, reflect.Complex128:
		return fmt.Sprint(v)
	case reflect.Array, reflect.Slice:
		return f.arrayAndSliceHelper(value)
	case reflect.Chan:
		// TODO: Use value.Type().ChanDir()
		return fmt.Sprintf("make(%s, %d)", value.Type(), value.Cap())
	case reflect.Func:
	case reflect.Interface:
	case reflect.Map:
		return f.mapHelper(value)
	case reflect.Ptr:
		return f.pointerHelper(value)
	case reflect.String:
		return strconv.Quote(v.(string))
	case reflect.Struct:
		return f.structHelper(value)
	case reflect.UnsafePointer:
	}
	if f.Panic {
		panic(fmt.Sprintf("Unable to handle type: %T", v))
	} else {
		return "nil"
	}
}

func (f Figurative) mapHelper(value reflect.Value) string {
	entries := make([]string, value.Len())
	for i, key := range value.MapKeys() {
		entries[i] = f.Literally(key.Interface()) + ": " + f.Literally(value.MapIndex(key).Interface())
	}
	return fmt.Sprintf("%s{%s}", f.typeName(value.Type()), strings.Join(entries, ", "))
}

func (f Figurative) arrayAndSliceHelper(value reflect.Value) string {
	entries := make([]string, value.Len())
	for i := 0; i < value.Len(); i++ {
		entries[i] = f.Literally(value.Index(i).Interface())
	}
	return fmt.Sprintf("%s{%s}", f.typeName(value.Type()), strings.Join(entries, ", "))
}

func (f Figurative) structHelper(value reflect.Value) string {
	// TODO: Allow better unnamed structs?
	if value.Type().Name() == "" {
		return "struct{}{}"
	}
	var (
		entries = make([]string, 0, value.NumField())
		typ     = value.Type()
	)
	for i := 0; i < value.NumField(); i++ {
		if typ.Field(i).PkgPath != "" {
			// Skip unexported fields
			continue
		}
		var (
			n = typ.Field(i).Name
			v = f.Literally(value.Field(i).Interface())
		)
		entries = append(entries, n+": "+v)
	}
	return fmt.Sprintf("%s{%s}", f.typeName(typ), strings.Join(entries, ", "))
}

func (f Figurative) pointerHelper(value reflect.Value) string {
	var (
		elem = value.Elem()
		kind = elem.Kind()
	)
	switch kind {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128,
		reflect.String:
		return fmt.Sprintf("%s%sPtr(%s)",
			f.Qualifier(reflect.TypeOf(f).PkgPath()),
			strings.Title(kind.String()),
			f.Literally(elem.Interface()))
	case reflect.Struct:
		return "&" + f.Literally(elem.Interface())
	}
	if f.Panic {
		panic(fmt.Sprintf("Unsupported type *%T", elem.Interface()))
	} else {
		return "nil"
	}
}
