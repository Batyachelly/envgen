.PHONY: install
install:
	go install ./cmd/envgen

.PHONY: lint
lint:
	golangci-lint run --config .golangci.yml ./...


.PHONY: run_example
run_example:
	go run cmd/envgen/main.go -target=example/config.go -structs=Config,AnotherConfig1,AnotherConfig2 -output_dir=example/generated