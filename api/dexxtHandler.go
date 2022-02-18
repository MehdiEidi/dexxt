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

// Pass token and sensible APIs through environment variables
const telegramApiBaseUrl string = "https://api.telegram.org/bot"
const telegramApiSendMessage string = "/sendMessage"
const telegramTokenEnv string = "TELEGRAM_BOT_TOKEN"

var telegramApi string = telegramApiBaseUrl + os.Getenv(telegramTokenEnv) + telegramApiSendMessage

// Update is a Telegram object that we receive every time an user interacts with the bot.
type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

// Implements the fmt.String interface to get the representation of an Update as a string.
func (u Update) String() string {
	return fmt.Sprintf("(update id: %d, message: %s)", u.UpdateId, u.Message)
}

// Message is a Telegram object that can be found in an update.
// Note that not all Update contains a Message. Update for an Inline Query doesn't.
type Message struct {
	Text     string   `json:"text"`
	Chat     Chat     `json:"chat"`
	Audio    Audio    `json:"audio"`
	Voice    Voice    `json:"voice"`
	Document Document `json:"document"`
}

// Implements the fmt.String interface to get the representation of a Message as a string.
func (m Message) String() string {
	return fmt.Sprintf("(text: %s, chat: %s, audio %s)", m.Text, m.Chat, m.Audio)
}

// Audio message has extra attributes
type Audio struct {
	FileId   string `json:"file_id"`
	Duration int    `json:"duration"`
}

// Implements the fmt.String interface to get the representation of an Audio as a string.
func (a Audio) String() string {
	return fmt.Sprintf("(file id: %s, duration: %d)", a.FileId, a.Duration)
}

// Voice Message can be summarized with similar attribute as an Audio message for our use case.
type Voice Audio

// Document Message refer to a file sent.
type Document struct {
	FileId   string `json:"file_id"`
	FileName string `json:"file_name"`
}

// Implements the fmt.String interface to get the representation of an Document as a string.
func (d Document) String() string {
	return fmt.Sprintf("(file id: %s, file name: %s)", d.FileId, d.FileName)
}

// A Chat indicates the conversation to which the Message belongs.
type Chat struct {
	Id int `json:"id"`
}

// Implements the fmt.String interface to get the representation of a Chat as a string.
func (c Chat) String() string {
	return fmt.Sprintf("(id: %d)", c.Id)
}

// Handler sends a message back to the chat with a punchline starting by the message provided by the user.
func Handler(w http.ResponseWriter, r *http.Request) {
	// Parse incoming request
	var update, err = parseTelegramRequest(r)
	if err != nil {
		log.Printf("error parsing update, %s", err.Error())
		return
	}

	// Send the punchline back to Telegram
	var telegramResponseBody, errTelegram = sendTextToTelegramChat(update.Message.Chat.Id, strings.ToLower(update.Message.Text))
	if errTelegram != nil {
		log.Printf("got error %s from telegram, response body is %s", errTelegram.Error(), telegramResponseBody)
	} else {
		log.Printf("successfully distributed to chat id %d", update.Message.Chat.Id)
	}
}

// parseTelegramRequest handles incoming update from the Telegram web hook
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

// sendTextToTelegramChat sends a text message to the Telegram chat identified by its chat Id
func sendTextToTelegramChat(chatId int, finglish string) (string, error) {
	text := getFarsi(finglish)

	log.Printf("Sending %s to chat_id: %d", text, chatId)
	response, err := http.PostForm(
		telegramApi,
		url.Values{
			"chat_id": {strconv.Itoa(chatId)},
			"text":    {text},
		})

	if err != nil {
		log.Printf("error when posting text to the chat: %s", err.Error())
		return "", err
	}
	defer response.Body.Close()

	var bodyBytes, errRead = io.ReadAll(response.Body)
	if errRead != nil {
		log.Printf("error in parsing telegram answer %s", errRead.Error())
		return "", err
	}
	bodyString := string(bodyBytes)
	log.Printf("Body of Telegram Response: %s", bodyString)

	return bodyString, nil
}

func getFarsi(finglish string) string {
	var farsi string

	for i, c := range finglish {
		switch c {
		case 'a':
			farsi += "ا"
		case 'b':
			farsi += "ب"
		case 'c':
			if peekChar(i, finglish) == "h" {
				farsi += "چ"
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
			} else {
				farsi += "س"
			}
		case 't':
			farsi += "ت"
		case 'u':
			farsi += "ی"
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
			farsi += string(c)
		}
	}

	return farsi
}

func peekChar(index int, str string) string {
	if index+1 < len(str) {
		return string(str[index+1])
	}
	return ""
}
