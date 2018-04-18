// Graceful shutdown is an example of how to gracefully shutdown an HTTP server.
// When requests are in-flight it seems less than ideal to simply stop an HTTP
// server. This is some code that should allow any long-running requests to
// terminate. The time that they are given to terminate may be configured.
//
// Environmental parameters:
//
// PORT defines the port on which the server should listen
//
// WAIT_TIMEOUT_SECONDS defines the timeout (in seconds) for long-running requests
// 
// It uses the context package and the Server Shutdown method introduced in Go 1.8.
//
package main
