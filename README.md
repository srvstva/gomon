### **gomon** - tcp based client server for running remote commands

`gomon` uses tcp socket to start and run remote commands. Any no of `gomon` client can connect to the server and request for running commands.

### Usage
---
```
➜  ~ gomon
usage: gomon <command> [<args>]
The most common gmon commands are
    serve    start the gomon server
  connect    connect to the gomon server
```
#### Starting the server
The server accepts a command line flag to run on a custom port. For details see `gomon serve -h`
```
gomon serve [-port=port]
```

#### Connecting to `gomon` server and running command
```
gomon connect -remotePort <port> -remoteHost <host/ip> -command <command>
```

### Examples
```
➜  ~ gomon serve
2019/06/28 11:12:16 server listening on [::]:7891

➜  ~ gomon connect -remoteHost localhost -remotePort 7891 -command 'uname -a'
2019/06/28 11:13:02 Connected to 127.0.0.1:7891
Linux Vostro-15-3568 4.15.0-52-generic #56-Ubuntu SMP Tue Jun 4 22:49:08 UTC 2019 x86_64 x86_64 x86_64 GNU/Linux

```

### How to build and install
-   Download or clone this [repository](https://github.com/srvstva/gomon)
-   `cd gomon && go build && go install`
-   Or run command `go get https://github.com/srvstva/gomon`
