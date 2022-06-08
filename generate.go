package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

// Init: Inisialisasi proses untuk file generate.go
func init() {
	rand.Seed(time.Now().UnixNano())
}

// GenerateFiles: fungsi yang digunakan untuk membuat file dummy dengan jumlah dan lebar konten sesuai dengan constant yang ditentukan
func GenerateFiles() {
	os.RemoveAll(tempPath)
	os.Mkdir(tempPath, os.ModePerm)

	// membuat file dummy dengan nama sesuai dengan urutan dan isi string acak sejumlah contentLength
	for i := 0; i < totalFile; i++ {
		fileName := filepath.Join(tempPath, fmt.Sprintf("file-%06d.txt", i+1))

		if err := ioutil.WriteFile(fileName, []byte(RandomString(contentLength)), os.ModePerm); err != nil {
			log.Println("ERROR WRITE FILE:", err.Error())
		}

		if counter := i + 1; counter%500 == 0 {
			log.Println(counter, "files created")
		}
	}

	log.Printf("%d of total files created", totalFile)
}
