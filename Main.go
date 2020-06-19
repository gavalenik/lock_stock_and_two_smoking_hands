package main

import (
    "os"
    "fmt"
    "log"
    "time"
    "bufio"
//    "reflect"  //fmt.Println(reflect.TypeOf(var))
    "context"
    "strings"
    "strconv"
//    "net/url"
    "net/http"
    "io/ioutil"
    "github.com/go-telegram-bot-api/telegram-bot-api"
    sdk "github.com/TinkoffCreditSystems/invest-openapi-go-sdk"
)

var (
    tele_bot *tgbotapi.BotAPI
    current_time = time.Now() //.Format("15:04:05") current_time.Add(24*time.Hour)
    balance_time = current_time.Add(24*time.Hour)
    index_time = current_time.Add(4*time.Hour)
    token = get_token_from_file("token")
    current_USD, current_EUR, current_RUB float64 = 0,0,0
    sp500, us30 float64
    stocks [20]string
    stock_prices [20]float64
)

const (
    timeout = 5*time.Second
)


type error interface {
    Error() string
}

func tele_initialization() {

// read token from a file
    tlg_tkn := get_token_from_file("telebot_token")
/*
//set proxy
    proxyUrl, err := url.Parse("http://103.235.30.54:8080")
    myClient := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
    bot, err := tgbotapi.NewBotAPIWithClient(tlg_tkn, "https://api.telegram.org/bot%s/%s", myClient)

*/
// bot initialization
    bot, err := tgbotapi.NewBotAPI(tlg_tkn)
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

func get_token_from_file(file_name string) string {
		info, err := ioutil.ReadFile(file_name)
		if err != nil {
				if strings.Contains(err.Error(), "no such file or directory") == true {
						log.Println("For successful result you need 3 files in catalog: token, telebot_token, yahoo_key\n")
				} else {
						log.Println("Error!", err)
				}
        fmt.Println()
        log.Panic("Critical Error!! The programme execution has been stopped\n")
		}
		return strings.TrimSpace(string(info))
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

func get_index_value(info string) (float64, int) {
    var start, finish int

    for i:=strings.Index(info, "regularMarketPrice"); i<len(info); i++ {
        if string(info[i]) == ":" {
            start = i+1
        } else if string(info[i]) == "," {
            finish = i
            break
        }
    }
    value, _ := strconv.ParseFloat(info[start:finish], 8)

    return value, finish
}

func getting_sp500_nasdaq() {

    //get indices once per 6 hours
    if time.Now().After(index_time) == true {
        var response string
        var split int

        url := "https://apidojo-yahoo-finance-v1.p.rapidapi.com/market/get-quotes?region=US&lang=en&symbols=%255EGSPC%252C%255EDJI%252CBAC%252CKC%253DF%252C002210.KS%252CIWM%252CAMECX"
      	req, _ := http.NewRequest("GET", url, nil)
      	req.Header.Add("x-rapidapi-host", "apidojo-yahoo-finance-v1.p.rapidapi.com")
      	req.Header.Add("x-rapidapi-key", get_token_from_file("yahoo_key"))

      	res, _ := http.DefaultClient.Do(req)
        defer res.Body.Close()
      	body, _ := ioutil.ReadAll(res.Body)
        response = string(body)

        sp500, split = get_index_value(response)
        us30, split = get_index_value(response[split:len(response)])

        var msg = "There's indices:\nsp500: "+strconv.FormatFloat(sp500,'f',2,64)+"\nnasdaq: "+strconv.FormatFloat(us30,'f',2,64)+"\n"
        fmt.Println(msg)

        index_time = time.Now().Add(6*time.Hour)
      }

}

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
                  log.Panic("Critical Error!! Undefined Currency in your f*cking account!! The programme execution has been stopped\n")
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

    var msg = "\nHey Bro! You have:\nUSD: "+strconv.FormatFloat(current_USD, 'f', 2, 64)+"\nEUR: "+strconv.FormatFloat(current_EUR, 'f', 2, 64)+"\nRUB: "+strconv.FormatFloat(current_RUB, 'f', 2, 64)+"\n"
    fmt.Println(msg)
    telegram(msg)
}

func balance_difference(client *sdk.RestClient) {

    //get balance difference once per day
    if time.Now().After(balance_time) == true {
        var old_USD float64 = current_USD
        var old_EUR float64 = current_EUR
        var old_RUB float64 = current_RUB

        current_USD, current_EUR, current_RUB = 0,0,0

        getting_current_balance(client)
        differnce_USD := old_USD - current_USD
        differnce_EUR := old_EUR - current_EUR
        differnce_RUB := old_RUB - current_RUB

        var msg = "Hey Bro! There's a difference between latest and current balances\nUSD: "+strconv.FormatFloat(differnce_USD,'f',2,64)+"\nEUR: "+strconv.FormatFloat(differnce_EUR,'f',2,64)+"\nRUB: "+strconv.FormatFloat(differnce_RUB,'f',2,64)+"\n"
        fmt.Println(msg)
        telegram (msg)

        balance_time = time.Now().Add(24*time.Hour)
    }
}
/*
func stock_info_request(client *sdk.RestClient, stocks [20]string) {

    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()

    log.Println("Getting stocks info\n")

    for i:=0; i<len(stocks); i++ {
        if stocks[i] != "" {
            instrument, err := client.InstrumentByFIGI(ctx, stocks[i])
	          if err != nil {
		            log.Fatalln(err)
	          }
	          log.Printf("%+v\n", instrument)
        } else {
            fmt.Printf()
            break
        }
    }
}
*/

func get_stock_orderbook(client *sdk.RestClient, stocks [20]string) {
    var deep int = 10

    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()

    for i:=0; i<len(stocks); i++ {
        if stocks[i] != "" {
            orderbook, err := client.Orderbook(ctx, deep, stocks[i])
            if err != nil {
                log.Fatalln(err)
            }
            //log.Printf("%+v\n", orderbook)
            stock_prices[i] = orderbook.LastPrice
        } else {
            fmt.Println()
            break
        }
    }
}

func get_stock_info(client *sdk.RestClient) {
    var i int = 0

// read stocks from a file
    file, err := os.Open("stock")
    if err != nil {
        fmt.Println(err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        stocks[i] = scanner.Text()
        i++
    }

    //stock_info_request(client, stocks) //nothing helpfull
    get_stock_orderbook(client, stocks)
}


//MAIN
func main() {
    log.Println("Let's get money!\n")

    tele_initialization() //telegram bot initialization  //telegram ("hello")   //sending message via telegram bot
    getting_sp500_nasdaq() //it works but we have limit for requests

    session := sdk.NewRestClient(token) //client for invest platform !!!      //session := sdk.NewSandboxRestClient(token) //client for Sandbox

    //getting_broker_accounts(session) //noting helpful, only contract number
    getting_current_balance(session)
    balance_difference(session)
    get_stock_info(session)

    // finish
    var msgf string = "Stocks info\n"
    for i:=0; i<len(stocks); i++ {
        if stocks[i] != "" {
            msgf = msgf+stocks[i]+": "+strconv.FormatFloat(stock_prices[i], 'f', 2, 64)+"\n"
        } else {
            break
        }
    }
    log.Println(msgf)
    log.Println("The End")
}

// TODO:
// add method to buy stocks
