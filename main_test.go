package main

import (
	"archive/zip"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func testRealMain(args []string) (int, string, string, error) {
	// Capture standard output and error
	oldOut := os.Stdout
	oldErr := os.Stderr

	stdoutReader, stdoutWriter, _ := os.Pipe()
	os.Stdout = stdoutWriter

	stderrReader, stderrWriter, _ := os.Pipe()
	os.Stderr = stderrWriter

	defer func() { os.Stdout = oldOut }()
	defer func() { os.Stderr = oldErr }()

	// Read output asynchronously so the pipes don't fill up
	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, stdoutReader)
		outC <- buf.String()
	}()

	errC := make(chan string)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, stderrReader)
		errC <- buf.String()
	}()

	// Run the main function
	status := realMain(args)

	// Close the pipes and return the results
	errOut := stdoutWriter.Close()
	errErr := stderrWriter.Close()
	if errOut != nil {
		return status, "", "", errOut
	}
	if errErr != nil {
		return status, "", "", errErr
	}

	outS := <-outC
	errS := <-errC
	return status, outS, errS, nil
}

func TestRealMainNoArgs(t *testing.T) {
	status, stdout, stderr, err := testRealMain([]string{"dirhasher"})
	if err != nil {
		t.Fatal(err)
	}
	if status != 1 {
		t.Fatalf("status %v", status)
	}

	expectedOut := ""
	expectedErr := "usage: dirhasher archive.zip|directory\n"
	if stdout != expectedOut {
		t.Errorf("stdout %q != %q", stdout, expectedOut)
	}
	if stderr != expectedErr {
		t.Errorf("stderr %q != %q", stderr, expectedErr)
	}
}

func TestRealMainNoSuchFile(t *testing.T) {
	status, stdout, stderr, err := testRealMain([]string{"dirhasher", "Jec#os}ai0ph"})
	if err != nil {
		t.Fatal(err)
	}
	if status != 1 {
		t.Fatalf("status %v", status)
	}

	expectedOut := ""
	expectedErr := "stat Jec#os}ai0ph: no such file or directory\n"
	if stdout != expectedOut {
		t.Errorf("stdout %q != %q", stdout, expectedOut)
	}
	if stderr != expectedErr {
		t.Errorf("stderr %q != %q", stderr, expectedErr)
	}
}

func TestRealMainDirectory(t *testing.T) {
	d, err := ioutil.TempDir("", "dirhasher-")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(d) }()

	if err := ioutil.WriteFile(filepath.Join(d, "foo"), []byte("foo contents"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := ioutil.WriteFile(filepath.Join(d, "bar"), []byte("bar contents"), 0644); err != nil {
		t.Fatal(err)
	}

	status, stdout, stderr, err := testRealMain([]string{"dirhasher", d})
	if err != nil {
		t.Fatal(err)
	}
	if status != 0 {
		t.Fatalf("status %v", status)
	}

	expectedOut := "h1:zNSU/Wy8yQDuLMTXCzfmK7DrPViDHSkkGBbFcIVxJ3A=\n"
	expectedErr := ""
	if stdout != expectedOut {
		t.Errorf("stdout %q != %q", stdout, expectedOut)
	}
	if stderr != expectedErr {
		t.Errorf("stderr %q != %q", stderr, expectedErr)
	}
}

func TestRealMainZip(t *testing.T) {
	f, err := ioutil.TempFile("", "dirhasher-")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Remove(f.Name()) }()

	z := zip.NewWriter(f)

	w, err := z.Create("foo")
	if err != nil {
		t.Fatal(err)
	}
	_, err = w.Write([]byte("foo contents"))
	if err != nil {
		t.Fatal(err)
	}

	w, err = z.Create("bar")
	if err != nil {
		t.Fatal(err)
	}
	_, err = w.Write([]byte("bar contents"))
	if err != nil {
		t.Fatal(err)
	}

	if err := z.Close(); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}

	status, stdout, stderr, err := testRealMain([]string{"dirhasher", f.Name()})
	if err != nil {
		t.Fatal(err)
	}
	if status != 0 {
		t.Fatalf("status %v", status)
	}

	expectedOut := "h1:zNSU/Wy8yQDuLMTXCzfmK7DrPViDHSkkGBbFcIVxJ3A=\n"
	expectedErr := ""
	if stdout != expectedOut {
		t.Errorf("stdout %q != %q", stdout, expectedOut)
	}
	if stderr != expectedErr {
		t.Errorf("stderr %q != %q", stderr, expectedErr)
	}
}
