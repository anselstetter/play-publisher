package assert

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

func Equals(t *testing.T, got any, want any, prefixes ...string) {
	t.Helper()
	prefix := strings.Join(prefixes, " ")
	if !reflect.DeepEqual(got, want) {
		t.Errorf("\n%s \ngot: %#v \nwant: %#v", prefix, got, want)
	}
}

func True(t *testing.T, test bool, prefixes ...string) {
	t.Helper()
	prefix := strings.Join(prefixes, " ")
	if !test {
		t.Errorf("\n%s \ntest was false", prefix)
	}
}

func NoError(t *testing.T, err error, prefixes ...string) {
	t.Helper()
	prefix := strings.Join(prefixes, " ")
	if err != nil {
		t.Errorf("\n%s \nunexpected error: %v", prefix, err)
	}
}

func AnyError(t *testing.T, err error, prefixes ...string) {
	t.Helper()
	prefix := strings.Join(prefixes, " ")
	if err == nil {
		t.Errorf("\n%s \nexpected error, but got nil", prefix)
	}
}

func IsError(t *testing.T, err error, target error, prefixes ...string) {
	t.Helper()
	prefix := strings.Join(prefixes, " ")
	if !errors.Is(err, target) {
		t.Errorf("\n%s \nexpected %v, but got %v", prefix, target, err)
	}
}
