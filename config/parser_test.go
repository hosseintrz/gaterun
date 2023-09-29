package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestReadCfg(t *testing.T) {
	folderName := "tmp"
	fileName := "cfg.json"
	configPath := filepath.Join(folderName, fileName)

	jsonConfig := []byte(`
	{
		"version": 3,
		"name": "gaterun api-gateway",
		"port": 8000,
		"cache_ttl": "3000s",
		"timeout": "4s",
		"host": "localhost",
		"endpoints": [
			{
				"endpoint": "/github",
				"method": "GET",
				"backends": [
					{
						"host": "https://api.github.com",
						"url_pattern": "/",
						"allow": [
							"authorizations_url",
							"code_search_url"
						]
					}
				]
			}	
		]
	}
	`)

	if err := os.Mkdir(folderName, os.ModePerm); err != nil {
		t.Fatalf("error creating new folder: %v\n", err)
	}

	if err := os.WriteFile(configPath, jsonConfig, 0644); err != nil {
		t.Fatalf("error writing file: %v\n", err)
	}

	cfg, err := Parse(configPath)
	if err != nil {
		t.Fatalf("error parsing config: %v", err)
	}

	if endpointCnt := len(cfg.Endpoints); endpointCnt != 1 {
		t.Fatalf("expected endpoints count to be %d but it's %d", 1, endpointCnt)
	}

	if cfg.Timeout != time.Duration(4*time.Second) {
		t.Fatalf("expected timeout to be %v but it's %v\n", 4*time.Second, cfg.Timeout)
	}

	if cfg.Port != 8000 {
		t.Fatalf("expected port %d but got %d\n", 8000, cfg.Port)
	}

	backends := cfg.Endpoints[0].Backends
	if backendsCnt := len(backends); backendsCnt != 1 {
		t.Fatalf("expected %d backends but got %d\n", 1, backendsCnt)
	}

	expectedHost := "https://api.github.com"
	if actualHost := backends[0].Host; actualHost != expectedHost {
		t.Fatalf("expected backend host to be %s but got %s\n", expectedHost, actualHost)
	}

	if err := os.RemoveAll(folderName); err != nil {
		t.Fatalf("error removing tmp json file")
	}

}
