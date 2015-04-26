package literally

import (
	"fmt"
	"path"
	"reflect"
	"strings"
)

func (f Figurative) typeName(typ reflect.Type) string {
	if typ.PkgPath() != "" {
		// Custom type
		return f.Qualifier(typ.PkgPath()) + typ.Name()
	}
	switch kind := typ.Kind(); kind {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128:
		return typ.Name()
	case reflect.Array:
		return fmt.Sprintf("[%d]%s", typ.Len(), f.typeName(typ.Elem()))
	case reflect.Chan:
		return "chan " + f.typeName(typ.Elem())
	case reflect.Func:
		var (
			in  = make([]string, typ.NumIn())
			out = make([]string, typ.NumOut())
		)
		for i := 0; i < typ.NumIn(); i++ {
			in[i] = f.typeName(typ.In(i))
		}
		for i := 0; i < typ.NumOut(); i++ {
			out[i] = f.typeName(typ.Out(i))
		}

		var (
			instring  = strings.Join(in, ", ")
			outstring = strings.Join(out, ", ")
		)
		switch len(out) {
		case 0:
			return fmt.Sprintf("func(%s)", instring)
		case 1:
			return fmt.Sprintf("func(%s) %s", instring, outstring)
		default:
			return fmt.Sprintf("func(%s) (%s)", instring, outstring)
		}
	case reflect.Interface:
		if typ.Name() != "" {
			return typ.Name()
		}
		return "interface{}"
	case reflect.Map:
		return fmt.Sprintf("map[%s]%s", f.typeName(typ.Key()), f.typeName(typ.Elem()))
	case reflect.Ptr:
		return "*" + f.typeName(typ.Elem())
	case reflect.Slice:
		return "[]" + f.typeName(typ.Elem())
	case reflect.String:
		return typ.Name()
	case reflect.Struct:
		if typ.Name() == "" {
			// TODO: Allow better unnamed structs?
			return "struct{}"
		}
		return f.Qualifier(typ.PkgPath()) + typ.Name()
	case reflect.UnsafePointer:
		return f.Qualifier("unsafe") + "Pointer"
	default:
		// Indicates change to reflect
		panic("Unable to get type name")
	}
}

// Qualifier returns the given package's package qualifier (with trailing ".") as
// determined by the Figurative's PkgNames
func (f Figurative) Qualifier(pkg string) string {
	q, ok := f.PkgNames[pkg]
	if !ok {
		_, q = path.Split(pkg)
	}
	if q == "" {
		return ""
	}
	return q + "."
}
