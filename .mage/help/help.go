package help

import (
	"fmt"
)

// Help - Show dependencies, useful tips.
func Help() {
	h := `Dependencies:
	- golang: os package install golang
	- go run mage.go tools

	Add 'export MAGEFILE_ENABLE_COLOR=1' to env for colors
`
	fmt.Println(h)
}
