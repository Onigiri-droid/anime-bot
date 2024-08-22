// main.go
package main

import (
<<<<<<< HEAD
    "encoding/json"
    "fmt"
    "os"
    "time"
=======
	"fmt"
	"os"
	"time"
>>>>>>> 7d58eae9377343d3a7cccc3eac076519df88520c

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

var (
<<<<<<< HEAD
    episodeTracker     = make(map[int]int)
    chatIDs            = []int64{}
    lastRequestTimes   = make(map[int64]time.Time)
    requestInterval    = 12 * time.Hour // Интервал между запросами свежей подборки
    episodeTrackerFile = "episode_tracker.json"
    chatIDsFile        = "chat_ids.json"
)

func main() {
    botToken := "5160413773:AAGyjpQbrAL-1hR6bnV8GwDY3ioIjxBVRzk"
=======
	episodeTracker   = make(map[int]int)
	chatIDs          = []int64{}
	lastRequestTimes = make(map[int64]time.Time)
	requestInterval  = 12 * time.Hour // Интервал между запросами свежей подборки
)

func main() {
	botToken := "5160413773:AAGyjpQbrAL-1hR6bnV8GwDY3ioIjxBVRzk"
>>>>>>> 7d58eae9377343d3a7cccc3eac076519df88520c

    err := loadEpisodeTracker()
    if err != nil {
        fmt.Println("Не удалось загрузить episode tracker: ", err)
    }

<<<<<<< HEAD
    // Загружаем данные из файлов
    loadEpisodeTracker()
    loadChatIDs()

    updates, _ := bot.UpdatesViaLongPolling(nil)
    bh, _ := th.NewBotHandler(bot, updates)
=======
    defer saveEpisodeTracker()
>>>>>>> 7d58eae9377343d3a7cccc3eac076519df88520c

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
<<<<<<< HEAD
    chatID := update.Message.Chat.ID

    // Проверяем, есть ли уже chatID в списке
    exists := false
    for _, id := range chatIDs {
        if id == chatID {
            exists = true
            break
        }
    }

    // Если chatID не существует, добавляем его в список
    if !exists {
        chatIDs = append(chatIDs, chatID)
        saveChatIDs()
    }

    // Отправляем приветственное сообщение и клавиатуру
    keyboard := tu.Keyboard(
        tu.KeyboardRow(
            tu.KeyboardButton("Свежая подборка"),
        ),
    ).WithResizeKeyboard()

    message := tu.Message(
        tu.ID(chatID),
        "Привет! Я бот, который сообщает о новинках аниме и позволяет подписаться на уведомления о выходе новых серий 📺 ✨\nДайте знать, что вам нравится, и я буду держать вас в курсе!",
    ).WithReplyMarkup(keyboard)

    _, _ = bot.SendSticker(
        tu.Sticker(
            tu.ID(chatID),
            tu.FileFromID("CAACAgIAAxkBAAEMkJdmqNnOUgkktshH0TJRYDmAGcb1wwACPwAD0DQ6Jwe4M9oNkfpONQQ"),
        ),
    )
=======
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
>>>>>>> 7d58eae9377343d3a7cccc3eac076519df88520c

	_, _ = bot.SendMessage(message)
}

func stopHandler(bot *telego.Bot, update telego.Update) {
    chatID := update.Message.Chat.ID

    // Удаляем chatID из списка chatIDs
    for i, id := range chatIDs {
        if id == chatID {
            chatIDs = append(chatIDs[:i], chatIDs[i+1:]...)
            saveChatIDs()
            break
        }
    }

    message := tu.Message(
        tu.ID(chatID),
        "Вы успешно отписались от уведомлений. Если захотите вернуться, просто напишите /start.",
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

    // Удаляем аниме, которых больше нет в списке
    for id := range episodeTracker {
        found := false
        for _, anime := range animes {
            if anime.Id == id {
                found = true
                break
            }
        }
        if !found {
            delete(episodeTracker, id)
        }
    }

    for _, anime := range animes {
        previousEpisodes, exists := episodeTracker[anime.Id]
        fmt.Printf("Проверка аниме с ID %d: предыдущие эпизоды = %d, текущие эпизоды = %d\n", anime.Id, previousEpisodes, anime.Episode)

        if !exists || anime.Episode > previousEpisodes {
            episodeTracker[anime.Id] = anime.Episode
<<<<<<< HEAD
            saveEpisodeTracker()
=======
            fmt.Printf("Обновлено: Аниме с ID %d теперь имеет %d эпизодов.\n", anime.Id, anime.Episode)
>>>>>>> 7d58eae9377343d3a7cccc3eac076519df88520c

            for _, chatID := range chatIDs {
                photoMessage := tu.Photo(
                    tu.ID(chatID),
                    tu.FileFromURL("https://shikimori.one"+anime.Image.Original),
                ).WithCaption(formatAnime(anime))

                _, err := bot.SendPhoto(photoMessage)
                if err != nil {
<<<<<<< HEAD
                    fmt.Printf("Failed to send message to chat %d: %v\n", chatID, err)
                    // Можно добавить логику для повторной отправки
=======
                    fmt.Println("Ошибка при отправке фото:", err)
                } else {
                    fmt.Printf("Уведомление отправлено в чат %d для аниме с ID %d\n", chatID, anime.Id)
>>>>>>> 7d58eae9377343d3a7cccc3eac076519df88520c
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

<<<<<<< HEAD
func saveEpisodeTracker() {
    file, err := os.Create(episodeTrackerFile)
    if err != nil {
        fmt.Println("Failed to save episode tracker:", err)
        return
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    if err := encoder.Encode(episodeTracker); err != nil {
        fmt.Println("Failed to encode episode tracker:", err)
    }
}

func loadEpisodeTracker() {
    file, err := os.Open(episodeTrackerFile)
    if err != nil {
        fmt.Println("Failed to load episode tracker:", err)
        return
    }
    defer file.Close()

    decoder := json.NewDecoder(file)
    if err := decoder.Decode(&episodeTracker); err != nil {
        fmt.Println("Failed to decode episode tracker:", err)
    }
}

func saveChatIDs() {
    file, err := os.Create(chatIDsFile)
    if err != nil {
        fmt.Println("Failed to save chat IDs:", err)
        return
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    if err := encoder.Encode(chatIDs); err != nil {
        fmt.Println("Failed to encode chat IDs:", err)
    }
}

func loadChatIDs() {
    file, err := os.Open(chatIDsFile)
    if err != nil {
        fmt.Println("Failed to load chat IDs:", err)
        return
    }
    defer file.Close()

    decoder := json.NewDecoder(file)
    if err := decoder.Decode(&chatIDs); err != nil {
        fmt.Println("Failed to decode chat IDs:", err)
    }
}
=======
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
>>>>>>> 7d58eae9377343d3a7cccc3eac076519df88520c
