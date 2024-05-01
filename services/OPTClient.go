package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/maikpro/mangadownloader/models"
)

func GetChapter(chapterNumber uint) (*models.Chapter, error) {
	// note: One-Piece-Tube does not provide Manga from Chapter 1 - 419
	if chapterNumber < 420 {
		return nil, errors.New("onepiece-tube.com does not provide Manga from Chapter 1 - 419")
	}

	url := fmt.Sprintf("https://onepiece-tube.com/manga/kapitel/%d", chapterNumber)
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
		return nil, err
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	scriptRawText := doc.Find("#app > script").First().Text()
	jsonString := strings.Replace(strings.Split(scriptRawText, "=")[1], ";", "", -1)
	var newjson models.Data
	json.Unmarshal([]byte(jsonString), &newjson)
	chapter := &newjson.Chapter
	chapter.Number = chapterNumber
	return chapter, nil
}

func GetPageImage(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error downloading file:", err)
		return nil, err
	}
	defer response.Body.Close()
	imageData, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return imageData, nil
}

func DownloadPage(url string, savePath string, filename string) (*string, error) {
	imageData, err := GetPageImage(url)
	if err != nil {
		return nil, err
	}

	// Get the extension from image
	// Get the file extension from the URL
	ext := filepath.Ext(url)

	// Remove the leading dot from the extension
	ext = ext[1:]
	filename = fmt.Sprintf("%s.%s", filename, ext)
	fullpath, err := SaveFile(savePath, filename, imageData)
	if err != nil {
		log.Println("Error saving image locally:", err)
		return nil, err
	}

	return &fullpath, err
}

type OPTListData struct {
	Arcs    []OPTArc   `json:"arcs"`
	Entries []OPTEntry `json:"entries"`
}

// {"id":41,"name":"Egghead Arc","min":1058,"max":1113}
type OPTArc struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Min  int    `json:"min"`
	Max  int    `json:"mAX"`
}

/* {
"id":1307,
"name":"Patt",
"number":1113,
"category_id":3,
"arc_id":41,
"specials_id":0,
"lang":"ger",
"pages":15,
"is_available":true,
"date":"26.04.2024",
"href":"https://onepiece-tube.com/manga/kapitel/1113/1"
},
*/type OPTEntry struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Number      int    `json:"number"`
	Pages       int    `json:"pages"`
	Date        string `json:"date"`
	IsAvailable bool   `json:"is_available"`
}

func GetArcList() ([]Arc, error) {
	res, err := http.Get("https://onepiece-tube.com/manga/kapitel-mangaliste")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
		return nil, err
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	scriptRawText := doc.Find("#main-content > script").First().Text()
	jsonString := strings.Replace(strings.Split(scriptRawText, "=")[1], ";", "", -1)
	var optListData OPTListData
	err = json.Unmarshal([]byte(jsonString), &optListData)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	arcs := mapIntoArcs(optListData)
	return arcs, nil
}

// on OPT the arc and chapter are splitted into seperated arrays (arcs and entries)
// therefore map them into a single struct to have Arc and entries together.
type Arc struct {
	Name  string
	Entry []OPTEntry
}

func mapIntoArcs(optListData OPTListData) []Arc {
	var arcs []Arc

	for _, optListArc := range optListData.Arcs {
		var arc Arc
		arc.Name = strings.TrimRight(optListArc.Name, " ")

		for _, optEntry := range optListData.Entries {
			min := optListArc.Min
			max := optListArc.Max

			if optEntry.Number <= max && optEntry.Number >= min {
				arc.Entry = append(arc.Entry, optEntry)
			}
		}

		arcs = append(arcs, arc)
	}
	return arcs
}
