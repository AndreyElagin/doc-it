package fsutils

import (
	"log"
	"os"
)

func CreateDirIfNotExist(dir string) error {
	_, err := os.ReadDir(dir)
	if err != nil {
		err := os.MkdirAll(dir, 0777)
		if err != nil {
			log.Fatal(err)
			return err
		}
	}
	return nil
}
