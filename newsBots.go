package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/bege13mot/simpleBot/pocket"
	"github.com/bege13mot/simpleBot/reddit"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func postMessage(bot *tgbotapi.BotAPI, list []string, message string) {
	for _, chat := range list {
		iChat, _ := strconv.ParseInt(chat, 10, 64)
		msg := tgbotapi.NewMessage(iChat, message)
		_, err := bot.Send(msg)
		if err != nil {
			log.Printf("postMessage Error: %v", err)
		}
	}
}

func forwardMessage(bot *tgbotapi.BotAPI, list []string, fromChatID int64, messageID int) {
	for _, chat := range list {
		iChat, _ := strconv.ParseInt(chat, 10, 64)
		msg := tgbotapi.NewForward(iChat, fromChatID, messageID)
		_, err := bot.Send(msg)
		if err != nil {
			log.Printf("forwardMessage Error: %v", err)
		}
	}
}

func main() {

	botToken := os.Getenv("TelegramBotToken")
	consumerKey := os.Getenv("CONSUMER_KEY")
	accessToken := os.Getenv("ACCESS_TOKEN")
	myID, err := strconv.Atoi(os.Getenv("MyID"))
	if err != nil {
		log.Fatalln("Can't parse MyID")
	}
	list := strings.Split(os.Getenv("Chats"), ",")
	port := ":" + os.Getenv("PORT")

	//Pictures
	numberOfPictures, err := strconv.Atoi(os.Getenv("NumberOfPictures"))
	if err != nil {
		log.Fatalln("Can't parse NumberOfPictures")
	}
	clientID := os.Getenv("RClientID")
	clientSecret := os.Getenv("RClientSecret")
	username := os.Getenv("RUsername")
	password := os.Getenv("RPassword")
	topics := strings.Split(os.Getenv("Topics"), ",")

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatalln("Can't login by Telegram API")
	}
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	// updates := bot.ListenForWebhook("/" + botToken)
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatalln("Get update error: ", err)
	}

	go http.ListenAndServe(port, nil)

	// Receive new updates
	for update := range updates {

		if update.Message.From.ID == myID {
			switch command := update.Message.Command(); command {
			case "":
				log.Printf("Chat ID: %d", update.Message.Chat.ID)
				forwardMessage(bot, list, int64(myID), update.Message.MessageID)

			case "post":
				for _, post := range pocket.RetrieveAndDelete(consumerKey, accessToken) {
					postMessage(bot, list, post)
				}

			case "picture":
				for _, pic := range reddit.GetRedditPictures(numberOfPictures, clientID, clientSecret, username, password, topics) {
					postMessage(bot, list, pic)
				}

			}

		}
	}
}
