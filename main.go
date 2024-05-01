package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/maikpro/mangadownloader/services"
)

type Config struct {
	IsList        bool
	ChapterNumber uint
	IsLocal       bool
	IsTelegram    bool
}

func main() {
	var config Config
	config.parseCLI()

	if config.IsList {
		arcs, err := services.GetArcList()
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		config.ChapterNumber = showArcList(arcs)
	}

	if config.ChapterNumber == 0 {
		log.Fatal("chapterNumber is not provided...")
		os.Exit(1)
	}
	log.Printf("You have selected %d", config.ChapterNumber)
	if config.IsLocal {
		var path *string
		action := func() {
			path = useLocally(&config.ChapterNumber)
		}
		spinner.New().Title(fmt.Sprintf("Downloading Chapter %d locally", config.ChapterNumber)).Action(action).Run()
		log.Printf("✅ Saved Chapter locally at %s\n", *path)
	}

	if config.IsTelegram {
		action := func() {
			useTelegram(&config.ChapterNumber)
		}
		spinner.New().Title(fmt.Sprintf("Sending Chapter %d to Telegram", config.ChapterNumber)).Action(action).Run()
		log.Println("✅ Sent Chapter to your Telegram Chat")
	}

	os.Exit(0)
}

func useLocally(chapterNumber *uint) *string {
	chapter, err := services.GetChapter(*chapterNumber)
	if err != nil {
		log.Fatalf("%s. Result: Couldn't get Chapter...", err)
		os.Exit(1)
	}

	var path *string
	for index, page := range chapter.Pages {
		path, err = services.DownloadPage(page.Url, fmt.Sprintf("./chapters/%d_%s", chapter.Number, chapter.Name), fmt.Sprintf("page_%d", index))
		if err != nil {
			log.Fatalf("%s. Result: Couldn't get Chapter...", err)
			os.Exit(1)
		}
	}
	return path
}

func (c *Config) parseCLI() {
	flag.UintVar(&c.ChapterNumber, "chapterNumber", 0, "ChapterNumber that has to be downloaded from onepiece-tube.com/manga/kapitel-mangaliste")
	flag.BoolVar(&c.IsLocal, "local", false, "set true if you want to download chapters useLocally on your computer. Don't forget setting the chapterNumber")
	flag.BoolVar(&c.IsList, "list", false, "show all current one piece manga chapters from onepiece-tube.com")
	flag.BoolVar(&c.IsTelegram, "telegram", false, "send chapter to your telegram chat")
	flag.Parse()
}

func useTelegram(chapterNumber *uint) {
	chapter, err := services.GetChapter(*chapterNumber)
	if err != nil {
		log.Fatalf("%s. Result: Couldn't get Chapter...", err)
		os.Exit(1)
	}

	err = services.SendChapter(*chapter)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func createOptions(arc services.Arc) []huh.Option[string] {
	var options []huh.Option[string]
	for _, chapter := range arc.Entry {
		if !chapter.IsAvailable {
			continue
		}
		option := huh.NewOption(fmt.Sprintf("%d - %s - [%s]", chapter.Number, chapter.Name, arc.Name), fmt.Sprintf("%d", chapter.Number))
		options = append(options, option)
	}

	return options
}

// try this: https://github.com/charmbracelet/bubbletea/blob/master/examples/list-default/main.go
func showArcList(arcs []services.Arc) uint {
	var options []huh.Option[string]
	var selectedChapter string

	for _, arc := range arcs {
		option := createOptions(arc)
		options = append(options, option...)
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				// Title(arc.Name).
				Options(options...).
				Value(&selectedChapter),
		),
	)
	err := form.Run()
	if err != nil {
		log.Fatal(err)
	}

	num, err := strconv.ParseUint(selectedChapter, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
	}

	return uint(num)
}
