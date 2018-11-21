package main

import (
	"fmt"
	"os"
)

// TailFolder holds the information about the folder and its files, needed for tailing.
type TailFolder struct {
	AbsolutePath string
	Files        map[string]uint64
}

// NewTailFolder makes a new TailFolder instance for the provided path.
// The path should point to a directory, otherwise a panic will occur.
func NewTailFolder(path string) TailFolder {
	dir := TailFolder{AbsolutePath: path, Files: make(map[string]uint64)}
	go dir.readFiles()
	return dir
}

// Path returns the absolute path to folder.
func (tailFolder TailFolder) Path() string {
	return tailFolder.AbsolutePath
}

// String returns the absolute path as string representation of this folder.
func (tailFolder TailFolder) String() string {
	return tailFolder.AbsolutePath
}

func (tailFolder *TailFolder) readFiles() {
	folder, err := os.Open(tailFolder.AbsolutePath)
	if err != nil {
		panic("Cannot open folder")
	}

	files, err := folder.Readdir(0)
	if err != nil {
		panic("Cannot read folder")
	}

	for _, file := range files {
		if !file.IsDir() {
			tailFolder.Files[file.Name()] = uint64(file.Size())
			fmt.Println("Adding file ", termBlue, file.Name(), termNormal, "for tracing. It is currently",
				termBlue, tailFolder.Files[file.Name()], termNormal, "long.")
		}
	}
}
