

# meserve

This is a simple HTTP server written in Go.
Currently, the program allows basic file browsing within a shared part of your filesystem.

## Basic Usage

The simplest way to run the server:

```./meserve```

Run server on specified port:

```./meserve --port 3000```


## To-do
- [x] - create alpha version of server
- [x] - create initial fileHandler
- [x] - create initial uploadHandler
- [x] - add way to select root HTTP directory
- [] - add tests
- [x] - split code into more packages
- [] - put managment of terminal args out of main

To-do terminal args:
- [] allow select port
- [] enable/disable upload
- [] select server root dir
