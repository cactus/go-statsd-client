go-statsd-client
================

[![Build Status](https://github.com/cactus/go-statsd-client/workflows/unit-tests/badge.svg)](https://github.com/cactus/go-statsd-client/actions)
[![GoDoc](https://godoc.org/github.com/cactus/go-statsd-client/statsd?status.png)](https://godoc.org/github.com/cactus/go-statsd-client/statsd)
[![Go Report Card](https://goreportcard.com/badge/cactus/go-statsd-client)](https://goreportcard.com/report/cactus/go-statsd-client)
[![License](https://img.shields.io/github/license/cactus/go-statsd-client.svg)](https://github.com/cactus/go-statsd-client/blob/master/LICENSE.md)

## About

A [StatsD][1] client (UDP) for Go.

## Docs

Viewable online at [godoc.org][2].

## Example

Some examples:

``` go
import (
    "log"

    "github.com/cactus/go-statsd-client/v5/statsd"
)

func main() {
    // First create a client config. Here is a simple config that sends one
    // stat per packet (for compatibility).
    config := &statsd.ClientConfig{
        Address: "127.0.0.1:8125",
        Prefix: "test-client",
    }

    /*
    // This one is for a client that re-resolves the hostname ever 30 seconds.
    // Useful if the address of a hostname changes frequently. Note that this
    // type of client has some additional locking overhead for safety.
    // As such, leave ResInetval as the zero value (previous exmaple) if you
    // don't specifically need this functionality.
    config := &statsd.ClientConfig{
        Address: "127.0.0.1:8125",
        Prefix: "test-client",
        ResInterval: 30 * time.Second,
    }

    // This one is for a buffered client, which sends multiple stats in one
    // packet, is recommended when your server supports it (better performance).
    config := &statsd.ClientConfig{
        Address: "127.0.0.1:8125",
        Prefix: "test-client",
        UseBuffered: true,
        // interval to force flush buffer. full buffers will flush on their own,
        // but for data not frequently sent, a max threshold is useful
        FlushInterval: 300*time.Millisecond,
    }

    // This one is for a buffered resolving client, which sends multiple stats
    // in one packet (like previous example), as well as re-resolving the
    // hostname every 30 seconds.
    config := &statsd.ClientConfig{
        Address: "127.0.0.1:8125",
        Prefix: "test-client",
        ResInterval: 30 * time.Second,
        UseBuffered: true,
        FlushInterval: 300*time.Millisecond,
    }

    // This one is an example of configuring "Tag" support
    // Supported formats are:
    //   InfixComma
    //   InfixSemicolon
    //   SuffixOctothorpe
    // The default, if not otherwise specified, is SuffixOctothorpe.
    config := &statsd.ClientConfig{
        Address: "127.0.0.1:8125",
        Prefix: "test-client",
        ResInterval: 30 * time.Second,
        TagFormat: statsd.InfixSemicolon,
    }
    */

    // Now create the client
    client, err := statsd.NewClientWithConfig(config)

    // and handle any initialization errors
    if err != nil {
        log.Fatal(err)
    }

    // make sure to close to clean up when done, to avoid leaks.
    defer client.Close()

    // Send a stat
    client.Inc("stat1", 42, 1.0)

    // Send a stat with "Tags"
    client.Inc("stat2", 41, 1.0, Tag{"mytag", "tagval"})
}
```

### Legacy Example

A legacy client creation method is still supported. This is retained so as not to break
or interrupt existing integrations.

``` go
import (
    "log"

    "github.com/cactus/go-statsd-client/v5/statsd"
)

func main() {
    // first create a client
    // The basic client sends one stat per packet (for compatibility).
    client, err := statsd.NewClient("127.0.0.1:8125", "test-client")

    // A buffered client, which sends multiple stats in one packet, is
    // recommended when your server supports it (better performance).
    // client, err := statsd.NewBufferedClient("127.0.0.1:8125", "test-client", 300*time.Millisecond, 0)

    // handle any errors
    if err != nil {
        log.Fatal(err)
    }
    // make sure to close to clean up when done, to avoid leaks.
    defer client.Close()

    // Send a stat
    client.Inc("stat1", 42, 1.0)
}
```


See [docs][2] for more info. There is also some additional example code in the
`test-client` directory.

## Contributors

See [here][4].

## Alternative Implementations

See the [statsd wiki][5] for some additional client implementations
(scroll down to the Go section).

## License

Released under the [MIT license][3]. See `LICENSE.md` file for details.


[1]: https://github.com/etsy/statsd
[2]: http://godoc.org/github.com/cactus/go-statsd-client/
[3]: http://www.opensource.org/licenses/mit-license.php
[4]: https://github.com/cactus/go-statsd-client/graphs/contributors
[5]: https://github.com/etsy/statsd/wiki#client-implementations
