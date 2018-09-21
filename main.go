package main

import (
	"fmt"
	"irregular_bot/model"
	"log"
	"net/http"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"gopkg.in/telegram-bot-api.v4"
)

//DatabaseLink link on database file,it fills on the init function
var DatabaseLink string

func main() {
	NewDbConnection := &model.AppDb{}
	NewDbConnection.SetDbPath(DatabaseLink)
	err := NewDbConnection.InitDb()
	if err != nil {
		log.Fatal(err.Error())
	}

	defer func() {
		err = NewDbConnection.CloseDb()
		if err != nil {
			log.Fatal(err.Error())
		}
	}()
	BotToken := os.Getenv("BotToken")
	WebhookURL := os.Getenv("WebhookURL")

	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		log.Fatal(err)
	}

	// bot.Debug = true
	log.Printf("Authorized on account %s\n", bot.Self.UserName)

	hookInfo, err := bot.GetWebhookInfo()
	if err != nil {
		log.Fatal(err)
	}

	if hookInfo.URL != "" {
		botResponse, err := bot.RemoveWebhook()
		if err != nil {
			log.Fatal(err)
		}
		log.Println(botResponse.Description)
	}

	// Use for check bot if WebHook is crash
	// u := tgbotapi.NewUpdate(0)
	// u.Timeout = 60
	// updates, err := bot.GetUpdatesChan(u)

	_, err = bot.SetWebhook(tgbotapi.NewWebhook(WebhookURL))
	if err != nil {
		panic(err)
	}

	updates := bot.ListenForWebhook("/")

	go http.ListenAndServe(":8080", nil)
	log.Println("start listen :8080")

	// Get all update from chsnel updates
	for update := range updates {
		if message, err := preapreResponse(NewDbConnection, update); err == nil {
			bot.Send(tgbotapi.NewMessage(
				update.Message.Chat.ID,
				message,
			))
		} else {
			bot.Send(tgbotapi.NewMessage(
				update.Message.Chat.ID,
				fmt.Sprintf("Error : %s", err.Error()),
			))
		}

	}
}
