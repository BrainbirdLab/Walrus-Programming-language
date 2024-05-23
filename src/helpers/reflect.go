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