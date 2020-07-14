package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/bege13mot/simpleBot/pocket"
	"github.com/bege13mot/simpleBot/reddit"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

func postMessages(bot *tgbotapi.BotAPI, list []string, message string) {
	for _, chat := range list {
		iChat, _ := strconv.ParseInt(chat, 10, 64)
		msg := tgbotapi.NewMessage(iChat, message)
		_, _ = bot.Send(msg)
	}
}

func main() {

	botToken := os.Getenv("TelegramBotToken")
	consumerKey := os.Getenv("CONSUMER_KEY")
	accessToken := os.Getenv("ACCESS_TOKEN")
	myID, _ := strconv.Atoi(os.Getenv("MyID"))
	list := strings.Split(os.Getenv("Chats"), ",")
	port := ":" + os.Getenv("PORT")

	bot, _ := tgbotapi.NewBotAPI(botToken)
	log.Printf("Authorized on account %s", bot.Self.UserName)

	//Pictures
	numberOfPictures, _ := strconv.Atoi(os.Getenv("NumberOfPictures"))
	clientID := os.Getenv("RClientID")
	clientSecret := os.Getenv("RClientSecret")
	username := os.Getenv("RUsername")
	password := os.Getenv("RPassword")
	topics := strings.Split(os.Getenv("Topics"), ",")

	updates := bot.ListenForWebhook("/" + botToken)
	// updates, err := bot.GetUpdatesChan(u)
	// if err != nil {
	// 	log.Printf("Get update error: ", err)
	// }

	go http.ListenAndServe(port, nil)

	// Receive new updates
	for update := range updates {

		if update.Message.From.ID == myID {
			switch command := update.Message.Command(); command {
			case "":
				log.Printf("Chat ID: %d", update.Message.Chat.ID)
				postMessages(bot, list, update.Message.Text)

			case "post":
				postMessages(bot, list, pocket.RetrieveAndDelete(consumerKey, accessToken))

			case "picture":
				pictures, err := reddit.GetRedditPictures(numberOfPictures, clientID, clientSecret, username, password, topics)
				if err == nil {
					for _, pic := range pictures {
						postMessages(bot, list, pic)
					}
				}

			}

		}
	}
}
