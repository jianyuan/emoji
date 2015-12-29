package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/codegangsta/cli"
)

type Emoji struct {
	Key      string   `json:"-"`
	Char     string   `json:"char"`
	Category string   `json:"category"`
	Keywords []string `json:"keywords"`
}

func makeKeywordLookUp(emojis []Emoji) map[string][]Emoji {
	kwdsMap := make(map[string][]Emoji)
	for _, emoji := range emojis {
		for _, kwd := range emoji.Keywords {
			kwdsMap[kwd] = append(kwdsMap[kwd], emoji)
		}
	}
	return kwdsMap
}

var keywordLookUp map[string][]Emoji

func init() {
	rand.Seed(time.Now().UTC().UnixNano())

	var data map[string]json.RawMessage
	if err := json.Unmarshal(MustAsset("emojis.json"), &data); err != nil {
		panic(err)
	}

	// Grab the keys
	var keys []string
	if err := json.Unmarshal(data["keys"], &keys); err != nil {
		panic(err)
	}

	// Parse emojis
	emojis := make([]Emoji, 0, len(keys))
	for _, key := range keys {
		var emoji Emoji
		if err := json.Unmarshal(data[key], &emoji); err != nil {
			panic(err)
		}
		emoji.Key = key
		emojis = append(emojis, emoji)
	}

	keywordLookUp = makeKeywordLookUp(emojis)
}

func main() {
	app := cli.NewApp()
	app.Name = "emoji"
	app.Usage = "find and copy emoji to clipboard"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "random, r",
			Usage: "select random emoji if multiple choices are available",
		},
	}
	app.Action = func(c *cli.Context) {
		if len(c.Args()) < 1 {
			cli.ShowAppHelp(c)
			return
		}

		query := strings.ToLower(strings.TrimSpace(c.Args()[0]))

		emojis, ok := keywordLookUp[query]
		if !ok {
			fmt.Println("emoji not found 😭")
			return
		}

		var emoji Emoji
		if len(emojis) > 1 {
			if c.Bool("random") {
				emoji = emojis[rand.Intn(len(emojis))]
			} else {
				for {
					for i, emj := range emojis {
						fmt.Printf("%d) %s\n", i+1, emj.Char)
					}

					var raw string
					fmt.Printf("choice [1-%d]: ", len(emojis))
					if _, err := fmt.Scanf("%s", &raw); err == nil {
						if choice, err := strconv.Atoi(raw); err == nil {
							if choice >= 1 && choice <= len(emojis) {
								emoji = emojis[choice-1]
								break
							}
						}
					}

					fmt.Println("invalid choice, please try again")
				}
			}
		} else {
			emoji = emojis[0]
		}

		if err := clipboard.WriteAll(emoji.Char); err != nil {
			panic(err)
		}
		fmt.Printf("copied %s\n", emoji.Char)
	}

	app.Run(os.Args)
}
