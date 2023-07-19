Changelog
=========

## head

## 5.1.0 2023-07-19
*   Fix for tag format in substatter. (GH-55)
*   Add support for Floats in several situations. (GH-57)
*   Add new ExtendedStatSender interface for the new Float methods.

## 5.0.0 2021-01-13
*   Add Tag support: suffix-octothorpe, infix-comma, infix-semicolon (GH-53)
*   Remove previously deprecated NoopClient. Use a nil `*Client` Statter as a
    replacement, if needed. Ex:
    ```
    var client Client
    // A nil *Client has noop behavior, so this is safe.
    // It will become a small overhead (just a couple function calls) noop.
    err = client.Inc("stat1", 42, 1.0)
    ```

## 4.0.0 2020-11-05
*   Fix go.mod versioning. (GH-51,GH-52)
*   Bump major version for go.mod change, just in an attempt to be safer
    for existing users.

## 3.2.1 2020-06-23
*   Export NewBufferedSenderWithSender for direct use where needed.

## 3.2.0 2019-09-21
*   A new client constructor with "config style" semantics.
    "legacy" client construction still supported, to retain backwards compat.
*   Add an optional re-resolving client configuration. This sets a schedule for
    having the client periodically re-resolve the addr to ip. This does add some
    overhead, so best used only when necessary.

## 3.1.1 2018-01-19
*   avoid some overhead by not using defer for two "hot" path funcs
*   Fix leak on sender create with unresolvable destination (GH-34).

## 3.1.0 2016-05-30
*   `NewClientWithSender(Sender, string) (Statter, error)` method added to
    enable building a Client from a prefix and an already created Sender.
*   Add stat recording sender in submodule statsdtest (GH-32).
*   Add an example helper stat validation function.
*   Change the way scope joins are done (GH-26).
*   Reorder some structs to avoid middle padding.

## 3.0.3 2016-02-18
*   make sampler function tunable (GH-24)

## 3.0.2 2016-01-13
*   reduce memory allocations
*   improve performance of buffered clients

## 3.0.1 2016-01-01
*   documentation typo fixes
*   fix possible race condition with `buffered_sender` send/close.

## 3.0.0 2015-12-04
*   add substatter support

## 2.0.2 2015-10-16
*   remove trailing newline in buffered sends to avoid etsy statsd log messages
*   minor internal code reorganization for clarity (no api changes)

## 2.0.1 2015-07-12
*   Add Set and SetInt funcs to support Sets
*   Properly flush BufferedSender on close (bugfix)
*   Add TimingDuration with support for sub-millisecond timing
*   fewer allocations, better performance of BufferedClient

## 2.0.0 2015-03-19
*   BufferedClient - send multiple stats at once
*   clean up godocs
*   clean up interfaces -- BREAKING CHANGE: for users who previously defined
    types as *Client instead of the Statter interface type.

## 1.0.1 2015-03-19
*   BufferedClient - send multiple stats at once

## 1.0.0 2015-02-04
*   tag a version as fix for GH-8
