//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/magefile/mage/sh"
)

func Test() error {
	err := os.MkdirAll(".tmp", 0750)
	if err != nil {
		return err
	}

	var (
		jsonOut     = ".tmp/gotestsum.json"
		coverageOut = ".tmp/coverage.out"
	)

	out, err := sh.Output("gotestsum", "--format", "pkgname-and-test-fails", "--jsonfile", jsonOut, "--", "-race", "-cover", "-coverprofile="+coverageOut, "./...")
	if err != nil {
		fmt.Println(out)
		return err
	}

	out, err = sh.Output("go", "tool", "cover", "-func="+coverageOut)
	if err != nil {
		return err
	}

	lines := strings.Split(out, "\n")
	totalLine := strings.Fields(lines[len(lines)-1])
	fmt.Printf("\nTotal test coverge: %s\n", totalLine[len(totalLine)-1])
	return nil
}
