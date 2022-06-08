package main

import (
	"crypto/md5"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// Proceed: fungsi yang digunakan untuk membaca file di folder temp, mengambil konten, melakukan
//			decode konten ke md5 dan merubahah nama file sesuai decode md5 isi dari konten tsb.
func Proceed() {
	var counterTotal, counterRenamed int

	err := filepath.WalkDir(tempPath, func(path string, d fs.DirEntry, err error) error {
		// jika ada permasalahan, langsung kembalikan
		if err != nil {
			return err
		}

		// jika menemukan sub dir, langsung kembalikan
		if d.IsDir() {
			return nil
		}

		// untuk menghitung jumlah file yang diketemukan
		counterTotal++

		// membaca konten file
		buf, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		sum := fmt.Sprintf("%x", md5.Sum(buf)) // konten >> ke >> md5

		// perubahan nama sumber file sesuai dengan hash md5 isi konten
		if err := os.Rename(path, filepath.Join(tempPath, fmt.Sprintf("file-%s.txt", sum))); err != nil {
			return err
		}

		counterRenamed++

		return nil
	})
	if err != nil {
		log.Println("Error Message:", err)
	}

	log.Printf("%d/%d files renamed", counterRenamed, counterTotal)
}
