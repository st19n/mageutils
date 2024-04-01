package build

import (
	"fmt"
	"strings"
)

// LdFlags
// Generates the -ldflags string
//   - packagePath: path of the package where the variables should be applied
//     e.g. app/internal/build
//   - ldFlags: key, value pairs of variables and values
func LdFlags(packagePath string, ldFlags map[string]string) string {
	flags := []string{}

	for key, val := range ldFlags {
		flags = append(flags, fmt.Sprintf("-X %s.%s=%s", packagePath, key, val))
	}

	return strings.Join(flags, " ")
}
