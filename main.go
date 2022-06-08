package main

import (
	"log"
	"math/rand"
	"time"
)

const (
	totalFile     int = 200000
	contentLength int = 5000

	tempPath     string = "./dummy"
	hashWorker   int    = 2
	renameWorker int    = 20
	buffer       int    = 10000
)

type FileInfo struct {
	FilePath  string // lokasi file
	Content   []byte // isi dari file
	Sum       string // md5 dari isi konten
	IsRenamed bool   // status sudah di rename atau belum
}

func main() {
	log.Println("Start")
	startTime := time.Now()

	// Membuat dummy data
	// GenerateFiles()

	// Normal Process:
	// Proceed() // 2, 1.5, 1.3, 1.4

	// Concurrency Process:
	ProceedConcurrency()

	log.Println("Done in ", time.Since(startTime).Seconds(), "seconds")
}

func ProceedConcurrency() {
	// fmt.Println("==> Monitoring Goroutine:", runtime.NumGoroutine())

	// pipeline 1: read files
	chanReadFiles := ReadFiles()

	// pipeline 2: create worker for sum content to hash md5
	chanGetMd5 := make([]<-chan FileInfo, hashWorker)
	for i := range chanGetMd5 {
		chanGetMd5[i] = GetSumMd5(chanReadFiles)
	}
	chanMergeFiles := MergeChanWorkerFileInfo(chanGetMd5...)

	// pipeline 3: create worker for rename file
	chanRenameFile := make([]<-chan FileInfo, renameWorker)
	for i := range chanRenameFile {
		chanRenameFile[i] = RenameFile(chanMergeFiles)
	}
	chanFullProcess := MergeChanWorkerFileInfo(chanRenameFile...)

	var counterTotal, counterRenamed int
	for file := range chanFullProcess {
		counterTotal++
		if file.IsRenamed {
			counterRenamed++
		}
	}
	log.Printf("%d / %d file renamed \n", counterRenamed, counterTotal)

	// fmt.Println("==> Monitoring Goroutine:", runtime.NumGoroutine())
}

// RandomString: digunakan untuk mendapatkan string acak dengan lebar sesuai nilai parameter.
func RandomString(length int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	random := make([]rune, length)
	for i := range random {
		random[i] = letters[rand.Intn(len(letters))]
	}

	return string(random)
}
