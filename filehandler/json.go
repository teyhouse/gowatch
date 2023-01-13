package filehandler

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Recursive file slice
var stmp []string

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// Iterate through directroys (recursive) - this will return files only (including the full path)
func walk(s string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}
	if !d.IsDir() {
		stmp = append(stmp, s)
	}
	return nil
}

// Since slices.Contain is not included in the current upstream-version of Go 1.19
func contains(elems []string, v string) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}

// Check if args contain --debug
func CheckDebug() bool {
	if arg := strings.Join(os.Args[1:], ""); strings.Contains(string(arg), "--debug") {
		return true
	}
	return false
}

func GetSettings() []string {
	data, err := os.ReadFile("settings.json")
	if err != nil {
		fmt.Print(err)
	}

	var result map[string]any
	json.Unmarshal([]byte(data), &result)
	files := result["files"].(map[string]any)

	var filelist = make([]string, 0)

	for _, value := range files {
		//Check if File or Directory
		fileInfo, err := os.Stat(value.(string))
		if err != nil {
			fmt.Print(err)
		}

		if fileInfo.IsDir() {
			// iterate through all files in directory
			filepath.WalkDir(value.(string), walk)
		} else {
			filelist = append(filelist, value.(string))
		}
	}

	//Merge recursive files
	filelist = append(filelist, stmp...)

	return filelist
}

func GetEventURI() string {
	data, err := os.ReadFile("event.json")
	if err != nil {
		fmt.Print(err)
	}

	var result map[string]any
	json.Unmarshal([]byte(data), &result)
	eventuri := result["http"].(map[string]any)
	return fmt.Sprint(eventuri["request"])
}

func GetHashes(filelist []string) sync.Map {
	if !FileExists("hashes.json") {
		file, err := os.Create("hashes.json")
		if err != nil {
			fmt.Printf("Error creating hashes.json: %s \n", err)
		}
		defer file.Close()
	}

	data, err := os.ReadFile("hashes.json")
	if err != nil {
		fmt.Print(err)
	}

	var res map[string]interface{}
	json.Unmarshal([]byte(data), &res)

	var savedhashes sync.Map

	for k, v := range res {
		//Check if Key (= filename) is still wanted, based on settings.json
		if contains(filelist, k) {
			savedhashes.Store(k, v)
		}
	}
	return savedhashes
}

func SaveHashes(savedhashes sync.Map) {
	m := make(map[string]string)

	savedhashes.Range(func(key, value interface{}) bool {
		m[key.(string)] = value.(string)
		return true
	})

	jsonStr, err := json.Marshal(m)

	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}

	_ = os.WriteFile("hashes.json", jsonStr, 0644)
}
