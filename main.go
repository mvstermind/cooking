package main

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

var (
	url string = "https://picsum.photos/1920/1080"
	wg  sync.WaitGroup
)

func main() {
	startTime := time.Now()
	if len(os.Args) < 1 {
		osArg := 1
		wg.Add(1)
		go downloadFiles(osArg, &wg)
		wg.Wait()
		fmt.Println("program ran in: ", time.Since(startTime))
		return
	}
	arg := os.Args[1]

	osArg, err := strconv.Atoi(arg)
	if err != nil {
		log.Println("couldn't convert os.arg to int: ", err)
		return
	}
	if argValid(arg) {
		wg.Add(1)
		go downloadFiles(osArg, &wg)
		wg.Wait()
		fmt.Println("program ran in: ", time.Since(startTime))
		return
	}
}

func downloadFiles(arg int, wg *sync.WaitGroup) {
	for i := 0; i < arg; i++ {
		defer wg.Done()
		resp := getImg()

		filePath := fmt.Sprintf("/Users/kneehead/Documents/wallpapers/image%v.jpg", fileCount())
		out, err := os.Create(filePath)
		if err != nil {
			log.Fatalf("Error while creating file: %v", err)
		}
		defer out.Close()

		// Copy the image data to the file
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			log.Fatalf("Error while writing to file: %v", err)
		}

		log.Println("Image successfully downloaded and saved to", filePath)
	}
}

func fileCount() int {
	count := 0
	files := "/Users/kneehead/Documents/wallpapers/"

	err := filepath.Walk(files, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if !info.IsDir() {
			count++
		}
		return nil
	})
	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", files, err)
	}
	return count
}

func getImg() *http.Response {
	r, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error while sending GET request: %v", err)
	}
	return r
}

func argValid(osArg string) bool {
	if osArg == "" {
		log.Println("put number of files to download")
		return false
	}
	return true
}
