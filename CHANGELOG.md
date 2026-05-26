# Changelog

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
