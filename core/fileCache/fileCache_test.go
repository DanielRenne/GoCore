package fileCache

import (
	"io/ioutil"
	"os"
	"testing"
)

func init() {
	//Inject the pool
	initializeGroupCache("127.0.0.1")
}

func TestStringGroupCache(t *testing.T) {
	SetString("/somePath", "test")

	val, err := GetString("/somePath")
	if err != nil {
		t.Errorf("Error at fileCache_test.TestStringGroupCache\nFailed to GetString():  %v", err.Error)
	}

	if val != "test" {
		t.Error("Error at fileCache_test.TestStringGroupCache\nFailed to return proper matching data")
	}
}

func TestHTMLFileGroupCache(t *testing.T) {

	tmpfile, err := ioutil.TempFile("", "groupCacheHTML.htm")
	if err != nil {
		t.Errorf("Error at fileCache_test.TestHTMLFileGroupCache\nFailed to Create Temp File:  %v", err.Error)
		return
	}

	defer os.Remove(tmpfile.Name()) // clean up
	t.Log(tmpfile.Name())

	if _, err := tmpfile.Write([]byte("testHTML")); err != nil {
		t.Errorf("Error at fileCache_test.TestHTMLFileGroupCache\nFailed to Write to Temp HTML File:  %v", err.Error)
		return
	}
	if err := tmpfile.Close(); err != nil {
		t.Errorf("Error at fileCache_test.TestHTMLFileGroupCache\nFailed to Close Temp HTML File:  %v", err.Error)
		return
	}

	val, err := GetHTMLFile(tmpfile.Name())

	if err != nil {
		t.Errorf("Error at fileCache_test.TestHTMLFileGroupCache\nFailed to GetHTMLFile():  %v", err.Error)
	}

	if val != "testHTML" {
		t.Error("Error at fileCache_test.TestHTMLFileGroupCache\nFailed to return proper matching data")
	}
}
