package golang

import "github.com/magefile/mage/sh"

// Tidy - Run go tidy
func Tidy() error {
	return sh.Run("go", "mod", "tidy")
}
