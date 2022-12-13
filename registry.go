package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"privateterraformregistry/internal/datamanager"
	"privateterraformregistry/internal/modules"
	"privateterraformregistry/internal/uploader"

	"github.com/gorilla/mux"
)

var dataDir = os.Getenv("DATA_DIR")

const (
	dataFile      = "data/data.json"
	maxUploadSize = 32 << 20
)

func getServiceDiscovery(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	type serviceDiscovery struct {
		ModulePath string `json:"modules.v1"`
	}

	json.NewEncoder(w).Encode(serviceDiscovery{
		ModulePath: "/terraform/modules/v1/",
	})
}

func listAvailableVersions(ms modules.Modules) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		namespace := vars["namespace"]
		system := vars["system"]
		name := vars["name"]

		type Version struct {
			Version string `json:"version"`
		}
		var availableVersions []Version

		for _, m := range ms.Modules {
			if m.Namespace == namespace && m.System == system && m.Name == name {
				availableVersions = append(availableVersions, Version{Version: m.Version})
			}
		}

		p := struct {
			Modules []struct {
				Versions []Version `json:"versions"`
			} `json:"modules"`
		}{
			Modules: []struct {
				Versions []Version `json:"versions"`
			}{
				{
					Versions: availableVersions,
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(p)
	}
}

func getDownloadPath(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.Header().Set("X-Terraform-Get", fmt.Sprintf("/modules/%s/%s/%s/%s", vars["namespace"], vars["system"], vars["name"], vars["version"]))
	w.WriteHeader(204)
}

func uploadOrDownloadModule(m *modules.Modules, d datamanager.DataManager) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			log.Println("Uploading Module")
			uploadModule(m, d, w, r)
			return
		}

		if r.Method == http.MethodGet {
			log.Println("Downloading Module")
			downloadModule(w, r)
			return
		}

		w.WriteHeader(405)
		w.Write([]byte("Method Not Allowed"))
	}
}

func uploadModule(ms *modules.Modules, datamanager datamanager.DataManager, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	m := modules.Module{
		Namespace: vars["namespace"],
		System:    vars["system"],
		Name:      vars["name"],
		Version:   vars["version"],
	}
	var err = m.Validate()

	if err != nil {
		log.Println(err)
		return
	}

	uploader := uploader.New(&uploader.Config{
		MaxUploadSize: maxUploadSize,
		DataDir:       dataDir,
	})
	err = uploader.Upload(r, m)

	if err != nil {
		log.Println(err)
		return
	}

	ms.Add(m)
	err = datamanager.Save()

	if err != nil {
		log.Println(err)
	}
}

func downloadModule(w http.ResponseWriter, r *http.Request) {
	// downloader := NewDownloader()
	// err = downloader.download(w, r))

	// if (err != nil) {
	// 	// do something
	// }
}

func main() {
	if dataDir == "" {
		dataDir = "/.privateterraformregistry/data"
	}

	ms := modules.Modules{}
	var datamanager = datamanager.New(&datamanager.Config{
		DataDir: dataDir,
	}, &ms)

	log.Print("Loading modules into memory from file")
	var err = datamanager.Load()

	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter().StrictSlash(true)

	log.Print("Registering Routes")
	router.HandleFunc("/.well-known/terraform.json", getServiceDiscovery)                    // terraform protocol
	router.HandleFunc("/v1/{namespace}/{name}/{system}/versions", listAvailableVersions(ms)) // terraform module protocol
	router.HandleFunc("/v1/{namespace}/{name}/{system}/{version}/download", getDownloadPath) // terraform module protocol
	router.HandleFunc("/modules/{namespace}/{name}/{system}/{version}", uploadOrDownloadModule(&ms, datamanager))

	log.Print("Server Ready")
	log.Fatal(http.ListenAndServe(":8080", router))
}
