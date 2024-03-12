package testutil

import (
	"testing"
)

func AssertNotNil(actual any, t *testing.T, objectName string) {
	if actual == nil {
		t.Errorf("Element %s should not be nil", objectName)
	}
}

func AssertNil(actual any, t *testing.T, objectName string) {
	if actual != nil {
		t.Errorf("Element %s should be nil, but has value %v", objectName, actual)
	}
}

func AssertEquals(expected any, actual any, t *testing.T, objectName string) {
	AssertNotNil(actual, t, objectName)
	if expected != actual {
		t.Errorf("Element %s is not as expected. expected: %v actual: %v", objectName, expected, actual)
	}
}

func AssertNotEquals(notExpected *any, actual *any, t *testing.T, objectName string) {
	AssertNotNil(actual, t, objectName)
	if *notExpected == *actual {
		t.Errorf("Element %s equals the unexpected. unexpected: %v actual: %v", objectName, *notExpected, *actual)
	}
}
