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
func GotestsumWithCoverage(ldFlags, files string, goTestArgs ...string) error {
	err := Gotestsum(ldFlags, files, goTestArgs...)
	if err != nil {
		return err
	}
	return GotestsumCoverage()
}

// Gotestsum - Run gotestsum.
func Gotestsum(ldFlags, files string, goTestArgs ...string) error {
	_, err := install.CreateDir("./tmp")
	if err != nil {
		return err
	}
	args := []string{
		"--format", "pkgname-and-test-fails",
		"--jsonfile", jsonOut,
		"--packages", files,

		"--",

		"-ldflags", ldFlags,
		"-race",
		"-cover",
		"-coverprofile="+testCoverageOut,
	}
	args = append(args, goTestArgs...)

	return sh.RunV(
		"gotestsum",
		args...,
	)
}

// GotestsumWatch - Run gotestsum with the --watch flag.
func GotestsumWatch(ldFlags string, goTestArgs ...string) error {
	args := []string{
		"--format", "testname",
		"--watch",

		"--",

		"-race",
		"-ldflags", ldFlags,
	}
	args = append(args, goTestArgs...)

	return sh.RunV(
		"gotestsum",
		args...,
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
