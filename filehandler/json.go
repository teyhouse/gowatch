package filehandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"sync"
)

// Explicit struct fields over map[string]interface{}: malformed or hostile
// JSON fails at parse time instead of panicking on a bad type assertion.
type settingsFile struct {
	Files map[string]string `json:"files"`
}

type eventFile struct {
	HTTP struct {
		Request string `json:"request"`
	} `json:"http"`
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// walk appends regular files (including their full path) to files.
// Old / dangling symbolic links are skipped via FileExists.
func walk(files *[]string) fs.WalkDirFunc {
	return func(s string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // unreadable entry: skip, don't abort the whole walk
		}
		if !d.IsDir() && FileExists(s) {
			*files = append(*files, s)
		}
		return nil
	}
}

// Check if args contain --debug
func CheckDebug() bool {
	return slices.Contains(os.Args[1:], "--debug")
}

func GetSettings() ([]string, error) {
	data, err := os.ReadFile("settings.json")
	if err != nil {
		return nil, fmt.Errorf("reading settings.json: %w", err)
	}

	var settings settingsFile
	if err := json.Unmarshal(data, &settings); err != nil {
		return nil, fmt.Errorf("parsing settings.json: %w", err)
	}

	filelist := make([]string, 0)
	for _, value := range settings.Files {
		fileInfo, err := os.Stat(value)
		if err != nil {
			fmt.Printf("Skipping %s: %s\n", value, err)
			continue
		}

		if fileInfo.IsDir() {
			// iterate through all (recursive) files of the directory
			if err := filepath.WalkDir(value, walk(&filelist)); err != nil {
				return nil, fmt.Errorf("walking directory %s: %w", value, err)
			}
		} else if FileExists(value) {
			filelist = append(filelist, value)
		}
	}

	return filelist, nil
}

func GetEventURI() (string, error) {
	data, err := os.ReadFile("event.json")
	if err != nil {
		return "", fmt.Errorf("reading event.json: %w", err)
	}

	var event eventFile
	if err := json.Unmarshal(data, &event); err != nil {
		return "", fmt.Errorf("parsing event.json: %w", err)
	}
	if event.HTTP.Request == "" {
		return "", fmt.Errorf("event.json: missing \"http.request\"")
	}
	return event.HTTP.Request, nil
}

// GetHashes loads hashes.json and keeps only entries whose key is still
// present in filelist.
func GetHashes(filelist []string) (*sync.Map, error) {
	var res map[string]string

	data, err := os.ReadFile("hashes.json")
	switch {
	case err == nil:
		if err := json.Unmarshal(data, &res); err != nil {
			return nil, fmt.Errorf("parsing hashes.json: %w", err)
		}
	case errors.Is(err, fs.ErrNotExist):
		res = map[string]string{} // first run: nothing saved yet
	default:
		return nil, fmt.Errorf("reading hashes.json: %w", err)
	}

	savedhashes := &sync.Map{}
	for k, v := range res {
		if slices.Contains(filelist, k) {
			savedhashes.Store(k, v)
		}
	}
	return savedhashes, nil
}

func SaveHashes(savedhashes *sync.Map) error {
	m := make(map[string]string)
	savedhashes.Range(func(key, value any) bool {
		m[key.(string)] = value.(string)
		return true
	})

	jsonStr, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("serializing hashes: %w", err)
	}
	if err := os.WriteFile("hashes.json", jsonStr, 0o600); err != nil {
		return fmt.Errorf("writing hashes.json: %w", err)
	}
	return nil
}
