# Configuration

Gomods has a set of default values for the config options and you can use it by
only adding the plugin name inside your config file.

Here is a example in Caddyfile format:

```
gomods.test {
  gomods
}
```

## Config Options

General options:

- `gobinary`: _string_ - Path to the Go binary (Default: `/usr/bin/go`)
- `workers`: _integer_ - Number of workers for fetching the modules (Default: `1`)
- `cache`: _object_ - Cache configuration object

Cache options:

- `type`: _string_ - Cache storage type (Default: `tmp`, Options: [`local`, `tmp`])
- `path`: _string_ - Cache storage path (Default: `/tmp`)

## Examples

Caddyfile Example:

```
gomods.test {
  gomods {
		gobinary /my/go/binary
		workers 5
		cache {
			type local
			path /my/cache/path
		}
	}
}
```

JSON Example:

```json
{
  "apps": {
    "http": {
      "servers": {
        "srv0": {
          "listen": [":443"],
          "routes": [
            {
              "match": [
                {
                  "host": ["gomods.test"]
                }
              ],
              "handle": [
                {
                  "handler": "subroute",
                  "routes": [
                    {
                      "handle": [
                        {
                          "Config": {
                            "cache": {
                              "path": "/my/cache/path",
                              "type": "local"
                            },
                            "gobinary": "/my/go/binary",
                            "workers": 5
                          },
                          "handler": "gomods"
                        }
                      ]
                    }
                  ]
                }
              ]
            }
          ]
        }
      }
    }
  }
}
```
