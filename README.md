# GoCore

A Golang solution of tools for building a full stack web application.

## Goals of the Project

Below are some targeted goals:

- Webserver Goals
  - Http & Https & HTTP 2.0 with Golang 1.9 and [gin-gonic/gin](https://github.com/gin-gonic/gin). [See GoCore/app documentation](https://pkg.go.dev/github.com/DanielRenne/GoCore/core/app))
    - Basic [configuration helpers with a gin-gonic server](https://pkg.go.dev/github.com/DanielRenne/GoCore/core/ginServer/#example_ConfigureGin)
    - Exposing [the gin-gonic Router engine](https://github.com/DanielRenne/GoCore/blob/master/doc/Basic_GinRouter.md) to your application for custom routes.
  - Setting up [dynamic routes through controller registration APIs (in the app/api package)](https://github.com/DanielRenne/GoCore/blob/master/doc/Controller_Registration_With_Api.md) and reflection to invoke your methods with interfaces
  - Websocket WSS and WS support through [github.com/gorilla/websocket](https://github.com/gorilla/websocket)
    - Replying to messages or Broadcasting to all
    - Giving you access to iterate your connected sockets
    - Managing deleted sockets and providing a way to safely iterate through connected sockets to publish messages to all or some sockets

---

- Database Goals. Provide Model/structs/ORM support and drivers for the following (for use with [GoCore full apps](https://github.com/DanielRenne/GoCore/blob/master/doc/GoCore_Full.md) or [standalone ORM](https://github.com/DanielRenne/GoCore/tree/master/core/dbServices/example)):
  - Supported databases:
    - MongoDB
    - BoltDB
  - Create SQL Schema (DDL) from JSON Configuration.
    - Generated golang structs and methods will also allow customization files to be injected inside your models/v1/model package
  - Create Golang ORM packages for RDBMS Transactions & Queries.
  - Create a [bootstrapping system](https://github.com/DanielRenne/GoCore/blob/master/doc/Bootstrap.md) to seed data in various configurations and data dumping formats
  - Recursive Joins with foreign and primary keys in mongo or bolt.
  - A pubsub store (core/store) for mongo or bolt to allow for interfaces to subscribe to changes in the database or to save changes to the database with either golang or a javascript client.

---

- General application toolbox and file utilities
  - Some basic crypto functions in the [github.com/DanielRenne/core/crypto](https://pkg.go.dev/github.com/DanielRenne/GoCore/core/crypto) package
  - Utility functions for versioning, hexadecimals, strings, human directory sizes, printable ascii/emojis, data type conversions inside of [github.com/DanielRenne/GoCore/core/extensions](https://pkg.go.dev/github.com/DanielRenne/GoCore/core/extensions)
  - Basic path helpers in [github.com/DanielRenne/GoCore/core/path](https://pkg.go.dev/github.com/DanielRenne/GoCore/core/path)
  - Utility functions for managing files and directories for getting all files in directories, copying and removing directories, reading files and also parsing interfaces into json, unGizipping files, and Untarring and Taring files natively inside of go with [github.com/DanielRenne/core/extensions/](https://pkg.go.dev/github.com/DanielRenne/GoCore/core/extensions) `fileExtensions.go` and `unix.go`
  - A simple logger with logging with colors, [goRoutine logging](https://pkg.go.dev/github.com/DanielRenne/GoCore/core/logger#example_GoRoutineLogger), [tailing files](https://pkg.go.dev/github.com/DanielRenne/GoCore/core/logger#Tail), and [measuring time of execution](https://pkg.go.dev/github.com/DanielRenne/GoCore/core/logger#example_TimeTrack) in [github.com/DanielRenne/core/logger](https://pkg.go.dev/github.com/DanielRenne/GoCore/core/logger) package
  - A ["core" package](https://pkg.go.dev/github.com/DanielRenne/GoCore/core) which currently has some debug Dump tools for dumping structs and variables to the console in a readable format
  - A [cron package](https://pkg.go.dev/github.com/DanielRenne/GoCore/core/cron) which sets up a ticker to run a function at a specified interval such as top of second (we never needed sub-second crons but can add it if you need it as we currently ticker every 100ms), minute, hour, day and helper functions to run at top of 15 minutes, 5 minutes [dynamically defined job](https://pkg.go.dev/github.com/DanielRenne/GoCore/core/cron#example-ExampleRegisterRecurring) and [even running only a job once](https://pkg.go.dev/github.com/DanielRenne/GoCore/core/cron#example-ExecuteOneTimeJob) (if successfully returned true)

---

- Channel management queues, pubsub functions, shell utilities and even a go worker package to assist spreading load
  - A [simple channel management queue system](https://pkg.go.dev/github.com/DanielRenne/GoCore/core/pubsub#example_Signal) for managing goRoutines and channels in [github.com/DanielRenne/GoCore/core/channels](https://pkg.go.dev/github.com/DanielRenne/GoCore/core/channels) package
  - pubsub package for [publishing to subscribers](https://pkg.go.dev/github.com/DanielRenne/GoCore/core/pubsub#example_Publish) in [github.com/DanielRenne/GoCore/core/pubsub](https://pkg.go.dev/github.com/DanielRenne/GoCore/core/pubsub) package
  - cmdExec package for [easily extracting sdtout and stderr and invoking system binaries in github.com/DanielRenne/GoCore/core/cmdExec](https://pkg.go.dev/github.com/DanielRenne/GoCore/core/cmdExec) package
  - A workqueue package simplifying how to append batches of function calls to N number of workers. Especially nice for long running tasks where you can either wait and block or be signaled when a job completes [github.com/DanielRenne/GoCore/core/workQueue](https://pkg.go.dev/github.com/DanielRenne/GoCore/core/workQueue)

---

- Atomic file locking functions (on many common types) for file system operations on thread safe files in [github.com/DanielRenne/GoCore/core/atomicTypes](https://pkg.go.dev/github.com/DanielRenne/GoCore/core/atomicTypes) package
  - Adds a Get() and Set() method with a mutex lock to the following types:
    - AtomicString
    - AtomicUInt16
    - AtomicUInt32
    - AtomicByteArray
    - AtomicFloat64
    - AtomicBoolArray
    - AtomicBool (ToggleTrue returning if changed to true)
    - AtomicInt (Add, Increment, Decrement)

---

## Get Started with GoCore

1.  To start a new project with go modules (after go 1.13) run the following steps in a new console window. Note, this just gets all packages indirectly and they will be removed in your app as you begin to use them.

```
go mod init yourProject/packageName
```

If you want to just play with all packages run:

    go get github.com/DanielRenne/GoCore

Otherwise [read the docs](https://pkg.go.dev/github.com/DanielRenne/GoCore) and see if anything adds value to your work and go get individual packages.

## Build GoCore Backend Only Webserver app

There are three options to start a webserver. GoCoreLite (just a gin-gonic server with a gorilla websocket where you pass the port you wish), GoCoreFull ( which assumes usages of our model and ORM with mongo or boltDB ), or GoCoreCreateApp (full front-end examples with a backend webserver).

- GoCore full docs are available at [here](https://github.com/DanielRenne/GoCore/blob/master/doc/GoCore_Full.md)

- GoCoreLite full docs are available [here](https://github.com/DanielRenne/GoCore/blob/master/doc/GoCore_Lite.md)

- GoCoreCreateApp full docs are available [here](https://github.com/DanielRenne/GoCore/blob/master/doc/FrontEnd_BackEnd.md)

## FAQ

### Why cant I go run my main in a full GoCore application

This is because in many cases for a full GoCore web app, we need to read a webConfig.json in your current directory so that you dont compile configurations inside your main.go and developers and servers can be reconfigured without recompiling. Go will compile into a tmp directory and we dont know where your webConfig.json is located. If you really want to use go run, you can pass `go run main.go $(pwd)` if you dont want to compile your web server which will pass the location where you webConfig.json sits next to your main program.

### If you decide not to use the web server functionality and want to try out some other helper utilities outlined in our goals, our main documentation for the codebase located here: [https://pkg.go.dev/github.com/DanielRenne/GoCore](https://pkg.go.dev/github.com/DanielRenne/GoCore)
