# Build Velociraptor for Windows 7

Velociraptor uses the latest Go toolchain to incorporate the latest
security patches and updates. However, Golang has dropped support for
unsupported operating systems at version 1.20:

https://go.dev/wiki/Windows

This script allows building a version of Velociraptor using the last
supported Go version 1.20. However, note the following caveats:

* To build under this unsupported Go version we had to freeze
  dependencies. Therefore this build includes known buggy and
  unsupported dependencies.

* This build may be insecure! since it includes unsupported
  dependencies.

* We typically update to the latest version of Velociraptor but it may
  be that in future we disable some feature (VQL plugins) that can not
  be easily updated.


NOTE: Do not use this build in a general deployment! Only use it for
deploying on deprecated, unsupported operating systems:

* Windows 7
* Windows 8, 8.1

As always, Windows XP is not supported.