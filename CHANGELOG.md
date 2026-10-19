# Changelog

## [0.5.7](https://github.com/lannuttia/terraform-provider-quay/compare/v0.5.6...v0.5.7) (2026-10-19)


### Bug Fixes

* replace deprecated archives.format with archives.formats ([422b658](https://github.com/lannuttia/terraform-provider-quay/commit/422b658d51f3b1b43583c405a8364076e00742e5))
* use literal path for subject-checksums in attest step ([1e4e796](https://github.com/lannuttia/terraform-provider-quay/commit/1e4e79682ccc63777f9bf84d7f29065321449266))
* use literal path for subject-checksums in attest step ([0346bfb](https://github.com/lannuttia/terraform-provider-quay/commit/0346bfb2d8f4c09c067f9672f17b93961b3b8273))
* use manifest mode by removing release-type from release-please action ([237f6fd](https://github.com/lannuttia/terraform-provider-quay/commit/237f6fd9e736114b00ad8f86434fdced47d8ea50))

## [0.5.6](https://github.com/lannuttia/terraform-provider-quay/compare/v0.5.5...v0.5.6) (2026-10-19)


### Bug Fixes

* use release-please manifest mode to respect draft and force-tag-creation config ([3f415fd](https://github.com/lannuttia/terraform-provider-quay/commit/3f415fd2573ae3f8410f4bcd7a1e6ceadd94ba0f))

## [0.5.5](https://github.com/lannuttia/terraform-provider-quay/compare/v0.5.4...v0.5.5) (2026-10-19)


### Bug Fixes

* enable force-tag-creation and goreleaser preflight for immutable releases ([3233e05](https://github.com/lannuttia/terraform-provider-quay/commit/3233e05a9641ee419ec37dec48733a0636ba6a8b))

## [0.5.4](https://github.com/lannuttia/terraform-provider-quay/compare/v0.5.3...v0.5.4) (2026-10-19)


### Bug Fixes

* switch to draft release flow for immutable release support ([95b5fca](https://github.com/lannuttia/terraform-provider-quay/commit/95b5fca29ebcf17c95bd39189851434167a354ef))

## [0.5.3](https://github.com/lannuttia/terraform-provider-quay/compare/v0.5.2...v0.5.3) (2026-10-19)


### Bug Fixes

* add missing if to build attestation action ([#39](https://github.com/lannuttia/terraform-provider-quay/issues/39)) ([a33fed9](https://github.com/lannuttia/terraform-provider-quay/commit/a33fed95fb26593d865a83e28e3400291f4efaa5))
* add pull-requests write permission to cd workflow ([b77945a](https://github.com/lannuttia/terraform-provider-quay/commit/b77945a2f7f206fd59e3730f117cfe53358a648e))

## [0.5.2](https://github.com/lannuttia/terraform-provider-quay/compare/v0.5.1...v0.5.2) (2026-10-18)


### Bug Fixes

* upgrade release-please-action to v5.0.0 ([#33](https://github.com/lannuttia/terraform-provider-quay/issues/33)) ([e998a5d](https://github.com/lannuttia/terraform-provider-quay/commit/e998a5d245208cafbee05ecd89f1302835a1b4e5))

## [0.5.1](https://github.com/lannuttia/terraform-provider-quay/compare/v0.5.0...v0.5.1) (2026-10-18)


### Bug Fixes

* use draft releases with force-tag-creation for Terraform Registry compatibility ([#30](https://github.com/lannuttia/terraform-provider-quay/issues/30)) ([40de26a](https://github.com/lannuttia/terraform-provider-quay/commit/40de26a1d46931f66cddc5090223dcb6c3fc6942))

## [0.5.0](https://github.com/lannuttia/terraform-provider-quay/compare/v0.4.2...v0.5.0) (2026-10-18)


### ⚠ BREAKING CHANGES

* replace event string with typed event blocks ([#22](https://github.com/lannuttia/terraform-provider-quay/issues/22))
* scope provider to quay_repository and quay_repository_notification only

### Features

* add configurable HTTP client timeout ([#6](https://github.com/lannuttia/terraform-provider-quay/issues/6)) ([5df8e77](https://github.com/lannuttia/terraform-provider-quay/commit/5df8e77c66bf5cbb75b71200c033544634857eb6))
* add import support to all resources ([086c290](https://github.com/lannuttia/terraform-provider-quay/commit/086c29057b28151f2bb2adb83176bf4965f1feb3))
* add quay_organization resource and data source ([a13ff36](https://github.com/lannuttia/terraform-provider-quay/commit/a13ff36b1214472d0190cef830c295b05780b098))
* replace event string with typed event blocks ([#22](https://github.com/lannuttia/terraform-provider-quay/issues/22)) ([ca99469](https://github.com/lannuttia/terraform-provider-quay/commit/ca99469f4c9b142e6fd7b77dd857fdeb265c4a9b))
* support terraform plan -generate-config-out for all resources ([eb4c24b](https://github.com/lannuttia/terraform-provider-quay/commit/eb4c24ba1ba562d0d1580188de8a30cd8135f163))


### Bug Fixes

* cleanup testcontainer after test ([d7c3974](https://github.com/lannuttia/terraform-provider-quay/commit/d7c3974b46f9d327746d2698678f7fdb735aa286))
* correct eventConfig JSON tag to event_config in CreateNotificationRequest ([719db24](https://github.com/lannuttia/terraform-provider-quay/commit/719db24d9e7586b55fa0ffd9abf98d5a4626b99a))
* correct RELEASE_PLEASE_TOKEN secret name in release workflow ([#10](https://github.com/lannuttia/terraform-provider-quay/issues/10)) ([bdf4c74](https://github.com/lannuttia/terraform-provider-quay/commit/bdf4c74d10e517764392d25b48a92ade337a8a19))
* only trigger release on Tests workflow completion ([#5](https://github.com/lannuttia/terraform-provider-quay/issues/5)) ([7855aa8](https://github.com/lannuttia/terraform-provider-quay/commit/7855aa8e3211c81854b100c19e43d00786ecd03e))
* remove draft release to fix release-please tag creation ([#29](https://github.com/lannuttia/terraform-provider-quay/issues/29)) ([a517725](https://github.com/lannuttia/terraform-provider-quay/commit/a517725c9108f59bc7f0063262729c10ae9a54e0))
* resolve golangci-lint staticcheck QF1008 in quayContainer ([ecfbd96](https://github.com/lannuttia/terraform-provider-quay/commit/ecfbd96c73ac5aac7b748b7dbfb538ef6e712a81))
* treat 403 as destroyed in post-destroy check helpers ([41fc6a4](https://github.com/lannuttia/terraform-provider-quay/commit/41fc6a4d10307eb8bfbdff72e25a1df1654827e5))
* trigger GoReleaser on release published instead of created ([#25](https://github.com/lannuttia/terraform-provider-quay/issues/25)) ([0efaba2](https://github.com/lannuttia/terraform-provider-quay/commit/0efaba2758a61a729695672c73d5520b91125ae7))
* trigger GoReleaser on tag push instead of release published ([#27](https://github.com/lannuttia/terraform-provider-quay/issues/27)) ([823ed7c](https://github.com/lannuttia/terraform-provider-quay/commit/823ed7cc6c92636f63afe9a2c96f95872f251f67))
* use strings.Split in splitImportID to reject import IDs with extra slashes ([2d08598](https://github.com/lannuttia/terraform-provider-quay/commit/2d085981c908a7f20225f1cc20110294f1f45fe1))
* validate visibility at client layer with typed constants ([d7f30d0](https://github.com/lannuttia/terraform-provider-quay/commit/d7f30d03ffec094160e170c88ade201a74edbc36))


### Documentation

* add CONTRIBUTING.md ([03711be](https://github.com/lannuttia/terraform-provider-quay/commit/03711beef653e037da54c5f60e494ad1d6ab05c9))
* correct v0.3.0 changelog entry to remove duplicate v0.2.2 changes ([#14](https://github.com/lannuttia/terraform-provider-quay/issues/14)) ([4c9c7a3](https://github.com/lannuttia/terraform-provider-quay/commit/4c9c7a371c4c886d29002a612d7fd7bd8140def2))
* document intentional camelCase/snake_case asymmetry in notification API ([108b3e0](https://github.com/lannuttia/terraform-provider-quay/commit/108b3e07a781797630c0b035a6a95a11f3355d05))
* generate Terraform Registry documentation ([877a901](https://github.com/lannuttia/terraform-provider-quay/commit/877a901475d82442906b68f8bd504aeb561abae1))
* note that pre-commit hooks require the dev shell to be active ([32188a2](https://github.com/lannuttia/terraform-provider-quay/commit/32188a2e87dd84c6c5076eae69b531a110f54ae0))
* update README to reflect current provider scope ([af108f0](https://github.com/lannuttia/terraform-provider-quay/commit/af108f093709bb43be2aa3ea1d9bed57312aa284))


### Maintenance

* scope provider to quay_repository and quay_repository_notification only ([cd581a1](https://github.com/lannuttia/terraform-provider-quay/commit/cd581a1f61aaefbdc88b0090b638fe228990c217))

## [0.4.2](https://github.com/lannuttia/terraform-provider-quay/compare/v0.4.1...v0.4.2) (2026-10-18)


### Bug Fixes

* trigger GoReleaser on tag push instead of release published ([#27](https://github.com/lannuttia/terraform-provider-quay/issues/27)) ([823ed7c](https://github.com/lannuttia/terraform-provider-quay/commit/823ed7cc6c92636f63afe9a2c96f95872f251f67))

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
