package main

import (
	"fmt"
	"gopkg.in/telegram-bot-api.v4"
	"irregular_bot/model"
)

func preapreResponse(connection *model.AppDb, fromChnel tgbotapi.Update) (message string, err error) {
	if fromChnel.Message.IsCommand() {
		err = fmt.Errorf("I dont now any command")
	} else if text := fromChnel.Message.Text; text != "" {
		var verb model.Verb
		switch text {
		default:
			verb, err = findVerb(connection, text)
			if err == nil {
				message = verb.String()
			}
		}
	} else {
		err = fmt.Errorf("You did something strange, did you try send me something amazing?")
	}
	
	return
}

func findVerb(connection *model.AppDb, text string) (model.Verb, error) {
	verb := model.Verb{}
	err := verb.Find(connection, text, "ru_RU")
	if err != nil {
		return verb, err
	}
	return verb, nil
}
