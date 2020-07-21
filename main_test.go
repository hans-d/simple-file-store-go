package filestore

import (
	"os"
	"path/filepath"
	"testing"
)

type Gopher struct {
	Type string `json:"type"`
}

const testFolder = "./test"

func TestMain(m *testing.M) {
	os.RemoveAll(testFolder)
	code := m.Run()
	os.RemoveAll(testFolder)
	os.Exit(code)
}

func TestNew(t *testing.T) {

	dataFolder := filepath.Join(testFolder, "new")

	if _, err := os.Stat(dataFolder); err == nil {
		t.Error("Expected nothing, got data folder")
	}

	if _, err := New(dataFolder, nil); err != nil {
		t.Errorf("Expected nothing, got %v", err)
	}

	if _, err := os.Stat(dataFolder); err != nil {
		t.Error("Expected data folder, got nothing")
	}

	if _, err := New(dataFolder, nil); err != nil {
		t.Errorf("Expected nothing, got %v", err)
	}
	if _, err := os.Stat(dataFolder); err != nil {
		t.Error("Expected data folder, got nothing")
	}
}

func TestWriteAndRead(t *testing.T) {

	type testCase struct {
		descr   string
		folder  string
		options *Options
	}
	testCases := []testCase{
		{descr: "simple", folder: "simple", options: nil},
		{descr: "yaml", folder: "yaml", options: &Options{Marshaler: &YAMLMarshaler{}}},
	}

	parent := "simpleread"
	key := "some-gopher"

	for _, tt := range testCases {
		dataFolder := filepath.Join(testFolder, tt.folder)
		f, err := New(dataFolder, tt.options)

		if err != nil {
			t.Errorf("Expected nothing for %s, got %v", tt.descr, err)
		}

		if err := f.Write(parent, key, Gopher{Type: "Simple"}); err != nil {
			t.Errorf("Write failed for %s: %s", tt.descr, err.Error())
		}

		gopher := Gopher{}
		if err := f.Read(parent, key, &gopher); err != nil {
			t.Errorf("Failed to read for %s: %s", tt.descr, err.Error())
		}
		if gopher.Type != "Simple" {
			t.Errorf("Expected Simple gopher for %s, got: %+v", tt.descr, gopher)
		}
	}
}
