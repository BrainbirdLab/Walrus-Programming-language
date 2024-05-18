package helpers

import (
	"fmt"
	"reflect"
)

func ExpectType[T any](ref any) T {
	expectedType := reflect.TypeOf((*T)(nil)).Elem()
	recievedType := reflect.TypeOf(ref)

	if recievedType != expectedType {
		panic(fmt.Sprintf("Expected %T but instead recived %T inside ExpectType[T](r)\n", expectedType, recievedType))
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