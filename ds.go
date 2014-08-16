package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/stacktic/dropbox"
)

func readCursor(filename string) (cursor string, err error) {
	cursorBytes, err := ioutil.ReadFile(filename)
	cursor = string(cursorBytes)
	return
}

func writeCursor(filename string, cursor string) (err error) {
	return ioutil.WriteFile(filename, []byte(cursor), 0644)
}

func createFolder(entry dropbox.DeltaEntry) (err error) {
	trimmedPath := strings.TrimLeft(entry.Entry.Path, "/")
	paths := strings.Split(trimmedPath, "/")
	localPath := path.Join(append([]string{"tmp"}, paths...)...)
	err = os.MkdirAll(localPath, 0755)
	return
}

func main() {
	token := os.Getenv("DROPBOX_TOKEN")
	db := dropbox.NewDropbox()
	db.SetAccessToken(token)

	cursor, err := readCursor("tmp/.dropbox")
	delta, err := db.Delta(cursor, "/")
	if err != nil {
		fmt.Println(err)
	}

	for _, entry := range delta.Entries {
		switch {
		case entry.Entry == nil:
			fmt.Println("delete")
		case entry.Entry.IsDir == false:
			fmt.Println("fetch file")
		default:
			createFolder(entry)
		}
	}

	//writeCursor("tmp/.dropbox", delta.Cursor)
}
