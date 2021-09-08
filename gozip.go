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

type fileEntry struct {
	rootPath   string
	parentPath string
	path       string
	isBareFile bool
}

func hasEntry(check fileEntry, feList *[]fileEntry) (found bool) {
	for _, fe := range *feList {
		if check.fullPath() == fe.fullPath() {
			return true
		}
	}
	return
}

func (fe *fileEntry) fullPath() (fullPath string) {
	fullPath = filepath.Join(fe.parentPath, fe.path)
	// fmt.Printf("%+v\n", fe)
	return
}

func (fe *fileEntry) archivePath() (archivePath string) {
	if fe.isBareFile {
		archivePath = filepath.Join(fe.parentPath, fe.path)
		if strings.HasPrefix(archivePath, "/") {
			archivePath = archivePath[1:]
		}
		return
	}
	archivePath = strings.Replace(fe.parentPath, fe.rootPath, "", 1)
	archivePath = filepath.Join(archivePath, fe.path)
	if strings.HasPrefix(archivePath, "/") && len(archivePath) > 1 {
		archivePath = archivePath[1:]
	}
	return
}

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

func walkAllFilesInDir(path string, fileEntries *[]fileEntry, errorMsgs *[]string) (err error) {
	// fmt.Println("path", path, "relative path", rel)
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

	rootPath, err := filepath.Abs(path)
	rootPath = filepath.Dir(rootPath)

	var curDir string
	var basePath string

	if !fileInfo.IsDir() {
		pth := path
		if strings.HasPrefix(path, "/") {
			pth = pth[1:]
		}
		fe := fileEntry{}
		fe.rootPath = rootPath

		abs, _ := filepath.Abs(path)
		parentPath := filepath.Dir(abs)

		fe.rootPath = filepath.Dir(rootPath)
		fe.parentPath = parentPath

		fe.path = fileInfo.Name()
		fe.isBareFile = true

		if !hasEntry(fe, fileEntries) {
			*fileEntries = append(*fileEntries, fe)
		}
		return
	}

	return filepath.Walk(path, func(path string, info os.FileInfo, e error) (err error) {
		if err != nil {
			*errorMsgs = append(*errorMsgs, e.Error())
			return err
		}

		if info.Name() == filepath.Base(path) {
			basePath, err = filepath.Abs(path)
			basePath = filepath.Dir(basePath)
		}
		if info.IsDir() {
			curDir = info.Name()
			curDir = filepath.Join(basePath, curDir)
		}
		// check if it is a regular file (not dir)
		if info.Mode().IsRegular() {
			fe := fileEntry{}
			fe.parentPath = curDir
			fe.rootPath = rootPath

			fe.path = filepath.Join(info.Name())
			// start with base path since it is a directory
			*fileEntries = append(*fileEntries, fe)
		}
		return
	})
}

func getZipfile(path string) (archive *os.File, err error) {
	_, err = os.Create(path)
	if err != nil {

	}

	if _, err = os.Stat(path); os.IsNotExist(err) {
		// Handle new file
		archive, err = os.Create(path)
		if err != nil {
			archive, err = os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
			if err != nil {
				return
			}
		}
	} else {
		archive, err = os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			return
		}
	}
	return
}

// archiveFiles make a zip archive and fill with information from list of fileEntries
func archiveFiles(zipFileName string, fileEntries []fileEntry) (err error) {
	var archive *os.File
	archive, err = getZipfile(zipFileName)
	if err != nil {
		fmt.Fprintln(os.Stderr, colour(brightRed, err.Error()))
		return
	}

	defer archive.Close()

	zipWriter := zip.NewWriter(archive)
	defer zipWriter.Close()

	for _, fileEntry := range fileEntries {
		file, err := os.Open(fileEntry.fullPath())
		if err != nil {
			fmt.Fprintln(os.Stderr, colour(brightRed, err.Error()))
			continue
		}
		defer file.Close()

		dest, err := zipWriter.Create(fileEntry.archivePath())
		if err != nil {
			fmt.Fprintln(os.Stderr, colour(brightRed, err.Error()))
			continue
		}
		if _, err := io.Copy(dest, file); err != nil {
			fmt.Fprintln(os.Stderr, colour(brightRed, err.Error()))
			continue
		}
	}

	return
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

	var fileEntries = []fileEntry{}
	var errorMsgs = []string{}

	if !args.Unzip {
		// Populate list of files
		for _, file := range fileEntries {
			walkAllFilesInDir(file.path, &fileEntries, &errorMsgs)
		}
		if len(fileEntries) == 0 {
			fmt.Fprintln(os.Stderr, colour(brightRed, "no valid files found"))
			os.Exit(1)
		}
	}

	// Exit for now
	if 1 == 1 {
		os.Exit(0)
	}

	var archive *os.File
	archive, err := getZipfile(args.Zipfile)
	if err != nil {
		fmt.Fprintln(os.Stderr, colour(brightRed, err.Error()))
		os.Exit(1)
	}

	defer archive.Close()

	archiveFiles(args.Zipfile, fileEntries)

	// zipWriter := zip.NewWriter(archive)
	// defer zipWriter.Close()

	// for _, path := range files {
	// 	file, err := os.Open(path.parentPath)
	// 	if err != nil {
	// 		fmt.Fprintln(os.Stderr, colour(brightRed, err.Error()))
	// 		continue
	// 	}
	// 	defer file.Close()

	// 	dest, err := zipWriter.Create(path.parentPath)
	// 	if err != nil {
	// 		fmt.Fprintln(os.Stderr, colour(brightRed, err.Error()))
	// 		continue
	// 	}
	// 	if _, err := io.Copy(dest, file); err != nil {
	// 		fmt.Fprintln(os.Stderr, colour(brightRed, err.Error()))
	// 		continue
	// 	}
	// }
}
