package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal(errors.New("Please provide a file or directory to compress."))
	}
	basePath := os.Args[1]
	createTarballs(basePath)
}

func createTarballs(pathToCheck string) {
	file, err := os.Stat(pathToCheck)
	if err != nil {
		log.Fatal(err)
	}

	if file.IsDir() {
		err = os.Chdir(pathToCheck)
		if err != nil {
			log.Fatal(err)
		}

		files, err := ioutil.ReadDir("./")
		if err != nil {
			log.Fatal(err)
		}

		for _, f := range files {
			if f.IsDir() {
				fmt.Println("compressing", file.Name())

				var b bytes.Buffer
				gz := gzip.NewWriter(&b)
				writer := tar.NewWriter(gz)

				addFilesToArchive(writer, f, pathToCheck)
				err = ioutil.WriteFile(f.Name()+".tar.gz", b.Bytes(), 0666)
				writer.Close()
				gz.Close()
				os.Chdir("../")
			}
		}
	} else {
		fmt.Println("compressing file", file.Name())

		err := os.Chdir(filepath.Dir(pathToCheck))
		if err != nil {
			log.Fatal(err)
		}

		var b bytes.Buffer
		gz := gzip.NewWriter(&b)
		writer := tar.NewWriter(gz)

		addFilesToArchive(writer, file, pathToCheck)
		err = ioutil.WriteFile(file.Name()+".tar.gz", b.Bytes(), 0666)
		writer.Close()
		gz.Close()

		fmt.Printf("archive %s created\n", file.Name()+".tar.gz")
	}
}

func statAndWriteHeader(tarWriter *tar.Writer, file os.FileInfo, curPath string) {
	var pathToCheck string
	if file.IsDir() {
		pathToCheck = curPath
	} else {
		tmp, _ := os.Stat(curPath)
		if !tmp.IsDir() {
			pathToCheck = curPath
		} else {
			pathToCheck = curPath + "/" + file.Name()
		}
	}
	f, err := os.Stat(pathToCheck)
	if err != nil {
		log.Fatal(err)
	}

	header, err := tar.FileInfoHeader(f, "")
	if err != nil {
		log.Fatal(err)
	}

	err = tarWriter.WriteHeader(header)
	if err != nil {
		log.Fatal(err)
	}
}

func addFilesToArchive(tarWriter *tar.Writer, file os.FileInfo, curPath string) {
	if file.IsDir() {
		statAndWriteHeader(tarWriter, file, curPath)

		err := os.Chdir(curPath + "/" + file.Name())
		if err != nil {
			log.Fatal(err)
		}

		files, err := ioutil.ReadDir(curPath + "/" + file.Name())
		if err != nil {
			log.Fatal(err)
		}

		for _, f := range files {
			addFilesToArchive(tarWriter, f, curPath+"/"+file.Name())
		}

		os.Chdir("../")
	} else {
		statAndWriteHeader(tarWriter, file, curPath)
		var pathToWrite string
		if strings.Compare(filepath.Base(curPath), file.Name()) == 0 {
			pathToWrite = curPath
		} else {
			pathToWrite = curPath + "/" + file.Name()
		}
		contents, err := ioutil.ReadFile(pathToWrite)
		if err != nil {
			log.Fatal(err)
		}

		_, err = tarWriter.Write(contents)
		if err != nil {
			log.Fatal(err)
		}
	}
}
