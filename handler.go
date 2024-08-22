package main

import (
    "time"
    "github.com/mymmrac/telego"
    "fmt"
    tu "github.com/mymmrac/telego/telegoutil"
)

func freshAnimeHandler(bot *telego.Bot, update telego.Update) {
    chatID := update.Message.Chat.ID
    currentTime := time.Now()

    // Проверка времени последнего запроса
    if lastRequestTime, exists := lastRequestTimes[chatID]; exists {
        if currentTime.Sub(lastRequestTime) < requestInterval {
            message := tu.Message(
                tu.ID(chatID),
                "Вы можете запросить свежую подборку раз в 12 часов ⏰.\nПопробуйте позже ⌛️",
            )
            _, _ = bot.SendMessage(message)
            return
        }
    }

    // Обновляем время последнего запроса
    lastRequestTimes[chatID] = currentTime

    animes, err := getAnimesFromShikimori()
    if err != nil {
        message := tu.Message(
            tu.ID(chatID),
            "Не удалось получить данные о аниме. Попробуйте позже.",
        )
        _, _ = bot.SendMessage(message)
        return
    }

    sortAnimesByScore(animes)

    for _, anime := range animes {
        photoMessage := tu.Photo(
            tu.ID(chatID),
            tu.FileFromURL("https://shikimori.one" + anime.Image.Original),
        ).WithCaption(formatAnime(anime))

        _, err := bot.SendPhoto(photoMessage)
        if err != nil {
            fmt.Printf("Failed to send message to chat %d: %v\n", chatID, err)
            // Логика повторной отправки может быть добавлена здесь
        }
    }
}
