Babyenv
=======

[![GoDoc Badge](https://godoc.org/github.com/magicnumbers/babylogger?status.svg)](http://godoc.org/github.com/magicnumbers/babyenv)

Package babyenv collects environment variables and places them in corresponding
struct fields. It aims to reduce the boilerplate in reading data from the
environment.

The struct should contain `env` tags indicating the names of corresponding
environment variables. The values of those environment variables will be then
collected and placed into the struct. If nothing is found, struct fields will
be given their default values (for example, `bool`s will be `false`).

Default values can also be provided in the `default` tag.


## Example

```go
package main

import (
    "fmt"
    "os"
    "github.com/magicnumbers/babyenv"
)

type config struct {
    Debug bool  `env:"DEBUG"`
    Port string `env:"PORT" default:"8000"`
    Workers int `env:"WORKERS" default:"16"`
}

func main() {
    os.Setenv("DEBUG", "true")
    os.Setenv("WORKERS", "4")

    var cfg config
    if err := babyenv.Parse(&cfg); err != nil {
        log.Fatalf("could not get environment vars: %v", err)
    }

    fmt.Printf("%b\n%s\n%d", cfg.Debug, cfg.Port, cfg.Workers)

    // Output:
    // true
    // 8000
    // 4
}
```


## Supported Types

Currently, only the following types are supported:

* string
* *string
* bool
* *bool
* int
* *int

Pull requests are welcome, especially for new types.


## Credit

This is entirely based on the [caarlos0][carlos]’s [env][carlosenv] package.
This one simply has a slightly different interface, and less functionality.

[carlos]: https://github.com/caarlos0
[carlosenv]: https://github.com/caarlos0/env


## LICENSE

MIT
