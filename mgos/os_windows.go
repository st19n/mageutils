//go:build windows
// +build windows

package mgos

// FileExt returns the default file extension based on the operating system.
func FileExt() string {
	return ".exe"
}
