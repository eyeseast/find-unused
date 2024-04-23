package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

var helptext = "Usage: find-unused lang.json ./src"

func main() {
	args := os.Args[1:]
	if len(args) < 2 {
		fmt.Println(helptext)
		os.Exit(0)
	}

	filename := args[0]
	dir := args[1]

	fmt.Println(filename, dir)

	lang, err := readJSON(filename)
	if err != nil {
		log.Fatal(err)
	}

	queue := make(chan string)
	go keys(lang, "", queue)

	for key := range queue {
		fmt.Println(key)
	}
}

func readJSON(filename string) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	file, err := os.Open(filename)
	if err != nil {
		return m, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	decoder.Decode(&m)
	return m, nil
}

func keys(object map[string]interface{}, prior string, channel chan string) {

	for k, v := range object {
		key := k
		if prior != "" {
			key = fmt.Sprintf("%s.%s", prior, k)
		}
		switch v.(type) {
		case string:
			channel <- key

		case map[string]interface{}:
			keys(v.(map[string]interface{}), key, channel)
		}
	}

	// we've reached the end of the original loop
	if prior == "" {
		close(channel)
	}
}
