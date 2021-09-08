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

// hasEntry check for duplicate absolute paths
func hasEntry(check fileEntry, feList *[]fileEntry) (found bool) {
	for _, fe := range *feList {
		if check.fullPath() == fe.fullPath() {
			return true
		}
	}
	return
}

// fullPath get full path for an entry
func (fe *fileEntry) fullPath() (fullPath string) {
	fullPath = filepath.Join(fe.parentPath, fe.path)
	// fmt.Printf("%+v\n", fe)
	return
}

// archivePath get path for archive for entry
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

// colour get colour output
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

// getFileForWriting get file for new or upened zip file
func getFileForWriting(path string) (archive *os.File, err error) {
	// _, err = os.Create(path)
	// if err != nil {
	// 	fmt.Fprintln(os.Stderr, colour(brightRed, err.Error()))
	// 	return
	// }

	if _, err = os.Stat(path); os.IsNotExist(err) {
		// Handle new file
		archive, err = os.Create(path)
		if err != nil {
			archive, err = os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
			if err != nil {
				fmt.Fprintln(os.Stderr, colour(brightRed, err.Error()))
				return
			}
		}
	} else {
		archive, err = os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			fmt.Fprintln(os.Stderr, colour(brightRed, err.Error()))
			return
		}
	}
	return
}

type zipFileEntry struct {
	name             string
	compressedSize   uint64
	uncompressedSize uint64
}

func printEntries(name string) {
	entries, err := fileList(name)
	if err != nil {
		fmt.Fprintln(os.Stderr, colour(brightRed, "no valid files found"))
		os.Exit(1)
	}
	for _, file := range entries {
		fmt.Printf("name %s compressed %d uncompressed %d\n", file.name, file.compressedSize, file.uncompressedSize)
	}
}

func fileList(name string) (entries []zipFileEntry, err error) {
	zf, err := zip.OpenReader(name)
	if err != nil {
		return
	}
	defer zf.Close()

	entries = make([]zipFileEntry, 0, 0)

	for _, file := range zf.File {
		entry := zipFileEntry{}
		entry.name = file.Name
		entry.compressedSize = file.CompressedSize64
		entry.uncompressedSize = file.UncompressedSize64
		entries = append(entries, entry)
	}

	return
}

// archiveFiles make a zip archive and fill with information from list of fileEntries
func archiveFiles(zipFileName string, fileEntries []fileEntry) (err error) {
	var archive *os.File
	archive, err = getFileForWriting(zipFileName)
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

/*
	The zip utility has a lot of options. It is not know at the time of this
	writing how much of the original utility will be implemented.
*/

var args struct {
	Unzip       bool     `arg:"-U" help:"unzip the archive"`
	File        string   `arg:"-f" help:"verbosity level"`
	List        bool     `arg:"l" help:"list entries in zip file"`
	Recursive   bool     `arg:"-r" help:"recursive"`
	Update      bool     `arg:"-u" help:"update existing"`
	Add         bool     `arg:"-a" help:"add if not existing"`
	Zipfile     string   `arg:"positional"`
	SourceFiles []string `arg:"positional"`
}

func main() {
	p := arg.MustParse(&args)
	fmt.Println(args.File)

	if args.Zipfile == "" {
		p.Fail(colour(brightRed, "no zipfile specified"))
	}

	if args.Unzip {
		if len(args.SourceFiles) > 0 {
			p.Fail(colour(brightRed, "can't specify source files witn unzip"))
		}
	}
	if !args.Unzip {
		if len(args.SourceFiles) == 0 {
			p.Fail(colour(brightRed, "no files to zip specified"))
		}
	}
	if !args.Unzip && !args.List {
		p.Fail(colour(brightRed, "either unzip or list must be specified"))
	}

	if args.List {
		entries, err := fileList(args.Zipfile)
		if err != nil {
			fmt.Fprintln(os.Stderr, colour(brightRed, "no valid files found"))
			os.Exit(1)
		}
		for _, file := range entries {
			fmt.Printf("name %s compressed %d uncompressed %d\n", file.name, file.compressedSize, file.uncompressedSize)
		}
	}

	var fileEntries = []fileEntry{}
	var errorMsgs = []string{}

	if !args.Unzip {
		// Populate list of files
		for _, path := range args.SourceFiles {
			walkAllFilesInDir(path, &fileEntries, &errorMsgs)
		}
		if len(fileEntries) == 0 {
			fmt.Fprintln(os.Stderr, colour(brightRed, "no valid files found"))
			os.Exit(1)
		}
	}

	// Show what has been found
	for _, fileEntry := range fileEntries {
		fmt.Printf("Entry %+v", fileEntry)
	}

	// Exit for now
	if 1 == 1 {
		os.Exit(0)
	}

	var archiveFile *os.File
	archiveFile, err := getFileForWriting(args.Zipfile)
	if err != nil {
		fmt.Fprintln(os.Stderr, colour(brightRed, err.Error()))
		os.Exit(1)
	}

	defer archiveFile.Close()

	archiveFiles(args.Zipfile, fileEntries)
}
