package filehash

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"

	"github.com/teyhouse/gowatch/filehandler"
)

func GetFileHash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		//log.Fatal(err)
		if filehandler.CheckDebug() {
			fmt.Printf("Error on file %s: %s\n", path, err)
		}
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		if filehandler.CheckDebug() {
			fmt.Println(err)
		}
	}
	return hex.EncodeToString(h.Sum(nil)), err
}
