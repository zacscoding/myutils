# My utils commands
this project is my command utils for dev :)  

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

---  

<div id="ssh_command"></div>

> ## SSH command  

```bash
$ myutils host add
```