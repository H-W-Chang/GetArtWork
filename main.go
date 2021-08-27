package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/H-W-Chang/GetArtWork/pkg/met"
)

func main() {
	if _, err := os.Stat("logs"); err != nil {
		if os.IsNotExist(err) {
			// not exist
			err := os.Mkdir("logs", os.ModeDir)
			if err != nil {
				log.Print(err)
			}
		} else {
			// other error
			log.Print(err)
		}
	}
	f, err := os.OpenFile(filepath.Join("logs", "GetArtWork.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	var artist, title string
	sourcePtr := flag.String("s", "", "source of art work, support: MET")
	artistPtr := flag.String("a", "", "artist")
	titlePtr := flag.String("t", "", "title")
	mediumPtr := flag.String("m", "", "medium type: Painting, Drawing, Print, Album, Tapestry, Printtest")
	downloadPathPtr := flag.String("p", "C:/images", "download path")
	downloadPtr := flag.Bool("d", false, "download or not")
	flag.Parse()
	artist = *artistPtr
	title = *titlePtr
	mediumType := *mediumPtr
	download := *downloadPtr
	downloadPath := *downloadPathPtr
	source := *sourcePtr
	log.Printf("artist: %s, title: %s medium: %s source: %s path: %s\n", artist, title, mediumType, source, downloadPath)
	switch source {
	case "MET":
		met.Search(artist, title, mediumType, download, downloadPath)
	default:
		log.Fatal("no such source")
	}
}
