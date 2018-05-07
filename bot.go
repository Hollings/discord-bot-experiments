package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	// "log"
	// "context"
	// "bytes"
	"net/http"
	"io/ioutil"
	// "encoding/json"
	// "github.com/zmb3/spotify"
	"strings"
	// "golang.org/x/oauth2/clientcredentials"
	// "github.com/markbates/goth"
	// "github.com/markbates/goth/providers/spotify"
	// "reflect"
	"regexp"
	"github.com/bwmarrin/discordgo"
	"github.com/tkanos/gonfig"
	// "golang.org/x/oauth2"

    "net/url"
    "strconv"
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
	// spotifyState = "abc123"
	// client spotify.Client
	// spotifyAuth spotify.Authenticator
	// spotifyCh = make(chan *spotify.Client)
	DiscordSecret = configuration.DiscordSecret
	configuration Configuration;
	//username should be pulled from conf

)
	
func init() {
	configuration = Configuration{}
	err := gonfig.GetConf("config/conf.json", &configuration)
	if err != nil {
		panic(err)
	}
}


func main() {

	// ConnectSpotify()
	//spotifyAddToPlaylist("hollingsxd","2z2WuA7x7Op9TvBoYh7y3w")

	dg, err := discordgo.New("Bot " + configuration.DiscordBotToken)
	if err != nil {
		fmt.Println("error creating bot: " , err)
		return
	}

	// messageCreate is a callback for MessageCreate event
	dg.AddHandler(messageCreate)

	// Open websocket to discy or die
	err = dg.Open()
	if err != nil {
		fmt.Println("error up opening websocket: " , err)
		return
	}

	// Wait for ctrl-c to end bot
	fmt.Println("bot is werking")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()

}

// This runs every time a message is received
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println(m.Content)
		if(getSongId(m.Content)!="") {
			addSongToPlaylist(configuration.SpotifyPlaylistId,getSongId(m.Content))
			s.MessageReactionAdd(m.ChannelID, m.ID, "🎵")
		}
		if(m.Content=="!playlist"){
			s.ChannelMessageSend(m.ChannelID, "https://open.spotify.com/user/"+configuration.UserID+"/playlist/"+configuration.SpotifyPlaylistId)

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
    // r.Header.Add("Authorization", "auth_token=\"XXXXXXX\"")
    r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
    r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

    resp, _ := client.Do(r)
	body, err := ioutil.ReadAll(resp.Body)
	bodyString := string(body) 
	fmt.Println(bodyString)
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
    // r.Header.Add("Authorization", "auth_token=\"XXXXXXX\"")
    r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
    r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
    resp, _ := client.Do(r)
	body, err := ioutil.ReadAll(resp.Body)
	bodyString := string(body) 
	fmt.Println(bodyString)
	fmt.Println(err)  
}

func addSongToPlaylist(playlistId string, songId string) {

	// This is bad
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
    fmt.Println(urlStr)
    client := &http.Client{}
    r, _ := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode())) // URL-encoded payload
    r.Header.Add("Authorization", "Bearer "+configuration.SpotifyAccessToken)
    resp, _ := client.Do(r)
	body, err := ioutil.ReadAll(resp.Body)
	bodyString := string(body) 
	fmt.Println(bodyString)
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
	// songName := spotifyFindTrack(songId)
	fmt.Println(songId)
	return songId

}

