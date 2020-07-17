package pocket

import (
	"log"
	"math/rand"
	"time"

	"github.com/mallipeddi/pocket"
)

//GetRandom receive random int
func GetRandom(limit int) int {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	return r.Intn(limit)
}

func getGreeting() string {
	first := []string{"Привет", "Доброе утро", "Шалом", "Мир вашему дому, Алоха"}
	second := []string{"человеки", "мешки с мясом", "котятки", "кожаные ...", "людишки"}

	return first[GetRandom(len(first))] + ", " + second[GetRandom(len(second))] + "!" + "\n\n"
}

//RetrieveAndDelete from Pocket
func RetrieveAndDelete(consumerKey string, accessToken string) []string {
	client := pocket.NewClientWithAccessToken(consumerKey, accessToken, "")
	req := pocket.NewRetrieveRequest().OnlyFavorited()
	m, err := client.Retrieve(req)
	if err != nil {
		log.Printf("error in retrieve: %s", err)
	}

	result := make([]string, 0, 10)

	result = append(result, getGreeting())
	if val, ok := m["list"].(map[string]interface{}); ok {

		for k, v := range val {
			url := v.(map[string]interface{})["given_url"]
			title := v.(map[string]interface{})["resolved_title"]
			result = append(result, url.(string)+" - "+title.(string))
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
	return result
}
