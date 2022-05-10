# Wharf core library changelog

This project tries to follow [SemVer 2.0.0](https://semver.org/).

<!--
	When composing new changes to this list, try to follow convention.

	The WIP release shall be updated just before adding the Git tag.
	From (WIP) to (YYYY-MM-DD), ex: (2021-02-09) for 9th of Febuary, 2021

	A good source on conventions can be found here:
	https://changelog.md/
-->

## v2.0.0 (WIP)

- BREAKING: Changed minor version of Go from 1.16 to 1.18. (#40)

- BREAKING: Changed module path from `github.com/iver-wharf/wharf-core` to
  `github.com/iver-wharf/wharf-core/v2`. (#40)

- Changed `env.Bind` to use generic constraints for compile-time assertions
  instead of runtime assertions. (#40)

- Changed quotation marks in `pkg/logger/consolepretty` to `“”` instead of
  backtick `` ` `` to result in fewer backslash escapes. (#43)

- Fixed shortened caller file name and scope being miscalculated due to the
  default ellipsis being 1 rune long but 3 bytes. It will now correctly treat
  `…` as only 1 rune when calculating the string shortening. (#44)

- Fixed scopes delimiter being written by `pkg/logger/consolepretty` when no
  scopes are used. (#45)

- Added `Config.DisableScope` in `pkg/logger/consolepretty`. (#45)

- Changed version of dependencies:

  - `github.com/fatih/color` from v1.12.0 to v1.13.0 (#46)
  - `github.com/gin-gonic/gin` from v1.7.1 to v1.7.7 (#46)
  - `github.com/mattn/go-colorable` from v0.1.8 to v0.1.12 (#46)
  - `github.com/spf13/viper` from v1.7.1 to v1.10.1 (#46)
  - `github.com/stretchr/testify` from v1.7.0 to v1.7.1 (#46)
  - `gorm.io/driver/postgres` from v1.1.0 to v1.3.1 (#46)
  - `gorm.io/gorm` from v1.21.11 to v1.23.3 (#46)

- Changed `DocsHome` in `pkg/problem` to `wharf.iver.com`. (#48)

## v1.3.0 (2021-11-30)

- Added `Event.WithFunc(func(Event) Event) Event` to `wharf-core/pkg/logger` to
  make it easier to reuse field inside a certain scope: (#29)

  ```go
  func someFunc(group, name string) {
    logArgs := func(ev logger.Event) logger.Event {
      return ev.
        WithString("group", group).
        WithString("name", name)
    }

    log.Debug().WithFunc(logArgs).Message("Foo bar.")
  }
  ```

- Added `consolepretty.Config.ScopeMinLengthAuto` which will pad scopes
  automatically based on the longest registered scope, or use the
  `ScopeMinLength` and `ScopeMaxLength` configs to get more fine grained
  control. The auto config is active by default. (#32)

- Added `consolepretty.Config.Ellipsis`, which defaults to the unicode ellipsis
  character `…`, which is used when trimming is applied by `CallerMaxLength` and
  `ScopeMaxLength`. (#33)

## v1.2.0 (2021-09-07)

- Added `wharf-core/pkg/cacertutil`, taken from `wharf-api/internal/httputils`,
  in preparation to delete it from the wharf-api, wharf-provider-github,
  wharf-provider-gitlab, and wharf-provider-azuredevops repos. (#27)

## v1.1.0 (2021-08-20)

- Changed default field colors in `pkg/logger/consolepretty` from yellow to
  white. (#23)

- Changed formatting in `pkg/logger/consolepretty` to have less whitespace
  around the scope and logging level, as well as adding padding and trimming to
  the caller so that it stays at a constant width for a given scope. (#23)

## v1.0.0 (2021-07-13)

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
