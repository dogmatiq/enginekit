# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog], and this project adheres to
[Semantic Versioning].

<!-- references -->

[Keep a Changelog]: https://keepachangelog.com/en/1.0.0/
[Semantic Versioning]: https://semver.org/spec/v2.0.0.html

## [Unreleased]

### Added

- Added `maps.NewFromSeq()` (and variants) which construct map types from
  `iter.Seq2` sequences.
- Added `sets.NewFromSeq()` (and variants) which construct set types from
  `iter.Seq` sequences.
- Added `sets.NewFromKeys()` (and variants) which construct set types from
  the keys of `iter.Seq2` sequences.
- Added `sets.NewFromValues()` (and variants) which construct set types from
  the values of `iter.Seq2` sequences.
- Added `Intersection()` method to all set types.

### Changed

- The results of `Clone()`, `Merge()`, `Select()` and `Project()` on any map
  type are now guaranteed to be non-nil.
- The results of `Clone()`, `Union()` and `Project()` on any set type are now
  guaranteed to be non-nil.

## [0.15.1] - 2024-10-02

### Removed

- Removed `nocopy` protection from collection types.

## [0.15.0] - 2024-10-02

### Added

- Added `sets.Proto` which is an unordered set of `proto.Message` values.
- Added `maps.Proto` which is an unordered map of `proto.Message` keys to
  arbitrary values.

### Changed

- **[BC]** Split `collection` package into separate `collections/maps` and
  `collections/sets` packages.

## [0.14.0] - 2024-09-30

### Added

- Added `collection.OrderedSet` and `UnorderedSet`.
- Added `collection.OrderedMap`.

### Removed

- **[BC]** Removed `uuidpb.OrderedSet` and `Map`, use `collections.OrderedSet`
  and `OrderedMap` instead.

## [0.13.0] - 2024-09-30

### Added

- Added `Marshaler.UnmarshalTypeFromMediaType()`.

### Removed

- Removed `Envelope.PortableName`. The `MediaType` field is now guaranteed to
  include the portable name as a parameter.
- Removed `Packet.PortableName()`.

## [0.12.2] - 2024-09-30

### Fixed

- Fixed error in `Marshaler.MarshalAs()` when passed a media-type that is
  unsupported because it does not have a `type` parameter. The implementation
  now correctly returns `false` instead.

## [0.12.1] - 2024-09-29

### Added

- Added `identitypb.Identity.Equal()`.
- Added `envelopepb.Packer` and `Transcoder` types.

## [0.12.0] - 2024-09-29

### Changed

- **[BC]** Renamed `uuidpb.FromString()` to `Parse()`.

### Added

- Added `uuidpb.MustParse()`.
- Added `marshaler` package as a replacement for `dogmatiq/marshalkit`.

## [0.11.1] - 2024-09-27

### Added

- Added JSON struct tags to stub types.

## [0.11.0] - 2024-09-25

### Added

- Added `Map.All()`, which returns an iterator that ranges over all key/value
  pairs in the map.
- Added`Map.Keys()` and `Values()` methods, which return iterators that range
  over the map's keys and values, respectively.
- Added `OrderedSet.All()`, which returns an iterator that ranges over all
  values in the set, in order.
- Added `Map.Len()` and `OrderedSet.Len()`.
- Added `protobuf/configpb` and `grpc/configgrpc` packages as a replacement for
  the `configspec` package from `dogmatiq/interopspec`.

### Changed

- Bumped minimum Go version to v1.23.
- **[BC]** `Map` is now a struct instead of an actual Go map. Iteration is
  provided by a new `All()` method that returns an iterator.
- **[BC]** `Set` has been renamed to `OrderedSet`, and is now a struct instead
  of a slice. Iteration is provided by a new `All()` method that returns an
  iterator.

### Removed

- **[BC]** Removed `MapKey` type.

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
[0.11.0]: https://github.com/dogmatiq/enginekit/releases/v0.11.0
[0.11.1]: https://github.com/dogmatiq/enginekit/releases/v0.11.1
[0.12.0]: https://github.com/dogmatiq/enginekit/releases/v0.12.0
[0.12.1]: https://github.com/dogmatiq/enginekit/releases/v0.12.1
[0.12.2]: https://github.com/dogmatiq/enginekit/releases/v0.12.2
[0.13.0]: https://github.com/dogmatiq/enginekit/releases/v0.13.0
[0.14.0]: https://github.com/dogmatiq/enginekit/releases/v0.14.0
[0.15.0]: https://github.com/dogmatiq/enginekit/releases/v0.15.0
[0.15.1]: https://github.com/dogmatiq/enginekit/releases/v0.15.1

<!-- version template
## [0.0.1] - YYYY-MM-DD

### Added
### Changed
### Deprecated
### Removed
### Fixed
### Security
-->
