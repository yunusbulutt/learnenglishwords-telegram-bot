package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	help = `Welcome to Learn English Words Bot.

Using the "/send" command, you can view 1 to 5 words that we randomly selected for you.

Example: "/send 3"`

	unknown_command = `Command is not found. Check "/help".`
	wrong_range     = `The argument must be between 1 and 5.`
)

type Words struct {
	Words []Word `json:"words"`
}

type Word struct {
	English string `json:"english"`
	Turkish string `json:"turkish"`
}

func main() {
	token := os.Getenv("token")
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	config := tgbotapi.NewUpdate(0)
	updates, err := bot.GetUpdatesChan(config)
	if err != nil {
		log.Fatalln(err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s]", update.Message.From.UserName)

		switch update.Message.Command() {
		case "help", "start":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, help)
			bot.Send(msg)
		case "send":
			query := update.Message.CommandArguments()
			count, err := strconv.Atoi(query)

			if err != nil || count <= 0 || count > 5 {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, wrong_range)
				bot.Send(msg)

				continue
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, getWord(count))
			bot.Send(msg)
		default:
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, unknown_command)
			bot.Send(msg)
		}
	}
}

func getWord(count int) string {

	jsonFile, err := os.Open("words.json")

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully Opened words.json")

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var words Words

	json.Unmarshal(byteValue, &words)

	var Word string

	var v [5]int

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < count; i++ {
		v[i] = rand.Intn(len(words.Words))
		Word += "🇬🇧 English: " + words.Words[v[i]].English + "\n"
		Word += "🇹🇷 Turkish: " + words.Words[v[i]].Turkish + "\n"
		Word += "➖➖➖➖➖➖➖➖➖➖" + "\n"
	}

	return Word
}
