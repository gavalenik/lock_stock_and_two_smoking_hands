package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "strings"
    "time"
		"os"
//		"github.com/Syfaro/telegram-bot-api"
)

var ()

const (
    timeout = time.Second * 3
    url     = "https://api-invest.tinkoff.ru/openapi/sandbox/sandbox/register"
		ChatID 	= 318841796
)


type error interface {
    Error() string
}

// TO DO. the problem with proxy
/*
func telegram(message string) {

		proxyUrl, err := url.Parse("http://169.57.1.85:80")
		myClient := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}

		// read token from a file
		data, err := ioutil.ReadFile("telebot_token")
		if err != nil {
				if strings.Contains(err.Error(), "no such file or directory") == true {
						fmt.Println("File with telegram bot token isn't exist")
				} else {
						fmt.Println("Error!", err)
				}
		}

		// bot initialization
		//bot, err := tgbotapi.NewBotAPI(strings.TrimSpace(string(data)))
		bot, err = tgbotapi.NewBotAPIWithClient(strings.TrimSpace(string(data)),"https://api.telegram.org/bot%s/%s", myClient)
		if err != nil {
				//log.Panic(err)
				fmt.Println("Error!", err)
		}

		bot.Debug = true
		msg := tgbotapi.NewMessage(ChatID, message)
		bot.Send(msg)
}
*/

func get_token_from_file() string {
		data, err := ioutil.ReadFile("token")
		if err != nil {
				if strings.Contains(err.Error(), "no such file or directory") == true {
						fmt.Println("You have to put file 'token' in catalog 'lock_stock_and_two_smoking_hands'")
						fmt.Println("Details here - https://tinkoffcreditsystems.github.io/invest-openapi/auth/")
				} else {
						fmt.Println("Error!", err)
				}
				fmt.Println("The programme execution has been stopped")
				os.Exit(0)
		}
		return strings.TrimSpace(string(data))
}

func register(token string){

		fmt.Println(token)
		client := &http.Client{
				Timeout: timeout,
		}

		req, err := http.NewRequest("POST", url, nil)
		if err != nil {
				log.Fatalf("Can't create register http request: %s", err)
		}

		req.Header.Add("Authorization", "Bearer "+token)
		resp, err := client.Do(req)
		if err != nil {
				log.Fatalf("Can't send register request: %s", err)
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
				log.Fatalf("Register, bad response code '%s' from '%s'", resp.Status, url)
		}

		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
				log.Fatalf("Can't read register response: %s", err)
		}

		type Register struct {
				TrackingID string `json:"trackingId"`
				Status     string `json:"status"`
		}

		var regResp Register
		err = json.Unmarshal(respBody, &regResp)
		if err != nil {
				log.Fatalf("Can't unmarshal register response: '%s' \nwith error: %s", string(respBody), err)
		}

		if strings.ToUpper(regResp.Status) != "OK" {
				log.Fatalf("Register failed, trackingId: '%s'", regResp.TrackingID)
		}

		fmt.Println("Register succeed")
}

func main() {
//		telegram ("hello")
    register(get_token_from_file())
		fmt.Println("the end")
}
