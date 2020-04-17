# My utils commands
this project is my command utils for dev :)  
this module is using leveldb for persistence in `$HOME/myutils` directory.

## Getting started  

There are two options to install `myutils`  

> Option1. Go get and install  

```bash
$ go get -u github.com/zacscoding/myutils/...
$ myutils
```  

> Option2. Go install (with git clone)  

```bash
$ mkdir -p $GOPATH/src/github.com/zacscoding
$ cd $GOPATH/src/github.com/zacscoding
$ git clone https://github.com/zacscoding/myutils.git
$ cd myutils/cmd/myutils
$ go install
```  

> Check myutils in $GOBIN

```bash
$ ls -la $GOBIN/myutils 
-rwxrwxr-x 1 app app 9417592 Nov  7 01:29 /home/app/go/bin/myutils
```

Don't forget to add `$GOPATH/bin` to ur `$PATH`.  

---  

## Commands  

- <a href="#host_command">host command</a>  
; host command is manage hosts such as save,update,get,remove.  
- <a href="#ssh_command">ssh command</a>  
; ssh command is utils for remote vm.

---  

<div id="host_command"></div>  

> ## Manage hosts command

```bash
$ ./bin/myutils host
NAME:
   myutils host - manage hosts such as add | get | gets | update | delete

USAGE:
   myutils host command [command options] [arguments...]

COMMANDS:
   import  Import hosts json file to local store.
   export  Export hosts json from local store
   add     Adds a host
   get     Get a host
   gets    Get hosts
   update  Update a host
   delete  Delete a host

OPTIONS:
   --help, -h  show help
```

> ## Example of hosts command  

>   

  

---  

<div id="ssh_command"></div>

> ## SSH command  

```bash
$ myutils host add
```