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
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal(errors.New("Please provide a directory to compress."))
	}
	baseDir := os.Args[1]
	createTarballs(baseDir)
}

func createTarballs(baseDir string) {
	err := os.Chdir(baseDir)
	if err != nil {
		log.Fatal(err)
	}

	files, err := ioutil.ReadDir("./")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if f.IsDir() {
			fmt.Println("compressing", f.Name())

			var b bytes.Buffer
			gz := gzip.NewWriter(&b)
			writer := tar.NewWriter(gz)

			addFilesToArchive(writer, f, baseDir)
			err = ioutil.WriteFile(f.Name()+".tar.gz", b.Bytes(), 0666)
			writer.Close()
			gz.Close()
			os.Chdir("../")
		}
	}
}

func statAndWriteHeader(tarWriter *tar.Writer, file os.FileInfo, curPath string) {
	var pathToCheck string
	if file.IsDir() {
		pathToCheck = curPath
	} else {
		pathToCheck = curPath + "/" + file.Name()
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
		contents, err := ioutil.ReadFile(curPath + "/" + file.Name())
		if err != nil {
			log.Fatal(err)
		}

		_, err = tarWriter.Write(contents)
		if err != nil {
			log.Fatal(err)
		}
	}
}
