package main

import (
	"fmt"
	"os"
)

type TailFolder struct {
	AbsolutePath string
	Files        map[string]uint64
}

func Create(path string) TailFolder {
	dir := TailFolder{AbsolutePath: path, Files: make(map[string]uint64)}
	go dir.readFiles()
	return dir
}

func (folder TailFolder) Path() string {
	return folder.AbsolutePath
}

func (folder TailFolder) String() string {
	return folder.AbsolutePath
}

func (this *TailFolder) readFiles() {
	folder, err := os.Open(this.AbsolutePath)
	if err != nil {
		panic("Cannot open folder")
	}

	files, err := folder.Readdir(0)
	if err != nil {
		panic("Cannot read folder")
	}

	for _, file := range files {
		if !file.IsDir() {
			this.Files[file.Name()] = uint64(file.Size())
			fmt.Println("Adding file ", TermBlue, file.Name(), TermNormal, "for tracing. It is currently",
				TermBlue, this.Files[file.Name()], TermNormal, "long.")
		}
	}
}
