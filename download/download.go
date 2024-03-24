package downloads

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"

	"github.com/magefile/mage/sh"
	"github.com/mholt/archiver/v3"
	"github.com/st19n/mageutils/mgos"
)

// remix of https://github.com/carolynvs/magex
//          https://github.com/uwu-tools/magex

type DownloadOptions struct {
	// URL is a Go template string for the URL to download. Required.
	// Available Template Variables:
	//   - {{.GOOS}}
	//   - {{.GOARCH}}
	//   - {{.EXT}}
	//   - {{.VERSION}}
	//   - {{.CLEANVERSION}} (equals VERSION, minus the 'v' prefix
	//   - {{.NAME}}
	URL string

	// Name of the binary, excluding OS specific file extension. Required.
	// Also set to {{.NAME}}
	Name string

	// Version to replace {{.VERSION}} in the URL template. Optional depending on whether or not the version is in the Url template
	Version string

	// OsReplacement maps from a GOOS to the os keyword used for the download. Optional, defaults to empty.
	OsReplacement map[string]string

	// ArchReplacement maps from a GOARCH to the arch keyword used for the download. Optional, defaults to empty.
	ArchReplacement map[string]string

	// ArchiveFilePath is the file path of the binary within the archive. Optional, defaults to `Name`
	ArchiveFilePath string
}

func DownloadBinary(destDir string, opts DownloadOptions) error {
	src, err := renderTemplate(opts.URL, opts)
	if err != nil {
		return err
	}

	fmt.Printf("downloading %s\n", src)

	// Download to a temp file
	tmpDir, err := os.MkdirTemp("", "mageutils")
	if err != nil {
		return fmt.Errorf("could not create temporary directory: %w", err)
	}
	defer os.RemoveAll(tmpDir)
	tmpFile := filepath.Join(tmpDir, filepath.Base(src))

	r, err := http.Get(src)
	if err != nil {
		return fmt.Errorf("could not resolve %s: %w", src, err)
	}
	defer r.Body.Close()
	if r.StatusCode >= 400 {
		return fmt.Errorf("error downloading %s (%d): %s", src, r.StatusCode, r.Status)
	}

	f, err := os.OpenFile(tmpFile, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0o755)
	if err != nil {
		return fmt.Errorf("could not open %s: %w", tmpFile, err)
	}
	defer f.Close()

	// Download to the temp file
	_, err = io.Copy(f, r.Body)
	if err != nil {
		return fmt.Errorf("error downloading %s: %w", src, err)
	}
	f.Close()

	// Make the binary executable
	err = os.Chmod(tmpFile, 0o755)
	if err != nil {
		return fmt.Errorf("could nog make %s executable: %w", tmpFile, err)
	}

	// Move it to the destination
	destPath := filepath.Join(destDir, opts.Name+mgos.FileExt())
	if err = sh.Copy(destPath, tmpFile); err != nil {
		return fmt.Errorf("error copying %s to %s: %w", tmpFile, destPath, err)
	}

	return nil
}

func DownloadArchiveExtractBinary(destDir string, opts DownloadOptions) error {
	src, err := renderTemplate(opts.URL, opts)
	if err != nil {
		return err
	}

	fmt.Printf("downloading archive %s\n", src)

	// Download to a temp file
	tmpDir, err := os.MkdirTemp("", "mageutils")
	if err != nil {
		return fmt.Errorf("could not create temporary directory: %w", err)
	}
	// defer os.RemoveAll(tmpDir)
	tarFile := filepath.Join(tmpDir, filepath.Base(src))
	binFile := filepath.Join(tmpDir, opts.Name+mgos.FileExt())

	r, err := http.Get(src)
	if err != nil {
		return fmt.Errorf("could not resolve %s: %w", src, err)
	}
	defer r.Body.Close()
	if r.StatusCode >= 400 {
		return fmt.Errorf("error downloading %s (%d): %s", src, r.StatusCode, r.Status)
	}

	f, err := os.OpenFile(tarFile, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0o755)
	if err != nil {
		return fmt.Errorf("could not open %s: %w", tarFile, err)
	}
	defer f.Close()

	// Download to the temp file
	_, err = io.Copy(f, r.Body)
	if err != nil {
		return fmt.Errorf("error downloading %s: %w", src, err)
	}
	f.Close()

	// extract archive
	archiveFilePath := opts.Name
	if opts.ArchiveFilePath != "" {
		archivePath, _err := renderTemplate(opts.ArchiveFilePath, opts)
		if _err != nil {
			fmt.Println(_err)
			return _err
		}
		archiveFilePath = archivePath
		binFile = filepath.Join(tmpDir, archiveFilePath+mgos.FileExt())
	}

	err = archiver.Extract(f.Name(), archiveFilePath, tmpDir)
	if err != nil {
		return fmt.Errorf("failed to extract archive %s: %w", f.Name(), err)
	}

	// Make the binary executable
	err = os.Chmod(binFile, 0o755)
	if err != nil {
		return fmt.Errorf("could nog make %s executable: %w", binFile, err)
	}

	// Move it to the destination
	destPath := filepath.Join(destDir, opts.Name+mgos.FileExt())
	if err = sh.Copy(destPath, binFile); err != nil {
		return fmt.Errorf("error copying %s to %s: %w", binFile, destPath, err)
	}

	return nil
}

// renderTemplate takes a Go templated string and expands template variables
// Available Template Variables:
// - {{.NAME}}
// - {{.GOOS}}
// - {{.GOARCH}}
// - {{.EXT}}
// - {{.VERSION}}
// - {{.CLEANVERSION}} (the {{.VERSION}} with `v` stripped).
func renderTemplate(templateString string, opts DownloadOptions) (string, error) {
	tmpl, err := template.New("url").Parse(templateString)
	if err != nil {
		return "", fmt.Errorf("error parsing %s as Go template: %w", opts.URL, err)
	}

	srcData := struct {
		NAME         string
		GOOS         string
		GOARCH       string
		EXT          string
		VERSION      string
		CLEANVERSION string
	}{
		NAME:         opts.Name,
		GOOS:         runtime.GOOS,
		GOARCH:       runtime.GOARCH,
		EXT:          mgos.FileExt(),
		VERSION:      opts.Version,
		CLEANVERSION: strings.Replace(opts.Version, "v", "", 1),
	}

	if overrideGoos, ok := opts.OsReplacement[runtime.GOOS]; ok {
		srcData.GOOS = overrideGoos
	}

	if overrideGoarch, ok := opts.ArchReplacement[runtime.GOARCH]; ok {
		srcData.GOARCH = overrideGoarch
	}

	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, srcData)
	if err != nil {
		return "", fmt.Errorf("error rendering %s as Go template with data: %#v: %w", opts.URL, srcData, err)
	}

	return buf.String(), nil
}
