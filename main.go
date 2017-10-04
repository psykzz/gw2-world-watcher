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
