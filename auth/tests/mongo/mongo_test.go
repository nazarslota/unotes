package mongo

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/udholdenhed/unotes/auth/pkg/mongo"
)

func TestBuildURI(t *testing.T) {
	type TestInput struct {
		Host     string
		Port     string
		Username string
		Password string
	}

	type TestExpected struct {
		URI string
	}

	type Test struct {
		TestInput
		TestExpected
	}

	tests := []Test{
		{
			TestInput: TestInput{
				Host:     "localhost",
				Port:     "27017",
				Username: "",
				Password: "",
			},
			TestExpected: TestExpected{
				URI: "mongodb://localhost:27017",
			},
		},
		{
			TestInput: TestInput{
				Host:     "localhost",
				Port:     "27017",
				Username: "username",
				Password: "",
			},
			TestExpected: TestExpected{
				URI: "mongodb://username@localhost:27017",
			},
		},
		{
			TestInput: TestInput{
				Host:     "localhost",
				Port:     "27017",
				Username: "username",
				Password: "password",
			},
			TestExpected: TestExpected{
				URI: "mongodb://username:password@localhost:27017",
			},
		},
		{
			TestInput: TestInput{
				Host:     "localhost",
				Port:     "",
				Username: "",
				Password: "",
			},
			TestExpected: TestExpected{
				URI: "mongodb+srv://localhost",
			},
		},
		{
			TestInput: TestInput{
				Host:     "localhost",
				Port:     "",
				Username: "username",
				Password: "",
			},
			TestExpected: TestExpected{
				URI: "mongodb+srv://username@localhost",
			},
		},
		{
			TestInput: TestInput{
				Host:     "localhost",
				Port:     "",
				Username: "username",
				Password: "password",
			},
			TestExpected: TestExpected{
				URI: "mongodb+srv://username:password@localhost",
			},
		},
	}

	for _, test := range tests {
		actual := mongo.BuildURI(
			test.TestInput.Host,
			test.TestInput.Port,
			test.TestInput.Username,
			test.TestInput.Password,
		)
		assert.Equal(t, test.TestExpected.URI, actual)
	}
}
