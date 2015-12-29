package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Emoji struct {
	Key      string   `json:"-"`
	Char     string   `json:"char"`
	Category string   `json:"category"`
	Keywords []string `json:"keywords"`
}

func main() {
	f, err := os.Open("emojis.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	dec := json.NewDecoder(f)

	var data map[string]json.RawMessage
	if err := dec.Decode(&data); err != nil {
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

	fmt.Println(emojis)
}
