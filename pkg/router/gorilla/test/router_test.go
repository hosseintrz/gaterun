package test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/hosseintrz/gaterun/config/models"
	"github.com/hosseintrz/gaterun/pkg/proxy"
	"github.com/hosseintrz/gaterun/pkg/router/gorilla"
)

func TestDefaultRouter(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	defaultPort := 9000

	cfg := models.ServiceConfig{
		Name: "test-models",
		Port: defaultPort,
		Endpoints: []*models.EndpointConfig{
			{
				Endpoint: "/users/{id}",
				Method:   "GET",
				Backends: []*models.BackendConfig{
					{},
				},
				Timeout: 5,
			},
			{
				Endpoint: "/users",
				Method:   "POST",
				Backends: []*models.BackendConfig{
					{},
				},
				Timeout: 5,
			},
			{
				Endpoint: "/users/{id}",
				Method:   "PUT",
				Backends: []*models.BackendConfig{
					{},
				},
				Timeout: 5,
			},
			{
				Endpoint: "/users/{id}",
				Method:   "PATCH",
				Backends: []*models.BackendConfig{
					{},
				},
				Timeout: 5,
			},
			{
				Endpoint: "/users/{id}",
				Method:   "DELETE",
				Backends: []*models.BackendConfig{
					{},
				},
				Timeout: 5,
			},
		},
	}

	routerFactory := gorilla.DefaultFactory(proxy.NoOpProxyFactory())
	r := routerFactory.NewWithContext(ctx)

	go func() {
		r.Run(cfg)
	}()

	time.Sleep(5 * time.Millisecond)

	for _, endpoint := range cfg.Endpoints {
		url := fmt.Sprintf("http://localhost:%d%s", defaultPort, endpoint.Endpoint)
		req, err := http.NewRequest(strings.ToTitle(endpoint.Method), url, http.NoBody)
		if err != nil {
			t.Errorf("error creating http request: %v\n", err)
			return
		}

		req.Header.Set("Content-Type", "application/json")

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Errorf("error in executing http request: %v\n", err)
			return
		}

		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			t.Errorf("error reading response body: %v\n", body)
		}

		expectedBody := `{"ok":true}`
		data := string(body)
		if data != expectedBody {
			t.Errorf("expected response body {%s} but got {%s}\n", expectedBody, data)
		}

		if cType := res.Header.Get("Content-Type"); cType != "application/json" {
			t.Errorf("epxected application/json contentType but got: %s - addr: %s \n", cType, url)
		}
		if gateRunVersion := res.Header.Get("GATERUN"); gateRunVersion != "1.0" {
			t.Errorf("expected gaterun version to be 1 but got %s - addr : %s\n", gateRunVersion, url)
		}

		if res.StatusCode != http.StatusOK {
			t.Errorf("expected statusCode to be 200 but got %d - addr: %s\n", res.StatusCode, url)
		}
	}

}
