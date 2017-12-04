package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"gopkg.in/telegram-bot-api.v4"
)

var (
	myID       int
	firstChat  int64
	secondChat int64
)

func MainHandler(resp http.ResponseWriter, _ *http.Request) {
	resp.Write([]byte("Hi there! I'm Bot!"))
}

func main() {
	botToken := os.Getenv("TelegramBotToken")
	strChats := os.Getenv("Chats")
	//strSecondChat := os.Getenv("SecondChat")
	strMyID := os.Getenv("MyID")

	bot, err := tgbotapi.NewBotAPI(botToken)
	myID, err := strconv.Atoi(strMyID)
	list := strings.Split(strChats, ",")
	//secondChat, err := strconv.ParseInt(strSecondChat, 10, 64)

	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	//updates, err := bot.GetUpdatesChan(u)
	updates := bot.ListenForWebhook("/" + botToken)

	// if err != nil {
	// 	log.Panic(err)
	// }

	http.HandleFunc("/", MainHandler)
	go http.ListenAndServe(":"+os.Getenv("PORT"), nil)

	// В канал updates будут приходить все новые сообщения.
	for update := range updates {

		if update.Message.From.ID == myID {
			// Создав структуру - можно её отправить обратно боту
			fmt.Println("Chat ID: ", update.Message.Chat.ID)
			for _, chat := range list {
				iChat, _ := strconv.ParseInt(chat, 10, 64)
				msg := tgbotapi.NewMessage(iChat, update.Message.Text)
				bot.Send(msg)
			}
		}
	}
}
