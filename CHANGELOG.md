# Changelog

## [Next]

## [1.4.0] — 2026-10-18

### Changes

- Add sorting to activity page table
- Show loading and error states in the icon picker
- Make table headings consistent
- Improve icon picker keyboard navigation, with Enter picking the highlighted icon
- Shut down gracefully on SIGINT and SIGTERM
- Fix icon search and vendor lookups failing in the Docker image
- Stop clearing a device's vendor when the lookup request fails

## [1.3.0] — 2026-09-26

### Changes

- Add dark theme
- Add theme switcher
- Upgrade to Bulma v1
- Add new fontset
- Align page headers on mobile
- Give each admin tab its own title
- Show timestamps slightly in the future as "just now"

## [1.2.1] — 2026-05-26

### Changes

- Fix session fixation attack
- Upgrade to Fiber v3
- Migrate container image from Alpine to Debian
- Disable services by default
- Remove old sessions on startup
- Remove unused `updated_at` columns
- Remove `prefixFS`
