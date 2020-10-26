// Implements the gotail application that can tail folders.
// It will tail all files in the specified folders, even new files that are newly created.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"unicode/utf8"

	"github.com/fsnotify/fsnotify"
)

const (
	termBlue   = "\x1b[1;34m"
	termNormal = "\x1b[0m"
)

var tailFolders = NewTailFolders()

var hlFlag = flag.String("hl", "", "Wil highlight any part of the logging that matches this regex.")
var hlRE *regexp.Regexp

func main() {
	flag.Parse()

	args := flag.Args()
	fmt.Printf("Flag highlight = [%v]\n", *hlFlag)

	if len(args) == 0 {
		fmt.Println("gotail INFO: USAGE: gotail [-hl 'regex'] <dir1> [<dir2>] [...]")
		fmt.Println("       -hl : high light regular expression - matching parts in the tailed logging is highlighted. Optional.")
		return
	}

	var err error
	hlRE, err = regexp.Compile(*hlFlag)
	if err != nil {
		fmt.Println("gotail ERROR: high light regular expression is not a valid regular expression.")
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

	eventsChannel, err := tailFolders.Watch()
	if err != nil {
		fmt.Printf("gotail ERROR: Cannot watch folders. Error: [%v]", err)
	}

	watch(eventsChannel)
}

func watch(watcher *fsnotify.Watcher) {
	for {
		// var e fsnotify.Event
		select {
		case e := <-watcher.Events:
			handleWatchEvent(e)
		case x := <-watcher.Errors:
			fmt.Printf("gotail ERROR: %v\n", x)
		}
	}
}

func handleWatchEvent(event fsnotify.Event) {
	// fmt.Printf("gotail INFO: Event: [%v] (name: %v, op: %v)\n", event, event.Name, event.Op)
	switch event.Op {
	case fsnotify.Create:
		handleNewFile(event.Name)
	case fsnotify.Write:
		handleWriteToFile(event.Name)
	}
}

func handleNewFile(newFile string) {
	//fmt.Println("Handle new file: ", newFile)
	tailFolders.AddFile(newFile)
}

func handleWriteToFile(writtenFile string) {
	//fmt.Println("Handle write to file: ", writtenFile)
	newLines := tailFolders.NewLines(writtenFile)
	//fmt.Printf("NewPart returned [%v]\n", newPart)

	tailFolder := findTailFolderForFile(writtenFile)
	// fmt.Printf("Tail folder for [%v] is [%v]\n", writtenFile, tailFolder)

	for _, line := range newLines {
		if utf8.ValidString(line) {
			// add scan for hl regex and print matching parts in blue
			if len(*hlFlag) > 0 {
				matched := hlRE.FindAllString(line, -1)
				i := -1
				line = hlRE.ReplaceAllStringFunc(line, func(src string) string {
					i = i + 1
					return termBlue + matched[i] + termNormal
				})
			}
			fmt.Printf("%v%v%v: %v\n", termBlue, tailFolder, termNormal, line)
		}
	}
}

func findTailFolderForFile(path string) string {
	return filepath.Base(filepath.Dir(path))
	//	dir := filepath.Dir(path)
	//	tailFolder := tailFolders[dir]
	//	fmt.Printf("gotail INFO: Found tailfolder [%v]\n", tailFolder)
	//	return tailFolder
}
