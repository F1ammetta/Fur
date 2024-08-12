package main

import (
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/ulikunitz/xz"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var unixIgnore = []string{"/proc", "/sys", "/dev", "/run", "/mnt"}

type Fs struct {
	Root FsEntry
}

type FsEntry struct {
	Name     string
	Children []FsEntry
}

type FileData struct {
	Path  string
	IsDir bool
}

func (f *FsEntry) Walk() {
	fmt.Println(f.Name)
	for _, child := range f.Children {
		child.Walk()
	}
}

func (f *Fs) Walk() {
	f.Root.Walk()
}

func (f *Fs) Serialize() error {
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

	err = encoder.Encode(f)

	if err != nil {
		return err
	}

	err = writer.Close()

	if err != nil {
		return err
	}

	return nil
}

func (f *Fs) Deserialize() error {
	filename := strings.ReplaceAll(root, string(os.PathSeparator), "_")
	filename = strings.ReplaceAll(filename, ":", "_")
	file, err := os.Open(filename + ".fs")

	if err != nil {
		return err
	}

	defer file.Close()

	reader, err := xz.NewReader(file)

	if err != nil {
		return err
	}

	decoder := gob.NewDecoder(reader)

	err = decoder.Decode(f)

	if err != nil {
		return err
	}

	return nil
}

func (f *Fs) AddEntry(path string) error {
	path = strings.TrimPrefix(path, root)
	parts := strings.Split(path, string(os.PathSeparator))
	if len(parts) == 0 {
		return errors.New("invalid path")
	}

	current := &f.Root

	for _, part := range parts {
		found := false

		for i, child := range current.Children {
			if child.Name == part {
				current = &current.Children[i]
				found = true
				break
			}
		}

		if !found {
			current.Children = append(current.Children, FsEntry{Name: part})
			current = &current.Children[len(current.Children)-1]
		}
	}

	return nil
}

func walkdir(wg *sync.WaitGroup, sem chan struct{}, root string, files chan<- FileData) {
	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() && root != path {
			for _, ignore := range unixIgnore {
				if strings.HasPrefix(path, ignore) {
					return filepath.SkipDir
				}
			}
			sem <- struct{}{}
			wg.Add(1)
			go func() {
				defer wg.Done()
				walkdir(wg, sem, path, files)

				<-sem
			}()

			return filepath.SkipDir
		}

		info, err := d.Info()

		if err != nil {
			fmt.Println(err)
		}

		if info.Mode().IsRegular() {
			files <- FileData{Path: path, IsDir: d.IsDir()}
		}

		return nil

	})
}
