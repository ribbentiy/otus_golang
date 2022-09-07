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
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
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
	invalidUser := User{
		ID:     "invalidLength",
		Name:   "TooOld",
		Age:    60,
		Email:  "invalidEmail",
		Role:   "notAdmin",
		Phones: []string{"123"},
		meta:   json.RawMessage("some message"),
	}

	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			App{Version: "12345"},
			nil,
		},
		{
			App{Version: "1234567"},
			ValidationErrors{ValidationError{
				Field: "Version",
				Err:   fmt.Errorf("%w: %v", ErrValidationError, ErrStringInvalidLen),
			}},
		},
		{
			User{
				ID:     "mx87HVXHNapvix9MAtmHCJCXqkGGPWiyp3xD",
				Name:   "Terry",
				Age:    40,
				Email:  "im@thehorse.yes",
				Role:   "admin",
				Phones: []string{"12312312312"},
				meta:   json.RawMessage("some message"),
			},
			nil,
		},
		{
			invalidUser,
			ValidationErrors{
				ValidationError{
					Field: "ID",
					Err:   fmt.Errorf("%w: %v", ErrValidationError, ErrStringInvalidLen),
				},
				ValidationError{
					Field: "Age",
					Err:   fmt.Errorf("%w: %v", ErrValidationError, ErrNumberMoreThanMax),
				},
				ValidationError{
					Field: "Email",
					Err:   fmt.Errorf("%w: %v", ErrValidationError, ErrStringInvalidRegexp),
				},
				ValidationError{
					Field: "Role",
					Err:   fmt.Errorf("%w: %v", ErrValidationError, ErrValueNotInList),
				},
				ValidationError{
					Field: "Phones",
					Err:   fmt.Errorf("%w: %v", ErrValidationError, ErrStringInvalidLen),
				},
			},
		},
		{
			struct {
				User User `validate:"nested"`
				App  App  `validate:"nested"`
			}{
				User: invalidUser,
				App: App{
					Version: "12345",
				},
			},
			ValidationErrors{
				ValidationError{
					Field: "User",
					Err: ValidationErrors{
						ValidationError{
							Field: "ID",
							Err:   fmt.Errorf("%w: %v", ErrValidationError, ErrStringInvalidLen),
						},
						ValidationError{
							Field: "Age",
							Err:   fmt.Errorf("%w: %v", ErrValidationError, ErrNumberMoreThanMax),
						},
						ValidationError{
							Field: "Email",
							Err:   fmt.Errorf("%w: %v", ErrValidationError, ErrStringInvalidRegexp),
						},
						ValidationError{
							Field: "Role",
							Err:   fmt.Errorf("%w: %v", ErrValidationError, ErrValueNotInList),
						},
						ValidationError{
							Field: "Phones",
							Err:   fmt.Errorf("%w: %v", ErrValidationError, ErrStringInvalidLen),
						},
					},
				},
			},
		},
	}

	for i, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tt.expectedErr, Validate(tt.in))
			_ = tt
		})
	}
}
