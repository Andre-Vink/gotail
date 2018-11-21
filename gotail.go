package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

const (
	TermBlue   = "\x1b[1;34m"
	TermNormal = "\x1b[0m"
)

/* All paths to tail. */
var tailFolders = make(map[string]TailFolder)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("gotail INFO: USAGE: gotail <dir1> [<dir2>] [...]")
		return
	}

	for _, path := range args {
		if absPath, err := filepath.Abs(path); err == nil {
			stat, _ := os.Stat(absPath)
			if stat == nil || !stat.IsDir() {
				fmt.Printf("gotail ERROR: [%v] does not exist or is not a directory.\n", absPath)
			} else {
				addTailFolder(absPath)
			}
		}
	}

	fmt.Printf("gotail INFO: Tailing folders %v\n", tailFolders)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Printf("gotail ERROR: Cannot create watcher. Aborting. Erro: [%v]", err)
		return
	}

	for _, path := range tailFolders {
		watcher.Add(path.Path())
	}

	watch(watcher)
}

func addTailFolder(folderPath string) {
	// TODO: only add when not already added (SET?)
	if _, exists := tailFolders[folderPath]; exists {
		fmt.Printf("gotail INFO: Folder [%v] will only be tailed once!\n", folderPath)
	} else {
		tailFolder := Create(folderPath)
		tailFolders[folderPath] = tailFolder
	}
}

func watch(watcher *fsnotify.Watcher) {
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
