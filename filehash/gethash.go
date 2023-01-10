package filehash

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"os"
)

func GetFileHash(path string) string {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}
	return hex.EncodeToString(h.Sum(nil))
}
