package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const TELEGRAM_API_BASE_URL string = "https://api.telegram.org/bot"
const BOT_TOKEN_ENV string = "TELEGRAM_BOT_TOKEN"

var telegramApi string = TELEGRAM_API_BASE_URL + os.Getenv(BOT_TOKEN_ENV) + TELEGRAM_API_BASE_URL

// Update is a Telegram object that we receive every time a user interacts with the bot.
type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

// String returns the representation of an Update as a string.
func (u Update) String() string {
	return fmt.Sprintf("(update id: %d, message: %s)", u.UpdateId, u.Message)
}

// Message is a Telegram object that can be found in an update.
type Message struct {
	Text string `json:"text"`
	Chat Chat   `json:"chat"`
}

// String returns representation of a Message as a string.
func (m Message) String() string {
	return fmt.Sprintf("(text: %s, chat: %s)", m.Text, m.Chat)
}

// Chat indicates the conversation to which the Message belongs.
type Chat struct {
	Id int `json:"id"`
}

// String returns representation of a Chat as a string.
func (c Chat) String() string {
	return fmt.Sprintf("(id: %d)", c.Id)
}

// Handler sends a message back to the chat.
func Handler(w http.ResponseWriter, r *http.Request) {
	// Parse incoming request
	update, err := parseTelegramRequest(r)
	if err != nil {
		log.Printf("error parsing update, %s", err.Error())
		return
	}

	telegramResponseBody, err := sendTextToTelegramChat(update.Message.Chat.Id, strings.ToLower(update.Message.Text))
	if err != nil {
		log.Printf("got error %s from telegram, response body is %s", err.Error(), telegramResponseBody)
	} else {
		log.Printf("successfully distributed to chat id %d", update.Message.Chat.Id)
	}
}

// parseTelegramRequest parses incoming update from the Telegram to Update.
func parseTelegramRequest(r *http.Request) (*Update, error) {
	var update Update

	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		log.Printf("could not decode incoming update %s", err.Error())
		return nil, err
	}

	if update.UpdateId == 0 {
		log.Printf("invalid update id, got update id = 0")
		return nil, errors.New("invalid update id of 0 indicates failure to parse incoming update")
	}

	return &update, nil
}

// sendTextToTelegramChat sends a text message to the Telegram chat identified by its chat Id.
func sendTextToTelegramChat(chatId int, finglish string) (string, error) {
	text := getFarsi(finglish)

	log.Printf("Sending %s to chat_id: %d", text, chatId)

	response, err := http.PostForm(telegramApi, url.Values{
		"chat_id": {strconv.Itoa(chatId)},
		"text":    {text},
	})
	if err != nil {
		log.Printf("error when posting text to the chat: %s", err.Error())
		return "", err
	}
	defer response.Body.Close()

	body, errRead := io.ReadAll(response.Body)
	if errRead != nil {
		log.Printf("error in parsing telegram answer %s", errRead.Error())
		return "", err
	}

	log.Printf("Body of Telegram Response: %s", string(body))

	return string(body), nil
}

// getFarsi constructs and returns the appropriate Farsi string for the given Finglish.
func getFarsi(finglish string) string {
	var farsi string

	for i := 0; i < len(finglish); i++ {
		switch finglish[i] {
		case 'a':
			farsi += "ا"
		case 'b':
			farsi += "ب"
		case 'c':
			if peekChar(i, finglish) == "h" {
				farsi += "چ"
				i++
			} else {
				farsi += "س"
			}

		case 'd':
			farsi += "د"
		case 'e':
			farsi += "ع"
		case 'f':
			farsi += "ف"
		case 'g':
			if peekChar(i, finglish) == "h" {
				farsi += "غ"
				i++
			} else {
				farsi += "گ"
			}

		case 'h':
			farsi += "ه"
		case 'i':
			farsi += "ی"
		case 'j':
			farsi += "ج"
		case 'k':
			if peekChar(i, finglish) == "h" {
				farsi += "خ"
				i++
			} else {
				farsi += "ک"
			}

		case 'l':
			farsi += "ل"
		case 'm':
			farsi += "م"
		case 'n':
			farsi += "ن"
		case 'o':
			farsi += "و"
		case 'p':
			farsi += "پ"
		case 'q':
			farsi += "ک"
		case 'r':
			farsi += "ر"
		case 's':
			if peekChar(i, finglish) == "h" {
				farsi += "ش"
				i++
			} else {
				farsi += "س"
			}

		case 't':
			farsi += "ت"
		case 'u':
			farsi += "و"
		case 'v':
			farsi += "و"
		case 'w':
			farsi += "و"
		case 'x':
			farsi += "خ"
		case 'y':
			farsi += "ی"
		case 'z':
			farsi += "ز"
		default:
			farsi += string(finglish[i])
		}
	}

	return farsi
}

// peekChar returns the next char in the given string if exists.
func peekChar(index int, str string) string {
	if index+1 < len(str) {
		return string(str[index+1])
	}
	return ""
}
