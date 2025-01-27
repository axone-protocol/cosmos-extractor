# Cosmos Extractor changelog

## 1.0.0 (2025-01-27)


### Features

* **cli:** add log_level option for customizable logging ([1cfc352](https://github.com/axone-protocol/cosmos-extractor/commit/1cfc3525ff36d1d6932ce594ff800ba839b20fbd))
* **cli:** set default log level to warn ([d4900c6](https://github.com/axone-protocol/cosmos-extractor/commit/d4900c67755e2613388f0907ccc5470c9803f8a9))
* **delegators:** add --hrp flag for Bech32 address delegators ([7f18288](https://github.com/axone-protocol/cosmos-extractor/commit/7f182884c1de5b477f904a34b01fd7e368c811b8))
* **delegators:** add --min-shares / --max-shares filters ([7f00c50](https://github.com/axone-protocol/cosmos-extractor/commit/7f00c500875f6995059dbe4ad6813e1d9f5d0a32))
* **delegators:** add chain metadata export ([d0ba97a](https://github.com/axone-protocol/cosmos-extractor/commit/d0ba97ae6009551574ecb02a2901969710813fb7))
* **delegators:** add delegators extraction ([4b9fa72](https://github.com/axone-protocol/cosmos-extractor/commit/4b9fa72d793343ab3cb8e9eddb4a3c40abb93c05))
* **docs:** add docs command to generate documentation ([8dc7a6a](https://github.com/axone-protocol/cosmos-extractor/commit/8dc7a6a6bc80e55539d126103b13bad20815782f))
* **extract:** introduce --output flag ([dbc87a1](https://github.com/axone-protocol/cosmos-extractor/commit/dbc87a1bc583a4bf19383c68599b1a46d04e88ed))
* **keeper:** add keepers management (auth, bank, staking) ([9d973b9](https://github.com/axone-protocol/cosmos-extractor/commit/9d973b9892026df66850576d81ef9ae0b705918f))


### Bug Fixes

* **delegators:** prevent duplicate addresses ([9ff3560](https://github.com/axone-protocol/cosmos-extractor/commit/9ff3560acb0d3243c0589ecbad6c9ad90fafa755))


### Performance Improvements

* **delegators:** significantly improve extraction performance ([03ec55e](https://github.com/axone-protocol/cosmos-extractor/commit/03ec55e3355353243b12af0460e90e8c1e308331))
