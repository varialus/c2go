c2go
====

Golang Translation of [c2go.py][1]

Status
------

**Rough Translation Complete**

 - c2go.py: Python global dictionaries translated to Go maps

**Roadmap**

 - Golang Rough Translation of c2go.py
 - Golang Rough Translation of pycparser using Golang [Yacc][2] and [golex][3] instead of PLY
 - Pass c2go.py Working Tests
 - Pass c2go.py Not Yet Working Tests
 - Parse Large Complex C Projects and Successfully Translate to Go

Setup
-----

 1. [Install and configure Go][4]
 2. export GOPATH=$HOME/go
 3. go get -d github.com/varialus/c2go
 4. cd ~/go/src/github.com/varialus/c2go/
 5. go run main.go

[1]:http://github.com/xyproto/c2go
[2]:http://golang.org/cmd/yacc/
[3]:http://github.com/cznic/golex
[4]:http://golang.org/doc/install
