package helpers

import (
	"fmt"
	"reflect"
)

func ExpectType[T any](ref any) T {

	expected := reflect.TypeOf((*T)(nil)).Elem()
	actual := reflect.TypeOf(ref)

	if expected!= actual {
		panic(fmt.Sprintf("expected type %s, got %s", expected, actual))
	}

	return ref.(T)
}

//func ExpectTypeAnyOf
// example expects the type to be either int or string or float, target is the value to be checked
func ExpectTypeAnyOf(target any, types ...interface{}) {
	for _, t := range types {
		if reflect.TypeOf(target) == reflect.TypeOf(t) {
			return
		}
	}
	panic(fmt.Sprintf("expected type %s, got %s", types, reflect.TypeOf(target)))
}

func TypesMatch(a any, b any) bool {
	return reflect.TypeOf(a) == reflect.TypeOf(b)
}

func TypesMatchT[T any](args ...any) bool {
	expectedType := reflect.TypeOf((*T)(nil)).Elem()

	for _, t := range args {
		if expectedType != reflect.TypeOf(t) {
			return false
		}
	}
	return true
}

func ContainsIn(items []string, targets ...string) bool {
	for _, target := range targets {
		for _, item := range items {
			if item == target {
				return true
			}
		}
	}
	return false
}