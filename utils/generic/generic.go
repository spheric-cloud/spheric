// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package generic

import (
	"fmt"
	"reflect"
	"strings"
)

// Identity is a function that returns its given parameters.
func Identity[E any](e E) E {
	return e
}

// Const produces a function that takes an argument and returns the original argument, ignoring the passed-in value.
func Const[E, F any](e E) func(F) E {
	return func(F) E {
		return e
	}
}

// Zero returns the zero value for the given type.
func Zero[E any]() E {
	var zero E
	return zero
}

func Cast[E any](v any) (E, error) {
	e, ok := v.(E)
	if !ok {
		return Zero[E](), fmt.Errorf("expected %T but got %T", e, v)
	}
	return e, nil
}

func MustCast[E any](v any) E {
	e, err := Cast[E](v)
	if err != nil {
		panic(err)
	}
	return e
}

func ReflectType[E any]() reflect.Type {
	ePtr := (*E)(nil)         // use a pointer to avoid initializing the entire type
	t := reflect.TypeOf(ePtr) // get the pointer type
	return t.Elem()           // return the element type (the actual requested type)
}

// Pointer returns a pointer for the given value.
func Pointer[E any](e E) *E {
	return &e
}

func PointerOrElse[E any](e *E, defaultFunc func() *E) *E {
	if e != nil {
		return e
	}
	return defaultFunc()
}

func PointerOr[E any](e *E, defaultPtr *E) *E {
	return PointerOrElse(e, func() *E {
		return defaultPtr
	})
}

func PointerOrNew[E any](e *E) *E {
	return PointerOrElse(e, func() *E {
		return new(E)
	})
}

// DerefOrElse returns the value e points to if it's non-nil. Otherwise, it returns the result of calling defaultFunc.
func DerefOrElse[E any](e *E, defaultFunc func() E) E {
	if e != nil {
		return *e
	}
	return defaultFunc()
}

// DerefOr returns the value e points to if it's non-nil. Otherwise, it returns the defaultValue.
func DerefOr[E any](e *E, defaultValue E) E {
	return DerefOrElse(e, func() E {
		return defaultValue
	})
}

// DerefOrZero returns the value e points to if it's non-nil. Otherwise, it returns the zero value for type E.
func DerefOrZero[E any](e *E) E {
	if e != nil {
		return *e
	}
	var zero E
	return zero
}

// TODO is a function to create holes when stubbing out more complex mechanisms.
//
// By default, it will panic with 'TODO: provide a value of type <type>' where <type> is the type of V.
// The panic message can be altered by passing in additional args that will be printed as
// 'TODO: <args separated by space>'
func TODO[V any](args ...any) V {
	var sb strings.Builder
	sb.WriteString("TODO: ")
	if len(args) > 0 {
		_, _ = fmt.Fprintln(&sb, args...)
	} else {
		_, _ = fmt.Fprintf(&sb, "provide a value of type %T", Zero[V]())
	}
	panic(sb.String())
}
