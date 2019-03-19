# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog], and this project adheres to
[Semantic Versioning].

<!-- references -->
[Keep a Changelog]: https://keepachangelog.com/en/1.0.0/
[Semantic Versioning]: https://semver.org/spec/v2.0.0.html

## [0.4.0] - 2019-03-19

### Changed

- **[BC]** Remove `logging.Logger` in favour of a simpler log message formatter

## [0.3.1] - 2019-03-07

### Added

- Add `message.Timeout` and `TimeoutMetaData`

## [0.3.0] - 2019-02-28

### Added

- Add `message.TypeContainer` interface
- Add `message.RoleMap` container

### Changed

- **[BC]** Change `config` package to use new `message.RoleMap` container

## [0.2.0] - 2019-02-26

### Added

- Add `config.HandlerConfig.ProducedMessageTypes()`

### Changed

- **[BC]** Remove `MessageTypes` field from handler configuration structs
- **[BC]** Numerous changes to `config.ApplicationConfig` to support produced message declarations
- **[BC]** Replace `config.HandlerConfig.CommandTypes()` and `EventTypes()` with `ConsumedMesssageTypes()`
- Improve application and handler name validation to match Dogma specification

## [0.1.0] - 2019-02-06

- Initial release

<!-- references -->
[Unreleased]: https://github.com/dogmatiq/enginekit
[0.1.0]: https://github.com/dogmatiq/enginekit/releases/tag/v0.1.0
[0.2.0]: https://github.com/dogmatiq/enginekit/releases/tag/v0.2.0
[0.3.0]: https://github.com/dogmatiq/enginekit/releases/tag/v0.3.0
[0.3.1]: https://github.com/dogmatiq/enginekit/releases/tag/v0.3.1
[0.4.0]: https://github.com/dogmatiq/enginekit/releases/tag/v0.4.0

<!-- version template
## [0.0.1] - YYYY-MM-DD

### Added
### Changed
### Deprecated
### Removed
### Fixed
### Security
-->
