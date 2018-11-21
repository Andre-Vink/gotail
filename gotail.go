package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"os"
	"path/filepath"
)

const (
	TERM_BLUE   = "\x1b[1;34m"
	TERM_NORMAL = "\x1b[0m"
)

/* All paths to tail. */
var foldersToTail []string

/* All files to trace. */
var tracedFiles = make(map[string]uint64)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("gotail INFO: USAGE: gotail <dir1> [<dir2>] ...")
		return
	}

	for _, path := range args {
		if absPath, err := filepath.Abs(path); err == nil {
			//fmt.Printf("%v, %v (%T)\n", absPath, err, absPath)
			stat, _ := os.Stat(absPath)
			if stat == nil || !stat.IsDir() {
				fmt.Printf("gotail ERROR: [%v] does not exist or is not a directory.\n", absPath)
			} else {
				addFolderToTrace(absPath)
			}
		}
	}

	fmt.Printf("gotail INFO: Tailing folders %v\n", foldersToTail)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Printf("gotail ERROR: Cannot create watcher. Aborting. Erro: [%v]", err)
		return
	}

	for _, path := range foldersToTail {
		watcher.Add(path)
	}

	for {
		var e fsnotify.Event
		select {
		case e = <-watcher.Events:
			fmt.Printf("gotail INFO: Event: [%v] (name: %v, op: %v)\n", e, e.Name, e.Op)
		case x := <-watcher.Errors:
			fmt.Printf("gotail INFO: Error: [%v]\n", x)
		}
	}
}

func addFolderToTrace(folderPath string) {
	// TODO: only add when not already added (SET?)
	foldersToTail = append(foldersToTail, folderPath)

	folder, err := os.Open(folderPath)
	if err != nil {
		panic("Cannot open folder")
	}

	files, err := folder.Readdir(0)
	if err != nil {
		panic("Cannot read folder")
	}

	for _, file := range files {
		if !file.IsDir() {
			tracedFiles[file.Name()] = uint64(file.Size())
			fmt.Println("Adding file ", TERM_BLUE, file.Name(), TERM_NORMAL, "for tracing. It is currently", TERM_BLUE, file.Size(), TERM_NORMAL, "long.")
		}
	}
}
