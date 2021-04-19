package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	b, err := gotgbot.NewBot(os.Getenv("BOT_TOKEN"))
	if err != nil {
		panic("failed to create new bot: " + err.Error())
	}

	// Create updater and dispatcher.
	updater := ext.NewUpdater(b, nil)
	dispatcher := updater.Dispatcher

	// Handlers for runnning commands.
	dispatcher.AddHandler(handlers.NewCommand("start", start))
	dispatcher.AddHandler(handlers.NewCommand("get", get))

	err = updater.StartPolling(b, &ext.PollingOpts{Clean: true})
	if err != nil {
		panic("failed to start polling: " + err.Error())
	}

	// log msg telling that bot has started
	fmt.Printf("%s has been started...!\nMade with ❤️ by @Divkix (@DivideProjects).\n", b.User.Username)

	// Idle, to keep updates coming in, and avoid bot stopping.
	updater.Idle()
}

func lorem_ipsum_gen() interface{} {

	api_url := "https://jsonplaceholder.typicode.com/todos/1"
	var raw map[string]interface{}
	response, err := http.Get(api_url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	responseString := string(responseData)

	in := []byte(responseString)
	if err := json.Unmarshal(in, &raw); err != nil {
		panic(err)
	}
	return raw
}

func start(ctx *ext.Context) error {
	// To ensure bot does not reply outside of private chats
	if ctx.EffectiveChat.Type != "private" {
		return nil
	}

	user_name := ctx.EffectiveUser.FirstName

	// Following string is replied to cmd user on /start
	MSG := "*Hi %v*,\n" +
		"I am a *Lorem Ipsum Generator Bot*\n" +
		"Brought to You with ❤️ By @DivideProjects"
	// For Checking either user joined channel or not
	ctx.EffectiveMessage.Reply(ctx.Bot, fmt.Sprintf(MSG, user_name), &gotgbot.SendMessageOpts{ParseMode: "Markdown"})
	return nil
}

func get(ctx *ext.Context) error {
	// To ensure bot does not reply outside of private chats
	if ctx.EffectiveChat.Type != "private" {
		return nil
	}
	quote := lorem_ipsum_gen()
	user := ctx.EffectiveUser
	channel_id, cerror := strconv.Atoi(os.Getenv("AUTH_GROUP_ID"))
	if cerror != nil {
		fmt.Println("Please Provide me a valid Channel/Supergroup ID")
	}
	member, eror := ctx.Bot.GetChatMember(int64(channel_id), user.Id)
	if eror != nil {
		ctx.Bot.SendMessage(ctx.EffectiveChat.Id, "Bot not admin in JoinCheck Channel", nil)
		return nil
	}

	// For Checking either user joined channel or not
	if member.Status == "member" || member.Status == "administrator" || member.Status == "creator" {
		_, err := ctx.EffectiveMessage.Reply(ctx.Bot, fmt.Sprintf("lorem ipsum: %v", quote), &gotgbot.SendMessageOpts{
			ParseMode: "Markdown",
		})
		if err != nil {
			fmt.Println("failed to send: " + err.Error())
		}
	} else {
		// An Error message replied to command user if he's not member of the JoinCheck Channel
		ctx.EffectiveMessage.Reply(ctx.Bot, fmt.Sprintf("*You must join %v to use me.*", os.Getenv("AUTH_GROUP")), &gotgbot.SendMessageOpts{ParseMode: "Markdown"})
	}
	return nil
}
