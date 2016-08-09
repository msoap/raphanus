package main

import (
	"fmt"
	"strings"

	"github.com/msoap/raphanus/client"
)

func main() {
	raph := raphanusclient.New()
	printKeys(raph)
	printLength(raph)
	saveIntKey(raph, "k1", 123)
	printIntKey(raph, "k1")
	removeKey(raph, "k1")
}

func printKeys(raph raphanusclient.Client) {
	allKeys, err := raph.Keys()
	if err != nil {
		fmt.Printf("Keys got error: %s\n", err)
		return
	}

	fmt.Printf("all keys: %s\n", strings.Join(allKeys, ", "))
}

func printLength(raph raphanusclient.Client) {
	length, err := raph.Length()
	if err != nil {
		fmt.Printf("Length got error: %s\n", err)
		return
	}

	fmt.Printf("Count of keys: %d\n", length)
}

func saveIntKey(raph raphanusclient.Client, key string, value int64) {
	err := raph.SetInt(key, value)
	if err != nil {
		fmt.Printf("SetInt got error: %s\n", err)
		return
	}

	fmt.Printf("Int value (%s: %d) saved\n", key, value)
}

func removeKey(raph raphanusclient.Client, key string) {
	err := raph.Remove(key)
	if err != nil {
		fmt.Printf("Remove got error: %s\n", err)
		return
	}

	fmt.Printf("Key %s removed\n", key)
}

func printIntKey(raph raphanusclient.Client, key string) {
	intVal, err := raph.GetInt(key)
	if err != nil {
		fmt.Printf("GetInt got error: %s\n", err)
		return
	}

	fmt.Printf("Key %s, integer value: %d\n", key, intVal)
}
