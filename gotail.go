package main

import (
	"fmt"
	"os"
	"path/filepath"
	"github.com/fsnotify/fsnotify"
)

var pathsToTail []string

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
	    		fmt.Printf("gotail ERROR: [%v] does not exist or is not a directory.\n", absPath);
			} else {
				// TODO: only add when not already added (SET?)
				pathsToTail = append(pathsToTail, absPath)
			}
		}
	}

	fmt.Printf("gotail INFO: Tailing folders %v\n", pathsToTail)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Printf("gotail ERROR: Cannot create watcher. Aborting. Erro: [%v]", err)
		return
	}

	for _, path := range pathsToTail {
		watcher.Add(path)
	}

	for {
		select {
		case e := <- watcher.Events:
			fmt.Printf("gotail INFO: Event: [%v]\n", e)
		case x := <- watcher.Errors:
			fmt.Printf("gotail INFO: Error: [%v]\n", x)
		}
	}
}
