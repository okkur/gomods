package gomods

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

var gomodsRegex = regexp.MustCompile("(list|info|mod|zip)")
var modVersionRegex = regexp.MustCompile("(.*)\\.(info|mod|zip)")

// DefaultGoBinaryPath is the default Golang binary installed on the machine
var DefaultGoBinaryPath = os.Getenv("GOROOT") + "/bin/go"

// Serve handles the incoming requests and serves the module files like .mod and etc
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
