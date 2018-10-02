package error

import (
	"fmt"
	"os"
)

// HasError checks if an error is present and handles any additional
// required steps
func HasError(err error) bool {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Shit just hit the fan: %v\n", err)
		return true
	}
	return false
}
