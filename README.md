# errcode

Extends go errors via interfaces to have:

- unique code (to understand its type)
- HTTP status code (to respond)

## Install

```shell
go get github.com/toanppp/errcode
```

## Documentation

https://pkg.go.dev/github.com/toanppp/errcode

## Example

```go
package main

import (
	"fmt"
	"github.com/toanppp/errcode"
	"net/http"
)

// Define error type with HTTP status code, message format and its unique code

var (
	InValidError = errcode.NewError(1001, http.StatusBadRequest, "Invalid %v")
)

func Repository() error {
	// Throw error on Repository
	return InValidError.WithArgs("username")
}

func Service() error {
	if err := Repository(); err != nil {
		// Wrap error in many times (on many functions) for easy tracking
		return fmt.Errorf("an error orcured in Repository: %w", err)
	}

	return nil
}

func main() {
	if err := Service(); err != nil {
		// Unwrap to get code, http status code and message to respond to the client
		if ec, ok := errcode.HardUnwrap(err); ok {
			fmt.Println(ec.Code())
			fmt.Println(ec.HTTPStatusCode())
			fmt.Println(ec.Message())
		}
	}
}
```
