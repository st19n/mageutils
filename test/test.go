package test

import (
	"fmt"
	"strings"

	"github.com/magefile/mage/sh"
	"github.com/st19n/mageutils/install"
)

const (
	jsonOut         = "./.tmp/test.json"
	testCoverageOut = "./.tmp/test_coverage"
)

// GotestsumWithCoverage - Run gotestsum and print the total test coverage.
func GotestsumWithCoverage(ldFlags string) error {
	err := Gotestsum(ldFlags)
	if err != nil {
		return err
	}
	return GotestsumCoverage()
}

// Gotestsum - Run gotestsum.
func Gotestsum(ldFlags string) error {
	_, err := install.CreateDir("./tmp")
	if err != nil {
		return err
	}
	return sh.RunV(
		"gotestsum",
		"--format", "pkgname-and-test-fails",
		"--jsonfile", jsonOut,

		"--",

		"-ldflags", ldFlags,
		"-race",
		"-cover",
		"-coverprofile="+testCoverageOut,
		"./...",
	)
}

// GotestsumWatch - Run gotestsum with the --watch flag.
func GotestsumWatch(ldFlags string) error {
	return sh.RunV(
		"gotestsum",
		"--format", "testname",
		"--watch",

		"--",

		"-race",
		"-ldflags", ldFlags,
	)
}

// GotestsumCoverage - Run coverage on a previous test run.
func GotestsumCoverage() error {
	_, err := install.CreateDir("./tmp")
	if err != nil {
		return err
	}
	out, err := sh.Output("go", "tool", "cover", fmt.Sprintf("-func=%s", testCoverageOut))
	if err != nil {
		return err
	}
	lines := strings.Split(out, "\n")
	totalLine := strings.Fields(lines[len(lines)-1])
	fmt.Printf("\nTotal test coverage: %s\n", totalLine[len(totalLine)-1])
	return nil
}

func CoverageReport() error {
	return sh.RunV("go", "tool", "cover", "-html="+testCoverageOut)
}
