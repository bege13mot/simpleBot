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
<<<<<<< HEAD
	for len(m) <= n {
=======

	for len(m) < n {
>>>>>>> 6d783d0... Optimize getPicture
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

	//retrive Pictures
	ln := len(topics)
	indexes := uniqRandom(ln, numberOfPictures)

	for _, i := range indexes {
		log.Println("reddit Topic: ", topics[i])
		wg.Add(1)
		go retrieveURL(pipe, cfg, topics[i], wg)
	}

	wg.Wait()
	close(pipe)

	result := make([]string, 0, numberOfPictures)

	for i := range pipe {
		result = append(result, i)
	}

	return result, nil
}
