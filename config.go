package gomods

import (
	"strconv"

	"github.com/caddyserver/caddy"
	"github.com/spf13/afero"
)

type Config struct {
	GoBinary string
	Workers  int
	Cache    Cache
	Fs       afero.Fs
}

type Cache struct {
	Enable bool
	Type   string
	Path   string
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
	conf.Fs = afero.NewOsFs()
	if conf.GoBinary == "" {
		conf.GoBinary = DefaultGoBinaryPath
	}
	if conf.Cache.Enable {
		if conf.Cache.Type == "" {
			conf.Cache.Type = DefaultGomodsCacheType
		}
		if conf.Cache.Path == "" {
			conf.Cache.Path = afero.GetTempDir(conf.Fs, "")
		}
	}
	if conf.Workers == 0 {
		conf.Workers = DefaultGomodsWorkers
	}
}

// ParseGomods parses the txtdirect config for gomods
func (conf *Config) ParseGomods(c *caddy.Controller) error {
	switch c.Val() {
	case "gobinary":
		conf.GoBinary = c.RemainingArgs()[0]

	case "workers":
		value, err := strconv.Atoi(c.RemainingArgs()[0])
		if err != nil {
			return c.ArgErr()
		}
		conf.Workers = value

	case "cache":
		conf.Cache.Enable = true
		c.NextArg()
		if c.Val() != "{" {
			break
		}
		for c.Next() {
			if c.Val() == "}" {
				break
			}
			err := conf.Cache.ParseCache(c)
			if err != nil {
				return err
			}
		}
	default:
		return c.ArgErr() // unhandled option for gomods
	}
	return nil
}

// ParseCache parses the txtdirect config for gomods cache
func (cache *Cache) ParseCache(c *caddy.Controller) error {
	switch c.Val() {
	case "type":
		cache.Type = c.RemainingArgs()[0]
	case "path":
		cache.Path = c.RemainingArgs()[0]
	default:
		return c.ArgErr() // unhandled option for gomods cache
	}
	return nil
}
