package main

import (
	"crypto/md5"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
)

// ReadFiles: fungsi yang digunakan untuk membaca file dengan output chan metadata file yang dibaca
//            fungsi ini berjalan secara concurrency
func ReadFiles() <-chan FileInfo {
	chanOut := make(chan FileInfo, buffer)

	go func() {
		defer close(chanOut)

		err := filepath.WalkDir(tempPath, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() {
				return nil
			}

			buf, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			chanOut <- FileInfo{
				FilePath: path,
				Content:  buf,
			}

			return nil
		})

		if err != nil {
			log.Println("Reading File Error:", err.Error())
		}
	}()

	return chanOut
}

// GetSumMd5: melakukan perhitungan hash md5 untuk konten file
func GetSumMd5(chanIn <-chan FileInfo) <-chan FileInfo {
	chanOut := make(chan FileInfo, buffer)

	go func() {
		defer close(chanOut)

		for file := range chanIn {
			file.Sum = fmt.Sprintf("%x", md5.Sum(file.Content))

			chanOut <- file
		}
	}()

	return chanOut
}

// RenameFile: digunakan untuk me-rename file
func RenameFile(chanIn <-chan FileInfo) <-chan FileInfo {
	chanOut := make(chan FileInfo, buffer)

	go func() {
		defer close(chanOut)

		for file := range chanIn {
			err := os.Rename(file.FilePath, filepath.Join(tempPath, fmt.Sprintf("file-%s.txt", file.Sum)))
			if err != nil {
				log.Println("Error Rename File:", err.Error())
			}
			file.IsRenamed = (err == nil)

			chanOut <- file
		}
	}()

	return chanOut
}

// MergeChanWorkerSumMd5: melakukan penggabungan chanel dari sekumpulan worker
func MergeChanWorkerFileInfo(chanels ...<-chan FileInfo) <-chan FileInfo {
	wg := new(sync.WaitGroup)
	chanOut := make(chan FileInfo, buffer)

	// penggabungan chanel ke chaOut
	wg.Add(len(chanels)) // []chan FileInfo

	for _, chanel := range chanels {
		go func(chanFile <-chan FileInfo) {
			defer wg.Done()

			for file := range chanFile {
				chanOut <- file
			}
		}(chanel)
	}

	go func() {
		wg.Wait()

		close(chanOut)
	}()

	return chanOut
}
