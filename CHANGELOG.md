# Wharf core library changelog

This project tries to follow [SemVer 2.0.0](https://semver.org/).

<!--
	When composing new changes to this list, try to follow convention.

	The WIP release shall be updated just before adding the Git tag.
	From (WIP) to (YYYY-MM-DD), ex: (2021-02-09) for 9th of Febuary, 2021

	A good source on conventions can be found here:
	https://changelog.md/
-->

## v1.0.0 (WIP)

- Added `pkg/app` with newly added `app.Version` struct together with an example
  of how to use with [gin-gonic/gin](https://github.com/gin-gonic/gin).
	Use this for versioning your built APIs and command-line programs. (#2)

- Added `pkg/config` with newly added `config.Config` interface and
  implementation to let you load configs from environment variables or YAML
  files. This is done via [spf13/viper](https://github.com/spf13/viper).
  (#4, #16)

- Added `pkg/problem` and `pkg/ginutil` to easily use the IETF RFC-7808
  compatible `problem.Response`, originally taken from
  [wharf-api](https://github.com/iver-wharf/wharf-api). (#9)

- Added utility to `pkg/problem` to identify and parse HTTP problem responses
  and made `problem.Response` conform to `error` and `fmt.Stringer` interfaces.
  (#12)

- Added `pkg/logger`, `pkg/logger/consolepretty`, and `pkg/logger/consolejson`
  as fast, low memory using, extensible, and highly customizable logging
  libraries. Heavily inspired by [rs/zerolog](https://github.com/rs/zerolog).
  (#10, #13, #15, #20)

- Added `pkg/logger` integration to Gin-Gonic and GORM inside `pkg/ginutil` and
  `pkg/gormutil`. (#11)

- Added new error response functions. (#17)

- Added `pkg/env` to bind environment variable of a range of different types to
  local variables. Usage is aimed towards backward compatability with old
  environment variables before the age of `pkg/config`, as the `pkg/config`
  already provides much better techniques for loading settings. (#18)
