package profiles

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

type Profile struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Endpoint   string `json:"endpoint"`
	Region     string `json:"region"`
	Bucket     string `json:"bucket"`
	AccessKey  string `json:"accessKey"`
	SecretKey  string `json:"secretKey,omitempty"`
	CDNURL     string `json:"cdnUrl"`
	PathStyle  bool   `json:"pathStyle"`
	HasSecret  bool   `json:"hasSecret,omitempty"`
	PublicBase string `json:"publicBase,omitempty"`
}

type PublicProfile struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Endpoint   string `json:"endpoint"`
	Region     string `json:"region"`
	Bucket     string `json:"bucket"`
	AccessKey  string `json:"accessKey"`
	HasSecret  bool   `json:"hasSecret"`
	CDNURL     string `json:"cdnUrl"`
	PathStyle  bool   `json:"pathStyle"`
	PublicBase string `json:"publicBase"`
}

type Store struct {
	mu       sync.RWMutex
	filePath string
	items    map[string]Profile
}

func NewStore(filePath string) (*Store, error) {
	if filePath == "" {
		dir, err := os.UserConfigDir()
		if err != nil {
			return nil, err
		}
		filePath = filepath.Join(dir, "minio-manager", "profiles.json")
	}

	store := &Store{
		filePath: filePath,
		items:    make(map[string]Profile),
	}

	if err := store.load(); err != nil {
		return nil, err
	}

	return store, nil
}

func (s *Store) List() []PublicProfile {
	s.mu.RLock()
	defer s.mu.RUnlock()

	list := make([]PublicProfile, 0, len(s.items))
	for _, item := range s.items {
		list = append(list, item.Public())
	}

	sort.Slice(list, func(i, j int) bool {
		return strings.ToLower(list[i].Name) < strings.ToLower(list[j].Name)
	})

	return list
}

func (s *Store) Get(id string) (Profile, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	item, ok := s.items[id]
	return item, ok
}

func (s *Store) Save(profile Profile) (PublicProfile, error) {
	profile = profile.Normalized()

	s.mu.Lock()
	defer s.mu.Unlock()

	if profile.SecretKey == "" {
		if existing, ok := s.items[profile.ID]; ok {
			profile.SecretKey = existing.SecretKey
		}
	}
	if err := profile.Validate(); err != nil {
		return PublicProfile{}, err
	}

	s.items[profile.ID] = profile
	return profile.Public(), s.persistLocked()
}

func (s *Store) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.items, id)
	return s.persistLocked()
}

func (s *Store) load() error {
	data, err := os.ReadFile(s.filePath)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return nil
	}
	return json.Unmarshal(data, &s.items)
}

func (s *Store) persistLocked() error {
	if err := os.MkdirAll(filepath.Dir(s.filePath), 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(s.items, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.filePath, data, 0o600)
}
