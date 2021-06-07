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
  files. This is done via [spf13/viper](https://github.com/spf13/viper). (#4)