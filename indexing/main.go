package main

import (
	"encoding/gob"
	"fmt"
	"indexing/trie"
	"os"
	"strings"
	"sync"
	"time"
)

// var root = "/home/fiammetta"
var root = "/mnt/stuff/A30s/DCIM/P"

type Flag int

const (
	Dir Flag = iota
	File
	None
)

func PrefixSearch(t *trie.RuneTrie, prefix string) []string {
	var paths []string
	t.Walk(func(key string, value interface{}) error {
		if strings.HasPrefix(strings.ToLower(key), strings.ToLower(prefix)) {
			paths = append(paths, key)
		}
		return nil
	})
	return paths

}

func ContainSearch(t *trie.RuneTrie, search string) []string {
	var paths []string
	t.Walk(func(key string, value interface{}) error {
		if strings.Contains(strings.ToLower(key), strings.ToLower(search)) {
			paths = append(paths, key)
		}
		return nil
	})
	return paths
}

func main() {

	gob.Register(Flag(0))

	fmt.Println("Importing File Tree...")
	tri, erro := DeserializeTrie()

	if erro != nil {
		fmt.Println(erro)
	}

	search := ".css"

	searchTime := time.Now()
	paths := ContainSearch(tri, search)

	for _, path := range paths {
		fmt.Println(path)
	}

	elapsedSearch := time.Since(searchTime)
	fmt.Printf("Time elapsed for search: %s\n", elapsedSearch)
	//
	os.Exit(0)

	// f := Fs{Root: FsEntry{Name: root}}
	files := make(chan FileData)

	t := trie.NewRuneTrie()
	t.Put(root, Dir)

	wgc := sync.WaitGroup{}
	wgc.Add(1)

	start_time := time.Now()
	go func() {
		defer wgc.Done()
		for file := range files {
			// f.AddEntry(file.Path)

			if file.IsDir {
				t.Put(file.Path, Dir)
			} else {
				t.Put(file.Path, File)
			}
		}
	}()

	wg := sync.WaitGroup{}
	wg.Add(1)

	sem := make(chan struct{}, 10000)

	go func() {
		defer wg.Done()
		walkdir(&wg, sem, root, files)
	}()

	wg.Wait()
	close(files)
	wgc.Wait()
	elapsed := time.Since(start_time)
	fmt.Printf("Time elapsed: %s\n", elapsed)

	err := SerializeTrie(t)

	if err != nil {
		fmt.Println(err)
	}

	// test trie for search
	// prefix := "/home/fiammetta/Downloads"
	// paths := PrefixSearch(t, prefix)

	// err = f.Serialize()
	//
	// if err != nil {
	// 	fmt.Println(err)
	// }
}
