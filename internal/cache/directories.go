package cache

import (
	"os"
	"path/filepath"
)

var CacheDir string
var TypeConverterCacheDir string
var TailwindCacheDir string

func SetupCacheDirectories() error {
	cacheDir, err := createCacheDirIfNotExists()
	if err != nil {
		return err
	}
	CacheDir = cacheDir
	TypeConverterCacheDir = filepath.Join(cacheDir, "typeconverter")
	TailwindCacheDir = filepath.Join(cacheDir, "tailwind")
	// Clean up type converter builds
	if err := os.RemoveAll(TypeConverterCacheDir); err != nil {
		return err
	}
	os.MkdirAll(TypeConverterCacheDir, os.ModePerm)
	// Clean up tailwind builds
	if err := os.RemoveAll(TailwindCacheDir); err != nil {
		return err
	}
	os.MkdirAll(TailwindCacheDir, os.ModePerm)

	// Remove react builds
	return os.RemoveAll(filepath.Join(cacheDir, "builds"))
}

// createCacheDirIfNotExists creates the cache directory if it doesn't exist
func createCacheDirIfNotExists() (string, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	baseCacheDir := filepath.Join(workingDir, ".gossr")
	err = os.MkdirAll(baseCacheDir, os.ModePerm)
	return baseCacheDir, err
}
