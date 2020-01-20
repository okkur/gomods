package gomods

import (
	"strconv"

	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/spf13/afero"
)

type Config struct {
	GoBinary string `json:"gobinary,omitempty"`
	Workers  int    `json:"workers,omitempty"`
	Cache    Cache  `json:"cache,omitempty"`
	fs       afero.Fs
}

type Cache struct {
	Enable bool
	Type   string `json:"type,omitempty"`
	Path   string `json:"path,omitempty"`
}

const (
	// DefaultGomodsCacheType is the default cache mode
	DefaultGomodsCacheType = "tmp"
	// DefaultGomodsWorkers is the default number of parallel workers
	DefaultGomodsWorkers = 1
)

// SetDefaults sets the default values for gomods config
// if the fields are empty
func (conf *Config) SetDefaults() {
	conf.fs = afero.NewOsFs()
	if conf.GoBinary == "" {
		conf.GoBinary = DefaultGoBinaryPath
	}
	if conf.Cache.Enable {
		if conf.Cache.Type == "" {
			conf.Cache.Type = DefaultGomodsCacheType
		}
		if conf.Cache.Path == "" {
			conf.Cache.Path = afero.GetTempDir(conf.fs, "")
		}
	}
	if conf.Workers == 0 {
		conf.Workers = DefaultGomodsWorkers
	}
}

// ParseGomods parses the txtdirect config for gomods
func (conf *Config) ParseGomods(d *caddyfile.Dispenser) error {
	switch d.Val() {
	case "gobinary":
		conf.GoBinary = d.RemainingArgs()[0]

	case "workers":
		value, err := strconv.Atoi(d.RemainingArgs()[0])
		if err != nil {
			return d.ArgErr()
		}
		conf.Workers = value

	case "cache":
		conf.Cache.Enable = true
		d.Next()

		if d.Val() != "{" {
			break
		}
		for d.Next() {
			if d.Val() == "}" {
				continue
			}
			err := conf.Cache.ParseCache(d)
			if err != nil {
				return err
			}
		}
	default:
		return d.ArgErr() // unhandled option for gomods
	}
	return nil
}

// ParseCache parses the txtdirect config for gomods cache
func (cache *Cache) ParseCache(d *caddyfile.Dispenser) error {
	switch d.Val() {
	case "type":
		cache.Type = d.RemainingArgs()[0]
	case "path":
		cache.Path = d.RemainingArgs()[0]
	default:
		return d.ArgErr() // unhandled option for gomods cache
	}
	return nil
}
