package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/mallipeddi/pocket"
	"github.com/turnage/graw/reddit"
	"gopkg.in/telegram-bot-api.v4"
)

func getRandom(limit int) int {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	return r.Intn(limit)
}

//Greeting does greeting string
func Greeting() string {
	first := []string{"Привет", "Доброе утро", "Шалом", "Мир вашему дому"}
	second := []string{"человеки", "мешки с мясом", "котятки", "кожаные ..."}

	return first[getRandom(len(first))] + ", " + second[getRandom(len(second))] + "!" + "\n\n"
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
			text += url.(string) + " - " + title.(string) + "\n" + "\n"
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

func retrieveURL(pipe chan<- string, bot reddit.Bot, topic string, wg *sync.WaitGroup) {
	defer wg.Done()

	log.Println("retrieveURL", topic)

	harvest, err := bot.Listing("/r/"+topic, "")
	if err != nil {
		log.Println("Reddit Topic Listing ERROR: ", err)
		return
	}

	runtime.Gosched()

	for _, post := range harvest.Posts[:5] {
		if strings.Contains(post.URL, "jpg") || strings.Contains(post.URL, "gif") {
			pipe <- post.URL
			break
		}
	}
	log.Println("retrieveURL", topic, "END")
}

func getRedditPictures() ([]string, error) {

	log.Println("Start getRedditPictures")

	clientID := os.Getenv("RClientID")
	clientSecret := os.Getenv("RClientSecret")
	username := os.Getenv("RUsername")
	password := os.Getenv("RPassword")
	topics := strings.Split(os.Getenv("Topics"), ",")

	cfg := reddit.BotConfig{
		Agent: "simpleBot",
		App: reddit.App{
			ID:       clientID,
			Secret:   clientSecret,
			Username: username,
			Password: password,
		},
	}

	rBot, error := reddit.NewBot(cfg)
	if error != nil {
		return nil, error
	}

	wg := &sync.WaitGroup{}
	pipe := make(chan string, 10)

	for _, topic := range topics {
		log.Println("reddit Topic: ", topic)
		wg.Add(1)
		go retrieveURL(pipe, rBot, topic, wg)
	}

	wg.Wait()
	close(pipe)

	result := make([]string, len(pipe))
	fmt.Println("!!!!, ", pipe)

	for i := range pipe {
		fmt.Println("ii: ", i)
		result = append(result, i)
	}

	ln := len(topics)

	rnd1, rnd2 := getRandom(ln), getRandom(ln)
	if rnd1 == rnd2 {
		rnd2 = getRandom(ln)
	}

	return []string{result[rnd1], result[rnd2]}, nil
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

		} else if update.Message.From.ID == myID && update.Message.Command() == "picture" {
			pictures, err := getRedditPictures()

			if err == nil {
				for _, pic := range pictures {
					PostMessages(bot, list, pic)
				}
			}
		}

	}
}
