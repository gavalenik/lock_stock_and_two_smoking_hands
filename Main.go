package main

import (
    "fmt"
    "log"
    "time"
//    "reflect"  //fmt.Println(reflect.TypeOf(var))
    "context"
    "strings"
    "net/url"
    "net/http"
    "io/ioutil"
//    "encoding/json"
    "github.com/go-telegram-bot-api/telegram-bot-api"
    sdk "github.com/TinkoffCreditSystems/invest-openapi-go-sdk"
)

var (
    tele_bot *tgbotapi.BotAPI
    current_time = time.Now().Format("15:04:05")
    token = get_token_from_file()
)

const (
    timeout = 5*time.Second
    general_url = "https://api-invest.tinkoff.ru/openapi/sandbox"
)


type error interface {
    Error() string
}

func tele_initialization() {
    //set proxy
    proxyUrl, err := url.Parse("http://195.154.62.117:5836")
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
		bot, err := tgbotapi.NewBotAPIWithClient(strings.TrimSpace(string(tlg_tkn)), "https://api.telegram.org/bot%s/%s", myClient)
		if err != nil {
				log.Panic(err)
    }
    tele_bot = bot
    tele_bot.Debug = true
}

func telegram(message string) {
    var ChatID = 318841796 //privat chat gavalenik

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
        log.Panic("Critical Error!! The programme execution has been stopped")
		}
		return strings.TrimSpace(string(tin_tkn))
}
/*
func getting_broker_accounts(client *sdk.RestClient) {
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()

    log.Println("Getting all broker's accounts")
    accounts, err := client.Accounts(ctx)
    if err != nil {
      log.Fatalln(err)
    }
    //log.Printf("%+v\n", accounts)

    //fmt.Println (accounts)
    //fmt.Println(reflect.TypeOf(accounts))
}
*/
func getting_all_assets(client *sdk.RestClient) {

    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()

    log.Println("Getting all assets")
    // Метод является совмещеним PositionsPortfolio и CurrenciesPortfolio
    portfolio, err := client.Portfolio(ctx, sdk.DefaultAccount)
    if err != nil {
      log.Fatalln(err)
    }
    log.Printf("%+v\n", portfolio) //type sdk.Portfolio
}


//MAIN
func main() {
    log.Println("Let's get money!")

//    tele_initialization() //telegram bot initialization
//    telegram ("hello")    //sending message via telegram bot
//    client := sdk.NewSandboxRestClient(token) //client for Sandbox
    session := sdk.NewRestClient(token) //client for invest platform !!!
//    getting_broker_accounts(session) //noting helpful, only contract number
    getting_all_assets(session)

		log.Println("The End")
}


/*TO DO
need to count unit in response "portfolio" 'getting_all_assets'
*/
