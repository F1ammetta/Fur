package main

import (
	"encoding/gob"
	"github.com/ulikunitz/xz"
	"indexing/trie"
	"os"
	"strings"
)

func SerializeTrie(t *trie.RuneTrie) error {
	filename := strings.ReplaceAll(root, string(os.PathSeparator), "_")
	filename = strings.ReplaceAll(filename, ":", "_")
	file, err := os.Create(filename + ".fs")

	if err != nil {
		return err
	}

	defer file.Close()

	writer, err := xz.NewWriter(file)

	if err != nil {
		return err
	}

	encoder := gob.NewEncoder(writer)

	err = encoder.Encode(t)

	if err != nil {
		return err
	}

	err = writer.Close()

	if err != nil {
		return err
	}

	return nil
}

func DeserializeTrie() (*trie.RuneTrie, error) {
	filename := strings.ReplaceAll(root, string(os.PathSeparator), "_")
	filename = strings.ReplaceAll(filename, ":", "_")
	file, err := os.Open(filename + ".fs")

	if err != nil {
		return nil, err
	}

	defer file.Close()

	reader, err := xz.NewReader(file)

	if err != nil {
		return nil, err
	}

	decoder := gob.NewDecoder(reader)

	t := &trie.RuneTrie{}

	err = decoder.Decode(t)

	if err != nil {
		return nil, err
	}

	return t, nil
}
