# Gomods

Caddy plugin that provides a caching Go modules proxy with your own domain.

 [![state](https://img.shields.io/badge/state-beta-blue.svg)]() [![release](https://img.shields.io/github/release/okkur/gomods.svg)](https://gomods.okkur.org/releases) [![license](https://img.shields.io/github/license/okkur/gomods.svg)](LICENSE)

**NOTE: This is a beta release, we do not consider it completely production ready yet. Use at your own risk.**

Gomods is a Caddy plugin that provides a caching Go modules proxy with your own domain.
It supports all the hosting services and VCS` that are supported by Go tools. It also provides local caching
and parallel workers to fetch and store Go modules.

## Using Gomods
Note: The `master` branch is using [Caddy v2](https://caddyserver.com), if you want to use Gomods with previous Caddy versions, check the caddy-v1 branch.

Gomods uses Go tools in the background for fetching the modules so there needs to be an installed version of Go on your machine.

For installing Gomods run the following command:
```
go get go.okkur.org/gomods/cmd/gomods
```

Then you should create a config file like this example:
```
gomods.test {
  gomods
}
```

The example above uses the default values for Go binary and number of parallel workers.
To customize these values add these fields to your config file:
```
gomods.test {
  gomods {
    gobinary /usr/bin/go
    workers 2
  }
}
```

To enable caching you should also add the `cache` field to the config:
```
gomods.test {
  gomods {
    cache
  }
}
```

Just like `gomods` itself, cache also uses its default values when not provided.
You can specify fields like `type` and `path` to customize caching:
```
gomods.test {
  gomods {
    cache {
      type local
      path /home/user/gomods_cache
    }
  }
}
```

For more information about the configuration options and the JSON config
example, check the [Configuration](/docs/configuration.md) page.

To run Gomods run the following command in the same directory that the config file is located:
```
$ gomods start
```


## Support
For detailed information on support options see our [support guide](/SUPPORT.md).

## Helping out
Best place to start is our [contribution guide](/CONTRIBUTING.md).

----

*Code is licensed under the [Apache License, Version 2.0](/LICENSE).*  
*Documentation/examples are licensed under [Creative Commons BY-SA 4.0](/docs/LICENSE).*  
*Illustrations, trademarks and third-party resources are owned by their respective party and are subject to different licensing.*

---

Copyright 2019 - The Gomods authors
