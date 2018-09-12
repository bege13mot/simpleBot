package reddit

import (
	"log"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/bege13mot/simpleBot/pocket"

	"github.com/turnage/graw/reddit"
)

func retrieveURL(pipe chan<- string, cfg reddit.BotConfig, topic string, wg *sync.WaitGroup) {
	defer wg.Done()

	log.Println("retrieveURL", topic)

	bot, error := reddit.NewBot(cfg)
	if error != nil {
		return
	}

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
}

//Get2RedditPictures return 2 random pictures
func Get2RedditPictures() ([]string, error) {

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

	wg := &sync.WaitGroup{}
	pipe := make(chan string, 10)

	for _, topic := range topics {
		log.Println("reddit Topic: ", topic)
		wg.Add(1)
		go retrieveURL(pipe, cfg, topic, wg)
	}

	wg.Wait()
	close(pipe)

	result := make([]string, 0)

	for i := range pipe {
		result = append(result, i)
	}

	ln := len(topics)

	rnd1, rnd2 := pocket.GetRandom(ln), pocket.GetRandom(ln)
	if rnd1 == rnd2 {
		rnd2 = pocket.GetRandom(ln)
	}

	log.Println("rnd1, rnd2, result", rnd1, rnd2, result)
	return []string{result[rnd1], result[rnd2]}, nil
}
