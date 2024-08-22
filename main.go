package main

import (
    "encoding/json"
    "fmt"
    "os"
    "time"

    "github.com/mymmrac/telego"
    th "github.com/mymmrac/telego/telegohandler"
    tu "github.com/mymmrac/telego/telegoutil"
)

var (
    episodeTracker     = make(map[int]int)
    chatIDs            = []int64{}
    lastRequestTimes   = make(map[int64]time.Time)
    requestInterval    = 12 * time.Hour // Интервал между запросами свежей подборки
    episodeTrackerFile = "episode_tracker.json"
    chatIDsFile        = "chat_ids.json"
)

func main() {
    botToken := "5160413773:AAGyjpQbrAL-1hR6bnV8GwDY3ioIjxBVRzk"

    bot, err := telego.NewBot(botToken, telego.WithDefaultDebugLogger())
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    // Загружаем данные из файлов
    loadEpisodeTracker()
    loadChatIDs()

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
    animes, err := getAnimesFromShikimori()
    if err != nil {
        fmt.Println("Failed to fetch animes:", err)
        return
    }

    sortAnimesByScore(animes)

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
        if !exists || anime.Episode > previousEpisodes {
            episodeTracker[anime.Id] = anime.Episode
            saveEpisodeTracker()

            for _, chatID := range chatIDs {
                photoMessage := tu.Photo(
                    tu.ID(chatID),
                    tu.FileFromURL("https://shikimori.one" + anime.Image.Original),
                ).WithCaption(formatAnime(anime))

                _, err := bot.SendPhoto(photoMessage)
                if err != nil {
                    fmt.Printf("Failed to send message to chat %d: %v\n", chatID, err)
                    // Можно добавить логику для повторной отправки
                }
            }
        }
    }
}

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
