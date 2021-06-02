package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

func main() {
	// Put Your Bot Token via ENV Vars
	b, err := gotgbot.NewBot(os.Getenv("BOT_TOKEN"), &gotgbot.BotOpts{
		Client:      http.Client{},
		GetTimeout:  gotgbot.DefaultGetTimeout,
		PostTimeout: gotgbot.DefaultPostTimeout,
	})
	if err != nil {
		panic("failed to create new bot: " + err.Error())
	}

	// Create updater and dispatcher.
	updater := ext.NewUpdater(nil)
	dispatcher := updater.Dispatcher

	// Handlers for runnning commands.
	dispatcher.AddHandler(handlers.NewCommand("start", start))
	dispatcher.AddHandler(handlers.NewCommand("get", get))

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
		return nil
	}

	user_name := ctx.EffectiveUser.FirstName

	// Following string is replied to cmd user on /start
	start_msg := "*Hi %v*,\n" +
		"I am a *Lorem Ipsum Generator Bot*\n" +
		"Brought to You with ❤️ By @DivideProjects"
	// For Checking either user joined channel or not
	msg.Reply(bot, fmt.Sprintf(start_msg, user_name), &gotgbot.SendMessageOpts{ParseMode: "Markdown"})
	return nil
}

func get(bot *gotgbot.Bot, ctx *ext.Context) error {
	chat := ctx.EffectiveChat
	user := ctx.EffectiveUser
	msg := ctx.EffectiveMessage

	// To ensure bot does not reply outside of private chats
	if chat.Type != "private" {
		return nil
	}
	quote := "Lorem Ipsum is simply dummy text of the printing and typesetting industry." +
		" Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an " +
		"unknown printer took a galley of type and scrambled it to make a type specimen book." +
		" It has survived not only five centuries, but also the leap into electronic typesetting," +
		" remaining essentially unchanged. It was popularised in the 1960s with the release of " +
		"Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing" +
		"software like Aldus PageMaker including versions of Lorem Ipsum."
	channel_id, cerror := strconv.Atoi(os.Getenv("AUTH_GROUP_ID"))
	if cerror != nil {
		log.Fatalln(cerror)
		log.Fatalln("Please Provide me a valid Channel/Supergroup ID")
	}
	member, eror := bot.GetChatMember(int64(channel_id), user.Id)
	if eror != nil {
		log.Fatalln(eror)
		bot.SendMessage(chat.Id, "Bot not admin in JoinCheck Channel", nil)
		return nil
	}

	// For Checking either user joined channel or not
	if member.Status == "member" || member.Status == "administrator" || member.Status == "creator" {
		_, err := msg.Reply(bot, fmt.Sprintf("*Sample Text:*\n%v", quote), &gotgbot.SendMessageOpts{
			ParseMode: "Markdown",
		})
		if err != nil {
			log.Fatalln(err)
			log.Fatalln("failed to send: " + err.Error())
		}
	} else {
		// An Error message replied to command user if he's not member of the JoinCheck Channel
		url, eror := bot.ExportChatInviteLink(int64(channel_id))
		if eror != nil {
			log.Fatalln(eror)
			bot.SendMessage(chat.Id, "I need invite rights in Channel to get the invite link!", nil)
			return nil
		}
		msg.Reply(bot, fmt.Sprintf("*You must join* [My Bots Channel](%v) *to use me.*", url), &gotgbot.SendMessageOpts{ParseMode: "Markdown", DisableWebPagePreview: true})
	}
	return nil
}
