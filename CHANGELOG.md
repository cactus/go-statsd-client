Changelog
=========

## head

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
