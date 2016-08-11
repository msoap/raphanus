package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/msoap/raphanus/common"
)

const (
	defaultAddress = raphanuscommon.DefaultHost + ":" + raphanuscommon.DefaultPort
	usageString    = "raphanus-server [options]\noptions:"
)

type config struct {
	address  string // like: "http://host:port"
	user     string // for HTTP basic authentication
	password string
	filename string // file for storage on disk
	syncTime int    // time in seconds between sync on disk
	logging  bool   // log calls
}

func getConfig() (cfg config) {
	flag.StringVar(&cfg.address, "address", defaultAddress, "address for bind server")
	flag.StringVar(&cfg.filename, "filename", "", "file name for storage on disk, '' - for work in-memory only")
	flag.IntVar(&cfg.syncTime, "sync-time", 0, "time in seconds between sync on disk")
	flag.BoolVar(&cfg.logging, "logging", false, "log calls")
	authUserPass := flag.String("auth", "", "user:password for enable HTTP basic authentication")

	showVersion := flag.Bool("version", false, "get version")
	flag.Usage = func() {
		fmt.Println(usageString)
		flag.PrintDefaults()
		os.Exit(0)
	}
	flag.Parse()

	if *showVersion {
		fmt.Println(raphanuscommon.Version)
		os.Exit(0)
	}

	if len(*authUserPass) > 0 {
		// TODO: allow ":" in password
		auth := strings.Split(*authUserPass, ":")
		if len(auth) == 2 && len(auth[0]) > 0 && len(auth[1]) > 0 {
			cfg.user, cfg.password = auth[0], auth[1]
		} else {
			fmt.Printf("Authentication user:password (%s) is not valid\n", *authUserPass)
			os.Exit(1)
		}
	}

	return cfg
}

func main() {
	api := newAPI(getConfig())
	api.run()
}
