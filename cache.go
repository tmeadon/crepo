package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type repoCache struct {
	Paths []string `json:"paths"`
}

func getCacheFilePath() (string, error) {
	configPath, err := getConfigFilePath()
	if err != nil {
		return "", err
	}
	cacheDir := filepath.Dir(configPath)
	return filepath.Join(cacheDir, "repos.json"), nil
}

func loadRepoCache() ([]string, error) {
	cachePath, err := getCacheFilePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(cachePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var cache repoCache
	if err := json.Unmarshal(data, &cache); err != nil {
		return nil, err
	}
	return cache.Paths, nil
}

func saveRepoCache(paths []string) error {
	cachePath, err := getCacheFilePath()
	if err != nil {
		return err
	}

	data, err := json.Marshal(repoCache{Paths: paths})
	if err != nil {
		return err
	}

	return os.WriteFile(cachePath, data, 0644)
}
