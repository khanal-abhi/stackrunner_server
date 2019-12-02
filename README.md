Stack Runner Server
==

## About
[Stack Runner Server](https://github.com/khanal-abhi/stackrunner_server) is the backend for the [Stack Runner](https://github.com/khanal-abhi/stack-runner), an extension for vs code that helps with stack builds for Haskell development. This server is required for the extension. Please make sure that the binary for the server is available in `$PATH` as `stackrunner_server`, or configure the settings with the proper path. Since this server is written in `go`, please make sure you have the latest version of `go` installed on your system.

<hr>

## Requirements
There are two main requirements for this extension:
- [go](https://golang.org/doc/install)
- [Haskell Stack](https://docs.haskellstack.org/en/stable/README)

<hr>

## Build instructions
To build this binary from source, you may follow the following steps:
````
$ git clone https://github.com/khanal-abhi/stackrunner_server
$ cd stackrunner_server
$ go build .
````

**Please make the binary available in your system path or configure the [Stack Runner](https://github.com/khanal-abhi/stack-runner) extension to point to the binary if planning to use this server with the extension**

<hr>

## Usage
To use the binary, you can run the command
```
$ stackrunner_server <path of stack project>
```

<hr>

## Issues
For any issues related to the server binary or installation, please report the bug here. For any issues related to the extension, please report the bug at [Stack Runner](https://github.com/khanal-abhi/stack-runner).