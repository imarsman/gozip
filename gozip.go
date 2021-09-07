package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/jwalton/gchalk"
	"github.com/rickb777/go-arg"
)

const (
	brightGreen = iota
	brightYellow
	brightBlue
	brightRed
	noColour // Can use to default to no colour output
)

func colour(colour int, input ...string) (output string) {
	str := fmt.Sprint(strings.Join(input, " "))
	str = strings.Replace(str, "  ", " ", -1)

	// Choose colour for output or none
	switch colour {
	case brightGreen:
		output = gchalk.BrightGreen(str)
	case brightYellow:
		output = gchalk.BrightYellow(str)
	case brightBlue:
		output = gchalk.BrightBlue(str)
	case brightRed:
		output = gchalk.BrightRed(str)
	default:
		output = str
	}

	return
}

/*
	The zip utility has a lot of options. It is not know at the time of this
	writing how much of the original utility will be implemented.
*/

func createZip(path string) (zipFile *os.File, err error) {
	zipFile, err = os.Create(path)
	if err != nil {
		return
	}
	return
}

var args struct {
	Unzip       bool     `arg:"-U" help:"unzip the archive"`
	File        string   `arg:"-f" help:"verbosity level"`
	Recursive   bool     `arg:"-r" help:"recursive"`
	Update      bool     `arg:"-u" help:"update existing"`
	Add         bool     `arg:"-a" help:"add if not existing"`
	Zipfile     string   `arg:"positional"`
	SourceFiles []string `arg:"positional"`
}

// zip a list of paths to files
// https://golang.cafe/blog/golang-zip-file-example.html
func zipFiles(zipFilePath string, paths []string) (err error) {
	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		panic(err)
	}
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	for _, currentPath := range paths {
		var file *os.File
		file, err = os.Open(currentPath)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		var zw io.Writer
		zw, err = zipWriter.Create(currentPath)
		if err != nil {
			return
		}
		if _, err := io.Copy(zw, file); err != nil {
			panic(err)
		}
	}

	return
}

func walkAllFilesInDir(path string, files *[]string, errorMsgs *[]string) (err error) {
	var file *os.File
	file, err = os.Open(path)
	if err != nil {
		*errorMsgs = append(*errorMsgs, err.Error())
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		*errorMsgs = append(*errorMsgs, err.Error())
		return
	}
	defer file.Close()

	if !fileInfo.IsDir() {
		*files = append(*files, path)
		return
	}

	return filepath.Walk(path, func(path string, info os.FileInfo, e error) error {
		if err != nil {
			*errorMsgs = append(*errorMsgs, e.Error())
			return err
		}

		// check if it is a regular file (not dir)
		if info.Mode().IsRegular() {
			path := info.Name()
			*files = append(*files, path)
		}
		return nil
	})
}

func main() {
	arg.MustParse(&args)
	fmt.Println(args.File)

	if args.Zipfile == "" {
		fmt.Fprintln(os.Stderr, colour(brightRed, "no zipfile specified"))
		os.Exit(1)
	}

	if args.Unzip {
		if len(args.SourceFiles) > 0 {
			fmt.Fprintln(os.Stderr, colour(brightRed, "files specified - invalid with unzip"))
			os.Exit(1)
		}
	}

	var files = []string{}
	var errorMsgs = []string{}

	if !args.Unzip {
		// Populate list of files
		for _, file := range files {
			walkAllFilesInDir(file, &files, &errorMsgs)
		}
		if len(files) == 0 {
			fmt.Fprintln(os.Stderr, colour(brightRed, "no valid files found"))
			os.Exit(1)
		}
	}

	// Exit for now
	if 1 == 1 {
		os.Exit(0)
	}

	var archive *os.File
	if _, err := os.Stat(args.Zipfile); os.IsNotExist(err) {
		// Handle new file
		archive, err = os.Create(args.Zipfile)
		if err != nil {
			return
		}
	} else {
		archive, err = os.OpenFile(args.Zipfile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			return
		}
	}
	defer archive.Close()

	zipWriter := zip.NewWriter(archive)
	defer zipWriter.Close()

	for _, path := range files {
		file, err := os.Open("test.csv")
		if err != nil {
			fmt.Fprintln(os.Stderr, colour(brightRed, err.Error()))
			continue
		}
		defer file.Close()

		dest, err := zipWriter.Create(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, colour(brightRed, err.Error()))
			continue
		}
		if _, err := io.Copy(dest, file); err != nil {
			fmt.Fprintln(os.Stderr, colour(brightRed, err.Error()))
			continue
		}
	}
}
