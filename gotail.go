package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

const (
	termBlue   = "\x1b[1;34m"
	termNormal = "\x1b[0m"
)

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
	if _, exists := tailFolders[folderPath]; exists {
		fmt.Printf("gotail INFO: Folder [%v] will only be tailed once!\n", folderPath)
	} else {
		tailFolder := NewTailFolder(folderPath)
		tailFolders[folderPath] = tailFolder
	}
}

func watch(watcher *fsnotify.Watcher) {
	for {
		var e fsnotify.Event
		select {
		case e = <-watcher.Events:
			handleWatchEvent(e)
		case x := <-watcher.Errors:
			fmt.Printf("gotail INFO: Error: [%v]\n", x)
		}
	}
}

func handleWatchEvent(event fsnotify.Event) {
	fmt.Printf("gotail INFO: Event: [%v] (name: %v, op: %v)\n", event, event.Name, event.Op)
	switch event.Op {
	case fsnotify.Create:
		handleNewFile(event.Name)
	case fsnotify.Write:
		handleWriteToFile(event.Name)
	}
}

func handleNewFile(path string) {
	tailFolder := findTailFolderForFile(path)
	fileName := filepath.Base(path)
	from, to := tailFolder.AddFile(fileName)
	fmt.Printf("gotail INFO: Positions returned (%v, %v)\n", from, to)
}

func handleWriteToFile(path string) {
	tailFolder := findTailFolderForFile(path)
	fileName := filepath.Base(path)
	from, to := tailFolder.Positions(fileName)
	fmt.Printf("gotail INFO: Positions returned (%v, %v)\n", from, to)
}

func findTailFolderForFile(path string) TailFolder {
	dir := filepath.Dir(path)
	tailFolder := tailFolders[dir]
	fmt.Printf("gotail INFO: Found tailfolder [%v]\n", tailFolder)
	return tailFolder
}
