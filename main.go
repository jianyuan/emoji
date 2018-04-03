package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/urfave/cli"
)

type Emoji struct {
	Emoji          string   `json:"emoji"`
	Description    string   `json:"description"`
	Category       string   `json:"category"`
	Aliases        []string `json:"aliases"`
	Tags           []string `json:"tags"`
	UnicodeVersion string   `json:"unicode_version"`
	IosVersion     string   `json:"ios_version"`
}

func (emj Emoji) TextCode() string {
	if len(emj.Aliases) > 0 {
		return fmt.Sprintf(":%s:", emj.Aliases[0])
	}
	return ""
}

func makeKeywordLookUp(emojis []Emoji) map[string][]*Emoji {
	m := make(map[string][]*Emoji)
	for _, emoji := range emojis {
		for _, alias := range emoji.Aliases {
			m[alias] = append(m[alias], &emoji)
		}
		for _, tag := range emoji.Tags {
			m[tag] = append(m[tag], &emoji)
		}
	}
	return m
}

var keywordLookUp map[string][]*Emoji

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
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
	app.Commands = []cli.Command{
		{
			Name: "generate-db",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "output",
					Value: "emojis.go",
				},
			},
			Action: func(c *cli.Context) error {
				src, err := generateDatabaseFile()
				if err != nil {
					return err
				}

				outputName := strings.ToLower(c.String("output"))
				return ioutil.WriteFile(outputName, src, 0644)
			},
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

		var emoji *Emoji
		if len(emojis) > 1 {
			if c.Bool("random") {
				emoji = emojis[rand.Intn(len(emojis))]
			} else {
				for {
					for i, emj := range emojis {
						if useText {
							fmt.Printf("%d) %s %s\n", i+1, emj.Emoji, emj.TextCode())
						} else {
							fmt.Printf("%d) %s\n", i+1, emj.Emoji)
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
			emojiCode = emoji.Emoji
		}

		if err := clipboard.WriteAll(emojiCode); err != nil {
			panic(err)
		}
		fmt.Printf("copied %s\n", emojiCode)
	}

	app.Run(os.Args)
}
