// anime.go
package main

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
)

type Anime struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Title       string    `json:"russian"`
	Image       ImageInfo `json:"image"`
	Score       string    `json:"score"`
	EpisodesAll int       `json:"episodes"`
	Episode     int       `json:"episodes_aired"`
}

type ImageInfo struct {
	Original string `json:"original"`
	Preview  string `json:"preview"`
	X96      string `json:"x96"`
	X48      string `json:"x48"`
}

func getCurrentSeason() string {
	year, month, _ := time.Now().Date()
	season := ""

	switch month {
	case time.December, time.January, time.February:
		season = "winter"
	case time.March, time.April, time.May:
		season = "spring"
	case time.June, time.July, time.August:
		season = "summer"
	case time.September, time.October, time.November:
		season = "fall"
	}

	return fmt.Sprintf("%s_%d", season, year)
}
// https://shikimori.one/api/animes?season=summer_2024&kind=tv&limit=99
func getAnimesFromShikimori() ([]Anime, error) {
	season := getCurrentSeason()
	apiURL := fmt.Sprintf("https://shikimori.one/api/animes?season=%s&kind=tv&limit=99", season)
	client := resty.New()
	var animes []Anime

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetResult(&animes).
		Get(apiURL)

	if err != nil {
		return nil, fmt.Errorf("Не удалось получить анимешки: %v", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("не удалось получить аниме: получен код состояния %d", resp.StatusCode())
	}

	return animes, nil
}

func sortAnimesByScore(animes []Anime) {
	sort.Slice(animes, func(i, j int) bool {
		scoreI, _ := strconv.ParseFloat(animes[i].Score, 64)
		scoreJ, _ := strconv.ParseFloat(animes[j].Score, 64)
		return scoreI > scoreJ
	})
}

func formatAnime(anime Anime) string {
	title := anime.Title
	if title == "" {
		title = anime.Name
	}

<<<<<<< HEAD
    episodesAll := fmt.Sprintf("%d", anime.EpisodesAll)
    if anime.EpisodesAll == 0 {
        episodesAll = "?"
    }
    return fmt.Sprintf("%s\nРейтинг: %s ⭐️\nСерии: %d из %s 📺\nСсылка: https://shikimori.one/animes/%d", title, anime.Score, anime.Episode, episodesAll, anime.Id)
}
=======
	episodesAll := fmt.Sprintf("%d", anime.EpisodesAll)
	if anime.EpisodesAll == 0 {
		episodesAll = "?"
	}
	return fmt.Sprintf("%s\nРейтинг: %s ⭐️\nСерии: %d из %s 📺\nСсылка: https://shikimori.one/animes/%d", title, anime.Score, anime.Episode, episodesAll, anime.Id)
}
>>>>>>> 7d58eae9377343d3a7cccc3eac076519df88520c
