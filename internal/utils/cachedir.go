package utils

import (
	"os"
	"path/filepath"
)

func CleanCacheDirectories() {
	cacheDir, err := createCacheDirIfNotExists()
	if err != nil {
		return
	}
	typeConverterCacheDir, _ := GetTypeConverterCacheDir()
	if err := os.RemoveAll(typeConverterCacheDir); err != nil {
		return
	}
	cssCacheDir, _ := GetCSSCacheDir()
	if err := os.RemoveAll(cssCacheDir); err != nil {
		return
	}
	// Remove GetServerBuildCacheDir
	if err := os.RemoveAll(filepath.Join(cacheDir, "builds")); err != nil {
		return
	}
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

// GetTypeConverterCacheDir returns the path to the type converter cache directory
func GetTypeConverterCacheDir() (string, error) {
	cacheDir, err := createCacheDirIfNotExists()
	if err != nil {
		return "", err
	}
	typeConverterCacheDir := filepath.Join(cacheDir, "typeconverter")
	err = os.MkdirAll(typeConverterCacheDir, os.ModePerm)
	return typeConverterCacheDir, err
}

// GetServerBuildCacheDir returns the path to the server build cache directory for the given route
func GetServerBuildCacheDir(fileName string) (string, error) {
	cacheDir, err := createCacheDirIfNotExists()
	if err != nil {
		return "", err
	}
	serverBuildCacheDir := filepath.Join(cacheDir, "builds")
	err = os.MkdirAll(serverBuildCacheDir, os.ModePerm)

	routeCacheDir := filepath.Join(serverBuildCacheDir, fileName)
	err = os.MkdirAll(routeCacheDir, os.ModePerm)
	return routeCacheDir, err
}

// GetCSSCacheDir returns the path to the server build cache directory for the given route
func GetCSSCacheDir() (string, error) {
	cacheDir, err := createCacheDirIfNotExists()
	if err != nil {
		return "", err
	}
	cssCacheDir := filepath.Join(cacheDir, "css_builds")
	err = os.MkdirAll(cssCacheDir, os.ModePerm)
	return cssCacheDir, err
}
