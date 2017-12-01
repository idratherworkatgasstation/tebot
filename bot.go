package main

import (
	"fmt"
	"log"
	"io/ioutil"
	"flag"
	"net/http"
	"math/rand"
	"os"
	"strings"

	"github.com/Syfaro/telegram-bot-api"
	"google.golang.org/api/youtube/v3"
	"google.golang.org/api/googleapi/transport"
	"github.com/PuerkitoBio/goquery"
)


var studyChoose = []tgbotapi.KeyboardButton{
	tgbotapi.KeyboardButton{Text: "Books"},
	tgbotapi.KeyboardButton{Text: "YouTube"},
}

var langChoose = []tgbotapi.KeyboardButton{
	tgbotapi.KeyboardButton{Text: "English"},
	tgbotapi.KeyboardButton{Text: "Русский"},
	tgbotapi.KeyboardButton{Text: "Deutsch"},
}

//Youtube search settings
var (
	query      = flag.String("query", "Golang", "Search term")
	maxResults = flag.Int64("max-results", 20, "Max YouTube results")
	videos = make(map[string]string)
)

const developerKey = ""  //Google API key
const WebhookURL = ""  //Heroku webhook
const botToken = ""  //Telegram token



func main() {

	flag.Parse()

	client := &http.Client{
		Transport: &transport.APIKey{Key: developerKey},
	}

	service, err := youtube.New(client)
	if err != nil {
		fmt.Println("Problems with Youtube client", err.Error())
	}

	call := service.Search.List("id, snippet").
		Q(*query).
		MaxResults(*maxResults)
	response, err := call.Do()
	if err != nil {
		log.Fatalf("Problems with search API", err.Error())
	}


	port := os.Getenv("PORT")
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		panic(err)
	}
	bot.Debug = true

	_ , err = bot.SetWebhook(tgbotapi.NewWebhook(WebhookURL))
	if err != nil {
		panic(err)
	}
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60


	updates := bot.ListenForWebhook("/")
	go http.ListenAndServe(":"+port, nil)

	for update := range updates {
		fmt.Println("receivedtext :", update.Message.Text)

		switch update.Message.Text {
		case "/start":
			chatID := update.Message.Chat.ID
			text := ("Hey, im GoGopher bot! Il help u with learning Go.")
			msg := tgbotapi.NewMessage(chatID, text)
			bot.Send(msg)

		case "/study":
			chatID := update.Message.Chat.ID
			text := ("Ok, what do u want?")
			msg := tgbotapi.NewMessage(chatID, text)
			msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(studyChoose)
			bot.Send(msg)

			for update := range updates {

				switch update.Message.Text {
				case "Books":
					getBooks()
					chatID := update.Message.Chat.ID
					text := ("Choose ur language:")
					msg := tgbotapi.NewMessage(chatID, text)
					msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(langChoose)
					bot.Send(msg)


					for update := range updates {
						switch update.Message.Text {
						case "English":
							en, err := ioutil.ReadFile("En.txt")
							if err != nil {
								fmt.Println(err.Error())
							}
							//Divide the message in half to send
							i := strings.Index(string(en), "Go Programming Blueprint")
							part1 := (en[:i])
							part2 := (en[i:])
							chatID := update.Message.Chat.ID
							text := (string(part1))
							text1 := (string(part2))
							msg := tgbotapi.NewMessage(chatID, text)
							msg1 := tgbotapi.NewMessage(chatID, text1)
							bot.Send(msg)
							bot.Send(msg1)

						case "Русский":
							ru, err := ioutil.ReadFile("Ru.txt")
							if err != nil {
								fmt.Println(err.Error())
							}
							chatID := update.Message.Chat.ID
							text := (string(ru))
							msg := tgbotapi.NewMessage(chatID, text)
							bot.Send(msg)

						case "Deutsch":
							de, err := ioutil.ReadFile("De.txt")
							if err != nil {
								fmt.Println(err.Error())
							}
							chatID := update.Message.Chat.ID
							text := (string(de))
							msg := tgbotapi.NewMessage(chatID, text)
							bot.Send(msg)
						}
						break
					}
					break

				case "YouTube":
					for _, item := range response.Items {
						videos[item.Snippet.Title] = "https://www.youtube.com/watch?v="+item.Id.VideoId
					}
					chatID := update.Message.Chat.ID
					text := (getVideo(videos))
					msg := tgbotapi.NewMessage(chatID, text)
					bot.Send(msg)
				}
				break
			}
			break

		default:
			chatID := update.Message.Chat.ID
			text := update.Message.Text
			reply := (text + "? Go is better!")
			msg := tgbotapi.NewMessage(chatID, reply)
			bot.Send(msg)
		}
	}
}
//Get random video from results
func getVideo(v map[string]string) string {
	i := int(float32(len(v)) * rand.Float32())
	for _, v := range v {
		if i == 0 {
			return v
		} else {
			i--
		}
	}
	panic("impossible")
}
//Get books list from github
func getBooks(){
	doc, err := goquery.NewDocument("https://github.com/golang/go/wiki/Books")
	if err != nil {
		fmt.Println(err.Error())
	}

	enbooks := doc.Find("#wiki-body > div > ul:nth-child(5)").Contents().Text()
	ef, err := os.Create("En.txt")
	if err != nil {
		fmt.Println("Problems with enbooks", err.Error())
	}
	ef.WriteString(enbooks)

	rubooks := doc.Find("#wiki-body > div > ul:nth-child(23)").Contents().Text()
	rf, err := os.Create("Ru.txt")
	if err != nil {
		fmt.Println("Problems with rubooks", err.Error())
	}
	rf.WriteString(rubooks)

	debooks := doc.Find("#wiki-body > div > ul:nth-child(19)").Contents().Text()
	df, err := os.Create("De.txt")
	if err != nil {
		fmt.Println("Broblems with debooks", err.Error())
	}
	df.WriteString(debooks)
}
