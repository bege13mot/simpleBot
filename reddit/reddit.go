package reddit

import (
	"log"
	"runtime"
	"strings"
	"sync"

	"github.com/bege13mot/simpleBot/pocket"

	"github.com/turnage/graw/reddit"
)

func getUniqRandom(max int, n int) []int {
	m := make(map[int]bool, n)

	for len(m) < n {
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

	log.Println("Retrieve URL", topic)

	bot, error := reddit.NewBot(cfg)
	if error != nil {
		return
	}

	harvest, err := bot.Listing("/r/"+topic, "")
	if err != nil {
		log.Fatalln("Reddit Topic Listing ERROR: ", err)
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

// GetRedditPictures return random pictures
func GetRedditPictures(numberOfPictures int, clientID string, clientSecret string, username string, password string, topics []string) []string {

	log.Println("Start getRedditPictures")

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
	if numberOfPictures > len(topics) {
		numberOfPictures = len(topics)
	}

	indexes := getUniqRandom(ln, numberOfPictures)

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

	return result
}
