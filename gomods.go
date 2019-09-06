package gomods

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/gomods/athens/pkg/download"
	"github.com/gomods/athens/pkg/download/addons"
	"github.com/gomods/athens/pkg/module"
	"github.com/gomods/athens/pkg/paths"
	"github.com/gomods/athens/pkg/stash"
	"github.com/gomods/athens/pkg/storage"
	"github.com/gomods/athens/pkg/storage/fs"
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

func (conf *Config) Serve(w http.ResponseWriter, r *http.Request) error {
	m := Module{}
	if err := m.ParseImportPath(r.URL.Path); err != nil {
		return fmt.Errorf("module url is empty")
	}

	dp, err := m.fetch(r, *conf)
	if err != nil {
		return err
	}

	switch m.FileExt {
	case "list":
		list, err := dp.List(r.Context(), m.Name)
		if err != nil {
			return err
		}
		_, err = w.Write([]byte(strings.Join(list, "\n")))
		if err != nil {
			return err
		}
		return nil
	case "info":
		info, err := dp.Info(r.Context(), m.Name, m.Version)
		if err != nil {
			return err
		}
		_, err = w.Write(info)
		if err != nil {
			return err
		}
		return nil
	case "mod":
		mod, err := dp.GoMod(r.Context(), m.Name, m.Version)
		if err != nil {
			return err
		}
		_, err = w.Write(mod)
		if err != nil {
			return err
		}
		return nil
	case "zip":
		zip, err := dp.Zip(r.Context(), m.Name, m.Version)
		if err != nil {
			return err
		}
		defer zip.Close()
		w.Write([]byte{})
		_, err = io.Copy(w, zip)
		if err != nil {
			return err
		}
		return nil
	case "latest":
		info, err := dp.Latest(r.Context(), m.Name)
		if err != nil {
			return err
		}
		json, err := json.Marshal(info)
		if err != nil {
			return err
		}
		_, err = w.Write(json)
		if err != nil {
			return err
		}
		return nil
	default:
		return fmt.Errorf("the requested file's extension is not supported")
	}
}

func (m Module) fetch(r *http.Request, c Config) (download.Protocol, error) {
	fetcher, err := module.NewGoGetFetcher(c.GoBinary, c.Fs)
	if err != nil {
		return nil, err
	}
	s, err := m.storage(c)
	if err != nil {
		return nil, err
	}
	dp := m.dp(fetcher, s, c)
	return dp, nil
}

func (m Module) storage(c Config) (storage.Backend, error) {
	switch c.Cache.Type {
	case "local":
		// Check if cache storage path exists, if not create it
		if _, err := os.Stat(c.Cache.Path); os.IsNotExist(err) {
			if err = os.MkdirAll(c.Cache.Path, os.ModePerm); err != nil {
				return nil, fmt.Errorf("couldn't create the cache storage directory on %s: %s", c.Cache.Path, err.Error())
			}
		}
		s, err := fs.NewStorage(c.Cache.Path, afero.NewOsFs())
		if err != nil {
			return nil, fmt.Errorf("could not create new storage from os fs (%s)", err)
		}
		return s, nil
	case "tmp":
		s, err := fs.NewStorage(c.Cache.Path, afero.NewOsFs())
		if err != nil {
			return nil, fmt.Errorf("could not create new storage from os fs (%s)", err)
		}
		return s, nil
	}
	return nil, fmt.Errorf("Invalid storage config for gomods")
}

func (m Module) dp(fetcher module.Fetcher, s storage.Backend, c Config) download.Protocol {
	lister := download.NewVCSLister(c.GoBinary, c.Fs)
	st := stash.New(fetcher, s, stash.WithPool(c.Workers), stash.WithSingleflight)
	dpOpts := &download.Opts{
		Storage: s,
		Stasher: st,
		Lister:  lister,
	}
	dp := download.New(dpOpts, addons.WithPool(c.Workers))
	return dp
}

// ParseImportPath parses the request path and exports the
// module's import path, module's version and file extension
func (m *Module) ParseImportPath(path string) error {
	if strings.Contains(path, "@latest") {
		pathLatest := strings.Split(path, "/@")
		m.Name, m.Version, m.FileExt = pathLatest[0][1:], "", pathLatest[1]
		if err := m.DecodeImportPath(); err != nil {
			return err
		}
		return nil
	}

	// First item in array is modules import path and the secondd item is version+extension
	pathSlice := strings.Split(path, "/@v/")
	if pathSlice[1] == "list" {
		m.Name, m.Version, m.FileExt = pathSlice[0][1:], "", "list"
		if err := m.DecodeImportPath(); err != nil {
			return err
		}
		return nil
	}

	versionExt := modVersionRegex.FindAllStringSubmatch(pathSlice[1], -1)[0]
	m.Name, m.Version, m.FileExt = pathSlice[0][1:], versionExt[1], versionExt[2]
	if err := m.DecodeImportPath(); err != nil {
		return err
	}

	return nil
}

// DecodeImportPath decodes the module's import path. For more information check
// https://github.com/golang/go/blob/master/src/cmd/go/internal/module/module.go#L375-L433
func (m *Module) DecodeImportPath() error {
	decoded, err := paths.DecodePath(m.Name)
	if err != nil {
		return err
	}
	m.Name = decoded
	return nil
}
