package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"privateterraformregistry/internal/modules"
	"testing"

	"github.com/gorilla/mux"
)

func TestServiceDiscovery(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/.well-known/terraform.json", nil)
	w := httptest.NewRecorder()

	getServiceDiscovery(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, _ := io.ReadAll(res.Body)
	if string(data) != `{"modules.v1":"/terraform/modules/v1/"}
` {
		t.Errorf(`expected {"modules.v1":"/terraform/modules/v1/"} got %v`, string(data))
	}
	header := res.Header.Get(`Content-Type`)
	if string(header) != "application/json" {
		t.Errorf(`expected Content-Type: application/json got Content-Type: %v`, header)
	}
}

func TestListAvailableVersions(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/v1/modules/hashicorp/consul/aws/versions", nil)
	w := httptest.NewRecorder()

	vars := map[string]string{
		"namespace": "hashicorp",
		"system":    "consul",
		"name":      "aws",
	}
	r = mux.SetURLVars(r, vars)

	listAvailableVersions(modules.Modules{
		Modules: []modules.Module{
			{Namespace: "hashicorp", System: "consul", Name: "aws", Version: "1.0.0"},
			{Namespace: "hashicorp", System: "consul", Name: "aws", Version: "1.1.0"},
			{Namespace: "hashicorp", System: "consul", Name: "aws", Version: "2.0.0"},
		},
	})(w, r)
	res := w.Result()
	defer res.Body.Close()
	data, _ := io.ReadAll(res.Body)
	if string(data) != `{"modules":[{"versions":[{"version":"1.0.0"},{"version":"1.1.0"},{"version":"2.0.0"}]}]}
` {
		t.Errorf(`expected {"modules":[{"versions":[{"version":"1.0.0"},{"version":"1.1.0"},{"version":"2.0.0"}]}]} got %v`, string(data))
	}
	header := res.Header.Get(`Content-Type`)
	if string(header) != "application/json" {
		t.Errorf(`expected Content-Type: application/json got Content-Type: %v`, header)
	}
}

func TestGetDownloadPath(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/v1/modules/hashicorp/consul/aws/0.0.1/download", nil)
	w := httptest.NewRecorder()

	vars := map[string]string{
		"namespace": "hashicorp",
		"system":    "consul",
		"name":      "aws",
		"version":   "0.0.1",
	}
	r = mux.SetURLVars(r, vars)

	getDownloadPath(w, r)
	res := w.Result()
	defer res.Body.Close()
	header := res.Header.Get("X-Terraform-Get")
	if string(header) != "/modules/hashicorp/consul/aws/0.0.1" {
		t.Errorf(`expected /modules/hashicorp/consul/aws/0.0.1 got %v`, header)
	}
	if res.StatusCode != 204 {
		t.Errorf(`expected status code 204 got %v`, res.StatusCode)
	}
}
