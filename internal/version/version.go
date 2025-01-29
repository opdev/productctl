// Package version contains all identifiable version information for the
// project.
package version

import "fmt"

var (
	baseName    = "productctl"
	projectName = "github.com/opdev/productctl"
	version     = "unknown"
	commit      = "unknown"
)

var Version = Info{
	BaseName: baseName,
	Name:     projectName,
	Version:  version,
	Commit:   commit,
}

type Info struct {
	BaseName string `json:"-"`
	Name     string `json:"name"`
	Version  string `json:"version"`
	Commit   string `json:"commit"`
}

func (vc *Info) String() string {
	return fmt.Sprintf("%s <commit: %s>", vc.Version, vc.Commit)
}
