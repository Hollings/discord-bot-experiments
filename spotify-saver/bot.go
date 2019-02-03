package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"net/http"
	"io/ioutil"
	"strings"
	"regexp"
	"github.com/bwmarrin/discordgo"
	"github.com/tkanos/gonfig"
    "net/url"
    "strconv"
    "encoding/json"

)

// I need to clean up this config code
type Configuration struct {
	SpotifyClientId string
	SpotifyCode string
	SpotifySecret string
	SpotifyRedirectURI string
	DiscordClientId string
	DiscordSecret string
	DiscordBotToken string
	SpotifyAccessToken string
	SpotifyRefreshToken string
	SpotifyTokenExpires string
	SpotifyPlaylistId string
	UserID string
}

var (

	DiscordSecret = configuration.DiscordSecret
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
type refreshTokenResp struct {
    Access_token   string      `json:"access_token"`
    Token_type   string      `json:"token_type"`
    Expires_in   int      `json:"expires_in"`
    Scope   string      `json:"scope"`

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
	refreshToken()
	// Wait for ctrl-c to end bot
	fmt.Println("bot is working")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()

}

// This function will be called (due to AddHandler above) when the bot receives
// the "ready" event from Discord.
func ready(s *discordgo.Session, event *discordgo.Ready) {
	// Set the playing status.
	s.UpdateStatus(0, "🎵")
}


// This runs every time a message is received
func messageCreate(spaghetti *discordgo.Session, m *discordgo.MessageCreate) {

		// If the message has a spotify link, add it to a playlist
		if(getSongId(m.Content)!="") {
			addSongToPlaylist(configuration.SpotifyPlaylistId,getSongId(m.Content))
			spaghetti.MessageReactionAdd(m.ChannelID, m.ID, "🎵")
		}

		// Returns the playlist that the bot has been adding to
		if(m.Content=="!playlist"){
			spaghetti.ChannelMessageSend(m.ChannelID, "https://open.spotify.com/user/"+configuration.UserID+"/playlist/"+configuration.SpotifyPlaylistId)
		}
}


//The Oauth token/refresh flow still needs to be automated. 
func getOuthTokens(){
  	apiUrl := "https://accounts.spotify.com"
    resource := "/api/token"
    data := url.Values{}
    data.Set("grant_type", "authorization_code")
    data.Add("code", configuration.SpotifyCode)
    data.Add("redirect_uri", configuration.SpotifyRedirectURI)
    data.Add("state", "abc123") //  I know...
    data.Add("client_id", configuration.SpotifyClientId)
    data.Add("client_secret", configuration.SpotifySecret)

    u, _ := url.ParseRequestURI(apiUrl)
    u.Path = resource
    urlStr := u.String() // 'https://api.com/user/'

    client := &http.Client{}
    r, _ := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode())) // URL-encoded payload
    r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
    r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

    resp, _ := client.Do(r)
	_, err := ioutil.ReadAll(resp.Body)
	// bodyString := string(body) 
	// fmt.Println(bodyString)
	fmt.Println(err)  
}

func refreshToken() {
	apiUrl := "https://accounts.spotify.com"
    resource := "/api/token"
    data := url.Values{}
    data.Set("grant_type", "refresh_token")
    data.Add("refresh_token", configuration.SpotifyRefreshToken)
    data.Add("redirect_uri", configuration.SpotifyRedirectURI)
    data.Add("state", "abc123")
    data.Add("client_id", configuration.SpotifyClientId)
    data.Add("client_secret", configuration.SpotifySecret)
    u, _ := url.ParseRequestURI(apiUrl)
    u.Path = resource
    urlStr := u.String() // 'https://api.com/user/'
    client := &http.Client{}
    r, _ := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode())) // URL-encoded payload
    r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
    r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
    resp, _ := client.Do(r)
	body, err := ioutil.ReadAll(resp.Body)
	bodyString := string(body) 
	fmt.Println(bodyString)
	fmt.Println(err)  
	fmt.Println("New AccessToken is" + bodyString)
	fmt.Println("---")
    var tok refreshTokenResp

	refreshErr := json.Unmarshal(body, &tok)
    if refreshErr != nil {
        fmt.Println(refreshErr)
        return
    }
    // fmt.Println(tok.Access_token)
    configuration.SpotifyAccessToken = tok.Access_token
    // fmt.Println("---")

    // res := refreshTokenResp{}
    // json.Unmarshal([]byte(bodyString), &res)
    // fmt.Println(res)

}


func addSongToPlaylist(playlistId string, songId string) {

	// This is ugly
	apiUrl := "https://api.spotify.com/"
    resource := "/v1/users/"+configuration.UserID+"/playlists/"+playlistId+"/tracks"
    fmt.Println(resource)
    data := url.Values{}
 	data.Set("position", "0")
    data.Add("uris", "spotify:track:"+songId)
    u, _ := url.ParseRequestURI(apiUrl)
    u.Path = resource
    urlStr := u.String() // 'https://api.com/user/'
    urlStr += "?position=0&uris=spotify:track:"+songId
    // fmt.Println(urlStr)
    client := &http.Client{}
    r, _ := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode())) // URL-encoded payload
    r.Header.Add("Authorization", "Bearer "+configuration.SpotifyAccessToken)
    resp, _ := client.Do(r)
	_, err := ioutil.ReadAll(resp.Body)
	// bodyString := string(body) 
	// fmt.Println(bodyString)
	fmt.Println(err)  


}

// Add song to spotify playlist
func getSongId(content string) string {

	// Stop interpreting playlist links as song links 
	if(strings.Contains(content, "playlist")){
		return ""
	}
	re := regexp.MustCompile("[a-zA-Z0-9]{22}")
	songId := re.FindString(content)
	fmt.Println(songId)
	return songId

}

