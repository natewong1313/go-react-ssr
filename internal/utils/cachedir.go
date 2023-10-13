package utils

import (
	"os"
	"path/filepath"
)

// GetCacheDir returns the path to the cache directory
func GetCacheDir() (string, error) {
	return createCacheDirIfNotExists()
}

// GetTypeConverterCacheDir returns the path to the type converter cache directory
func GetTypeConverterCacheDir() (string, error) {
	cacheDir, err := createCacheDirIfNotExists()
	if err != nil {
		return "", err
	}
	typeConverterCacheDir := filepath.Join(cacheDir, "type_converter")
	if err := os.RemoveAll(typeConverterCacheDir); err != nil {
		return "", err
	}
	err = os.MkdirAll(typeConverterCacheDir, os.ModePerm)
	return typeConverterCacheDir, err
}

// createCacheDirIfNotExists creates the user cache directory if it doesn't exist
func createCacheDirIfNotExists() (string, error) {
	userCacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	gossrCacheDirPath := filepath.Join(userCacheDir, "gossr")
	err = os.MkdirAll(gossrCacheDirPath, os.ModePerm)
	return gossrCacheDirPath, err
}
