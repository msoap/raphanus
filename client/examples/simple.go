package main

import (
	"fmt"
	"strings"

	"github.com/msoap/raphanus/client"
)

func main() {
	raph := raphanusclient.New()

	saveIntKey(raph, "k1", 123)
	saveIntKey(raph, "k2", 777)
	incrDecrIntKey(raph, "k2")
	printIntKey(raph, "k1")
	updateIntKey(raph, "k1", 321)
	printIntKey(raph, "k1")

	printKeys(raph)
	printLength(raph)

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

func updateIntKey(raph raphanusclient.Client, key string, value int64) {
	err := raph.UpdateInt(key, value)
	if err != nil {
		fmt.Printf("UpdateInt got error: %s\n", err)
		return
	}

	fmt.Printf("Int value (%s: %d) updated\n", key, value)
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

func incrDecrIntKey(raph raphanusclient.Client, key string) {

	if err := raph.IncrInt(key); err != nil {
		fmt.Printf("IncrInt got error: %s\n", err)
		return
	}
	printIntKey(raph, key)

	if err := raph.DecrInt(key); err != nil {
		fmt.Printf("DecrInt got error: %s\n", err)
		return
	}
	printIntKey(raph, key)
}
