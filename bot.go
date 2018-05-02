package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"log"
	"context"
	"bytes"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/zmb3/spotify"
	"strings"
	// "reflect"
	"regexp"
	"github.com/bwmarrin/discordgo"
	"github.com/tkanos/gonfig"
	"golang.org/x/oauth2/clientcredentials"

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
}

var (
	spotifyState = "abc123"
	client spotify.Client
	spotifyAuth spotify.Authenticator
	spotifyCh = make(chan *spotify.Client)
	DiscordSecret = configuration.DiscordSecret
	configuration Configuration;
	//username should be pulled from conf
	userID string = "hollingsxd"

)
	
func init() {
	configuration = Configuration{}
	err := gonfig.GetConf("config/conf.json", &configuration)
	if err != nil {
		panic(err)
	}
	os.Setenv("SPOTIFY_ID", configuration.SpotifyClientId)
	os.Setenv("SPOTIFY_SECRET", configuration.SpotifySecret)
	spotifyAuth = spotify.NewAuthenticator(configuration.SpotifyRedirectURI, spotify.ScopeUserReadCurrentlyPlaying, spotify.ScopeUserReadPlaybackState, spotify.ScopeUserModifyPlaybackState, spotify.ScopePlaylistModifyPrivate, spotify.ScopePlaylistModifyPublic)
}


func main() {

	ConnectSpotify()
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

// This runs every time a message sis received
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println(m.Content)
	if (checkValidLink(m.Content)) {
		//https://open.spotify.com/user/hollingsxd/playlist/2z2WuA7x7Op9TvBoYh7y3w?si=wG4dJ_KhTfWbl7nJ94pMLA
		// yes, this playlist ID shouldn't be here. I'll fix it later
		spotifyAddToPlaylist(userID,"2z2WuA7x7Op9TvBoYh7y3w", getSongId(m.Content))
		s.MessageReactionAdd(m.ChannelID, m.ID, "🎵")

	}

}

 func GetTokensScope(tokUrl string, clientId string, secret string) (string,error){
        body := bytes.NewBuffer([]byte("grant_type=client_credentials&client_id="+clientId+"&client_secret="+secret+"&response_type=token"))
        req, err := http.NewRequest("POST",tokUrl,body)
        req.Header.Set("Content-Type","application/x-www-form-urlencoded")  
        client := &http.Client{}
        resp, err := client.Do(req)
        if err != nil {
            return "",err
        }

        defer resp.Body.Close()
        rsBody, err := ioutil.ReadAll(resp.Body)
        type WithScope struct {
            Scope string `json:"scope"`
        }
        var dat WithScope
        err = json.Unmarshal(rsBody,&dat)
        if err != nil {
            return "",err
        }

        return dat.Scope,err
    }

func ConnectSpotify(){
	// fmt.Println(spotify.TokenURL)
// 
	config := &clientcredentials.Config{
		ClientID:     os.Getenv("SPOTIFY_ID"),
		ClientSecret: os.Getenv("SPOTIFY_SECRET"),
		TokenURL:     spotify.TokenURL,
		Scopes: []string{"playlist-modify-private","playlist-modify-public"},
	}
	// fmt.Println(config)
	token, err := config.Token(context.Background())
	if err != nil {
		log.Fatalf("couldn't get token: %v", err)
	}

	// fmt.Println(GetTokensScope(spotify.TokenURL, os.Getenv("SPOTIFY_ID"), os.Getenv("SPOTIFY_SECRET")))
	// fmt.Println(token.AccessToken)

	// I can't figure out how to get scopes working, so I load the oauth token from config. This will expire in 3600 seconds so its WRONG
	token.AccessToken = configuration.SpotifyCode
	// fmt.Println(reflect.TypeOf(token))
	client = spotifyAuth.NewClient(token)



}

func spotifyFindTrack(trackId string) string{
			// search for playlists and albums containing "holiday"
	results, err := client.GetTrack(spotify.ID(trackId))

	if err != nil {
		return "Not Found"
		log.Fatal(err)
	}
	fmt.Println(results.Name)

	return results.Name
} 

func spotifyAddToPlaylist(userId string, playlistId string, trackId string) bool {


	// client := spotify.Authenticator{}.NewClient(token)

	results, err := client.AddTracksToPlaylist(userId,  spotify.ID(playlistId), spotify.ID(trackId)) 

	if err != nil {
		fmt.Println(err)
		return false
	}
	fmt.Println(results)
	return true
}


// Check if this message contains a spotify song
func checkValidLink(content string) bool {

	if( strings.Contains(content, "open.spotify.com")) {
		return true
	}
	return false
}

//https://open.spotify.com/track/76r8BBGixz8suvdjcxMze3?si=IuXCz7JeTvu2xFAodJMa9Q

// Add song to spotify playlist
func getSongId(content string) string {

	re := regexp.MustCompile("[a-zA-Z0-9]{22}")
	songId := re.FindString(content)
	// songName := spotifyFindTrack(songId)
	return songId

}

