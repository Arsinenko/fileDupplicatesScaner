package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func GetMd5Hash(filepath string) string {
	file, _ := os.Open(filepath)
	defer file.Close()
	hash := md5.New()
	io.Copy(hash, file)
	return hex.EncodeToString(hash.Sum(nil))
}

func main() {
	bySize := make(map[int64][]string)
	var root string
	if len(os.Args) > 1 {
		root = os.Args[1]
	} else {
		root, _ = os.Getwd()
	}
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && strings.HasPrefix(d.Name(), ".") {
			return fs.SkipDir
		}
		if !d.IsDir() {
			info, _ := d.Info()
			bySize[info.Size()] = append(bySize[info.Size()], path)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Ошибка обхода: %v\n", err)
	}
	for size, files := range bySize {
		if len(files) < 2 || size == 0 {
			continue
		}
		duplicates := make(map[string][]string)
		for _, path := range files {
			hash := GetMd5Hash(path)
			duplicates[hash] = append(duplicates[hash], path)
		}

		for hash, paths := range duplicates {
			if len(paths) > 1 {
				fmt.Printf("Дубликаты (%d байт), хеш %s:\n - %s\n",
					size, hash, strings.Join(paths, "\n - "))
			}
		}
	}

}
