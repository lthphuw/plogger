package plogger

import (
	"reflect"

	"github.com/lthphuw/pdefault"
)

// Setter function.
type Setter[T any] func(*T) error

// Lister for list all setter functions.
type Lister[T any] interface {
	List() []Setter[T]
}

func newOptions[T any](opts ...Lister[T]) (*T, error) {
	args := new(T)

	for _, opt := range opts {
		if opt == nil || reflect.ValueOf(opt).IsNil() {
			// Do nothing if the option is nil or if opt is nil.
			continue
		}

		for _, setArgs := range opt.List() {
			if setArgs == nil {
				continue
			}

			if err := setArgs(args); err != nil {
				return nil, err
			}
		}
	}

	// Apply defaults
	if err := pdefault.Defaults(args); err != nil {
		return nil, err
	}

	return args, nil
}
