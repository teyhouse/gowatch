package watcher

import (
	"fmt"
	"sync"

	"github.com/teyhouse/gowatch/filehandler"
	"github.com/teyhouse/gowatch/filehash"
	"github.com/teyhouse/gowatch/logger"
)

// List of Files you want to watch from settings.json
var filelist = make([]string, 0)

// Generated Hashes for all files
var hashlist sync.Map

// Saved List of Hashes from hashes.json
var savedhashes sync.Map

func getSettings() {
	filelist = filehandler.GetSettings()
}

func getStoreHashes() {
	//filelist is passed in order to prevent old items being added again
	savedhashes = filehandler.GetHashes(filelist)
}

func getHashes() {
	var wg sync.WaitGroup

	for _, key := range filelist {
		wg.Add(1)

		go func(key string) {
			//Get Filehash from current iteration
			hash := filehash.GetFileHash(key)

			//Add filename:hash to hashlist
			hashlist.Store(key, hash)

			//Check if hash is in savedhashes - add otherwise
			value, found := savedhashes.LoadOrStore(key, hash)
			if found {
				if filehandler.CheckDebug() {
					fmt.Printf("Found: %s:%s\n", key, value)
				}

				//Check if hash still matches with stored hash
				if hash == value {
					//fmt.Println("Hashes match.")
				} else {
					if filehandler.CheckDebug() {
						fmt.Printf("Hashes don't match: %s:%s\n", key, value)
					}
					//log.Printf("File-Change detected: %s:%s", key, value)
					message := fmt.Sprintf("File-Change detected: %s:%s", key, value)
					logger.Log(message)
					logger.LogHTTP(message)
					savedhashes.Store(key, hash)
				}
			} else {
				if filehandler.CheckDebug() {
					fmt.Printf("Not found: %s:%s\n", key, value)
				}
				//savedhashes.Store(key, hash) //Not necessary since LoadorStore does both
			}
			wg.Done()
		}(key)
	}
	wg.Wait()
}

func SaveHashes() {
	filehandler.SaveHashes(savedhashes)
}

func Watch() {

	getSettings()
	getStoreHashes()
	getHashes()

	//Iterate synced map - debug output
	if filehandler.CheckDebug() {
		fmt.Println("\n\nIterating over saved hashes:")
		savedhashes.Range(func(key, value interface{}) bool {
			fmt.Printf("%s:%s\n", key, value)
			return true
		})
	}

	SaveHashes()
	fmt.Printf("DONE - checked %s files.\n", fmt.Sprint(len(filelist)))
}
