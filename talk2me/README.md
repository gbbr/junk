talk2me
=======

Server & Client chat via raw TCP through XML (warning! this is not secure) 

Start server first: `go run server.go`  
Optionally, an argument may be specified to define the host and port to listen on, for example `./server -addr=HOST:PORT`. By default, it listens on localhost port 1234.


Then, start client: `go run client.go`  
Messages will be sent to all connected clients and are read from standard I/O. Optionally the `addr` flag is available to conenct to a custom host/port set-up on the server.
