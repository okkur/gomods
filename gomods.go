package gomods

import (
	"net/http"
	"os"
	"regexp"

	"github.com/gomods/athens/pkg/download"
	"github.com/gomods/athens/pkg/module"
	"github.com/gomods/athens/pkg/storage"
	"github.com/spf13/afero"
)

type Module struct {
	Name    string
	Version string
	FileExt string
}

type ModuleHandler interface {
	fetch(r *http.Request, c Config) (*storage.Version, error)
	storage(c Config) (storage.Backend, error)
	dp(fetcher module.Fetcher, s storage.Backend, fs afero.Fs) download.Protocol
}

var gomodsRegex = regexp.MustCompile("(list|info|mod|zip)")
var modVersionRegex = regexp.MustCompile("(.*)\\.(info|mod|zip)")
var DefaultGoBinaryPath = os.Getenv("GOROOT") + "/bin/go"
