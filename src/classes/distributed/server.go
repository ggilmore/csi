package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/peterbourgon/ff/v3/ffcli"
)

type server struct {
	mu    sync.RWMutex
	cache map[string]string

	storagePath string
	address     string
}

func NewServer(storagePath, address string) (*server, error) {
	data, err := os.ReadFile(storagePath)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("loading data from %q: %s", storagePath, err)
		}

		return &server{
			cache:       make(map[string]string),
			storagePath: storagePath,
			address:     address,
		}, nil
	}

	var cache map[string]string
	err = json.Unmarshal(data, &cache)
	if err != nil {
		return nil, fmt.Errorf("parsing data from %q: %s", storagePath, err)
	}

	return &server{
		cache:       cache,
		storagePath: storagePath,
		address:     address,
	}, nil
}

func serverCommand() (*ffcli.Command, error) {
	var (
		port        int
		hostname    string
		storagePath string
	)

	fs := flag.NewFlagSet("server", flag.ExitOnError)

	fs.IntVar(&port, "port", defaultPort, "the port that the server will listening on")
	fs.StringVar(&hostname, "hostname", defaultHostname, "the hostname of the server")

	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to discover user's home directory: %w", err)
	}
	fs.StringVar(&storagePath, "storagePath", filepath.Join(home, "kv_storage.json"), "the file path used to store the key/value data")

	cmd := &ffcli.Command{
		Name:       "server",
		ShortUsage: fmt.Sprintf("%s server [-port N] [-hostname foo.com]", os.Args[0]),
		ShortHelp:  "start an instance of the key/value server",
		FlagSet:    fs,
		Exec: func(ctx context.Context, args []string) error {
			address := fmt.Sprintf("%s:%d", hostname, port)

			s, err := NewServer(storagePath, address)
			if err != nil {
				return fmt.Errorf("initializing server: %s", err)
			}

			err = s.Run()
			if err != nil {
				return fmt.Errorf("running server: %w", err)
			}

			return nil
		},
	}

	return cmd, nil
}

func (s *server) Run() error {
	http.HandleFunc("/v1/store", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			s.handleGet(w, r)
		case http.MethodPost:
			s.handleSet(w, r)
		default:
			message := fmt.Sprintf("method must be one of %s or %s, got %s", http.MethodGet, http.MethodPost, r.Method)
			http.Error(w, message, http.StatusBadRequest)
			return
		}
	})

	log.Printf("listening on %q", s.address)

	err := http.ListenAndServe(s.address, nil)
	if err != nil {
		return fmt.Errorf("listening on %q: %s", s.address, err)
	}

	return nil
}

func (s *server) handleGet(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	key := params.Get("key")

	if key == "" {
		http.Error(w, "key query parameter must be specified", http.StatusBadRequest)
		return
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	value, found := s.cache[key]
	if !found {
		http.Error(w, "key %q not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	_, err := w.Write([]byte(value))
	if err != nil {
		log.Printf("failed to write error response for set %q=%q, err: %s", key, value, err)
		return
	}
}

func (s *server) handleSet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, fmt.Sprintf("method must be POST, got %s", r.Method), http.StatusBadRequest)
		return
	}

	params := r.URL.Query()
	key := params.Get("key")

	if key == "" {
		http.Error(w, "key query parameter must be specified", http.StatusBadRequest)
		return
	}

	value := params.Get("value")

	s.mu.Lock()
	defer s.mu.Unlock()

	// save old cache value in case something goes wrong
	oldValue, found := s.cache[key]
	s.cache[key] = value

	err := s.flushStorage()
	if err != nil {
		// restore old cache value
		if found {
			s.cache[key] = oldValue
		} else {
			delete(s.cache, key)
		}

		http.Error(w, fmt.Sprintf("flushing cache to storage: %s", err), http.StatusInternalServerError)
		return
	}
}

func (s *server) flushStorage() error {
	bs, err := json.Marshal(s.cache)
	if err != nil {
		return fmt.Errorf("marshalling cache: %s", err)
	}

	err = os.WriteFile(s.storagePath, bs, 0666)
	if err != nil {
		return fmt.Errorf("writing to %q: %s", s.storagePath, err)
	}

	return nil
}
