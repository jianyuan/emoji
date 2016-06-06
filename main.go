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

func (emj Emoji) TextCode() string {
	return fmt.Sprintf(":%s:", emj.Key)
}

func makeKeywordLookUp(emojis []Emoji) map[string][]Emoji {
	kwdsMap := make(map[string][]Emoji)
	for _, emoji := range emojis {
		kwdsMap[emoji.Key] = append(kwdsMap[emoji.Key], emoji)
		for _, kwd := range emoji.Keywords {
			kwdsMap[kwd] = append(kwdsMap[kwd], emoji)
		}
	}
	return kwdsMap
}

var keywordLookUp map[string][]Emoji

func init() {
	rand.Seed(time.Now().UTC().UnixNano())

	var data map[string]Emoji
	if err := json.Unmarshal(MustAsset("emojis.json"), &data); err != nil {
		panic(err)
	}

	// Parse emojis
	emojis := make([]Emoji, 0, len(data))
	for key, emoji := range data {
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
		cli.BoolFlag{
			Name:  "text, t",
			Usage: "return emoji as ascii code instead of unicode",
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

		useText := c.Bool("text")

		var emoji Emoji
		if len(emojis) > 1 {
			if c.Bool("random") {
				emoji = emojis[rand.Intn(len(emojis))]
			} else {
				for {
					for i, emj := range emojis {
						if useText {
							fmt.Printf("%d) %s %s\n", i+1, emj.Char, emj.TextCode())
						} else {
							fmt.Printf("%d) %s\n", i+1, emj.Char)
						}
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

		var emojiCode string
		if useText {
			emojiCode = emoji.TextCode()
		} else {
			emojiCode = emoji.Char
		}

		if err := clipboard.WriteAll(emojiCode); err != nil {
			panic(err)
		}
		fmt.Printf("copied %s\n", emojiCode)
	}

	app.Run(os.Args)
}
