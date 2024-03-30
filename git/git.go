package git

import "github.com/magefile/mage/sh"

// Tag returns the git tag for the current branch or "" if none.
func Tag() string {
	s, _ := sh.Output("bash", "-c", "git describe --tags 2>/dev/null")
	return s
}

// Hash returns the git hash for the current repo or "" if none.
func Hash() string {
	hash, _ := sh.Output("git", "rev-parse", "--short", "HEAD")
	return hash
}

// Branch returns the current git branch.
func Branch() string {
	branch, _ := sh.Output("git", "rev-parse", "--abbrev-ref", "HEAD")
	return branch
}
