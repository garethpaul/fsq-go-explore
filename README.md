# fsq-go-explore

<!-- README-OVERVIEW-IMAGE -->
![Project overview](docs/readme-overview.svg)

## Overview

`garethpaul/fsq-go-explore` is a Go project. Foursquare GoLang Explore Sample

This README is based on the checked-in source, manifests, scripts, and repository metadata on the `master` branch. The project language mix found during review was: Go (14).

## Repository Contents

- `README.md` - project overview and local usage notes
- `fsq` - source or example code
- `limiter` - source or example code
- `SECURITY.md` - security reporting and disclosure guidance
- `static` - source or example code
- `templates` - source or example code
- `VISION.md` - project direction and maintenance guardrails

Additional scan context:

- Source directories: fsq, limiter, static, templates
- Dependency and build manifests: none detected
- Entry points or build surfaces: main.go
- Test-looking files: no obvious test files detected

## Getting Started

### Prerequisites

- Git

### Setup

```bash
git clone https://github.com/garethpaul/fsq-go-explore.git
cd fsq-go-explore
```

The setup commands above are derived from repository files. Legacy mobile, Python, or JavaScript samples may require older SDKs or package versions than a modern workstation uses by default.

## Running or Using the Project

- Run `go run .` or build the module with `go build ./...`.

## Testing and Verification

- No dedicated automated test command was identified from the checked-in files. Verify changes by running the relevant build or manually exercising the sample.

When the required SDK or runtime is unavailable, use static checks and source review first, then verify on a machine that has the matching platform toolchain.

## Configuration and Secrets

- Detected references to Foursquare, Twitter. Keep API keys, OAuth credentials, tokens, and account-specific values in local configuration only.

## Security and Privacy Notes

- Review changes touching authentication or token handling; examples from the scan include auth.go, limiter/config/config.go, limiter/limiter.go, main.go.
- Review changes touching external API calls or credential-adjacent configuration; examples from the scan include auth.go, edit.go, fsq/api.go, fsq/common.go, and 6 more.
- Review changes touching network requests, sockets, or service endpoints; examples from the scan include auth.go, fsq/api.go, fsq/common.go, fsq/venue.go, and 3 more.
- Review changes touching file, media, JSON, XML, CSV, OCR, or data parsing; examples from the scan include auth.go, fsq/api.go, fsq/common.go, fsq/keys.go, and 3 more.

## Maintenance Notes

- See `SECURITY.md` for vulnerability reporting and safe research guidance.
- See `VISION.md` for project direction and contribution guardrails.

## Contributing

Keep changes small and tied to the project that is already present in this repository. For code changes, document the toolchain used, avoid committing generated dependency directories or local configuration, and update this README when setup or verification steps change.

## Existing Project Notes

Prior README summary:

> fsq-go-explore Foursquare-Go-Explore =================== Getting Started ---------- Here are some steps to get you started. 1. Get your developer keys
