# mdk
Cloud Computing 15619 - Twitter Analytics Web Service

## mdk_go
1. Install git on ubuntu
2. Clone this project
```
$ cd mdk
```
```
$ tar -C /usr/local -xzf go1.5.1.linux-amd64.tar.gz
```
```
$ export PATH=$PATH:/usr/local/go/bin
```
```
$ GOPATH=$(pwd)/mdk_go
```
```
$ export GOBIN=$GOPATH/bin
```
```
$ cd $GOPATH/src/mdk
```
```
$ if above does not word, add the following to ~/.bashrc and source ~/.bashrc
  export GOROOT=/usr/local/go
  export GOPATH=$HOME/mdk/mdk_go
  export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
```
```
$ go get
```
```
$ cd $GOPATH
```
```
$ sudo ./bin/mdk
```
