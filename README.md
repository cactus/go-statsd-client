go-statsd-client
================

[![Build Status](https://travis-ci.org/cactus/go-statsd-client.png?branch=master)](https://travis-ci.org/cactus/go-statsd-client)
[![GoDoc](https://godoc.org/github.com/cactus/go-statsd-client/statsd?status.png)](https://godoc.org/github.com/cactus/go-statsd-client/statsd)

## About

A [StatsD][1] client for Go.

## Docs

Viewable online at [godoc.org][2].

## Example

``` go
import (
    "log"

    "github.com/cactus/go-statsd-client/statsd"
)

func main() {
    // first create a client
    client, err := statsd.NewClient("127.0.0.1:8125", "test-client")
    // handle any errors
    if err != nil {
        log.Fatal(err)
    }
    // make sure to clean up
    defer client.Close()

    // Send a stat
    client.Inc("stat1", 42, 1.0)
}
```

See [docs][2] for more info.

## Contributors

See [here][4].

## Alternative Implementations

See the [statsd wiki][5] for some additional client implementations
(scroll down to the Go section).

## License

Released under the [MIT license][3]. See `LICENSE.md` file for details.


[1]: https://github.com/etsy/statsd
[2]: http://godoc.org/github.com/cactus/go-statsd-client/statsd
[3]: http://www.opensource.org/licenses/mit-license.php
[4]: https://github.com/cactus/go-statsd-client/graphs/contributors
[5]: https://github.com/etsy/statsd/wiki#client-implementations
