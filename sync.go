package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/stacktic/dropbox"
)

type DropboxSync struct {
	Dropbox   *dropbox.Dropbox
	Directory string
}

func (ds *DropboxSync) readCursor() (cursor string, err error) {
	cursorBytes, err := ioutil.ReadFile(path.Join(ds.Directory, ".dropbox"))
	cursor = string(cursorBytes)
	return
}

func (ds *DropboxSync) writeCursor(cursor string) (err error) {
	return ioutil.WriteFile(path.Join(ds.Directory, ".dropbox"), []byte(cursor), 0644)
}

func (ds *DropboxSync) createFolder(entry dropbox.DeltaEntry) error {
	trimmedPath := entry.Entry.Path[1:]
	paths := strings.Split(trimmedPath, "/")
	localPath := path.Join(append([]string{ds.Directory}, paths...)...)
	return os.MkdirAll(localPath, 0755)
}

func (ds *DropboxSync) fetchFile(entry dropbox.DeltaEntry) error {
	src := entry.Entry.Path
	rev := entry.Entry.Revision
	dst := ds.Directory + entry.Entry.Path

	return ds.Dropbox.DownloadToFile(src, dst, rev)
}

func (ds *DropboxSync) remove(entry dropbox.DeltaEntry) error {
	dst := ds.Directory + entry.Path
	return os.Remove(dst)
}

func (ds *DropboxSync) Sync() error {
	cursor, err := ds.readCursor()
	delta, err := ds.Dropbox.Delta(cursor, "/")
	if err != nil {
		fmt.Println(err)
	}

	for _, entry := range delta.Entries {
		switch {
		case entry.Entry == nil:
			ds.remove(entry)
		case entry.Entry.IsDir == false:
			ds.fetchFile(entry)
		default:
			ds.createFolder(entry)
		}
	}

	ds.writeCursor(delta.Cursor)
	return nil
}

func NewDropboxSync(token string) *DropboxSync {
	db := dropbox.NewDropbox()
	db.RootDirectory = "auto"
	db.SetAccessToken(token)

	ds := &DropboxSync{
		Dropbox:   db,
		Directory: "tmp",
	}
	return ds
}

func main() {
	token := os.Getenv("DROPBOX_TOKEN")
	ds := NewDropboxSync(token)
	ds.Sync()
}
