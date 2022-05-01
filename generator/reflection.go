package generator

import (
	"fmt"
	"log"
	"reflect"

	"github.com/peterzeller/go-fun/zero"
)

// ErrUnsupportedType is returned if a type is not supported by a ReflectionGenFun
var ErrUnsupportedType = fmt.Errorf("unsupported type")

// ReflectionGenFun is a function to create a generator for the given type.
// Returns error ErrUnsupportedType if the type is not supported and other errors if the generator cannot be created.
type ReflectionGenFun func(t reflect.Type, opts *ReflectionGeneratorOptions) (UntypedGenerator, error)

// ReflectionGeneratorOptions contains options for reflection-based generators.
type ReflectionGeneratorOptions struct {
	generators []ReflectionGenFun
}

// Register a generator in the options.
// Generators that are added later will overwrite previous generators.
func (r *ReflectionGeneratorOptions) Register(genFunc ReflectionGenFun) {
	r.generators = append(r.generators, genFunc)
}

// RegisterConstructor registers a constructor for a custom type.
func (r *ReflectionGeneratorOptions) RegisterConstructor(constructorFun interface{}) {
	g, err := r.generatorFromConstructor(constructorFun)
	if err != nil {
		panic(err)
	}
	returnType := reflect.TypeOf(constructorFun).Out(0)
	f := func(t reflect.Type, opts *ReflectionGeneratorOptions) (UntypedGenerator, error) {
		if !returnType.AssignableTo(t) {
			log.Printf("Constructor did not match: %v != %v", returnType, t)
			return nil, ErrUnsupportedType
		}
		return g, nil
	}
	r.Register(f)
}

func (r *ReflectionGeneratorOptions) generatorFromConstructor(constructorFun interface{}) (UntypedGenerator, error) {
	v := reflect.ValueOf(constructorFun)
	t := reflect.TypeOf(constructorFun)
	if t.Kind() != reflect.Func {
		return nil, fmt.Errorf("constructorFun is not a function")
	}
	switch t.NumOut() {
	case 1:
		// ok
	case 2:
		errType := t.Out(1)
		if errType != reflect.TypeOf(error(nil)) {
			return nil, fmt.Errorf("constructorFun returns two values, but second is not error")
		}
	default:
		return nil, fmt.Errorf("constructorFun must return a single value, or a single value and an error")
	}
	// if we have no parameters, generate constant generator
	if t.NumIn() == 0 {
		results := v.Call([]reflect.Value{})
		return ToUntyped(Constant(results[0].Interface())), nil
	}

	// create generator for arguments
	gens := make([]Generator[reflect.Value], t.NumIn())
	for i := 0; i < t.NumIn(); i++ {
		g, err := buildGenerator(r, t.In(i))
		if err != nil {
			return nil, fmt.Errorf("cannot create generator for parameter %d: %w", i, err)
		}
		gens[i] = Map(ToTypedGenerator[interface{}](g),
			func(v interface{}) reflect.Value {
				return reflect.ValueOf(v)
			})
	}
	// combine gens

	g := Map(gens[0], func(value reflect.Value) []reflect.Value {
		values := make([]reflect.Value, 1, len(gens))
		values[0] = value
		return values
	})
	for i := 1; i < len(gens); i++ {
		g = Zip(g, gens[i], func(ar []reflect.Value, v reflect.Value) []reflect.Value {
			return append(ar, v)
		})
	}
	cg := FilterMap(g,
		func(args []reflect.Value) (interface{}, bool) {
			results := v.Call(args)
			if len(results) == 2 {
				if !results[1].IsNil() {
					return reflect.Value{}, false
				}
			}
			return results[0].Interface(), true
		})
	return ToUntyped(cg), nil
}

// ReflectionGenDefaultOpts returns default options for reflection-based generators.
// It contains generators for basic types, slices, maps, and structs.
func ReflectionGenDefaultOpts() *ReflectionGeneratorOptions {
	r := &ReflectionGeneratorOptions{}
	r.Register(reflectionGenBasicTypes)
	return r
}

// ReflectionGen uses reflection to return a generator for the generic type T.
// To support custom types, pass in ReflectionGeneratorOptions with registered generators for those types.
func ReflectionGen[T any](opts *ReflectionGeneratorOptions) Generator[T] {
	typ := reflect.TypeOf(zero.Value[T]())
	gen, err := buildGenerator(opts, typ)
	if err != nil {
		panic(err)
	}
	return ToTypedGenerator[T](gen)
}

func buildGenerator(opts *ReflectionGeneratorOptions, typ reflect.Type) (UntypedGenerator, error) {
	for i := len(opts.generators) - 1; i >= 0; i-- {
		genFunc := opts.generators[i]
		gen, err := genFunc(typ, opts)
		log.Printf("Trying generator %v for %v: %v", i, typ, err)
		if err == nil {
			return gen, nil
		}
	}
	return nil, fmt.Errorf("no generator found for type %s", typ)
}

func reflectionGenBasicTypes(t reflect.Type, opts *ReflectionGeneratorOptions) (UntypedGenerator, error) {
	switch t.Kind() {
	case reflect.Invalid:
		return nil, ErrUnsupportedType
	case reflect.Bool:
		return ToUntyped(Bool()), nil
	case reflect.Int:
		return ToUntyped(Int()), nil
	case reflect.Int8:
		return nil, ErrUnsupportedType
	case reflect.Int16:
		return nil, ErrUnsupportedType
	case reflect.Int32:
		return nil, ErrUnsupportedType
	case reflect.Int64:
		return nil, ErrUnsupportedType
	case reflect.Uint:
		return nil, ErrUnsupportedType
	case reflect.Uint8:
		return nil, ErrUnsupportedType
	case reflect.Uint16:
		return nil, ErrUnsupportedType
	case reflect.Uint32:
		return nil, ErrUnsupportedType
	case reflect.Uint64:
		return nil, ErrUnsupportedType
	case reflect.Uintptr:
		return nil, ErrUnsupportedType
	case reflect.Float32:
		return nil, ErrUnsupportedType
	case reflect.Float64:
		return nil, ErrUnsupportedType
	case reflect.Complex64:
		return nil, ErrUnsupportedType
	case reflect.Complex128:
		return nil, ErrUnsupportedType
	case reflect.Array:
		return nil, ErrUnsupportedType
	case reflect.Chan:
		return nil, ErrUnsupportedType
	case reflect.Func:
		return nil, ErrUnsupportedType
	case reflect.Interface:
		return nil, ErrUnsupportedType
	case reflect.Map:
		return nil, ErrUnsupportedType
	case reflect.Pointer:
		return nil, ErrUnsupportedType
	case reflect.Slice:
		return sliceGenerator(t, opts)
	case reflect.String:
		return ToUntyped(String()), nil
	case reflect.Struct:
		return structGenerator(t, opts)
	case reflect.UnsafePointer:
		return nil, ErrUnsupportedType
	default:
		return nil, ErrUnsupportedType
	}
}

func sliceGenerator(t reflect.Type, opts *ReflectionGeneratorOptions) (UntypedGenerator, error) {
	if t.Kind() != reflect.Slice {
		return nil, fmt.Errorf("not a slice: %v", t)
	}
	// get slice element type
	elemType := t.Elem()
	// create generator for slice element type
	elemGen, err := buildGenerator(opts, elemType)
	if err != nil {
		return nil, fmt.Errorf("generating element generator for %v slice: %w", elemType, err)
	}
	sliceGen := Map(Slice(ToTypedGenerator[interface{}](elemGen)),
		func(ar []interface{}) interface{} {
			r := reflect.MakeSlice(t, len(ar), len(ar))
			for i, v := range ar {
				r.Index(i).Set(reflect.ValueOf(v))
			}
			return r.Interface()
		})
	return ToUntyped(sliceGen), nil
}

func structGenerator(t reflect.Type, opts *ReflectionGeneratorOptions) (UntypedGenerator, error) {
	fields := reflect.VisibleFields(t)
	if len(fields) == 0 {
		zeroVal := reflect.New(t).Elem()
		return ToUntyped(Constant[interface{}](zeroVal)), nil
	}

	fieldGens := make([]UntypedGenerator, len(fields))
	for i, field := range fields {
		g, err := reflectionGenBasicTypes(field.Type, opts)
		if err != nil {
			return nil, fmt.Errorf("error creating generator for struct %s, field %s: %s", t.Name(), field.Name, err)
		}
		fieldGens[i] = g
	}
	g := Map(ToTypedGenerator[interface{}](fieldGens[0]), func(value interface{}) []interface{} {
		values := make([]interface{}, 1, len(fields))
		values[0] = value
		return values
	})
	for i := 1; i < len(fieldGens); i++ {
		g = Zip(g, ToTypedGenerator[interface{}](fieldGens[i]), func(ar []interface{}, v interface{}) []interface{} {
			return append(ar, v)
		})
	}
	sGen := Map(g, func(fieldValues []interface{}) interface{} {
		val := reflect.New(t).Elem()
		for i, v := range fieldValues {
			val.Field(i).Set(reflect.ValueOf(v))
		}
		return val.Interface()
	})
	return ToUntyped(sGen), nil
}
