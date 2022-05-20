# Wharf core library

[![Codacy Badge](https://app.codacy.com/project/badge/Grade/7eddf4e7be814c2f9fb2b6efbec69fc9)](https://www.codacy.com/gh/iver-wharf/wharf-core/dashboard?utm_source=github.com\&utm_medium=referral\&utm_content=iver-wharf/wharf-core\&utm_campaign=Badge_Grade)
[![Go Reference](https://pkg.go.dev/badge/github.com/iver-wharf/wharf-core.svg)](https://pkg.go.dev/github.com/iver-wharf/wharf-core)

Utility Go code used by numerous other Wharf components.

```sh
go get -u github.com/iver-wharf/wharf-core/v2
```

Mantra of this repository is to include code that will be used in more than 1
other repository, and does not solve any particular use case. It's more for
common code. It holds code that does not solve any particular problems that’s
specific for the different component’s domains.

Instead, it is a place of common utility code. What you will find in this
utility repository is Go code that features:

- :heavy_check_mark: Reading configuration from files and/or environment
  variables

- :heavy_check_mark: Logging in a unified manner

- :heavy_check_mark: Serving common endpoints such as GET /version

- :heavy_check_mark: Common HTTP JSON models, such as the IETF RFC-7807
  problem response.

What you **will not** find in this repository:

- :x: Parsing .wharf-ci.yml files

- :x: Abstractions over Kubernetes

- :x: Abstractions over AMQP
  (already found in [iver-wharf/messagebus-go](https://github.com/iver-wharf/messagebus-go))

- :x: Wharf API database or HTTP JSON models

## Dependencies

- YAML library [gopkg.in/yaml.v2 (github.com/go-yaml/yaml)](https://github.com/go-yaml/yaml)
- Configuration library [spf13/viper](https://github.com/spf13/viper)
- Web framework [github.com/gin-gonic/gin](https://github.com/gin-gonic/gin)
- Database ORM library [gorm.io/gorm](https://gorm.io/)
- Terminal coloring library [github.com/fatih/color](https://github.com/fatih/color)

## Development

1. Install Go 1.18 or later: <https://golang.org/>

2. Install the [swaggo/swag](https://github.com/swaggo/swag) CLI globally:

   ```sh
   # Run this outside of any Go module, including this repository, to not
   # have `go get` update the go.mod file.
   $ cd ..

   $ go get -u github.com/swaggo/swag
   ```

3. Generate the swaggo files (this has to be redone each time the swaggo
   documentation comments has been altered):

   ```sh
   # Navigate back to this repository
   $ cd wharf-api

   # Generate the files into docs/
   $ swag
   ```

4. Start hacking with your favorite tool. For example VS Code, GoLand,
   Vim, Emacs, or whatnot.

## Tests

Requires Go 1.16 or later to be installed: <https://golang.org/>

```sh
go test -v ./...
```

## Linting

```sh
make deps # download linting dependencies

make lint

make lint-go # only lint Go code
make lint-md # only lint Markdown files
```

Some errors can be fixed automatically. Keep in mind that this updates the
files in place.

```sh
make lint-fix

make lint-fix-go # only lint and fix Go files
make lint-fix-md # only lint and fix Markdown files
```

---

Maintained by [Iver](https://www.iver.com/en).
Licensed under the [MIT license](./LICENSE).
