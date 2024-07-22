package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:ss|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			User{
				ID:     "00001",
				Age:    41,
				Email:  "user41@test.mail",
				Role:   "admin",
				Phones: []string{"88005353535"},
			},
			ErrLength,
		},

		{
			App{
				Version: "0.1.9",
			},
			nil,
		},

		{
			App{
				Version: "1.0.2387",
			},
			ErrLength,
		},

		{
			Token{
				[]byte{0, 0, 0, 0, 1},
				[]byte{0, 0, 0, 0, 2},
				[]byte{0, 0, 0, 0, 3},
			},
			nil,
		},

		{
			Response{
				Code: 200,
				Body: "Ok",
			},
			nil,
		},

		{
			Response{
				Code: 503,
				Body: "Server error",
			},
			ErrNotIn,
		},

		{
			nil,
			ErrNotStruct,
		},
	}

	for i, tt := range tests {
		expextedErrors := &tt.expectedErr
		structs := &tt.in
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			validationErrors := Validate(structs)
			require.ErrorAs(t, validationErrors, expextedErrors)
		})
	}
}
