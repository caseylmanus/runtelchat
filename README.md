# runtelchat

runtelchat is a chat server made for a coding exersize. This is a multi-socket TCP chat server that will accept tcp connections on the configured port and relay messages between multiple tcp clients.

runtelhat is a play on words of the telnet client meant to be used as a redumentory client

The main entry point is in /cmd/runtelchat.

## Getting Started

Once you either clone this repository or download the .zip file, navigate to {project root}/cmd/runtelchat and perform go build.  Then you may execute ./runtelchat 

### Prerequisites

Only requires Go and a configured GOPATH

### Configuration

A default configuration file is provided, if no configuration file exists, then a configuration identical to this default file will be used. The config file should be named config.json and be in the same directory as the compiled application.  This app uses a .json configuration file.

#### Configuration Items
* host - IP or Hostname of the machine to accept connections on
* ports - An array of ports to use
* logFilePath - location of the log file


## Running the tests

This is a normal Go project so tests are in the root directory and you make execute them with "go test". I'll admit the test suite is pretty thin at this point. 


## Dependencies and Third Party Code
Right now all dependencies are vendored into the vendor folder, so no action should be needed for them.

### Lumberjack
This lib is used for rolling logging, simplifing logging
[GoDoc](https://godoc.org/gopkg.in/natefinch/lumberjack.v2)

### pkg/errors
The errors package is used to provide error wrapping
[GoDoc](https://godoc.org/github.com/pkg/errors)

