package main

import (
    "fmt"
    "log"
    "time"
//    "reflect"  //fmt.Println(reflect.TypeOf(var))
    "context"
    "strings"
    "strconv"
    "net/url"
    "net/http"
    "io/ioutil"
    "github.com/go-telegram-bot-api/telegram-bot-api"
    sdk "github.com/TinkoffCreditSystems/invest-openapi-go-sdk"
)

var (
    tele_bot *tgbotapi.BotAPI
    current_time = time.Now().Format("15:04:05")
    token = get_token_from_file()
    current_USD float64 = 0
    current_EUR float64 = 0
    current_RUB float64 = 0
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

func getting_current_balance(client *sdk.RestClient) {
/* type HINT
  type PositionBalance struct {
    FIGI                      string         `json:"figi"`
    Ticker                    string         `json:"ticker"`
    ISIN                      string         `json:"isin"`
    InstrumentType            InstrumentType `json:"instrumentType"`
    Balance                   float64        `json:"balance"`
    Blocked                   float64        `json:"blocked"`
    Lots                      int            `json:"lots"`
    ExpectedYield             MoneyAmount    `json:"expectedYield"`
    AveragePositionPrice      MoneyAmount    `json:"averagePositionPrice"`
    AveragePositionPriceNoNkd MoneyAmount    `json:"averagePositionPriceNoNkd"`
    Name                      string         `json:"name"`
  }
*/
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()
    log.Println("Getting current Balance")

// get current money in valuable pappers
  	positions, err := client.PositionsPortfolio(ctx, sdk.DefaultAccount)
  	if err != nil {
  		log.Fatalln(err)
  	}
  	//log.Printf("%+v\n", positions) //type []sdk.PositionBalance

    for i := 0; i < len(positions); i++ {
        if positions[i].ISIN != "" { //work only with pappers
            switch positions[i].ExpectedYield.Currency {
                case "USD":
                  current_USD = current_USD + (positions[i].Balance*positions[i].AveragePositionPrice.Value)
                  if positions[i].ExpectedYield.Value != 0 {
                      current_USD = current_USD + positions[i].ExpectedYield.Value
                  }
                case "EUR":
                  current_EUR = current_EUR + (positions[i].Balance*positions[i].AveragePositionPrice.Value)
                  if positions[i].ExpectedYield.Value != 0 {
                      current_EUR = current_EUR + positions[i].ExpectedYield.Value
                  }
                case "RUB":
                  current_RUB = current_RUB + (positions[i].Balance*positions[i].AveragePositionPrice.Value)
                  if positions[i].ExpectedYield.Value != 0 {
                      current_RUB = current_RUB + positions[i].ExpectedYield.Value
                  }
                default:
                  log.Panic("Critical Error!! Undefined Currency in your f*cking account!! The programme execution has been stopped")
            }
        }
    }

// get current money in cash
    positionCurrencies, err := client.CurrenciesPortfolio(ctx, sdk.DefaultAccount)
    if err != nil {
        log.Fatalln(err)
    }
    //log.Printf("%+v\n", positionCurrencies)

    for i := 0; i < len(positionCurrencies); i++ {
      switch positionCurrencies[i].Currency {
          case "USD":
              current_USD = current_USD + (positionCurrencies[i].Balance-positionCurrencies[i].Blocked)
          case "EUR":
              current_EUR = current_EUR + (positionCurrencies[i].Balance-positionCurrencies[i].Blocked)
          case "RUB":
              current_RUB = current_RUB + (positionCurrencies[i].Balance-positionCurrencies[i].Blocked)
          default:
              log.Panic("Critical Error!! Undefined Currency in your f*cking account!! The programme execution has been stopped")
      }
    }

    var msg = "Hey Bro! You have:\nUSD: "+strconv.FormatFloat(current_USD, 'f', 2, 64)+"\nEUR: "+strconv.FormatFloat(current_EUR, 'f', 2, 64)+"\nRUB: "+strconv.FormatFloat(current_RUB, 'f', 2, 64)
    fmt.Println(msg)
    telegram (msg)
}


//MAIN
func main() {
    log.Println("Let's get money!")

    tele_initialization() //telegram bot initialization
    //telegram ("hello")    //sending message via telegram bot

    //session := sdk.NewSandboxRestClient(token) //client for Sandbox
    session := sdk.NewRestClient(token) //client for invest platform !!!

    //getting_broker_accounts(session) //noting helpful, only contract number
    getting_current_balance(session)
		log.Println("The End")
}
