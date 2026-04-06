package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"time"
)

// Cache provides local file-based caching for HTTP responses
type Cache struct {
	Dir string
	TTL time.Duration
}

// DefaultTTL is the default cache time-to-live
const DefaultTTL = 24 * time.Hour

// New creates a new cache with the default directory (~/.cache/fedramp-browser)
func New() (*Cache, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	cacheDir := filepath.Join(homeDir, ".cache", "fedramp-browser")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return nil, err
	}

	return &Cache{
		Dir: cacheDir,
		TTL: DefaultTTL,
	}, nil
}

// keyToFilename converts a URL or key to a safe filename
func (c *Cache) keyToFilename(key string) string {
	hash := sha256.Sum256([]byte(key))
	return hex.EncodeToString(hash[:16]) + ".json"
}

// Path returns the full path to the cache file for a key
func (c *Cache) Path(key string) string {
	return filepath.Join(c.Dir, c.keyToFilename(key))
}

// Get retrieves data from cache if it exists and is not expired
func (c *Cache) Get(key string) ([]byte, bool) {
	path := c.Path(key)

	info, err := os.Stat(path)
	if err != nil {
		return nil, false
	}

	// Check if cache is expired
	if time.Since(info.ModTime()) > c.TTL {
		return nil, false
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, false
	}

	return data, true
}

// Set stores data in the cache
func (c *Cache) Set(key string, data []byte) error {
	path := c.Path(key)
	return os.WriteFile(path, data, 0644)
}

// Clear removes all cached files
func (c *Cache) Clear() error {
	entries, err := os.ReadDir(c.Dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			_ = os.Remove(filepath.Join(c.Dir, entry.Name()))
		}
	}
	return nil
}
