package fmt

import (
	"fmt"

	"github.com/magefile/mage/sh"
)

// Fmt - Run gofmt linters
func Fmt() error {
	if err := sh.RunV("gofumpt", "-l", "-w", "."); err != nil {
		fmt.Printf("ERROR: running gofumpt: %v\n", err)
		return err
	}
	if err := sh.RunV("golangci-lint", "run", "--fix", "--enable-only", "gci"); err != nil {
		fmt.Printf("ERROR: running gci: %v\n", err)
		return err
	}
	return nil
}
