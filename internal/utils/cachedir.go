package utils

import (
	"os"
	"path/filepath"
)

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

// GetServerBuildCacheDir returns the path to the server build cache directory for the given route
func GetServerBuildCacheDir(routeName string) (string, error) {
	cacheDir, err := createCacheDirIfNotExists()
	if err != nil {
		return "", err
	}
	serverBuildCacheDir := filepath.Join(cacheDir, "builds")
	err = os.MkdirAll(serverBuildCacheDir, os.ModePerm)

	routeCacheDir := filepath.Join(serverBuildCacheDir, routeName)
	err = os.MkdirAll(routeCacheDir, os.ModePerm)
	return routeCacheDir, err
}
