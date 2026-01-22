.PHONY: example build-plugin run-host

example: build-plugin run-host

build-plugin:
	@echo "Building example plugin..."
	cd examples/plugin && GOOS=wasip1 GOARCH=wasm go build -buildmode=c-shared -o plugin.wasm .

run-host:
	@echo "Running example host runtime..."
	cd examples/host-runtime && go run . ../plugin/plugin.wasm