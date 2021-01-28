Babyenv
=======

[![GoDoc Badge](https://godoc.org/github.com/meowgorithm/babylogger?status.svg)](http://godoc.org/github.com/meowgorithm/babyenv)

Package babyenv collects environment variables and places them in corresponding
struct fields. It aims to reduce the boilerplate in reading data from the
environment.

The struct should contain `env` tags indicating the names of corresponding
environment variables. The values of those environment variables will be then
collected and placed into the struct. If nothing is found, struct fields will
be given their default values (for example, `bool`s will be `false`).

```go
type config struct {
    Name string `env:"NAME"`
}
```

Default values can also be provided in the `default` tag.

```go
    type config struct {
        Name string `env:"NAME" default:"Jane"`
    }
```

A 'required' flag can also be set in the following format:

```go
    type config struct {
        Name string `env:"NAME,required"`
    }
```

If a required flag is set the 'default' tag will be ignored.


## Example

```go
package main

import (
    "fmt"
    "os"
    "github.com/meowgorithm/babyenv"
)

type config struct {
    Debug   bool   `env:"DEBUG"`
    Port    string `env:"PORT" default:"8000"`
    Workers int    `env:"WORKERS" default:"16"`
    Name    string `env:"NAME,required"`
}

func main() {
    os.Setenv("DEBUG", "true")
    os.Setenv("WORKERS", "4")
    os.Setenv("NAME", "Jane")

    var cfg config
    if err := babyenv.Parse(&cfg); err != nil {
        log.Fatalf("could not get environment vars: %v", err)
    }

    fmt.Printf("%b\n%s\n%d\n%s", cfg.Debug, cfg.Port, cfg.Workers, cfg.Name)

    // Output:
    // true
    // 8000
    // 4
    // Jane
}
```


## Supported Types

Currently, only the following types are supported:

* `string`
* `bool`
* `int`
* `int64`
* `[]byte`/`[]uint8`
* `*string`
* `*bool`
* `*int`
* `*int64`
* `*[]byte`/`*[]uint8`

Pull requests are welcome, especially for new types.


## Credit

This is entirely based on [caarlos0][carlos]’s [env][carlosenv] package.
This one simply has a slightly different interface, and less functionality.

[carlos]: https://github.com/caarlos0
[carlosenv]: https://github.com/caarlos0/env


## LICENSE

MIT
