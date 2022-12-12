# k8sh

Simple POSIX inspired shell for Kubernetes operations.

## Getting Started

1. go run cmd/k8sh/k8sh.go

## TODO:
- expansion and parsing

  [ ] variable expansion within quotes
  
  [ ] variable assignment
  
  [ ] read only or internal symbols (can't be set)
  
    - this will likely be replaced by a shell context

  [ ] tilde expansion
  
    - delay until a better concept of $HOME exists

  [ ] posix parameter expansion
  
  [ ] functions
  
  [ ] command argument parser

- builtins

  [ ] ls basic, '-l' only
  
  [ ] cd basic
  
  [ ] watch
  
  [ ] stat (dumps object)
  
  [ ] mkdir (creates namespace)
  
  [ ] cp
  
  [ ] mv
  
  [ ] rm
  
  [ ] history
  
    - save each command
    - no arg print history
    - arg integer to replay command

- plugins

  [ ] add new commands

- shell

  [ ] shell context for state
  
    - CWD is managed and guarded by shell

  [ ] each command is given a shell handle
  
  [ ] shell holds client and can be constructed with a mock
