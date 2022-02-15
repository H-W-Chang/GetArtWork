package met

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/dustin/go-humanize"
)

const searchAPI = "https://collectionapi.metmuseum.org/public/collection/v1/search"
const objectAPI = "https://collectionapi.metmuseum.org/public/collection/v1/objects"

func Search(artist, title, mediumType string, download bool, downloadPath string) {
	var wg sync.WaitGroup
	client := &http.Client{}
	req, err := http.NewRequest("GET", searchAPI, nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	q := req.URL.Query()
	q.Add("hasImages", "true")
	if artist != "" {
		q.Add("q", artist)
		q.Add("artistOrCulture", "true")
	}
	if title != "" {
		q.Add("title", title)
	}
	req.URL.RawQuery = q.Encode()
	log.Println(req.URL.String())
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var objmap map[string]interface{}
	err = json.Unmarshal(body, &objmap)
	if err != nil {
		log.Fatal(err)
	}
	for _, id := range objmap["objectIDs"].([]interface{}) {
		wg.Add(1)
		s := fmt.Sprintf("%d", int(id.(float64)))
		go getObject(s, mediumType, downloadPath, download, &wg)
		break
	}
	wg.Wait()
}

func getObject(id, mediumType, downloadPath string, download bool, wg *sync.WaitGroup) {
	defer wg.Done()
	url := objectAPI + "/" + id
	objmap := getObjectInfo(url)
	objectName := objmap["objectName"].(string)
	if objectName != mediumType {
		return
	}
	imagePath := objmap["primaryImage"].(string)
	log.Println("id:", id, "medium:", mediumType, "path:", imagePath)
	title := objmap["title"].(string)
	objectDate := objmap["objectDate"].(string)
	medium := objmap["medium"].(string)
	artist := objmap["artistDisplayName"].(string)
	r, _ := regexp.Compile(`\([0-9.x ]+`)
	dimensions := objmap["dimensions"].(string)
	if dimensions != "" {
		dimensions = strings.ReplaceAll(r.FindString(objmap["dimensions"].(string)), " ", "")[1:]
	}
	filename := artist + "_" + title + "_" + objectDate + "_" + medium + "_" + dimensions
	filename = strings.ReplaceAll(filename, ":", "")
	filename = strings.ReplaceAll(filename, "/", "")
	filename = strings.ReplaceAll(filename, "\"", "")
	log.Print(filename)
	if download {
		downloadImage(filename, artist, imagePath, downloadPath)
	}
	return
}

func getObjectInfo(url string) (objmap map[string]interface{}) {
	resp, err := http.Get(url)
	if err != nil {
		log.Print(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	// var objmap map[string]string
	err = json.Unmarshal(body, &objmap)
	if err != nil {
		log.Fatal(err)
	}
	return
}

// WriteCounter counts the number of bytes written to it. It implements to the io.Writer interface
// and we can pass this into io.TeeReader() which will report progress on each write cycle.
type WriteCounter struct {
	Total uint64
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	return n, nil
}

func (wc WriteCounter) PrintProgress() {
	// Clear the line by using a character return to go back to the start and remove
	// the remaining characters by filling it with spaces
	fmt.Printf("\r%s", strings.Repeat(" ", 35))

	// Return again and print current status of download
	// We use the humanize package to print the bytes in a meaningful way (e.g. 10 MB)
	fmt.Printf("\rDownloading... %s complete", humanize.Bytes(wc.Total))
}

func downloadImage(filename, artist, imagePath, downloadPath string) {
	err := os.MkdirAll(filepath.Join(downloadPath, artist), os.ModeDir)
	if err != nil {
		log.Print(err)
	}
	out, err := os.Create(filepath.Join(downloadPath, artist, filename) + ".jpg")
	if err != nil {
		log.Print(err)
	}
	defer out.Close()
	resp, err := http.Get(imagePath)
	if err != nil {
		log.Print(err)
		return
	}
	defer resp.Body.Close()

	counter := &WriteCounter{}
	// io.Copy(out, resp.Body)
	if _, err = io.Copy(out, io.TeeReader(resp.Body, counter)); err != nil {
		out.Close()
		return
	}
}
