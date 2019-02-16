package main

import (
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "net/http"
    // "io/ioutil"
    "strings"
    // "regexp"
    "github.com/bwmarrin/discordgo"
    "github.com/tkanos/gonfig"
    // "net/url"
    "strconv"
    "encoding/json"
    // "math/rand"
    "log"
    // "bufio"
    // "database/sql"
    // "time"
    // _ "github.com/mattn/go-sqlite3"
)

// I need to clean up this config code
type Configuration struct {
    DiscordClientId string
    DiscordSecret string
    DiscordBotToken string
}



type AutoGenerated struct {
    Quote struct {
        Symbol                string      `json:"symbol"`
        CompanyName           string      `json:"companyName"`
        PrimaryExchange       string      `json:"primaryExchange"`
        Sector                string      `json:"sector"`
        CalculationPrice      string      `json:"calculationPrice"`
        Open                  float64     `json:"open"`
        OpenTime              int64       `json:"openTime"`
        Close                 float64     `json:"close"`
        CloseTime             int64       `json:"closeTime"`
        High                  float64     `json:"high"`
        Low                   float64     `json:"low"`
        LatestPrice           float64     `json:"latestPrice"`
        LatestSource          string      `json:"latestSource"`
        LatestTime            string      `json:"latestTime"`
        LatestUpdate          int64       `json:"latestUpdate"`
        LatestVolume          int         `json:"latestVolume"`
        IexRealtimePrice      interface{} `json:"iexRealtimePrice"`
        IexRealtimeSize       interface{} `json:"iexRealtimeSize"`
        IexLastUpdated        interface{} `json:"iexLastUpdated"`
        DelayedPrice          float64     `json:"delayedPrice"`
        DelayedPriceTime      int64       `json:"delayedPriceTime"`
        ExtendedPrice         int         `json:"extendedPrice"`
        ExtendedChange        float64     `json:"extendedChange"`
        ExtendedChangePercent float64     `json:"extendedChangePercent"`
        ExtendedPriceTime     int64       `json:"extendedPriceTime"`
        PreviousClose         float64     `json:"previousClose"`
        Change                float64     `json:"change"`
        ChangePercent         float64     `json:"changePercent"`
        IexMarketPercent      interface{} `json:"iexMarketPercent"`
        IexVolume             interface{} `json:"iexVolume"`
        AvgTotalVolume        int         `json:"avgTotalVolume"`
        IexBidPrice           interface{} `json:"iexBidPrice"`
        IexBidSize            interface{} `json:"iexBidSize"`
        IexAskPrice           interface{} `json:"iexAskPrice"`
        IexAskSize            interface{} `json:"iexAskSize"`
        MarketCap             int64       `json:"marketCap"`
        PeRatio               float64     `json:"peRatio"`
        Week52High            float64     `json:"week52High"`
        Week52Low             float64     `json:"week52Low"`
        YtdChange             float64     `json:"ytdChange"`
    } `json:"quote"`
}

var (
    DiscordSecret = configuration.DiscordSecret
    Stocks map[string]map[string]int
    configuration Configuration;
    banStrings = []string{"get rekt noob scrub","ur mums gay lol","ligma","ur dad gay","ur moms fat"} 
)
    
func init() {
    configuration = Configuration{}
    err := gonfig.GetConf("config/conf.json", &configuration)
    if err != nil {
        panic(err)
    }
}

func connectSqlite(){
        // db, err := sql.Open("sqlite3", "./foo.db")
  //       checkErr(err)
}

func main() {


    dg, err := discordgo.New("Bot " + configuration.DiscordBotToken)
    if err != nil {
        fmt.Println("error creating bot: " , err)
        return
    }

    // messageCreate is a callback for MessageCreate event
    dg.AddHandler(messageCreate)
    dg.AddHandler(ready)

    // Open websocket to discy or die
    err = dg.Open()
    if err != nil {
        fmt.Println("error up opening websocket: " , err)
        return
    }
    // refreshToken()
    // Wait for ctrl-c to end bot
    fmt.Println("bot is werking")
    sc := make(chan os.Signal, 1)
    signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
    <-sc

    dg.Close()

}
// This function will be called (due to AddHandler above) when the bot receives
// the "ready" event from Discord.
func ready(s *discordgo.Session, event *discordgo.Ready) {
    // Set the playing status.
    s.UpdateStatus(0, "hey")
}

func getStocks(username string)map[string]int{

    data := map[string]map[string]int{"Hollings":{"TSLA":3,"SPY":2,"AIEQ":2,"NIO":3}}
    return data[username]
}
// This runs every time a message is received
func messageCreate(session *discordgo.Session, m *discordgo.MessageCreate) {
        
        fmt.Println(m.Content)
        
        mentionsMe := false
        for _, us := range m.Mentions {
            if us.ID == session.State.User.ID {
                mentionsMe = true
                break
            }
        }
        
        if(strings.Contains(m.Content,"ligma")&&(m.Content!="What's ligma?")){
            session.ChannelMessageSend(m.ChannelID, "What's ligma?")
        }
        if(strings.Contains(m.Content,"updog")&&(m.Content!="What's updog?")){
            session.ChannelMessageSend(m.ChannelID, "What's updog?")
        }

        if(mentionsMe){
            stock := ""
            command := ""
            // Set a minimum length for the random response
            split :=  strings.Split(m.Content, " ")
            log.Println(len(split))
            if(len(split) == 2){
                    stock = split[1]
            }else if(len(split) == 3){
                stock = split[2]
                command = split[1]
            }else{
                command = "status"
            }
            // log.Println()

            if(command == "add") {
                session.ChannelMessageSend(m.ChannelID, "Added " + stock)
                return
            }

            if(command == "status") {
                data := getStocks("Hollings")
                // log.Println(data["F"])

                response := ""
                totalPrice := float64(0)

                for k, v := range data { 
                    _, price := getStockPrice(k)
                    totalSinglePrice := price * float64(v)
                    totalPrice += totalSinglePrice
                    // s := strconv.Itoa(v)
                    if len(k) == 3{
                        k += " "
                    }
                    response += strconv.Itoa(v) + "\t" + k + ":\t$" + FloatToString(totalSinglePrice) + "\n"
                }

                response = "```Total Portfolio for @Hollings: $" + FloatToString(totalPrice) +"\n ----- \n"+ response + "```"
                session.ChannelMessageSend(m.ChannelID, response)
                return
            }

            if(stock != ""){

                companyName, stockPrice := getStockPrice(stock)
                previousClose := FloatToString(stockPrice)
                finalString := companyName +": " + previousClose
    
                session.ChannelMessageSend(m.ChannelID, finalString)
                return
            }
            
        }


        if(m.Content=="test"){
            session.MessageReactionAdd(m.ChannelID, m.ID, "🇹")
            session.MessageReactionAdd(m.ChannelID, m.ID, "🇪")
            session.MessageReactionAdd(m.ChannelID, m.ID, "🇸")
            session.MessageReactionAdd(m.ChannelID, m.ID, "🔨")
        }
    }

func getStockPrice(stock string) (string,float64) {

    // response, err := http.Get("https://api.iextrading.com/1.0/stock/tsla/batch?types=quote&range=1m&last=10")
    
    // QueryEscape escapes the phone string so
    // it can be safely placed inside a URL query

    url := fmt.Sprintf("https://api.iextrading.com/1.0/stock/%s/batch?types=quote&range=1m&last=10", stock)

    // Build the request
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        log.Fatal("NewRequest: ", err)
        return "",0.0
    }

    // For control over HTTP client headers,
    // redirect policy, and other settings,
    // create a Client
    // A Client is an HTTP client
    client := &http.Client{}

    // Send the request via a client
    // Do sends an HTTP request and
    // returns an HTTP response
    resp, err := client.Do(req)
    if err != nil {
        log.Fatal("Do: ", err)
        return "",0.0
    }

    // Callers should close resp.Body
    // when done reading from it
    // Defer the closing of the body
    defer resp.Body.Close()

    // Fill the record with the data from the JSON
    var quote AutoGenerated

    // Use json.Decode for reading streams of JSON data
    if err := json.NewDecoder(resp.Body).Decode(&quote); err != nil {
        log.Println(err)
    }
   
    return quote.Quote.CompanyName, quote.Quote.PreviousClose
}

func FloatToString(input_num float64) string {
    // to convert a float number to a string
    return strconv.FormatFloat(input_num, 'f', 2, 64)
}
