package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/zhanglp92/plugins/imports/internal"
)

var (
	writeToFile   bool
	updateComment bool
	_testing      bool
)

func init() {

	// the flag parsing in this init function is necessary for the command binary to execute
	// using a bool to skip flag parsing while running unit tests as a workaround
	// without this the unit test cases will fail at flag.Parse() and will exit
	if _testing {
		return
	}

	flag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, "Usage of %s:\n  %s [-whc] <file|dir>\n", os.Args[0], os.Args[0])
		flag.PrintDefaults()
	}

	flag.BoolVar(&writeToFile, "w", false, "write result to (source) file instead of stdout")
	flag.BoolVar(&updateComment, "c", false, "update comment to node")
	flag.Parse()
}

func main() {
	var err error
	if len(flag.Args()) < 1 {
		err = processRaw(os.Stdin, os.Stdout)
	} else {
		err = processPath(flag.Arg(0))
	}
	if err != nil {
		log.Printf("error occurred : %v", err)
	}
}

func processPath(path string) error {
	stat, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("unable to stat %s : %v", path, err)
	}

	if stat.IsDir() {
		return processDir(path)
	}

	return processFile(path)
}

func processDir(dirpath string) error {
	return filepath.Walk(dirpath, func(itemPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// skip vendor directories
		if info.IsDir() && info.Name() == "vendor" {
			return filepath.SkipDir
		}

		if path.Ext(itemPath) != ".go" {
			return nil
		}

		// at this point, item is a go source file
		return processFile(itemPath)
	})
}

func processFile(filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error while opening file : %v", err)
	}
	defer func() {
		_ = f.Close()
	}()

	fileData, err := ioutil.ReadAll(f)
	if err != nil {
		return fmt.Errorf("unable to read data in file : %v", err)
	}

	newData, err := process(fileData)
	if err != nil {
		return fmt.Errorf("error while process file : %v", err)
	}

	// if data is unchanged, stop here
	if bytes.Equal(newData, fileData) {
		return nil
	}

	if writeToFile {
		log.Printf("old=[%d]; new=[%d] : %s", len(fileData), len(newData), filePath)
		if err := ioutil.WriteFile(filePath, newData, os.ModePerm); err != nil {
			return fmt.Errorf("error while writing to file : %v", err)
		}
	} else {
		log.Printf("--- %s\n", filePath)
		fmt.Printf("%s", newData)
	}

	return nil
}

func processRaw(in io.Reader, out io.Writer) error {
	src, err := ioutil.ReadAll(in)
	if err != nil {
		return err
	}

	res, err := process(src)
	if err != nil {
		return err
	}

	_, err = out.Write(res)
	return err
}

func process(src []byte) ([]byte, error) {
	return internal.Process(src, updateComment)
}
