package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
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

	// fmt.Println(filename, dir)

	lang, err := readJSON(filename)
	if err != nil {
		log.Fatal(err)
	}

	queue := make(chan string)
	go keys(lang, "", queue)

	unused := make([]string, 0)
	for key := range queue {
		found, err := find(key, dir)
		if (err) != nil {
			log.Fatal(err)
		}

		if !found {
			unused = append(unused, key)
		}
	}

	fmt.Println("Unused keys:")
	fmt.Println(strings.Join(unused, "\n"))
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

// read a file and check if it contains a key
func checkFile(filename string, key string) (bool, error) {
	b, err := os.ReadFile(filename)
	if err != nil {
		return false, err
	}

	s := string(b)
	return strings.Contains(s, key), nil
}

// walk a directory until we find a file containing key
func find(key string, dir string) (bool, error) {
	// fmt.Printf("Finding key: %s\n", key)
	found := false
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		// go to next
		if d.IsDir() {
			return nil
		}

		if strings.HasPrefix(d.Name(), ".") {
			return nil
		}

		if err != nil {
			return err
		}

		found, err = checkFile(path, key)
		if err != nil {
			return err
		}

		// if we found the key in a file, we're done
		if found {
			// log.Printf("Found %s in file %s", key, path)
			return fs.SkipAll
		}

		return nil
	})

	return found, err
}
