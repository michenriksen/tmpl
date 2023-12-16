package testutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ErrorAssertion is a function that asserts an error.
//
// The function is passed the test context and the error to assert. If the error
// does not match the assertion, the test is failed and function returns false.
type ErrorAssertion func(testing.TB, error) bool

// AssertErrorIs asserts that the given error is the target error.
func AssertErrorIs(target error) ErrorAssertion {
	return func(tb testing.TB, err error) bool {
		tb.Helper()
		return assert.ErrorIs(tb, err, target)
	}
}

// AssertErrorContains asserts that the given error contains the provided
// substring.
func AssertErrorContains(substr string) ErrorAssertion {
	return func(tb testing.TB, err error) bool {
		tb.Helper()
		return assert.ErrorContains(tb, err, substr)
	}
}

// RequireErrorIs is like [AssertErrorIs], but fails the test immediately if
// the assertion fails.
func RequireErrorIs(target error) ErrorAssertion {
	return func(tb testing.TB, err error) bool {
		tb.Helper()

		if !AssertErrorIs(target)(tb, err) {
			tb.FailNow()
			return false
		}

		return true
	}
}

// RequireErrorContains is like [AssertErrorContains], but fails the test
// immediately if the assertion fails.
func RequireErrorContains(substr string) ErrorAssertion {
	return func(tb testing.TB, err error) bool {
		tb.Helper()

		if !AssertErrorContains(substr)(tb, err) {
			tb.FailNow()
			return false
		}

		return true
	}
}
