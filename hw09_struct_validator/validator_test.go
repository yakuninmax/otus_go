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
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Fruit struct {
		Color string `validate:"len:xx"`
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
		in            interface{}
		expectedError error
	}{
		{
			User{
				ID:     "00001",
				Age:    41,
				Email:  "user41@test.mail",
				Role:   "admin",
				Phones: []string{"88005353535"},
			},
			ValidationErrors{
				{Field: "ID", Err: ErrLength},
			},
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
			ValidationErrors{
				{Field: "Version", Err: ErrLength},
			},
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
			ValidationErrors{
				{Field: "Code", Err: ErrNotIn},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			err := Validate(tt.in)
			require.Equal(t, tt.expectedError, err)
		})
	}
}

func TestCommonErrors(t *testing.T) {
	tests := []struct {
		in            interface{}
		expectedError string
	}{
		{
			Fruit{
				Color: "red",
			},
			`strconv.Atoi: parsing "xx": invalid syntax`,
		},
		{
			nil,
			`value is not a structure`,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			err := Validate(tt.in)
			require.Contains(t, err.Error(), tt.expectedError)
		})
	}
}
