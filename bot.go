package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"log"
	"context"
	// "net/http"
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
	userID string = "hollingsxd"
	 html = `
Logged in
`

)
	
func init() {
	configuration = Configuration{}
	err := gonfig.GetConf("config/conf.json", &configuration)
	if err != nil {
		panic(err)
	}
	os.Setenv("SPOTIFY_ID", configuration.SpotifyClientId)
	os.Setenv("SPOTIFY_SECRET", configuration.SpotifySecret)
	spotifyAuth = spotify.NewAuthenticator(configuration.SpotifyRedirectURI, spotify.ScopeUserReadCurrentlyPlaying, spotify.ScopeUserReadPlaybackState, spotify.ScopeUserModifyPlaybackState)

}

func ConnectSpotify(){

	config := &clientcredentials.Config{
		ClientID:     os.Getenv("SPOTIFY_ID"),
		ClientSecret: os.Getenv("SPOTIFY_SECRET"),
		TokenURL:     spotify.TokenURL,
	}
	token, err := config.Token(context.Background())
	if err != nil {
		log.Fatalf("couldn't get token: %v", err)
	}

	client = spotify.Authenticator{}.NewClient(token)

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


func main() {
	ConnectSpotify()
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
	if (checkValidLink(m.Content)) {

		addToPlaylist(m.Content)
		fmt.Println(s.MessageReactionAdd(m.ChannelID, m.ID, "👍"))

	}

}


// Check if this message contains a spotify song
func checkValidLink(content string) bool {

	if( strings.Contains(content, "open.spotify.com")) {
		return true
	}
	return false
}


// Add song to spotify playlist
func addToPlaylist(content string) string {

	re := regexp.MustCompile("[a-zA-Z0-9]{22}")
	songId := re.FindString(content)
	songName := spotifyFindTrack(songId)
	return songName

}