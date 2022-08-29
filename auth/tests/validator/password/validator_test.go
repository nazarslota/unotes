package password

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/udholdenhed/unotes/auth/pkg/validator/password"
)

func TestNewDefaultValidationFuncs(t *testing.T) {
	type TestInput struct {
		MinPassLength int
		MaxPassLength int
	}

	type TestExpected struct {
		Err error
	}

	type Test struct {
		TestInput
		TestExpected
	}

	tests := []Test{
		{
			TestInput:    TestInput{MinPassLength: 0, MaxPassLength: 0},
			TestExpected: TestExpected{Err: nil},
		},
		{
			TestInput:    TestInput{MinPassLength: 4, MaxPassLength: 0},
			TestExpected: TestExpected{Err: nil},
		},
		{
			TestInput:    TestInput{MinPassLength: 0, MaxPassLength: 8},
			TestExpected: TestExpected{Err: nil},
		},
		{
			TestInput:    TestInput{MinPassLength: 4, MaxPassLength: 8},
			TestExpected: TestExpected{Err: nil},
		},
		{
			TestInput:    TestInput{MinPassLength: 4, MaxPassLength: 4},
			TestExpected: TestExpected{Err: password.ErrInvalidValidationFuncsConfig},
		},
		{
			TestInput:    TestInput{MinPassLength: 8, MaxPassLength: 4},
			TestExpected: TestExpected{Err: password.ErrInvalidValidationFuncsConfig},
		},
	}

	for _, test := range tests {
		_, actual := password.NewDefaultValidationFuncs(test.MinPassLength, test.MaxPassLength)
		assert.Equalf(t, test.TestExpected.Err, actual,
			"MinPassLength: %d. MaxPassLength: %d.", test.MinPassLength, test.MaxPassLength,
		)
	}
}

func TestDefaultValidationFuncs_ValidatePasswordLength(t *testing.T) {
	TestNewDefaultValidationFuncs(t)

	type TestInput struct {
		MinPassLength int
		MaxPassLength int
		Password      string
	}

	type TestExpected struct {
		Err error
	}

	type Test struct {
		TestInput
		TestExpected
	}

	tests := []Test{
		{
			TestInput: TestInput{
				MinPassLength: 4, MaxPassLength: 8,
				Password: "1234",
			},
			TestExpected: TestExpected{Err: nil},
		},
		{
			TestInput: TestInput{
				MinPassLength: 4, MaxPassLength: 8,
				Password: "12345678",
			},
			TestExpected: TestExpected{Err: nil},
		},
		{
			TestInput: TestInput{
				MinPassLength: 4, MaxPassLength: 8,
				Password: "123",
			},
			TestExpected: TestExpected{Err: password.ErrPasswordIsTooShort},
		},
		{
			TestInput: TestInput{
				MinPassLength: 4, MaxPassLength: 8,
				Password: "123456789",
			},
			TestExpected: TestExpected{Err: password.ErrPasswordIsTooLong},
		},
	}

	for _, test := range tests {
		funcs, _ := password.NewDefaultValidationFuncs(test.MinPassLength, test.MaxPassLength)
		actual := funcs.ValidatePasswordLength(test.Password)
		assert.Equalf(t, test.Err, actual,
			"Password: %s. Exepected: %v. Actual: %v", test.Password, test.TestExpected.Err, actual,
		)
	}
}
