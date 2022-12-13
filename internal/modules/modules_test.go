package modules

import (
	"testing"
)

func TestAddAppendsModule(t *testing.T) {
	ms := Modules{}

	ms.Add(Module{
		Namespace: "hashicorp",
		System:    "consul",
		Name:      "aws",
		Version:   "0.1.0",
	})

	if len(ms.Modules) != 1 {
		t.Error("Module not appended to Modules")
	}
}

func TestAddCannotAddDuplicateModule(t *testing.T) {
	ms := Modules{}

	ms.Add(Module{
		Namespace: "hashicorp",
		System:    "consul",
		Name:      "aws",
		Version:   "0.1.0",
	})
	ms.Add(Module{
		Namespace: "hashicorp",
		System:    "consul",
		Name:      "aws",
		Version:   "0.1.0",
	})

	if len(ms.Modules) == 1 {
		t.Error("Module not appended to Modules")
	}
}
