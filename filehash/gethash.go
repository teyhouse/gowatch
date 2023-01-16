package filehash

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"

	"github.com/teyhouse/gowatch/filehandler"
	"github.com/teyhouse/gowatch/logger"
)

func GetFileHash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		message := fmt.Sprintf("Error on file %s: %s\n", path, err)
		if filehandler.CheckDebug() {
			fmt.Print(message)
		}
		logger.Log(message)
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
