package install

import (
	"runtime"

	download "github.com/st19n/mageutils/download"
)

func Tools(targetDir string, tools map[string]string) error {
	if version, ok := tools["air"]; ok {
		if err := Air(targetDir, version); err != nil {
			return err
		}
	}
	if version, ok := tools["gofumpt"]; ok {
		if err := Gofumpt(targetDir, version); err != nil {
			return err
		}
	}
	if version, ok := tools["golangci-lint"]; ok {
		if err := GolangciLint(targetDir, version); err != nil {
			return err
		}
	}
	if version, ok := tools["gosimports"]; ok {
		if err := Gosimports(targetDir, version); err != nil {
			return err
		}
	}
	if version, ok := tools["gotestsum"]; ok {
		if err := Gotestsum(targetDir, version); err != nil {
			return err
		}
	}
	if version, ok := tools["misspell"]; ok {
		if err := Misspell(targetDir, version); err != nil {
			return err
		}
	}
	if version, ok := tools["swag"]; ok {
		if err := Swag(targetDir, version); err != nil {
			return err
		}
	}
	if version, ok := tools["templ"]; ok {
		if err := Templ(targetDir, version); err != nil {
			return err
		}
	}

	return nil
}

func Air(targetDir, version string) error {
	return download.DownloadBinary(targetDir, download.DownloadOptions{
		Name:    "air",
		URL:     "https://github.com/cosmtrek/air/releases/download/{{.VERSION}}/air_{{.CLEANVERSION}}_{{.GOOS}}_{{.GOARCH}}",
		Version: version,
	})
}

func Gofumpt(targetDir, version string) error {
	return download.DownloadBinary(targetDir, download.DownloadOptions{
		Name:    "gofumpt",
		URL:     "https://github.com/mvdan/gofumpt/releases/download/{{.VERSION}}/gofumpt_{{.VERSION}}_{{.GOOS}}_{{.GOARCH}}",
		Version: version,
	})
}

func GolangciLint(targetDir, version string) error {
	return download.DownloadArchiveExtractBinary(targetDir, download.DownloadOptions{
		Name:            "golangci-lint",
		URL:             "https://github.com/golangci/golangci-lint/releases/download/{{.VERSION}}/golangci-lint-{{.CLEANVERSION}}-{{.GOOS}}-{{.GOARCH}}.tar.gz",
		Version:         version,
		ArchiveFilePath: "golangci-lint-{{.CLEANVERSION}}-{{.GOOS}}-{{.GOARCH}}/{{.NAME}}",
	})
}

func Gosimports(targetDir, version string) error {
	return download.DownloadArchiveExtractBinary(targetDir, download.DownloadOptions{
		Name:    "gosimports",
		URL:     "https://github.com/rinchsan/gosimports/releases/download/{{.VERSION}}/gosimports_{{.CLEANVERSION}}_{{.GOOS}}_{{.GOARCH}}.tar.gz",
		Version: version,
	})
}

func Gotestsum(targetDir, version string) error {
	return download.DownloadArchiveExtractBinary(targetDir, download.DownloadOptions{
		Name:            "gotestsum",
		URL:             "https://github.com/gotestyourself/gotestsum/releases/download/{{.VERSION}}/gotestsum_{{.CLEANVERSION}}_{{.GOOS}}_{{.GOARCH}}.tar.gz",
		Version:         version,
		ArchiveFilePath: "{{.NAME}}",
	})
}

func Misspell(targetDir, version string) error {
	return download.DownloadArchiveExtractBinary(targetDir, download.DownloadOptions{
		Name:            "misspell",
		URL:             "https://github.com/client9/misspell/releases/download/{{.VERSION}}/misspell_{{.CLEANVERSION}}_{{.GOOS}}_{{.GOARCH}}.tar.gz",
		Version:         version,
		ArchiveFilePath: "{{.NAME}}",
		OsReplacement: map[string]string{
			"darwin": "mac",
		},
		ArchReplacement: map[string]string{
			"amd64": "64bit",
		},
	})
}

func Swag(targetDir, version string) error {
	arch := runtime.GOARCH
	if runtime.GOOS == "darwin" {
		if runtime.GOARCH == "amd64" {
			arch = "x86_64"
		}
	}
	return download.DownloadArchiveExtractBinary(targetDir, download.DownloadOptions{
		Name:            "swag",
		URL:             "https://github.com/swaggo/swag/releases/download/{{.VERSION}}/swag_{{.CLEANVERSION}}_{{.GOOS}}_{{.GOARCH}}.tar.gz",
		Version:         version,
		ArchiveFilePath: "{{.NAME}}",
		OsReplacement: map[string]string{
			"linux":   "Linux",
			"windows": "Windows",
			"darwin":  "Darwin",
		},
		ArchReplacement: map[string]string{
			"amd64": arch,
			"386":   "i386",
		},
	})
}

func Templ(targetDir, version string) error {
	return download.DownloadArchiveExtractBinary(targetDir, download.DownloadOptions{
		Name:            "templ",
		URL:             "https://github.com/a-h/templ/releases/download/{{.VERSION}}/templ_{{.GOOS}}_{{.GOARCH}}.tar.gz",
		Version:         version,
		ArchiveFilePath: "{{.NAME}}",
		OsReplacement: map[string]string{
			"linux":   "Linux",
			"windows": "Windows",
			"darwin":  "Darwin",
		},
		ArchReplacement: map[string]string{
			"amd64": "x86_64",
			"386":   "i386",
		},
	})
}
