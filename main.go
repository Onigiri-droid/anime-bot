// main.go
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

var (
	episodeTracker   = make(map[int]int)
	chatIDs          = []int64{}
	lastRequestTimes = make(map[int64]time.Time)
	requestInterval  = 12 * time.Hour // Интервал между запросами свежей подборки
)

func main() {
	botToken := "5160413773:AAGyjpQbrAL-1hR6bnV8GwDY3ioIjxBVRzk"

    err := loadEpisodeTracker()
    if err != nil {
        fmt.Println("Не удалось загрузить episode tracker: ", err)
    }

    defer saveEpisodeTracker()

	bot, err := telego.NewBot(botToken, telego.WithDefaultDebugLogger())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	updates, _ := bot.UpdatesViaLongPolling(nil)
	bh, _ := th.NewBotHandler(bot, updates)

	defer bh.Stop()
	defer bot.StopLongPolling()

	bh.Handle(startHandler, th.CommandEqual("start"))
	bh.Handle(freshAnimeHandler, th.TextEqual("Свежая подборка"))

	go func() {
		for {
			checkForNewEpisodes(bot)
			time.Sleep(1 * time.Hour)
		}
	}()

	bh.Start()
}

func startHandler(bot *telego.Bot, update telego.Update) {
	chatID := tu.ID(update.Message.Chat.ID)
	chatIDs = append(chatIDs, update.Message.Chat.ID)

	keyboard := tu.Keyboard(
		tu.KeyboardRow(
			tu.KeyboardButton("Свежая подборка"),
		),
	).WithResizeKeyboard()

	message := tu.Message(
		chatID,
		"Привет! Я бот, который сообщает о новинках аниме и позволяет подписаться на уведомления о выходе новых серий 📺 ✨\nДайте знать, что вам нравится, и я буду держать вас в курсе!",
	).WithReplyMarkup(keyboard)

	_, _ = bot.SendSticker(
		tu.Sticker(
			chatID,
			tu.FileFromID("CAACAgIAAxkBAAEMkJdmqNnOUgkktshH0TJRYDmAGcb1wwACPwAD0DQ6Jwe4M9oNkfpONQQ"),
		),
	)

	_, _ = bot.SendMessage(message)
}

func checkForNewEpisodes(bot *telego.Bot) {
    fmt.Println("Проверка на наличие новых эпизодов...")
    animes, err := getAnimesFromShikimori()
    if err != nil {
        fmt.Println("Не удалось получить аниме:", err)
        return
    }

    sortAnimesByScore(animes)
    fmt.Println("Аниме успешно получены и отсортированы")

    for _, anime := range animes {
        previousEpisodes, exists := episodeTracker[anime.Id]
        fmt.Printf("Проверка аниме с ID %d: предыдущие эпизоды = %d, текущие эпизоды = %d\n", anime.Id, previousEpisodes, anime.Episode)

        if !exists || anime.Episode > previousEpisodes {
            episodeTracker[anime.Id] = anime.Episode
            fmt.Printf("Обновлено: Аниме с ID %d теперь имеет %d эпизодов.\n", anime.Id, anime.Episode)

            for _, chatID := range chatIDs {
                photoMessage := tu.Photo(
                    tu.ID(chatID),
                    tu.FileFromURL("https://shikimori.one"+anime.Image.Original),
                ).WithCaption(formatAnime(anime))

                _, err := bot.SendPhoto(photoMessage)
                if err != nil {
                    fmt.Println("Ошибка при отправке фото:", err)
                } else {
                    fmt.Printf("Уведомление отправлено в чат %d для аниме с ID %d\n", chatID, anime.Id)
                }
            }
        }
    }

    err = saveEpisodeTracker()
    if err != nil {
        fmt.Println("Не удалось сохранить episode tracker:", err)
    } else {
        fmt.Println("episodeTracker успешно сохранен")
    }
}

func simulateAnimeUpdate(bot *telego.Bot) {
    fmt.Println("Запуск симуляции обновления аниме...")
    time.Sleep(30 * time.Second) // Ждем 30 секунд перед обновлением

    animeID := 49785
    newEpisodeCount := 8

    previousCount, exists := episodeTracker[animeID]
    if !exists {
        previousCount = 0
    }

    episodeTracker[animeID] = newEpisodeCount
    fmt.Printf("Аниме с ID %d обновлено: теперь %d эпизодов.\n", animeID, newEpisodeCount)

    if newEpisodeCount > previousCount {
        for _, chatID := range chatIDs {
            photoMessage := tu.Photo(
                tu.ID(chatID),
                tu.FileFromURL("https://shikimori.one/system/animes/original/49785.jpg"),
            ).WithCaption(fmt.Sprintf("Аниме обновлено! Теперь %d серий.", newEpisodeCount))

            _, err := bot.SendPhoto(photoMessage)
            if err != nil {
                fmt.Println("Ошибка при отправке уведомления:", err)
            } else {
                fmt.Printf("Уведомление отправлено в чат %d для аниме с ID %d\n", chatID, animeID)
            }
        }
    }

    err := saveEpisodeTracker()
    if err != nil {
        fmt.Println("Не удалось сохранить episode tracker:", err)
    } else {
        fmt.Println("episodeTracker успешно сохранен после симуляции обновления")
    }
}




// func freshAnimeHandler(bot *telego.Bot, update telego.Update) {
//     chatID := update.Message.Chat.ID
//     currentTime := time.Now()

//     if lastRequestTime, exists := lastRequestTimes[chatID]; exists {
//         if currentTime.Sub(lastRequestTime) < requestInterval {
//             message := tu.Message(
//                 tu.ID(chatID),
//                 "Вы можете запросить свежую подборку раз в 12 часов ⏰.\nПопробуйте позже ⌛️",
//             )
//             _, _ = bot.SendMessage(message)
//             return
//         }
//     }

//     lastRequestTimes[chatID] = currentTime

//     handleAnimes(bot, []int64{chatID})
// }

// func checkForNewEpisodes(bot *telego.Bot) {
//     handleAnimes(bot, chatIDs)
// }