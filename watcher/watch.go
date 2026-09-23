package watcher

import (
	"fmt"
	"runtime"
	"sync"

	"github.com/teyhouse/gowatch/filehandler"
	"github.com/teyhouse/gowatch/filehash"
	"github.com/teyhouse/gowatch/logger"
)

// Generated Hashes for all files (current run)
var hashlist sync.Map

// Saved List of Hashes from hashes.json
var savedhashes *sync.Map

func getHashes(filelist []string) {
	var wg sync.WaitGroup

	// Counting semaphore sized to the CPU count, acquired BEFORE the
	// goroutine is spawned, so concurrency is truly bounded and no more
	// than NumCPU file handles are open at once.
	sem := make(chan struct{}, runtime.NumCPU())

	for _, key := range filelist {
		sem <- struct{}{}

		wg.Go(func() {
			defer func() { <-sem }()

			hash, err := filehash.GetFileHash(key)
			if err != nil {
				message := fmt.Sprintf("Error on file %s: %s", key, err)
				if filehandler.CheckDebug() {
					fmt.Println(message)
				}
				logger.Log(message)
				return
			}

			hashlist.Store(key, hash)

			// Check if hash is in savedhashes - add otherwise
			value, found := savedhashes.LoadOrStore(key, hash)
			if !found {
				if filehandler.CheckDebug() {
					fmt.Printf("Not found: %s:%v\n", key, value)
				}
				return
			}

			if filehandler.CheckDebug() {
				fmt.Printf("Found: %s:%s\n", key, value)
			}

			if hash == value {
				return // unchanged
			}

			if filehandler.CheckDebug() {
				fmt.Printf("Hashes don't match: %s:%s\n", key, value)
			}
			message := fmt.Sprintf("File-Change detected: %s:%s", key, value)
			logger.Log(message)
			logger.LogHTTP(message)
			savedhashes.Store(key, hash)
		})
	}
	wg.Wait()
}

func SaveHashes() error {
	return filehandler.SaveHashes(savedhashes)
}

func Watch() error {
	filelist, err := filehandler.GetSettings()
	if err != nil {
		return err
	}

	// filelist is passed in order to prevent old items being added again
	savedhashes, err = filehandler.GetHashes(filelist)
	if err != nil {
		return err
	}

	getHashes(filelist)

	// Iterate synced map - debug output
	if filehandler.CheckDebug() {
		fmt.Println("\n\nIterating over saved hashes:")
		savedhashes.Range(func(key, value any) bool {
			fmt.Printf("%s:%s\n", key, value)
			return true
		})
	}

	if err := SaveHashes(); err != nil {
		return err
	}
	fmt.Printf("DONE - checked %d files.\n", len(filelist))
	return nil
}
