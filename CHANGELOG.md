# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog], and this project adheres to
[Semantic Versioning].

<!-- references -->
[Keep a Changelog]: https://keepachangelog.com/en/1.0.0/
[Semantic Versioning]: https://semver.org/spec/v2.0.0.html

## [Unreleased]

### Added

- **[BC]** Add `HandlerKey` field to handler-related error types
- **[BC]** Add `Message` field to handler-related error types
- **[BC]** Add `type.TypeContainer.Each()`
- Add `type.RoleMap.Each()`
- Add `type.TypeSet.Each()`
- Add `handler.UnexpectedMessageError`

## [0.6.0] - 2019-08-01

### Added

- Add `config.ApplicationConfig.ApplicationKey`
- Add `config.ApplicationConfig.HandlersByName`
- Add `config.Aggregate.HandlerKey`
- Add `config.Process.HandlerKey`
- Add `config.Projection.HandlerKey`
- Add `config.Integration.HandlerKey`

### Changed

- **[BC]** Replace configurer `Name()` methods with `Identity()`
- **[BC]** Rename `config.ApplicationConfig.Handlers` to `HandlersByName`
- **[BC]** `config.ApplicationConfig.Consumers` now maps to handler configs instead of names
- **[BC]** `config.ApplicationConfig.Producers` now maps to handler configs instead of names

## [0.5.2] - 2019-06-17

### Added

- Add `marshaling` package, which marshals messages and state in various formats

## [0.5.1] - 2019-06-10

### Added

- Add `handler.Type.Validate()`
- Add `message.Correlation.Validate()`
- Add `message.Direction.Validate()`
- Add `message.Role.Validate()`

## [0.5.0] - 2019-04-17

### Added

- Add `SchedulesTimeoutType()` to `ProcessConfigurer` implementation

### Removed

- **[BC]** Remove `message.MetaData`
- **[BC]** Remove `message.Timeout`
- **[BC]** Remove `message.TimeoutMetaData`

### Changed

- **[BC]** Require Dogma v0.4.0

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
[0.5.0]: https://github.com/dogmatiq/enginekit/releases/tag/v0.5.0
[0.5.1]: https://github.com/dogmatiq/enginekit/releases/tag/v0.5.1
[0.5.2]: https://github.com/dogmatiq/enginekit/releases/tag/v0.5.2
[0.6.0]: https://github.com/dogmatiq/enginekit/releases/tag/v0.6.0

<!-- version template
## [0.0.1] - YYYY-MM-DD

### Added
### Changed
### Deprecated
### Removed
### Fixed
### Security
-->
