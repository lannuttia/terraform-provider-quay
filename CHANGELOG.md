# Changelog

## [0.4.1](https://github.com/lannuttia/terraform-provider-quay/compare/v0.4.0...v0.4.1) (2026-10-18)


### Bug Fixes

* trigger GoReleaser on release published instead of created ([#25](https://github.com/lannuttia/terraform-provider-quay/issues/25)) ([0efaba2](https://github.com/lannuttia/terraform-provider-quay/commit/0efaba2758a61a729695672c73d5520b91125ae7))

## [0.4.0](https://github.com/lannuttia/terraform-provider-quay/compare/v0.3.1...v0.4.0) (2026-10-18)


### ⚠ BREAKING CHANGES

* replace event string with typed event blocks ([#22](https://github.com/lannuttia/terraform-provider-quay/issues/22))

### Features

* replace event string with typed event blocks ([#22](https://github.com/lannuttia/terraform-provider-quay/issues/22)) ([ca99469](https://github.com/lannuttia/terraform-provider-quay/commit/ca99469f4c9b142e6fd7b77dd857fdeb265c4a9b))

## [0.3.1](https://github.com/lannuttia/terraform-provider-quay/compare/v0.3.0...v0.3.1) (2026-10-18)


### Documentation

* correct v0.3.0 changelog entry to remove duplicate v0.2.2 changes ([#14](https://github.com/lannuttia/terraform-provider-quay/issues/14)) ([4c9c7a3](https://github.com/lannuttia/terraform-provider-quay/commit/4c9c7a371c4c886d29002a612d7fd7bd8140def2))

## [0.3.1](https://github.com/lannuttia/terraform-provider-quay/compare/v0.3.0...v0.3.1) (2026-10-18)


### Documentation

* correct v0.3.0 changelog entry to remove duplicate v0.2.2 changes ([#14](https://github.com/lannuttia/terraform-provider-quay/issues/14)) ([4c9c7a3](https://github.com/lannuttia/terraform-provider-quay/commit/4c9c7a371c4c886d29002a612d7fd7bd8140def2))

## [0.3.0](https://github.com/lannuttia/terraform-provider-quay/compare/v0.2.2...v0.3.0) (2026-10-17)

No user-facing changes. This release was generated due to a release-please tracking issue caused by a missing v0.2.2 git tag.

## [0.2.2](https://github.com/lannuttia/terraform-provider-quay/compare/v0.2.1...v0.2.2) (2026-10-17)


### Features

* add configurable HTTP client timeout ([#6](https://github.com/lannuttia/terraform-provider-quay/issues/6)) ([5df8e77](https://github.com/lannuttia/terraform-provider-quay/commit/5df8e77c66bf5cbb75b71200c033544634857eb6))


### Bug Fixes

* cleanup testcontainer after test ([d7c3974](https://github.com/lannuttia/terraform-provider-quay/commit/d7c3974b46f9d327746d2698678f7fdb735aa286))
* correct eventConfig JSON tag to event_config in CreateNotificationRequest ([719db24](https://github.com/lannuttia/terraform-provider-quay/commit/719db24d9e7586b55fa0ffd9abf98d5a4626b99a))
* correct RELEASE_PLEASE_TOKEN secret name in release workflow ([#10](https://github.com/lannuttia/terraform-provider-quay/issues/10)) ([bdf4c74](https://github.com/lannuttia/terraform-provider-quay/commit/bdf4c74d10e517764392d25b48a92ade337a8a19))
* only trigger release on Tests workflow completion ([#5](https://github.com/lannuttia/terraform-provider-quay/issues/5)) ([7855aa8](https://github.com/lannuttia/terraform-provider-quay/commit/7855aa8e3211c81854b100c19e43d00786ecd03e))
* treat 403 as destroyed in post-destroy check helpers ([41fc6a4](https://github.com/lannuttia/terraform-provider-quay/commit/41fc6a4d10307eb8bfbdff72e25a1df1654827e5))
* use strings.Split in splitImportID to reject import IDs with extra slashes ([2d08598](https://github.com/lannuttia/terraform-provider-quay/commit/2d085981c908a7f20225f1cc20110294f1f45fe1))
* validate visibility at client layer with typed constants ([d7f30d0](https://github.com/lannuttia/terraform-provider-quay/commit/d7f30d03ffec094160e170c88ade201a74edbc36))


### Documentation

* add CONTRIBUTING.md ([03711be](https://github.com/lannuttia/terraform-provider-quay/commit/03711beef653e037da54c5f60e494ad1d6ab05c9))
* document intentional camelCase/snake_case asymmetry in notification API ([108b3e0](https://github.com/lannuttia/terraform-provider-quay/commit/108b3e07a781797630c0b035a6a95a11f3355d05))
* note that pre-commit hooks require the dev shell to be active ([32188a2](https://github.com/lannuttia/terraform-provider-quay/commit/32188a2e87dd84c6c5076eae69b531a110f54ae0))
* update README to reflect current provider scope ([af108f0](https://github.com/lannuttia/terraform-provider-quay/commit/af108f093709bb43be2aa3ea1d9bed57312aa284))
