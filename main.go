package main

import (
	"os"

	"github.com/gudmundur/ds/ds"
)

func main() {
	token := os.Getenv("DROPBOX_TOKEN")
	ds := ds.NewDropboxSync(token)
	ds.Sync()
}
