# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres
to [Semantic Versioning](https://semver.org/spec/v2.0.0.html) (after version 0.0.5).

## [0.11.0] - 2024-06-16
### Added
- new `--ignoreTables` option ([PR #59](https://github.com/KarnerTh/mermerd/pull/59))

### Changed
- go 1.22 is now used

### Fixed

## [0.10.0] - 2023-11-21
### Added
- Support relationship labels ([PR #50](https://github.com/KarnerTh/mermerd/pull/50))
- Support multiple key constraints on a single attribute ([PR #52](https://github.com/KarnerTh/mermerd/pull/52))
- Add unique constraint to key column ([PR #53](https://github.com/KarnerTh/mermerd/pull/53))
- Add sqlite support ([PR #55](https://github.com/KarnerTh/mermerd/pull/55))

## [0.9.0] - 2023-08-06
### Changed
- Sort constraints ([Issue #44](https://github.com/KarnerTh/mermerd/issues/44))
- Add NOT NULL constraint to description ([Issue #42](https://github.com/KarnerTh/mermerd/issues/42)

## [0.8.1] - 2023-07-14
### Fixed
- Sort column names ([Issue #40](https://github.com/KarnerTh/mermerd/issues/40))

## [0.8.0] - 2023-05-30
### Changed
- Table names are now sorted in mermaid file ([Issue #34](https://github.com/KarnerTh/mermerd/issues/34))

### Fixed
- Allow not unique constraint names for postgres ([Issue #36](https://github.com/KarnerTh/mermerd/issues/36))

## [0.7.1] - 2023-05-12
### Fixed
- Fix escaping of quote marks in descriptions ([PR #33](https://github.com/KarnerTh/mermerd/pull/33))

## [0.7.0] - 2023-04-06
### Added
- Support column comments ([PR #32](https://github.com/KarnerTh/mermerd/pull/32))

### Changed
- **Breaking change:** In order to support different descriptions the `--showEnumValues` flag was replaced
by `--showDescriptions enumValues` (for details see [PR #32](https://github.com/KarnerTh/mermerd/pull/32))

## [0.6.1] - 2023-03-09
### Fixed
- Fixed wrong table name in constraints if schema prefix was used

## [0.6.0] - 2023-03-08
### Added
- Support schema prefix ([Issue #30](https://github.com/KarnerTh/mermerd/issues/30))

### Changed
- updated dependencies

## [0.5.0] - 2022-12-28
### Added
- Support enum description ([Issue #15](https://github.com/KarnerTh/mermerd/issues/15))
- Support multiple schemas ([Issue #23](https://github.com/KarnerTh/mermerd/issues/23))

## [0.4.1] - 2022-09-28
### Fixed
- Fix wrong column format for `is_primary` ([Issue #24](https://github.com/KarnerTh/mermerd/issues/24))

## [0.4.0] - 2022-09-11
### Added
- Support variable expansion ([Issue #16](https://github.com/KarnerTh/mermerd/issues/16))
- Added support for attribute keys ([Issue #20](https://github.com/KarnerTh/mermerd/issues/20))

### Changed
- mermerd will now by default add the attribute key if applicable (PK or FK). If this is undesired, it can be
  disabled by the `--omitAttributeKeys` flag or the `omitAttributeKeys` config (example is in the readme).

## [0.3.0] - 2022-09-02
### Added
- Added --selectedTables switch ([PR #12](https://github.com/KarnerTh/mermerd/pull/12))
- Added support for MSSQL ([Issue #13](https://github.com/KarnerTh/mermerd/issues/13))

### Changed
- go 1.19 is now used
- updated dependencies

### Fixed
- Fixed some typos and documentation

## [0.2.1] - 2022-06-03
### Fixed
- Embed the template file into the binary ([Issue #10](https://github.com/KarnerTh/mermerd/issues/10))

## [0.2.0] - 2022-06-01
### Added
- A `--debug` flag/config to show debug information
- A `--omitConstraintLabels` flag/config to toggle the new constraint labels

### Changed
- The column name is now displayed as the constraint label (can be switched off)

### Fixed
- Sub query for constraints returned multiple items  ([Issue #8](https://github.com/KarnerTh/mermerd/issues/8))

## [0.1.0] - 2022-04-15
### Added
- Mermerd is available via the go tools

### Changed
- go 1.18 is now used

### Fixed
- MySQL query fix for constraints ([Issue #7](https://github.com/KarnerTh/mermerd/issues/7))

## [0.0.5] - 2022-03-17
### Added
- New config: allow surrounding output with mermerd backticks ([PR #4](https://github.com/KarnerTh/mermerd/pull/4))

## [0.0.4] - 2022-03-14
### Added
- Licence

### Fixed
- Do not require a global configuration file

## [0.0.3] - 2022-03-12
### Added
- Possibility to opt in for all tables
- Start mermerd with a predefined run config
- Add version command
- Show version number in intro header

### Changed
- Improved help command output
- Exit with error code 1 on failure
- Fully POSIX-compliant flags (including short & long versions)
- the parameter for the connection string suggestions (previously `connectionStrings`) was renamed to
  `connectionStringSuggestions`
- the flag `-ac` was replaced with `--showAllConstraints`

### Removed
- `.mermerd` configuration file is not automatically created on first use anymore

## [0.0.2] - 2022-01-30
### Added
- Configurable suggestions for connection string input

### Changed
- improved one to many constraint detection for mysql
- improved one to many constraint detection for postgres

## [0.0.1] - 2022-01-17
### Added
- Initial release of mermerd

[0.10.0]: https://github.com/KarnerTh/mermerd/releases/tag/v0.10.0

[0.9.0]: https://github.com/KarnerTh/mermerd/releases/tag/v0.9.0

[0.8.1]: https://github.com/KarnerTh/mermerd/releases/tag/v0.8.1

[0.8.0]: https://github.com/KarnerTh/mermerd/releases/tag/v0.8.0

[0.7.1]: https://github.com/KarnerTh/mermerd/releases/tag/v0.7.1

[0.7.0]: https://github.com/KarnerTh/mermerd/releases/tag/v0.7.0

[0.6.1]: https://github.com/KarnerTh/mermerd/releases/tag/v0.6.1

[0.6.0]: https://github.com/KarnerTh/mermerd/releases/tag/v0.6.0

[0.5.0]: https://github.com/KarnerTh/mermerd/releases/tag/v0.5.0

[0.4.1]: https://github.com/KarnerTh/mermerd/releases/tag/v0.4.1

[0.4.0]: https://github.com/KarnerTh/mermerd/releases/tag/v0.4.0

[0.3.0]: https://github.com/KarnerTh/mermerd/releases/tag/v0.3.0

[0.2.1]: https://github.com/KarnerTh/mermerd/releases/tag/v0.2.1

[0.2.0]: https://github.com/KarnerTh/mermerd/releases/tag/v0.2.0

[0.1.0]: https://github.com/KarnerTh/mermerd/releases/tag/v0.1.0

[0.0.5]: https://github.com/KarnerTh/mermerd/releases/tag/v0.0.5

[0.0.4]: https://github.com/KarnerTh/mermerd/releases/tag/v0.0.4

[0.0.3]: https://github.com/KarnerTh/mermerd/releases/tag/v0.0.3

[0.0.2]: https://github.com/KarnerTh/mermerd/releases/tag/v0.0.2

[0.0.1]: https://github.com/KarnerTh/mermerd/releases/tag/v0.0.1
