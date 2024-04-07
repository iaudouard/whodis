package store

import (
	"fmt"
	"log"
	"os"
	"time"
)

// saves data to disk atomically
func saveToDisk(data []byte, path string) error {
	tmp := fmt.Sprintf("tmp-%d.wdb", time.Now().Unix())

	file, err := os.OpenFile(tmp, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0644)
	if err != nil {
		return err
	}
	defer func() {
		file.Close()
		os.Remove(tmp)
	}()

	n, err := file.Write(data)
	log.Printf("wrote %d bytes", n)
	if err != nil {
		return err
	}
	err = file.Sync()
	if err != nil {
		return err
	}

	return os.Rename(tmp, path)
}
