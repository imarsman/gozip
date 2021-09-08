package main

import (
	"bytes"
	"os/exec"
	"testing"

	"github.com/matryer/is"
)

//                Tests and benchmarks
// -----------------------------------------------------
// benchmark
//   go test -run=XXX -bench=. -benchmem
// Get allocation information and pipe to less
//   go build -gcflags '-m -m' ./*.go 2>&1 |less
// Run all tests
//   go test -v
// Run one test and do allocation profiling
//   go test -run=XXX -bench=IterativeISOTimestampLong -gcflags '-m' 2>&1 |less
// Run a specific test by function name pattern
//  go test -run=TestParsISOTimestamp
//
//  go test -run=XXX -bench=.
//  go test -bench=. -benchmem -memprofile memprofile.out -cpuprofile cpuprofile.out
//  go tool pprof -http=:8080 memprofile.out
//  go tool pprof -http=:8080 cpuprofile.out

func runCmd(command string) error {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()

	return err
}

func cleanup() error {
	err := runCmd("script/reset.sh")

	return err
}

func TestStart(t *testing.T) {
	is := is.New(t)
	err := cleanup()
	is.NoErr(err)
}

func TestRunCmd(t *testing.T) {
	is := is.New(t)

	err := runCmd("script/reset.sh")
	is.NoErr(err)
}

func TestGetFiles(t *testing.T) {
	is := is.New(t)

	var fileEntries = []fileEntry{}
	var errorMsgs = []string{}

	walkAllFilesInDir("./sample", &fileEntries, &errorMsgs)

	for _, fe := range fileEntries {
		t.Logf("%s", fe.fullPath())
	}

	err := runCmd("script/reset.sh")
	is.NoErr(err)
}

func TestMakeZip(t *testing.T) {
	is := is.New(t)

	path := "./sample"
	var fileEntries = []fileEntry{}
	var errorMsgs = []string{}
	// walk through dir or single file to get new entries
	walkAllFilesInDir(path, &fileEntries, &errorMsgs)

	path = "/Users/ian/git/gozip/sample/orig/1.txt"
	// walk through dir or single file to get new entries
	walkAllFilesInDir(path, &fileEntries, &errorMsgs)

	err := archiveFiles("test/archive.zip", fileEntries)
	is.NoErr(err)
}

func TestListZip(t *testing.T) {
	is := is.New(t)

	path := "test/archive.zip"

	entries, err := fileList(path)
	is.NoErr(err)

	for _, file := range entries {
		t.Logf("%+v", file)
	}
}

func TestEnd(t *testing.T) {
	is := is.New(t)
	err := cleanup()
	is.NoErr(err)
}
