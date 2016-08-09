package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/msoap/raphanus/common"
)

const (
	version        = "0.1"
	defaultAddress = raphanuscommon.DefaultHost + ":" + raphanuscommon.DefaultPort
	usageString    = "raphanus-server [options]\noptions:"
)

type config struct {
	address string // like: "http://host:port"
}

func getConfig() (cfg config) {
	flag.StringVar(&cfg.address, "address", defaultAddress, "address for bind server")
	showVersion := flag.Bool("version", false, "get version")
	flag.Usage = func() {
		fmt.Println(usageString)
		flag.PrintDefaults()
		os.Exit(0)
	}
	flag.Parse()

	if *showVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	return cfg
}

func main() {
	api := newAPI(getConfig())
	api.run()
}
