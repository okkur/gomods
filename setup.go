package gomods

import (
	"fmt"
	"net/http"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

type Gomods struct {
	Config Config
}

func init() {
	caddy.RegisterModule(Gomods{})
	httpcaddyfile.RegisterHandlerDirective("gomods", parse)
}

func parse(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var g Gomods
	err := g.UnmarshalCaddyfile(h.Dispenser)
	return g, err
}

// UnmarshalCaddyfile implements caddyfile.Unmarshaler.
func (g *Gomods) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		if d.Val() == "gomods" {
			d.Next()
			if d.Val() != "{" {
				break
			}
			d.Next()
		}
		if err := g.Config.ParseGomods(d); err != nil {
			return fmt.Errorf("[gomods]: Couldn't parse the config: %s", err.Error())
		}
	}
	g.Config.SetDefaults()
	return nil
}

// CaddyModule returns the Caddy module information.
func (Gomods) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		Name: "http.handlers.gomods",
		New:  func() caddy.Module { return new(Gomods) },
	}
}

// ServeHTTP implements caddyhttp.MiddlewareHandler.
func (g Gomods) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	return g.Config.Serve(w, r)
}
