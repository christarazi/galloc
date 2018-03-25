# galloc
An implementation of `malloc(3)` in the Go programming language.

**Disclaimer**: this was created for educational purposes.

### Dependencies

This program is intended to be used on a Unix system because it depends on the
[`unix`](https://godoc.org/golang.org/x/sys/unix) package.

### Usage

Simply download / clone the repo then run:

```shell
$ cd path/to/repo
$ go build && ./malloc
```

## TODO

 - Add implementation for `free(3)`
 - Add implementation for `realloc(3)`
 - Add implementation for `calloc(3)`
 - Allow for use in C using cgo

## License

This program is free software, distributed under the terms of the [GNU] General
Public License as published by the Free Software Foundation, version 3 of the
License (or any later version).  For more information, see the file LICENSE.
