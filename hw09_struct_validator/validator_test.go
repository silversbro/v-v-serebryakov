package hw09structvalidator

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

//nolint:funlen
func TestValidate(t *testing.T) {
	tests := []struct {
		name         string
		input        interface{}
		expectedErr  error
		expectErrors int
	}{
		{
			name:        "not a struct",
			input:       "just a string",
			expectedErr: ErrNotStruct,
		},
		{
			name:        "nil input",
			input:       nil,
			expectedErr: ErrNotStruct,
		},
		{
			name:  "valid struct with no validation tags",
			input: struct{ Name string }{Name: "test"},
		},
		{
			name: "valid struct with validation",
			input: struct {
				Name string `validate:"len:5"`
				Age  int    `validate:"min:18|max:99"`
			}{
				Name: "Alice",
				Age:  25,
			},
		},
		{
			name: "invalid string length",
			input: struct {
				Name string `validate:"len:5"`
			}{
				Name: "Bob",
			},
			expectedErr:  ErrValidationLength,
			expectErrors: 1,
		},
		{
			name: "invalid int min value",
			input: struct {
				Age int `validate:"min:18"`
			}{
				Age: 15,
			},
			expectedErr:  ErrValidationMin,
			expectErrors: 1,
		},
		{
			name: "invalid regexp",
			input: struct {
				Email string `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
			}{
				Email: "invalid-email",
			},
			expectedErr:  ErrValidationRegexp,
			expectErrors: 1,
		},
		{
			name: "invalid 'in' validation for string",
			input: struct {
				Status string `validate:"in:active,pending,completed"`
			}{
				Status: "deleted",
			},
			expectedErr:  ErrValidationIn,
			expectErrors: 1,
		},
		{
			name: "invalid 'in' validation for int",
			input: struct {
				Code int `validate:"in:200,404,500"`
			}{
				Code: 403,
			},
			expectedErr:  ErrValidationIn,
			expectErrors: 1,
		},
		{
			name: "multiple validation errors",
			input: struct {
				Name  string `validate:"len:5|regexp:^[A-Z]"`
				Age   int    `validate:"min:18|max:99"`
				Email string `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
			}{
				Name:  "bob",
				Age:   15,
				Email: "invalid-email",
			},
			expectErrors: 4,
		},
		{
			name: "invalid validator syntax",
			input: struct {
				Name string `validate:"len"`
			}{
				Name: "test",
			},
			expectedErr:  ErrInvalidValidator,
			expectErrors: 1,
		},
		{
			name: "unsupported field type",
			input: struct {
				Data map[string]string `validate:"len:5"`
			}{
				Data: make(map[string]string),
			},
			expectedErr:  ErrUnsupportedType,
			expectErrors: 1,
		},
		{
			name: "slice validation",
			input: struct {
				IDs []int `validate:"min:1"`
			}{
				IDs: []int{1, 0, 2},
			},
			expectErrors: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.input)

			if tt.expectedErr == nil && tt.expectErrors == 0 {
				require.NoError(t, err)
				return
			}

			require.Error(t, err)

			if tt.expectedErr != nil {
				var valErr ValidationErrors
				if errors.As(err, &valErr) {
					require.True(t, valErr.Is(tt.expectedErr),
						"expected error %v, got %v", tt.expectedErr, valErr)
				} else {
					require.ErrorIs(t, err, tt.expectedErr)
				}
			}

			if tt.expectErrors > 0 {
				var valErr ValidationErrors
				require.True(t, errors.As(err, &valErr))
				require.Len(t, valErr, tt.expectErrors)
			}
		})
	}
}

func TestValidationErrors(t *testing.T) {
	errs := ValidationErrors{
		{Field: "Name", Err: ErrValidationLength},
		{Field: "Age", Err: ErrValidationMin},
	}

	require.True(t, errs.Is(ErrValidationLength))
	require.True(t, errs.Is(ErrValidationMin))
	require.False(t, errs.Is(ErrValidationMax))
	require.False(t, errs.Is(nil))
}

func TestValidationErrorWrapping(t *testing.T) {
	t.Run("wrapped validation errors", func(t *testing.T) {
		err := Validate(struct {
			Name string `validate:"len:5"`
		}{
			Name: "Bob",
		})

		var valErr ValidationErrors
		require.True(t, errors.As(err, &valErr))
		require.Len(t, valErr, 1)
		require.ErrorIs(t, valErr[0].Err, ErrValidationLength)
		require.Contains(t, valErr[0].Err.Error(), "must be 5")
	})

	t.Run("wrapped program errors", func(t *testing.T) {
		err := Validate(struct {
			Data float64 `validate:"min:1"`
		}{
			Data: 0.5,
		})

		var valErr ValidationErrors
		require.True(t, errors.As(err, &valErr))
		require.Len(t, valErr, 1)
		require.ErrorIs(t, valErr[0].Err, ErrUnsupportedType)
	})
}
