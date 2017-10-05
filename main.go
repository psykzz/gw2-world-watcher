package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/mitsuse/pushbullet-go"
	"github.com/mitsuse/pushbullet-go/requests"
)

const (
	Low      = "Low"
	Medium   = "Medium"
	High     = "High"
	VeryHigh = "VeryHigh"
	Full     = "Full"
)

type World struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Population string `json:"population"`
}

func getWorldsJson(url string) (worlds []World, err error) {
	var resp *http.Response
	if resp, err = http.Get(url); err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data []byte
	if data, err = ioutil.ReadAll(resp.Body); err != nil {
		return nil, err
	}

	if err = json.Unmarshal(data, &worlds); err != nil {
		return nil, err
	}
	return
}
type Payload struct {
	Notification struct {
		Title       string `json:"title"`
		Body        string `json:"body"`
		Icon        string `json:"icon"`
		ClickAction string `json:"click_action"`
	} `json:"notification"`
	To string `json:"to"`
}
func sendNotification() {
	
	// TODO: Fill out the rest of this payload.
	data := Payload{
		notification: Notification{
			Title: "GW2 Watcher" 
			Body: "Server change status!" 
		}
		To: ""
	}
	
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Unable to marshall notification JSON");
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "https://fcm.googleapis.com/fcm/send", body)
	if err != nil {
		fmt.Println("Unable to create message for FCM")
	}
	req.Header.Set("Authorization", fmt.Sprintf("key=%s", os.GetEnv('FCM_SERVER_KEY')))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Unable to send message to FCM")
	}
	defer resp.Body.Close()

}

func main() {
	token := os.Getenv("PUSHBULLET_TOKEN")
	pb := pushbullet.New(token)

	status := "unknown"

	t := time.NewTicker(time.Millisecond * 5000)
	for {

		worlds, err := getWorldsJson("https://api.guildwars2.com/v2/worlds?ids=all")
		if err != nil {
			fmt.Println("Error getting GW2 Worlds JSON.")
			panic(err.Error())
		}

		var world = worlds[30] // Far shiverpeaks
		if world.Population != status {

			var message = fmt.Sprintf("GW2 Watcher: FSP Pop status changed to %s from %s.", world.Population, status)
			fmt.Println(message)

			n := requests.NewNote()
			n.Title = "GW2 Watcher"
			n.Body = message
			if _, err := pb.PostPushesNote(n); err != nil {
				fmt.Println("Error pushing notification.")
				fmt.Fprintf(os.Stderr, "error: %s\n", err)
			}

			status = world.Population
		}

		<-t.C
	}
}
