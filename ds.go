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

func createFolder(entry dropbox.DeltaEntry) error {
	trimmedPath := entry.Entry.Path[1:]
	paths := strings.Split(trimmedPath, "/")
	localPath := path.Join(append([]string{"tmp"}, paths...)...)
	return os.MkdirAll(localPath, 0755)
}

func fetchFile(db *dropbox.Dropbox, entry dropbox.DeltaEntry) error {
	src := entry.Entry.Path
	rev := entry.Entry.Revision
	dst := "tmp/" + entry.Entry.Path

	return db.DownloadToFile(src, dst, rev)
}

func remove(entry dropbox.DeltaEntry) error {
	dst := "tmp" + entry.Path
	return os.Remove(dst)
}

func main() {
	token := os.Getenv("DROPBOX_TOKEN")
	db := dropbox.NewDropbox()
	db.RootDirectory = "auto"
	db.SetAccessToken(token)

	cursor, err := readCursor("tmp/.dropbox")
	delta, err := db.Delta(cursor, "/")
	if err != nil {
		fmt.Println(err)
	}

	for _, entry := range delta.Entries {
		switch {
		case entry.Entry == nil:
			remove(entry)
		case entry.Entry.IsDir == false:
			fetchFile(db, entry)
		default:
			createFolder(entry)
		}
	}

	writeCursor("tmp/.dropbox", delta.Cursor)
}
