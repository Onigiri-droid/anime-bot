package main

import (
	"fmt"
	
    "github.com/mymmrac/telego"
    tu "github.com/mymmrac/telego/telegoutil"
)

func freshAnimeHandler(bot *telego.Bot, update telego.Update) {
    chatID := tu.ID(update.Message.Chat.ID)

    animes, err := getAnimesFromShikimori()
    if err != nil {
        message := tu.Message(
            chatID,
            "Не удалось получить данные о аниме. Попробуйте позже.",
        )
        _, _ = bot.SendMessage(message)
        return
    }

    sortAnimesByScore(animes)

    for i, anime := range animes {
        if i >= 50 {
            break
        }

        photoMessage := tu.Photo(
            chatID,
            tu.FileFromURL("https://shikimori.one" + anime.Image.Original),
        ).WithCaption(formatAnime(anime)).
            WithReplyMarkup(tu.InlineKeyboard(
                tu.InlineKeyboardRow(
                    tu.InlineKeyboardButton("Подписаться").WithCallbackData(fmt.Sprintf("subscribe_%d", anime.Id)),
                ),
            ))

        _, _ = bot.SendPhoto(photoMessage)
    }
}
