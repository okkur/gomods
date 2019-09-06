package gomods

import (
	"net/http"

	"github.com/caddyserver/caddy"
	"github.com/caddyserver/caddy/caddy/caddymain"
	"github.com/caddyserver/caddy/caddyhttp/httpserver"
)

func main() {
	caddymain.EnableTelemetry = false
	caddymain.Run()
}

func init() {
	caddy.RegisterPlugin("gomods", caddy.Plugin{
		ServerType: "http",
		Action:     setup,
	})
	// TODO: hardcode directive after stable release into Caddy
	httpserver.RegisterDevDirective("gomods", "")
}

func parse(c *caddy.Controller) (Config, error) {
	var config Config

	for c.Next() {
		if c.Val() == "gomods" {
			c.Next() // skip directive name
		}

		config.ParseGomods(c)
	}

	config.SetDefaults()

	return config, nil
}

func setup(c *caddy.Controller) error {
	config, err := parse(c)
	if err != nil {
		return err
	}

	// Add handler to Caddy
	cfg := httpserver.GetConfig(c)
	mid := func(next httpserver.Handler) httpserver.Handler {
		return Gomods{
			Next:   next,
			Config: config,
		}
	}
	cfg.AddMiddleware(mid)

	return nil
}

type Gomods struct {
	Next   httpserver.Handler
	Config Config
}

func (rd Gomods) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	if err := rd.Config.Serve(w, r); err != nil {
		if err.Error() == "option disabled" {
			return rd.Next.ServeHTTP(w, r)
		}
		return http.StatusInternalServerError, err
	}

	return 0, nil
}
