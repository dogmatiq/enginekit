# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog], and this project adheres to
[Semantic Versioning].

<!-- references -->

[Keep a Changelog]: https://keepachangelog.com/en/1.0.0/
[Semantic Versioning]: https://semver.org/spec/v2.0.0.html

## [0.10.3] - 2024-08-21

### Added

- Added the `enginetest/stubs` package as a replacement for the deprecated
  `github.com/dogmatiq/dogma/fixtures` package.

## [0.10.2] - 2024-04-08

### Added

- Added `uuidpb.Map.Has()`.

### Fixed

- Fixed unsigned integer overflow in `uuidpb.Compare()`.

## [0.10.1] - 2024-04-08

### Added

- Added `uuidpb.MapKey.Format()`, `String()` and `DapperString()` methods.

## [0.10.0] - 2024-03-27

- Initial release.

> [!NOTE]
> Releases v0.9.0 and prior where part of a prior "incarnation" of the
> `enginekit` package name. These versions are still cached on go.pkg.dev, but
> are unrelated to this repository.

<!-- references -->

[Unreleased]: https://github.com/dogmatiq/enginekit
[0.10.0]: https://github.com/dogmatiq/enginekit/releases/v0.10.0
[0.10.1]: https://github.com/dogmatiq/enginekit/releases/v0.10.1
[0.10.2]: https://github.com/dogmatiq/enginekit/releases/v0.10.2
[0.10.3]: https://github.com/dogmatiq/enginekit/releases/v0.10.3

<!-- version template
## [0.0.1] - YYYY-MM-DD

### Added
### Changed
### Deprecated
### Removed
### Fixed
### Security
-->
