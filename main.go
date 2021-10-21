package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

func main() {
	// Put Your Bot Token via ENV Vars
	b, err := gotgbot.NewBot(
		os.Getenv("BOT_TOKEN"),
		&gotgbot.BotOpts{
			Client:      http.Client{},
			GetTimeout:  gotgbot.DefaultGetTimeout,
			PostTimeout: gotgbot.DefaultPostTimeout,
		},
	)
	if err != nil {
		panic("failed to create new bot: " + err.Error())
	}

	// Create updater and dispatcher.
	updater := ext.NewUpdater(nil)
	dispatcher := updater.Dispatcher

	// Handlers for runnning commands.
	dispatcher.AddHandler(handlers.NewCommand("start", start))
	dispatcher.AddHandler(handlers.NewCommand("run", run))

	err = updater.StartPolling(b, &ext.PollingOpts{DropPendingUpdates: true})
	if err != nil {
		log.Fatalf("[Polling] Failed to start polling: %v\n", err)
	} else {
		log.Println("[Polling] Started Polling...!")
	}

	// log msg telling that bot has started
	fmt.Printf("%s has been started...!\nMade with ❤️ by @DivideProjects\n", b.User.Username)

	// Idle, to keep updates coming in, and avoid bot stopping.
	updater.Idle()
}

func start(bot *gotgbot.Bot, ctx *ext.Context) error {
	msg := ctx.EffectiveMessage
	// To ensure bot does not reply outside of private chats
	if ctx.EffectiveChat.Type != "private" {
		return ext.EndGroups
	}

	user_name := ctx.EffectiveUser.FirstName

	// Following string is replied to cmd user on /start
	start_msg := "*Hi %v*,\n" +
		"I am a Simple Telegram made using [Go](https://go.dev)*\n" +
		"Brought to You with ❤️ By @DivideProjects"
	// For Checking either user joined channel or not
	msg.Reply(bot, fmt.Sprintf(start_msg, user_name), &gotgbot.SendMessageOpts{ParseMode: "Markdown"})
	return ext.EndGroups
}

func run(bot *gotgbot.Bot, ctx *ext.Context) error {
	chat := ctx.EffectiveChat
	msg := ctx.EffectiveMessage

	// To ensure bot does not reply outside of private chats
	if chat.Type != "private" {
		msg.Reply(bot, "This command only works in private chats!", nil)
		return ext.EndGroups
	}
	text := "This command does nothing, you can build your bot by looking my source code here:\nhttps://github.com/DivideProjects/GoLangTgBot"

	msg.Reply(bot, text, &gotgbot.SendMessageOpts{ParseMode: "Markdown", DisableWebPagePreview: false})

	return nil
}
