# docker-utils

Collection of utility commands that extend Docker CLI.

In order to use, you need to have godep utility. If you don't have it already, type `go get github.com/tools/godep` into your terminal. After that, type:
```
go get -d github.com/tkopczynski/docker-utils
cd $GOPATH/src/github.com/tkopczynski/docker-utils
godep go install .
```
  
Now docker-utils should be available on PATH. Type `docker-utils -h` to see usage options.
