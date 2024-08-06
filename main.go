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
    episodeTracker     = make(map[int]int)
    chatIDs            = []int64{}
    lastRequestTimes   = make(map[int64]time.Time)
    requestInterval    = 12 * time.Hour // Интервал между запросами свежей подборки
)

func main() {
    botToken := "5160413773:AAFUgxwGQhWyFbSB2eECW3506AyTNgkXPEQ"

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
    animes, err := getAnimesFromShikimori()
    if err != nil {
        fmt.Println("Failed to fetch animes:", err)
        return
    }

    sortAnimesByScore(animes)

    for _, anime := range animes {
        previousEpisodes, exists := episodeTracker[anime.Id]
        if !exists || anime.Episode > previousEpisodes {
            episodeTracker[anime.Id] = anime.Episode
            for _, chatID := range chatIDs {
                photoMessage := tu.Photo(
                    tu.ID(chatID),
                    tu.FileFromURL("https://shikimori.one" + anime.Image.Original),
                ).WithCaption(formatAnime(anime))

                _, _ = bot.SendPhoto(photoMessage)
            }
        }
    }
}
