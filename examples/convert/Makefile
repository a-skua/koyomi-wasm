SRC := $(shell find src -name "*.rs")
CONFIG := Cargo.toml Cargo.lock
WIT := $(shell find wit -name "*.wit")

PACKAGE_NAME := $(shell cargo metadata --format-version=1 --no-deps | jq -r '.packages[0].name')
TARGET := wasm32-wasip2
RELEASE_BUILD := target/$(TARGET)/release/$(PACKAGE_NAME).wasm

.PHONY: build
build: bin/$(PACKAGE_NAME).wasm

.PHONY: fmt
fmt:
	cargo fmt

.PHONY: clean
clean:
	cargo clean
	rm -rf bin deps wit/deps

wit/deps: wkg.toml wit/*.wit
	wkg wit fetch

deps/a-skua/koyomi.wasm:
	@mkdir -p $(dir $@)
	wkg oci pull -o $@ ghcr.io/a-skua/koyomi:0.1.1-component

deps/example/convert.wasm: $(RELEASE_BUILD)
	@mkdir -p $(dir $@)
	cp $< $@

bin/$(PACKAGE_NAME).wasm: compose.wac deps/a-skua/koyomi.wasm deps/example/convert.wasm
	@mkdir -p $(dir $@)
	wac compose compose.wac -o $@

$(RELEASE_BUILD): wit/deps $(CONFIG) $(SRC) $(WIT)
	cargo build --release --target $(TARGET)
