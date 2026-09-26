# Changelog

## [1.3.0] - 2026-09-26

### Added

- Add dark theme
- Add theme switcher

### Changed

- Upgrade to Bulma v1
- Add new fontset
- Align page headers on mobile
- Give each admin tab its own title

### Fixed

- Show timestamps slightly in the future as "just now"

## [1.2.1] - 2026-05-26

### Fixed

- Fix session fixation attack

### Changed

- Upgrade to Fiber v3
- Migrate container image from Alpine to Debian
- Disable services by default
- Remove old sessions on startup

### Removed

- Remove unused `updated_at` columns
- Remove `prefixFS`
