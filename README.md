# Graceful Shutdown

[![Build status](https://travis-ci.org/mramshaw/Graceful-Shutdown.svg?branch=master)](https://travis-ci.org/mramshaw/Graceful-Shutdown)
[![Go Report Card](https://goreportcard.com/badge/github.com/mramshaw/Graceful-Shutdown?style=flat-square)](https://goreportcard.com/badge/github.com/mramshaw/Graceful-Shutdown)
[![GoDoc](https://godoc.org/github.com/mramshaw/Graceful-Shutdown?status.svg)](https://godoc.org/github.com/mramshaw/Graceful-Shutdown)
[![GitHub release](https://img.shields.io/github/release/mramshaw/Graceful-Shutdown.svg?style=flat-square)](https://github.com/mramshaw/Graceful-Shutdown/releases)

Gracefully shut down a Golang web server


## Motivation

Shutting down a server when requests are in-flight really doesn't sound ideal.

Kubernetes generally has a gradual shutdown of its containers (30 seconds or so),
so gracefully shutting down a server may not be absolutely necessary in a cloud
(meaning Kubernetes) environment, but even so it's probably still a good practice.

The `Server.Shutdown` and `Server.Close` methods were added to Golang in
[Go 1.8](https://golang.org/doc/go1.8#http_shutdown).

This code was heavily inspired by the
[Gorilla/mux](https://github.com/Gorilla/mux#graceful-shutdown) code.

#### Go versions prior to 1.8

For Go versions prior to __1.8__ there is [graceful](https://github.com/tylerb/graceful).

Full circle: __graceful__ recommends using the [Shutdown](https://golang.org/pkg/net/http/#Server.Shutdown)
approach listed here for Golang versions __1.8__ and later.

This particular
[caveat from graceful](https://github.com/tylerb/graceful#important-things-to-note-when-setting-timeout-to-0)
is well worth a read.


## Principles

I am gradually coming around to the [12-Factor](https://12factor.net/config) way
of doing things, accordingly our runtime variables `PORT` and `WAIT_TIMEOUT_SECONDS`
will be defined as environment variables (this has significant advantages in
Docker container deployments, as well as when working with `docker-compose`).

In practice, this looks like this:

    $ PORT=8080 WAIT_TIMEOUT_SECONDS=30 go run graceful_shutdown.go 

We will use the __log__ package to give us our actual event timings.


## Testing

Almost all of the code came from __Gorilla/mux__, but how to test that it works?

What I came up with can probably be improved, but it works so I'm happy with
it for now. What we are really trying to do is model an inconsiderate __client__,
however we can actually force the client to be inconsiderate on the __server__
side, which seems to be a little more satisfactory from a cohesiveness point
of view.

We will create a `/timer` endpoint that will simply `sleep` for a minute.

The plan then will be to open a browser to our __/timer__ route and then hit
__Ctrl-C__ during this minute. We should then see our shutdown period gracefully
enforced.

We will specify port 8080, which means the our browser endpoint will be:

    http://localhost:8080/timer


## Request in-flight

WAIT_TIMEOUT_SECONDS not specified, one request in-flight:

    $ PORT=8080 go run graceful_shutdown.go 
    2018/03/14 23:40:59 Invalid WAIT_TIMEOUT_SECONDS, setting to 15 seconds
    2018/03/14 23:40:59 Listening on http://localhost:8080 ...
    ^C2018/03/14 23:41:15 Shutdown request (Ctrl-C) caught
    2018/03/14 23:41:15 http: Server closed
    2018/03/14 23:41:30 Shutting down ...
    $

Server shuts down after 15 seconds, as expected. Good!


## Smoke test

WAIT_TIMEOUT_SECONDS not specified, no requests in-flight:

    $ PORT=8080 go run graceful_shutdown.go 
    2018/03/14 23:42:01 Invalid WAIT_TIMEOUT_SECONDS, setting to 15 seconds
    2018/03/14 23:42:01 Listening on http://localhost:8080 ...
    ^C2018/03/14 23:42:03 Shutdown request (Ctrl-C) caught
    2018/03/14 23:42:03 Shutting down ...
    $

Server shuts down immediately. Fine.


## Longer timeout, request in-flight

WAIT_TIMEOUT_SECONDS specified, no requests in-flight:

    $ PORT=8080 WAIT_TIMEOUT_SECONDS=30 go run graceful_shutdown.go 
    2018/03/14 23:43:09 Listening on http://localhost:8080 ...
    ^C2018/03/14 23:43:11 Shutdown request (Ctrl-C) caught
    2018/03/14 23:43:11 Shutting down ...
    $

Server again shuts down immediately. As expected.


## Longer timeout

WAIT_TIMEOUT_SECONDS specified, one request in-flight:

    $ PORT=8080 WAIT_TIMEOUT_SECONDS=30 go run graceful_shutdown.go 
    2018/03/14 23:43:15 Listening on http://localhost:8080 ...
    ^C2018/03/14 23:43:21 Shutdown request (Ctrl-C) caught
    2018/03/14 23:43:21 http: Server closed
    2018/03/14 23:43:51 Shutting down ...
    $

Server shuts down after 30 seconds, as expected. Excellent!


## To Do

- [x] Add notes for __graceful__ for Go versions prior to __1.8__
- [ ] Investigate edge cases (browser shut down while server timing-out, timeout longer then server sleep)
- [ ] Research the subject more carefully (Just For Func, etc)
- [ ] Figure out how to automate the testing


## Credits

The server code was heavily inspired by the excellent Gorilla/mux:

    https://github.com/Gorilla/mux#graceful-shutdown
