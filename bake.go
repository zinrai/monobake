package main

import (
	"encoding/json"
	"fmt"
	"os"
)

const DefaultBakeFile = "docker-bake.json"

// BakeFile represents a docker-bake.json structure.
type BakeFile struct {
	Target map[string]json.RawMessage `json:"target"`
}

// LoadBakeFile loads and parses a Bake file from the given path.
func LoadBakeFile(path string) (*BakeFile, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open bake file: %w", err)
	}
	defer f.Close()

	var b BakeFile
	if err := json.NewDecoder(f).Decode(&b); err != nil {
		return nil, fmt.Errorf("failed to parse bake file: %w", err)
	}

	return &b, nil
}

// HasTarget checks if the target exists in the Bake file.
func (b *BakeFile) HasTarget(name string) bool {
	_, ok := b.Target[name]
	return ok
}
