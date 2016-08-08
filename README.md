Raphanus - simple Redis-like in-memory cache
--------------------------------------------

## Install

From source:

    go get -u github.com/msoap/raphanus
    # build server & cli
    cd $GOPATH/src/github.com/msoap/raphanus/server && go build -o $GOPATH/bin/raphanus-server
    cd $GOPATH/src/github.com/msoap/raphanus/cli && go build -o $GOPATH/bin/raphanus-cli

Download from github [releases](https://github.com/msoap/raphanus/releases)

From Docker hub:

    docker pull msoap/raphanus

## Run server

    raphanus-server [-port NNNN]

as Docker container:

    docker run --name raphanus --publish 8771:8771 --detach msoap/raphanus

## Use as embed DB

```Go
import (
    "github.com/msoap/raphanus"
)

func main() {
    raph := raphanus.New()
    raph.SetStr("key", "value")
    v, err := raph.GetStr("key")
    if err == raphanus.ErrKeyNotExists {
        ...
    }

    raph.UnderLock(func () {
        v, err := raph.GetStr("k1")
        if err != nil {
            return
        }
        raph.SetStr("k1", v + " updated")
    })

    multiKeys := []string{"k1", "k2"}
    multiVal := []string{}
    raph.UnderRLock(func () {
        for _, k := multiKeys {
            v, err := raph.GetStr(k)
            if err != nil {
                multiVal = multiVal[:]
                return
            }
            multiVal = append(multiVal, v)
        }
    })
}
```

## Use client library for connect with server

```Go
import (
    raphanusclient "github.com/msoap/raphanus/client"
)

func main() {
    raph := raphanusclient.New()
    raph := raphanusclient.NewWithAddr("http://localhost:8770")
    raph.SetStr("key", "value")
    v, err := raph.GetStr("key")
    if err == raphanus.ErrKeyNotExists {
        ...
    }
}
```
