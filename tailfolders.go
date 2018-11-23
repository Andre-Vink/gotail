package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"os"
	"path/filepath"
)

// TailFolders holds the information about the folders and their files, needed for tailing.
type TailFolders struct {
	folders []string
	files   map[string]uint64
}

// NewTailFolders makes a new TailFolders instance.
// The path should point to a directory, otherwise a panic will occur.
func NewTailFolders() TailFolders {
	return TailFolders{files: make(map[string]uint64)}
}

// Folders returns all folders that are tailed or going to be tailed.
func (tailFolders TailFolders) Folders() []string {
	return tailFolders.folders
}

// AddFolder adds a folder for tailing.
func (tailFolders *TailFolders) AddFolder(folder string) {
	tailFolders.folders = append(tailFolders.folders, folder)
	tailFolders.readFiles(folder)
}

// AddFile adds a file for tailing.
func (tailFolders *TailFolders) AddFile(file string) {
	tailFolders.files[file] = 0
}

// Watch starts watching the folders.
func (tailFolders TailFolders) Watch() *fsnotify.Watcher {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Printf("gotail ERROR: Cannot create watcher. Aborting. Error: [%v]", err)
		//return nil
	}

	for _, folder := range tailFolders.folders {
		watcher.Add(folder)
	}

	return watcher
}

// NewPart returns the new part of the file.
func (tailFolders TailFolders) NewPart(fileName string) string {
	lastPosition := tailFolders.files[fileName]
	if fileInfo, err := os.Stat(fileName); err == nil {
		currentFileSize := uint64(fileInfo.Size())
		if currentFileSize < lastPosition { // fileName was truncated
			lastPosition = 0
		}
		if file, err := os.Open(fileName); err == nil {
			defer file.Close()
			file.Seek(int64(lastPosition), 0)
			bufferSize := currentFileSize - lastPosition
			var buffer = make([]byte, bufferSize)
			n, err := file.Read(buffer)
			fmt.Println("Read file [", fileName, "] returned", n, "bytes and error", err)
			return string(buffer)
		}
	}
	return ""
}

// Path returns the absolute path to folder.
//func (tailFolder TailFolders) Path() string {
//	return tailFolder.folders
//}

// String returns the absolute path as string representation of this folder.
//func (tailFolder TailFolders) String() string {
//	return tailFolder.folders
//}

func (tailFolders TailFolders) readFiles(folder string) {
	openFolder, err := os.Open(folder)
	if err != nil {
		panic(fmt.Sprintf("Cannot open folder [%v].", folder))
	}

	files, err := openFolder.Readdir(0)
	if err != nil {
		panic(fmt.Sprintf("Cannot read folder [%v].", folder))
	}

	for _, file := range files {
		if !file.IsDir() {
			absFilePath := filepath.Join(folder, file.Name())
			tailFolders.files[absFilePath] = uint64(file.Size())
			fmt.Println("Adding file ", termBlue, absFilePath, termNormal, "for tracing. It is currently",
				termBlue, tailFolders.files[absFilePath], termNormal, "long.")
		}
	}
}

// Positions returns the last echo'd position as from and the current size of the file in the to.
// If the file is not found it will panic.
//func (tailFolders TailFolders) Positions(fileName string) (from uint64, to uint64) {
//	value, ok := tailFolders.files[fileName]
//	if ok {
//		fmt.Printf("gotail TRACE: found file [%v] with position [%v]\n", fileName, value)
//		return value, 0
//	}
//	panic(fmt.Sprintf("Did not find file [%v]", fileName))
//}
