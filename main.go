package main

import (
    "fmt"
    "os"

    "github.com/mymmrac/telego"
    th "github.com/mymmrac/telego/telegohandler"
    tu "github.com/mymmrac/telego/telegoutil"
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

    bh.Start()
}

func startHandler(bot *telego.Bot, update telego.Update) {
    chatID := tu.ID(update.Message.Chat.ID)

    keyboard := tu.Keyboard(
        tu.KeyboardRow(
            tu.KeyboardButton("Свежая подборка"),
            tu.KeyboardButton("Подписки"),
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
