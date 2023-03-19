package errcode

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"testing"
)

func randomUniqueCode() int {
	code := rand.Int()
	if _, ok := codes[code]; ok {
		return randomUniqueCode()
	}

	return code
}

func TestNewError(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("duplicate code: not panic")
		}
	}()

	code := randomUniqueCode()
	_ = NewError(code, http.StatusInternalServerError, "Internal Server Error")
	_ = NewError(code, http.StatusBadRequest, "Bad Request")
}

func TestUnwrap(t *testing.T) {
	err := errors.New("error")
	ec := NewError(randomUniqueCode(), http.StatusInternalServerError, "Error")

	tests := []struct {
		name string
		arg  error
		want error
	}{
		{
			name: "unwrap",
			arg:  err,
			want: err,
		},
		{
			name: "unwrap with code",
			arg:  ec,
			want: ec,
		},
		{
			name: "wrap",
			arg:  fmt.Errorf("wrap: %w", fmt.Errorf("wrap: %w", err)),
			want: err,
		},
		{
			name: "wrap with code",
			arg:  fmt.Errorf("wrap: %w", fmt.Errorf("wrap: %w", ec)),
			want: ec,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Unwrap(tt.arg); !errors.Is(got, tt.want) {
				t.Errorf("Unwrap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHardUnwrap(t *testing.T) {
	err := errors.New("error")
	ec := NewError(randomUniqueCode(), http.StatusInternalServerError, "Error")

	tests := []struct {
		name string
		arg  error
		want error
	}{
		{
			name: "unwrap",
			arg:  err,
			want: nil,
		},
		{
			name: "unwrap with code",
			arg:  ec,
			want: ec,
		},
		{
			name: "wrap",
			arg:  fmt.Errorf("wrap: %w", fmt.Errorf("wrap: %w", err)),
			want: nil,
		},
		{
			name: "wrap with code",
			arg:  fmt.Errorf("wrap: %w", fmt.Errorf("wrap: %w", ec)),
			want: ec,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotOK := HardUnwrap(tt.arg)

			if tt.want == nil {
				if gotOK {
					t.Errorf("HardUnwrap() = %v, want nil", got)
				}

				return
			}

			if !gotOK || !errors.Is(got, tt.want) {
				t.Errorf("HardUnwrap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ExampleNewError() {
	err := NewError(1, http.StatusInternalServerError, "Internal Server Error")
	fmt.Println(err)
	// Output: 1 - Internal Server Error
}

func ExampleError_HTTPStatusCode() {
	err := NewError(2, http.StatusNotFound, "Not Found")
	fmt.Println(err.HTTPStatusCode())
	// Output: 404
}

func ExampleError_Message() {
	err := NewError(3, http.StatusBadRequest, "Invalid %v", "username")
	fmt.Println(err.Message())
	// Output: Invalid username
}

func ExampleUnwrap() {
	// Define error
	ec := NewError(4, http.StatusBadRequest, "Invalid %v")

	// Wrap error in many times (on many functions) for easy tracking
	err := fmt.Errorf("service: %w", fmt.Errorf("repo: %w", ec.WithArgs("username")))

	// Unwrap to get code, http status code and message to respond to the client
	u := Unwrap(err)

	fmt.Println(err)
	fmt.Println(u)
	// Output:
	// service: repo: 4 - Invalid username
	// 4 - Invalid username
}

func ExampleHardUnwrap() {
	// Define error
	ec := NewError(5, http.StatusBadRequest, "Invalid %v")

	// Wrap error in many times (on many functions) for easy tracking
	err := fmt.Errorf("service: %w", fmt.Errorf("repo: %w", ec.WithArgs("username")))

	// Unwrap to get code, http status code and message to respond to the client
	u, ok := HardUnwrap(err)

	fmt.Println(err)
	fmt.Println(u)
	fmt.Println(ok)
	fmt.Println(u.Code())
	fmt.Println(u.HTTPStatusCode())
	fmt.Println(u.Message())
	// Output:
	// service: repo: 5 - Invalid username
	// 5 - Invalid username
	// true
	// 5
	// 400
	// Invalid username
}
