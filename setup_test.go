package gomods

import (
	"bytes"
	"testing"

	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/spf13/afero"
)

func Test_parse(t *testing.T) {
	tests := []struct {
		configFile string
		config     Config
	}{
		{
			`
			gomods
			`,
			Config{
				GoBinary: DefaultGoBinaryPath,
				Workers:  DefaultGomodsWorkers,
			},
		},
		{
			`
			gomods {
				cache
			}
			`,
			Config{
				GoBinary: DefaultGoBinaryPath,
				Workers:  DefaultGomodsWorkers,
				Cache: Cache{
					Enable: true,
					Type:   DefaultGomodsCacheType,
					Path:   "/tmp/txtdirect/gomods",
				},
			},
		},
		{
			`
			gomods {
				gobinary /my/go/binary
				cache {
					type local
					path /my/cache/path
				}
			}
			`,
			Config{
				GoBinary: "/my/go/binary",
				Workers:  DefaultGomodsWorkers,
				Cache: Cache{
					Enable: true,
					Type:   "local",
					Path:   "/my/cache/path",
				},
			},
		},
		{
			`
			gomods {
				gobinary /my/go/binary
				workers 5
				cache {
					type local
					path /my/cache/path
				}
			}
			`,
			Config{
				GoBinary: "/my/go/binary",
				Cache: Cache{
					Enable: true,
					Type:   "local",
					Path:   "/my/cache/path",
				},
				Workers: 5,
			},
		},
	}
	for _, test := range tests {
		buf := bytes.NewBuffer([]byte(test.configFile))
		block, err := caddyfile.Parse("Caddyfile", buf)
		if err != nil {
			t.Errorf("Couldn't read the config: %s", err.Error())
		}

		// Extract the config tokens from the server blocks
		var tokens []caddyfile.Token
		for _, segment := range block[0].Segments {
			for _, token := range segment {
				tokens = append(tokens, token)
			}
		}

		d := caddyfile.NewDispenser(tokens)
		g := &Gomods{}

		if err := g.UnmarshalCaddyfile(d); err != nil {
			t.Errorf("Couldn't parse the config: %s", err.Error())
		}

		// Fs field gets filled by default when parsing the config
		test.config.fs = g.Config.fs
		// Set the default cache path for expected config if cache type is tmp
		if g.Config.Cache.Type == "tmp" {
			test.config.Cache.Path = afero.GetTempDir(test.config.fs, "")
		}

		if g.Config != test.config {
			t.Errorf("Expected config to be %+v, but got %+v", test.config, g.Config)
		}
	}
}
