package datamanager

import (
	"os"
	"privateterraformregistry/internal/modules"
	"testing"
)

var dataDir = os.Getenv("DATA_DIR")

const (
	testFileContents = `{"modules":[{"Namespace":"hashicorp","System":"consul","Name":"aws","Version":"0.1.0"}]}`
)

func TestLoadPopulatesModules(t *testing.T) {
	_, err := os.OpenFile(dataDir+"/modules.json", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(dataDir+"/modules.json", []byte(testFileContents), 0644)
	if err != nil {
		t.Fatal(err)
	}

	ms := modules.Modules{}
	dm := New(&Config{
		DataDir: dataDir,
	}, &ms)
	dm.Load()

	if len(ms.Modules) < 1 {
		t.Fatal("No modules loaded")
	}

	if ms.Modules[0].Namespace != "hashicorp" {
		t.Errorf("Expected Namespace to be hashicorp got %v", ms.Modules[0].Namespace)
	}

	if ms.Modules[0].System != "consul" {
		t.Errorf("Expected System to be consul got %v", ms.Modules[0].Namespace)
	}

	if ms.Modules[0].Name != "aws" {
		t.Errorf("Expected Name to be aws got %v", ms.Modules[0].Namespace)
	}

	if ms.Modules[0].Version != "0.1.0" {
		t.Errorf("Expected Version to be 0.1.0 got %v", ms.Modules[0].Version)
	}
}

func TestSaveUpdatesModules(t *testing.T) {
	os.Remove(dataDir + "/modules.json")

	ms := modules.Modules{}

	ms.Add(modules.Module{
		Namespace: "hashicorp",
		System:    "consul",
		Name:      "aws",
		Version:   "0.1.0",
	})

	dm := New(&Config{
		DataDir: dataDir,
	}, &ms)
	err := dm.Save()
	if err != nil {
		t.Error(err)
	}

	_, err = os.Stat(dataDir)
	if err != nil {
		t.Errorf("Unable to retrieve file info for %s", dataDir)
	}
}
