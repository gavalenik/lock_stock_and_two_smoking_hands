package index


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
    var response string
    var split int

    url := "https://apidojo-yahoo-finance-v1.p.rapidapi.com/market/get-quotes?region=US&lang=en&symbols=%255EGSPC%252C%255EDJI%252CBAC%252CKC%253DF%252C002210.KS%252CIWM%252CAMECX"
  	req, _ := http.NewRequest("GET", url, nil)
  	req.Header.Add("x-rapidapi-host", "apidojo-yahoo-finance-v1.p.rapidapi.com")
  	req.Header.Add("x-rapidapi-key", get_token_from_file("yaahoo_key"))

  	res, _ := http.DefaultClient.Do(req)
    defer res.Body.Close()
  	body, _ := ioutil.ReadAll(res.Body)
    response = string(body)
    fmt.Println(res.body[1])

    sp500, split = get_index_value(response)
    us30, split = get_index_value(response[split:len(response)])

    fmt.Println()
    fmt.Print("sp500 index: ")
    fmt.Println(sp500)
    fmt.Print("nasdaq index: ")
    fmt.Println(us30)
    fmt.Println()
}
