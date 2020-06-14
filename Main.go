package main

import (
    "os"
    "fmt"
    "log"
    "time"
    "strings"
    "net/url"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "github.com/go-telegram-bot-api/telegram-bot-api"
)

var ()

const (
    timeout = time.Second * 3
    tin_url = "https://api-invest.tinkoff.ru/openapi/sandbox/sandbox/register"
)


type error interface {
    Error() string
}

func telegram(message string) {
    var (
      ChatID = 318841796 //privat chat gavalenik
      apiEndpoint = "https://api.telegram.org/bot"
    )

    //set proxy
    proxyUrl, err := url.Parse("http://185.25.207.165:3128")
    myClient := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}

		// read token from a file
		tlg_tkn, err := ioutil.ReadFile("telebot_token")
		if err != nil {
				if strings.Contains(err.Error(), "no such file or directory") == true {
						fmt.Println("File with telegram bot token isn't exist")
				} else {
						fmt.Println("Error!", err)
				}
		}

		// bot initialization
		tele_bot, err := tgbotapi.NewBotAPIWithClient(strings.TrimSpace(string(tlg_tkn)), apiEndpoint+"%s/%s", myClient)
		if err != nil {
				log.Panic(err)
		}
    tele_bot.Debug = true

		tele_bot.Send(tgbotapi.NewMessage(int64(ChatID), message))
}

func get_token_from_file() string {
		tin_tkn, err := ioutil.ReadFile("token")
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
		return strings.TrimSpace(string(tin_tkn))
}

func register(token string){

		fmt.Println(token)
		client := &http.Client{
				Timeout: timeout,
		}

		req, err := http.NewRequest("POST", tin_url, nil)
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
				log.Fatalf("Register, bad response code '%s' from '%s'", resp.Status, tin_url)
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


//MAIN
func main() {
    telegram ("hello") //sending message via telegram bot
    register(get_token_from_file())
		fmt.Println("the end")
}

/*
TO DO
1 - quite message to telegram
2 - don't stop execution cause proxy problem
3 - let forward with general tin apiEndpoint
4 - to find stable proxy
5 - sending message to telegram show up the message in console. need to turn off console

*/
