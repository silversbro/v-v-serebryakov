package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"testing"
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
		name        string
		in          interface{}
		expectedErr error
	}{
		{
			name: "valid user",
			in: User{
				ID:     "123456789012345678901234567890123456",
				Name:   "John Doe",
				Age:    25,
				Email:  "john@example.com",
				Role:   "admin",
				Phones: []string{"12345678901", "10987654321"},
			},
			expectedErr: nil,
		},
		{
			name: "invalid user - multiple errors",
			in: User{
				ID:     "short",
				Name:   "John Doe",
				Age:    17,
				Email:  "invalid-email",
				Role:   "user",
				Phones: []string{"123", "456"},
			},
			expectedErr: ValidationErrors{
				{Field: "ID", Err: errors.New("length must be 36")},
				{Field: "Age", Err: errors.New("must be >= 18")},
				{Field: "Email", Err: errors.New("must match regexp ^\\w+@\\w+\\.\\w+$")},
				{Field: "Role", Err: errors.New("must be one of [admin stuff]")},
				{Field: "Phones[0]", Err: errors.New("length must be 11")},
				{Field: "Phones[1]", Err: errors.New("length must be 11")},
			},
		},
		{
			name: "valid app",
			in: App{
				Version: "1.0.0",
			},
			expectedErr: nil,
		},
		{
			name: "invalid app - version length",
			in: App{
				Version: "1.0.0.1",
			},
			expectedErr: ValidationErrors{
				{Field: "Version", Err: errors.New("length must be 5")},
			},
		},
		{
			name: "valid response",
			in: Response{
				Code: 200,
				Body: "OK",
			},
			expectedErr: nil,
		},
		{
			name: "invalid response code",
			in: Response{
				Code: 400,
				Body: "Bad Request",
			},
			expectedErr: ValidationErrors{
				{Field: "Code", Err: errors.New("must be one of [200 404 500]")},
			},
		},
		{
			name:        "invalid input - not a struct",
			in:          "not a struct",
			expectedErr: errors.New("input must be a struct or pointer to struct"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.in)

			if tt.expectedErr == nil {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
				return
			}

			if err == nil {
				t.Errorf("expected error %v, got nil", tt.expectedErr)
				return
			}

			switch expected := tt.expectedErr.(type) {
			case ValidationErrors:
				actual, ok := err.(ValidationErrors)
				if !ok {
					t.Errorf("expected ValidationErrors, got %T", err)
					return
				}

				if len(actual) != len(expected) {
					t.Errorf("expected %d errors, got %d", len(expected), len(actual))
					return
				}

				for i := range expected {
					if actual[i].Field != expected[i].Field || actual[i].Err.Error() != expected[i].Err.Error() {
						t.Errorf("error %d: expected %v, got %v", i, expected[i], actual[i])
					}
				}

			case error:
				if err.Error() != expected.Error() {
					t.Errorf("expected error %v, got %v", expected, err)
				}
			}
		})
	}
}
