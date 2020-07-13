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

func uniqRandom(max int, n int) []int {
	m := make(map[int]bool, n)
	for len(m) <= n {
		x := pocket.GetRandom(max)
		m[x] = true
	}

	result := make([]int, 0, n)

	for k := range m {
		result = append(result, k)
	}

	return result
}
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

//GetRedditPictures return random pictures
func GetRedditPictures(numberOfPictures int) ([]string, error) {

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

	indexes := uniqRandom(ln, numberOfPictures)
	answer := make([]string, 0, numberOfPictures)

	for _, v := range indexes {
		answer = append(answer, result[v])
	}

	log.Println("Indexes: ", indexes)
	return answer, nil
}
