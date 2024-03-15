package testutil

import (
	"reflect"
	"testing"
)

func AssertNotNil(actual any, t *testing.T, objectName string) {
	if isNil(actual) {
		t.Errorf("Element %s should not be nil", objectName)
	}
}

func AssertNil(actual any, t *testing.T, objectName string) {
	if !isNil(actual) {
		t.Errorf("Element %s should be nil, but has value '%v'", objectName, actual)
	}
}

func AssertEquals(expected any, actual any, t *testing.T, objectName string) {
	AssertNotNil(actual, t, objectName)
	if expected != actual {
		t.Errorf("Element %s is not as expected. expected: '%v' actual: '%v'", objectName, expected, actual)
	}
}

func AssertNotEquals(notExpected any, actual any, t *testing.T, objectName string) {
	AssertNotNil(actual, t, objectName)
	if notExpected == actual {
		t.Errorf("Element %s equals the unexpected. unexpected: '%v' actual: '%v'", objectName, notExpected, actual)
	}
}

func isNil(toCheck any) bool {
	v := reflect.ValueOf(toCheck)
	k := v.Kind()

	switch k {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Pointer, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		return v.IsNil()
	default:
		return toCheck == nil
	}
}
