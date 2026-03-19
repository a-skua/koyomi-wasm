# koyomi-wasm

A WebAssembly Component (WASI Preview 2) for converting between Western calendar (西暦) and Japanese era calendar (和暦).

Built with [TinyGo](https://tinygo.org/) and [goark/koyomi](https://github.com/goark/koyomi).

## WIT Interface

```wit
package a-skua:koyomi@0.1.0;

interface convert {
    enum month { january, february, ... december }
    enum era { meiji, taisho, showa, heisei, reiwa }

    resource western-date {
        constructor(year: s32, month: month, day: u8);
        year: func() -> s32;
        month: func() -> month;
        day: func() -> u8;
        to-wareki: func() -> result<wareki-date, string>;
        to-string: func() -> string;  // "2019-05-01"
    }

    resource wareki-date {
        constructor(era: era, year: s32, month: month, day: u8);
        era: func() -> era;
        year: func() -> s32;
        month: func() -> month;
        day: func() -> u8;
        to-seireki: func() -> result<western-date, string>;
        to-string: func() -> string;  // "令和元年5月1日"
    }
}
```

## Requirements

- [Go 1.25](https://go.dev/)
- [TinyGo 0.40+](https://tinygo.org/)
- [wkg](https://github.com/bytecodealliance/wasm-pkg-tools)

## Build

```sh
make build
```

## Examples

See [examples/convert](examples/convert) for a CLI tool that converts between Western and Japanese era dates.

```sh
wasmtime run bin/convert.wasm -s 2019-05-01
# 令和元年5月1日

wasmtime run bin/convert.wasm -w 令和元年5月1日
# 2019-05-01
```

## License

[Apache License 2.0](LICENSE)
