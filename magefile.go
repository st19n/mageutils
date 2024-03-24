//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"

	//mage:import
	_ "github.com/st19n/mageutils/.mage/fmt"
	//mage:import
	_ "github.com/st19n/mageutils/.mage/golang"
	//mage:import
	_ "github.com/st19n/mageutils/.mage/help"

	"github.com/st19n/mageutils/install"
)

var tools = map[string]string{
	"gofumpt":       "v0.6.0",
	"misspell":      "v0.3.4",
	"golangci-lint": "v1.57.1",
	"gotestsum":     "v1.11.0",
}

// Test - Run go tests
func Test() error {
	err := os.MkdirAll(".tmp", 0o750)
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

// Tools - Install required tools
func Tools() error {
	toolBinDir, err := filepath.Abs("./.tmp/bin")
	if err != nil {
		return err
	}
	err = os.MkdirAll(toolBinDir, 0o750)
	if err != nil {
		return err
	}

	return install.Tools(toolBinDir, tools)
}

// Lint - Run golangci-lint
func Lint() error {
	mg.Deps(misspell)
	return sh.RunV("golangci-lint", "run", "--fix")
}

func misspell() error {
	return sh.RunV("misspell", "-error", ".")
}
