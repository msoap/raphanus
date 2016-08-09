package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/msoap/raphanus/client"
)

func main() {
	raph := raphanusclient.New()
	printKeys(raph)
	removeKey(raph, "k1")
}

func printKeys(raph raphanusclient.Client) {
	allKeys, err := raph.Keys()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("all keys: %s\n", strings.Join(allKeys, ", "))
}

func removeKey(raph raphanusclient.Client, key string) {
	err := raph.Remove(key)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Key %s removed\n", key)
}
