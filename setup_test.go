package gomods

import (
	"testing"

	"github.com/caddyserver/caddy"
	"github.com/spf13/afero"
)

func Test_parse(t *testing.T) {
	type args struct {
		c *caddy.Controller
	}
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
				cache {
					type local
					path /my/cache/path
				}
				workers 5
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
		c := caddy.NewTestController("http", test.configFile)
		config, err := parse(c)
		if err != nil {
			t.Error(err)
		}

		// Fs field gets filled by default when parsing the config
		test.config.Fs = config.Fs
		// Set the default cache path for expected config if cache type is tmp
		if config.Cache.Type == "tmp" {
			test.config.Cache.Path = afero.GetTempDir(test.config.Fs, "")
		}

		if config != test.config {
			t.Errorf("Expected config to be %+v, but got %+v", test.config, config)
		}
	}
}
