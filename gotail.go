package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"os"
	"path/filepath"
)

const (
	termBlue   = "\x1b[1;34m"
	termNormal = "\x1b[0m"
)

var tailFolders = NewTailFolders()

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
				tailFolders.AddFolder(absPath)
			}
		}
	}

	fmt.Printf("gotail INFO: Tailing folders %v\n", tailFolders.Folders())

	eventsChannel := tailFolders.Watch()

	watch(eventsChannel)
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

//func addTailFolder(folderPath string) {
//	if _, exists := tailFolders[folderPath]; exists {
//		fmt.Printf("gotail INFO: Folder [%v] will only be tailed once!\n", folderPath)
//	} else {
//		tailFolder := NewTailFolders(folderPath)
//		tailFolders[folderPath] = tailFolder
//	}
//}

func handleWatchEvent(event fsnotify.Event) {
	fmt.Printf("gotail INFO: Event: [%v] (name: %v, op: %v)\n", event, event.Name, event.Op)
	switch event.Op {
	case fsnotify.Create:
		handleNewFile(event.Name)
	case fsnotify.Write:
		handleWriteToFile(event.Name)
	}
}

func handleNewFile(newFile string) {
	fmt.Println("Handle new file: ", newFile)
	//tailFolder := findTailFolderForFile(path)
	//fileName := filepath.Base(path)
	//from, to := tailFolder.AddFile(fileName)
	//fmt.Printf("gotail INFO: Positions returned (%v, %v)\n", from, to)
}

func handleWriteToFile(writtenFile string) {
	fmt.Println("Handle write to file: ", writtenFile)
	//	tailFolder := findTailFolderForFile(path)
	//	fileName := filepath.Base(path)
	//	from, to := tailFolder.Positions(fileName)
	//	fmt.Printf("gotail INFO: Positions returned (%v, %v)\n", from, to)
}

//func findTailFolderForFile(path string) TailFolders {
//	dir := filepath.Dir(path)
//	tailFolder := tailFolders[dir]
//	fmt.Printf("gotail INFO: Found tailfolder [%v]\n", tailFolder)
//	return tailFolder
//}
