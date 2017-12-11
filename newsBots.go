package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mallipeddi/pocket"
	"gopkg.in/telegram-bot-api.v4"
)

//Greeting does greeting string
func Greeting() string {
	first := []string{"Привет", "Доброе утро", "Шалом", "Мир вашему дому"}
	second := []string{"человеки", "мешки с мясом", "котятки", "кожаные ..."}
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	return first[r.Intn(4)] + ", " + second[r.Intn(4)] + "!" + "\n\n"
}

//PostMessages send messages to chats
func PostMessages(bot *tgbotapi.BotAPI, list []string, message string) {
	for _, chat := range list {
		iChat, _ := strconv.ParseInt(chat, 10, 64)
		msg := tgbotapi.NewMessage(iChat, message)
		bot.Send(msg)
	}
}

//RetrieveAndDelete from Pocket
func RetrieveAndDelete(consumerKey string, accessToken string) (message string) {
	client := pocket.NewClientWithAccessToken(consumerKey, accessToken, "")
	req := pocket.NewRetrieveRequest().OnlyFavorited()
	m, err := client.Retrieve(req)
	if err != nil {
		log.Printf("error in retrieve: %s", err)
	}

	text := Greeting()
	if val, ok := m["list"].(map[string]interface{}); ok {

		for k, v := range val {
			url := v.(map[string]interface{})["given_url"]
			title := v.(map[string]interface{})["resolved_title"]
			text += url.(string) + " - " + title.(string) + "\n"
			//Delete item from Pocket
			req := new(pocket.ModifyRequest)
			action := pocket.Action{Kind: pocket.ActionDelete, Params: map[string]string{"item_id": k}}
			req.AddAction(action)
			m, err := client.Modify(req)
			if err != nil {
				log.Printf("error in modify: %s", err)
			}
			log.Printf("modify response: %s\n", m)
		}
	}
	return text
}

func main() {
	botToken := os.Getenv("TelegramBotToken")
	consumerKey := os.Getenv("CONSUMER_KEY")
	accessToken := os.Getenv("ACCESS_TOKEN")
	myID, _ := strconv.Atoi(os.Getenv("MyID"))
	list := strings.Split(os.Getenv("Chats"), ",")

	bot, _ := tgbotapi.NewBotAPI(botToken)
	log.Printf("Authorized on account %s", bot.Self.UserName)

	updates := bot.ListenForWebhook("/" + botToken)
	// updates, err := bot.GetUpdatesChan(u)
	// if err != nil {
	// 	log.Printf("Get update error: ", err)
	// }

	go http.ListenAndServe(":"+os.Getenv("PORT"), nil)

	// Receive new updates
	for update := range updates {
		if update.Message.From.ID == myID && update.Message.Command() == "" {
			log.Printf("Chat ID: %d", update.Message.Chat.ID)
			PostMessages(bot, list, update.Message.Text)

		} else if update.Message.From.ID == myID && update.Message.Command() == "post" {

			message := RetrieveAndDelete(consumerKey, accessToken)
			PostMessages(bot, list, message)
		}
	}
}
