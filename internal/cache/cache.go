package cache

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

// UploadCache gerencia o estado local de uploads bem-sucedidos
type UploadCache struct {
	Entries map[string]string `json:"entries"` // mapping: FileHash -> URL
	path    string
	mu      sync.RWMutex
}

// NewCache cria ou carrega o cache local do diretório do mangá.
func NewCache(mangaDir string) *UploadCache {
	c := &UploadCache{
		Entries: make(map[string]string),
		path:    filepath.Join(mangaDir, ".manga_cache.json"),
	}
	c.Load()
	return c
}

// Load lê o cache do disco, se existir.
func (c *UploadCache) Load() {
	c.mu.Lock()
	defer c.mu.Unlock()

	data, err := os.ReadFile(c.path)
	if err == nil {
		json.Unmarshal(data, &c.Entries)
	}
}

// Save escreve o cache no disco.
func (c *UploadCache) Save() error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	data, err := json.MarshalIndent(c.Entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(c.path, data, 0600)
}

// Set adiciona uma nova entrada ao cache.
func (c *UploadCache) Set(hash string, url string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Entries[hash] = url
}

// Get recupera a URL do cache se ela existir.
func (c *UploadCache) Get(hash string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	url, exists := c.Entries[hash]
	return url, exists
}

// HashFile gera um SHA256 do arquivo para identificar se a imagem já foi enviada.
func HashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
