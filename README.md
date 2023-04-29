# envgen
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
