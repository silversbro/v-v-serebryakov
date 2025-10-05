package hw09structvalidator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name          string
		input         interface{}
		expectedError error
	}{
		{
			name:          "not a struct",
			input:         "string",
			expectedError: ErrNotStruct,
		},
		{
			name:          "not a struct with pointer",
			input:         func() *string { s := "test"; return &s }(),
			expectedError: ErrNotStruct,
		},
		{
			name: "valid struct without tags",
			input: struct {
				Name string
				Age  int
			}{
				Name: "John",
				Age:  30,
			},
			expectedError: nil,
		},
		{
			name: "pointer to struct",
			input: &struct {
				Name string `validate:"len:5"`
			}{
				Name: "Alice",
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.input)
			if tt.expectedError == nil {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, tt.expectedError)
			}
		})
	}
}

func TestValidateString(t *testing.T) {
	type User struct {
		Name     string `validate:"len:5"`
		Email    string `validate:"regexp:^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"`
		Role     string `validate:"in:admin,user,moderator"`
		Password string `validate:"len:8|regexp:[A-Z]|regexp:[a-z]|regexp:[0-9]"`
	}

	tests := []struct {
		name          string
		user          User
		expectedError bool
		checkError    func(t *testing.T, err error)
	}{
		{
			name: "valid user",
			user: User{
				Name:     "Alice",
				Email:    "test@example.com",
				Role:     "user",
				Password: "Pass1234",
			},
			expectedError: false,
		},
		{
			name: "invalid name length",
			user: User{
				Name:     "Bob",
				Email:    "test@example.com",
				Role:     "user",
				Password: "Pass1234",
			},
			expectedError: true,
			checkError: func(t *testing.T, err error) {
				t.Helper()
				var valErr ValidationErrors
				require.ErrorAs(t, err, &valErr)
				require.Len(t, valErr, 1)
				require.ErrorIs(t, valErr[0].Err, ErrValidationLength)
			},
		},
		{
			name: "invalid email regexp",
			user: User{
				Name:     "Alice",
				Email:    "invalid-email",
				Role:     "user",
				Password: "Pass1234",
			},
			expectedError: true,
			checkError: func(t *testing.T, err error) {
				t.Helper()
				var valErr ValidationErrors
				require.ErrorAs(t, err, &valErr)
				require.ErrorIs(t, valErr[0].Err, ErrValidationRegexp)
			},
		},
		{
			name: "invalid role",
			user: User{
				Name:     "Alice",
				Email:    "test@example.com",
				Role:     "guest",
				Password: "Pass1234",
			},
			expectedError: true,
			checkError: func(t *testing.T, err error) {
				t.Helper()
				var valErr ValidationErrors
				require.ErrorAs(t, err, &valErr)
				require.ErrorIs(t, valErr[0].Err, ErrValidationIn)
			},
		},
		{
			name: "multiple validation errors",
			user: User{
				Name:     "Bob",
				Email:    "invalid-email",
				Role:     "guest",
				Password: "weak",
			},
			expectedError: true,
			checkError: func(t *testing.T, err error) {
				t.Helper()
				var valErr ValidationErrors
				require.ErrorAs(t, err, &valErr)
				require.Greater(t, len(valErr), 1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.user)
			if tt.expectedError {
				require.Error(t, err)
				if tt.checkError != nil {
					tt.checkError(t, err)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateInt(t *testing.T) {
	type Config struct {
		Port    int `validate:"min:1|max:65535"`
		Timeout int `validate:"in:5,10,15,30"`
		Level   int `validate:"min:1|max:10"`
	}

	tests := []struct {
		name          string
		config        Config
		expectedError bool
		checkError    func(t *testing.T, err error)
	}{
		{
			name: "valid config",
			config: Config{
				Port:    8080,
				Timeout: 10,
				Level:   5,
			},
			expectedError: false,
		},
		{
			name: "port too small",
			config: Config{
				Port:    0,
				Timeout: 10,
				Level:   5,
			},
			expectedError: true,
			checkError: func(t *testing.T, err error) {
				t.Helper()
				var valErr ValidationErrors
				require.ErrorAs(t, err, &valErr)
				require.ErrorIs(t, valErr[0].Err, ErrValidationMin)
			},
		},
		{
			name: "port too large",
			config: Config{
				Port:    70000,
				Timeout: 10,
				Level:   5,
			},
			expectedError: true,
			checkError: func(t *testing.T, err error) {
				t.Helper()
				var valErr ValidationErrors
				require.ErrorAs(t, err, &valErr)
				require.ErrorIs(t, valErr[0].Err, ErrValidationMax)
			},
		},
		{
			name: "invalid timeout value",
			config: Config{
				Port:    8080,
				Timeout: 25,
				Level:   5,
			},
			expectedError: true,
			checkError: func(t *testing.T, err error) {
				t.Helper()
				var valErr ValidationErrors
				require.ErrorAs(t, err, &valErr)
				require.ErrorIs(t, valErr[0].Err, ErrValidationIn)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.config)
			if tt.expectedError {
				require.Error(t, err)
				if tt.checkError != nil {
					tt.checkError(t, err)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateSlice(t *testing.T) {
	type Form struct {
		Tags      []string `validate:"len:3"`
		IDs       []int    `validate:"min:1"`
		Statuses  []string `validate:"in:active,pending,completed"`
		EmptyTags []string `validate:"len:2"`
	}

	tests := []struct {
		name          string
		form          Form
		expectedError bool
		errorCount    int
	}{
		{
			name: "valid slices",
			form: Form{
				Tags:      []string{"abc", "def", "ghi"},
				IDs:       []int{5, 10, 15},
				Statuses:  []string{"active", "pending"},
				EmptyTags: []string{},
			},
			expectedError: false,
		},
		{
			name: "invalid string slice elements",
			form: Form{
				Tags:     []string{"ab", "def", "ghij"},
				IDs:      []int{5, 10, 15},
				Statuses: []string{"active", "invalid"},
			},
			expectedError: true,
			errorCount:    3,
		},
		{
			name: "invalid int slice elements",
			form: Form{
				Tags:     []string{"abc", "def", "ghi"},
				IDs:      []int{0, -5, 10},
				Statuses: []string{"active"},
			},
			expectedError: true,
			errorCount:    2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.form)
			if tt.expectedError {
				require.Error(t, err)
				var valErr ValidationErrors
				require.ErrorAs(t, err, &valErr)
				if tt.errorCount > 0 {
					require.Len(t, valErr, tt.errorCount)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidationErrors(t *testing.T) {
	type TestStruct struct {
		Field1 string `validate:"len:3"`
		Field2 int    `validate:"min:10"`
		Field3 string `validate:"in:yes,no"`
	}

	testStruct := TestStruct{
		Field1: "abcd",
		Field2: 5,
		Field3: "maybe",
	}

	err := Validate(testStruct)
	require.Error(t, err)

	var valErr ValidationErrors
	require.ErrorAs(t, err, &valErr)
	require.Len(t, valErr, 3)

	// Test Error() method
	errorStr := valErr.Error()
	require.Contains(t, errorStr, "Field1")
	require.Contains(t, errorStr, "Field2")
	require.Contains(t, errorStr, "Field3")

	// Test Is() method
	require.True(t, valErr.Is(ErrValidationLength))
	require.True(t, valErr.Is(ErrValidationMin))
	require.True(t, valErr.Is(ErrValidationIn))
	require.False(t, valErr.Is(ErrNotStruct))
}

func TestInvalidValidators(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		check func(t *testing.T, err error)
	}{
		{
			name: "invalid validator syntax",
			input: struct {
				Field string `validate:"len"`
			}{},
			check: func(t *testing.T, err error) {
				t.Helper()
				require.ErrorIs(t, err, ErrInvalidValidator)
			},
		},
		{
			name: "unsupported field type",
			input: struct {
				Field float64 `validate:"min:1.5"`
			}{},
			check: func(t *testing.T, err error) {
				t.Helper()
				require.ErrorIs(t, err, ErrUnsupportedType)
			},
		},
		{
			name: "invalid regexp pattern",
			input: struct {
				Field string `validate:"regexp:[invalid"`
			}{},
			check: func(t *testing.T, err error) {
				t.Helper()
				require.ErrorIs(t, err, ErrInvalidRegexp)
			},
		},
		{
			name: "invalid int value in validator",
			input: struct {
				Field int `validate:"min:not-a-number"`
			}{},
			check: func(t *testing.T, err error) {
				t.Helper()
				require.ErrorIs(t, err, ErrInvalidIntValue)
			},
		},
		{
			name: "invalid len value",
			input: struct {
				Field string `validate:"len:not-a-number"`
			}{},
			check: func(t *testing.T, err error) {
				t.Helper()
				require.ErrorIs(t, err, ErrInvalidLenValue)
			},
		},
		{
			name: "unknown string validator",
			input: struct {
				Field string `validate:"unknown:value"`
			}{},
			check: func(t *testing.T, err error) {
				t.Helper()
				require.ErrorContains(t, err, "unknown validator unknown for string")
			},
		},
		{
			name: "unknown int validator",
			input: struct {
				Field int `validate:"unknown:value"`
			}{},
			check: func(t *testing.T, err error) {
				t.Helper()
				require.ErrorContains(t, err, "unknown validator unknown for int")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.input)
			require.Error(t, err)
			tt.check(t, err)
		})
	}
}

func TestPrivateFields(t *testing.T) {
	type PrivateTest struct {
		PublicField  string `validate:"len:5"`
		privateField string `validate:"len:3"`
	}

	testStruct := PrivateTest{
		PublicField:  "valid",
		privateField: "invalid",
	}

	err := Validate(testStruct)
	require.NoError(t, err) // private field should be ignored
}

func TestComplexValidators(t *testing.T) {
	type UserProfile struct {
		Username   string `validate:"len:3|regexp:^[a-z]+$|in:admin,user"`
		Age        int    `validate:"min:18|max:120|in:18,21,30,40"`
		PIN        string `validate:"len:4|regexp:^[0-9]+$"`
		MultiCheck string `validate:"len:2|in:ab,cd,ef|regexp:^[a-z]+$"`
	}

	tests := []struct {
		name          string
		profile       UserProfile
		expectedError bool
	}{
		{
			name: "valid complex profile",
			profile: UserProfile{
				Username:   "admin",
				Age:        21,
				PIN:        "1234",
				MultiCheck: "ab",
			},
			expectedError: false,
		},
		{
			name: "multiple validator failures",
			profile: UserProfile{
				Username:   "ADMIN",
				Age:        15,
				PIN:        "12A4",
				MultiCheck: "xyz",
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.profile)
			if tt.expectedError {
				require.Error(t, err)
			}
		})
	}
}
