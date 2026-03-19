# Example: convert

A CLI tool that converts between Western calendar (西暦) and Japanese era calendar (和暦) using the `a-skua:koyomi/convert` component.

## Requirements

- [Rust](https://www.rust-lang.org/) (wasm32-wasip2 target)
- [wkg](https://github.com/bytecodealliance/wasm-pkg-tools)
- [wac](https://github.com/bytecodealliance/wac)
- [wasmtime](https://wasmtime.dev/)

## Setup

Configure the wkg registry for the `a-skua` namespace:

```sh
wkg config --edit
```

```toml
[namespace_registries.a-skua]
registry = "ghcr.io"

[namespace_registries.a-skua.metadata]
preferredProtocol = "oci"
oci = { registry = "ghcr.io" }
```

## Build

```sh
make build
```

## Usage

Western to Japanese era (西暦 → 和暦):

```sh
wasmtime run bin/convert.wasm -- -s 2019-05-01
# 令和元年5月1日
```

Japanese era to Western (和暦 → 西暦):

```sh
wasmtime run bin/convert.wasm -- -w 令和元年5月1日
# 2019-05-01

wasmtime run bin/convert.wasm -- -w 令和1年5月1日
# 2019-05-01
```
