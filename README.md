# Cosmos Extractor

> CLI tool for extracting diverse data from Cosmos chain snapshots

[![version](https://img.shields.io/github/v/release/axone-protocol/cosmos-extractor?style=for-the-badge&logo=github)](https://github.com/axone-protocol/cosmos-extractor/releases)
[![lint](https://img.shields.io/github/actions/workflow/status/axone-protocol/cosmos-extractor/lint.yml?branch=main&label=lint&style=for-the-badge&logo=github)](https://github.com/axone-protocol/cosmos-extractor/actions/workflows/lint.yml)
[![build](https://img.shields.io/github/actions/workflow/status/axone-protocol/cosmos-extractor/build.yml?branch=main&label=build&style=for-the-badge&logo=github)](https://github.com/axone-protocol/cosmos-extractor/actions/workflows/build.yml)
[![test](https://img.shields.io/github/actions/workflow/status/axone-protocol/cosmos-extractor/test.yml?branch=main&label=test&style=for-the-badge&logo=github)](https://github.com/axone-protocol/cosmos-extractor/actions/workflows/test.yml)
[![codecov](https://img.shields.io/codecov/c/github/axone-protocol/cosmos-extractor?style=for-the-badge&token=6NL9ICGZQS&logo=codecov)](https://codecov.io/gh/axone-protocol/cosmos-extractor)
[![conventional commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg?style=for-the-badge&logo=conventionalcommits)](https://conventionalcommits.org)
[![semantic-release](https://img.shields.io/badge/%20%20%F0%9F%93%A6%F0%9F%9A%80-semantic--release-e10079.svg?style=for-the-badge)](https://github.com/semantic-release/semantic-release)
[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-2.1-4baaaa.svg?style=for-the-badge)](https://github.com/okp4/.github/blob/main/CODE_OF_CONDUCT.md)
[![License](https://img.shields.io/badge/License-BSD_3--Clause-blue.svg?style=for-the-badge)](https://opensource.org/licenses/BSD-3-Clause)

## Purpose

This CLI tool offers a straightforward way to extract diverse data from Cosmos chains using the application file store.

Initially developed for Axone’s internal needs, it provides a versatile foundation for addressing various data extraction requirements.

## Features

- Export chain store information.
- Export delegators and their delegations.

## Usage example

```bash
# Export delegators and their delegations to a CSV file
$ cosmos-extractor extract delegators ./data/bitsong/data ./output --chain-name bitsong
```
