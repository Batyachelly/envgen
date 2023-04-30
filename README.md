# envgen

[![codebeat badge](https://codebeat.co/badges/a58abcf6-7138-4a6e-905c-a6b09b650deb)](https://codebeat.co/projects/github-com-batyachelly-envgen-main)
[![Go Report Card](https://goreportcard.com/badge/github.com/Batyachelly/envgen)](https://goreportcard.com/report/github.com/Batyachelly/envgen)
[![Go Doc](https://godoc.org/github.com/Batyachelly/envgen?status.svg)](https://godoc.org/github.com/Batyachelly/envgen)
[![Release](https://img.shields.io/github/v/release/Batyachelly/envgen.svg?style=flat-square)](https://github.com/Batyachelly/envgen/releases)

Utility to parse golang env structures into .env files.

## Getting started

1. Download swag by using:
```shell
go install github.com/Batyachelly/envgen/cmd/envgen@latest
```

2. Run `envgen` in the project's root folder which contains a config files.
```shell
envgen -target=example/config.go -structs=Config,AnotherConfig1,AnotherConfig2 -output_dir example/generated
```
