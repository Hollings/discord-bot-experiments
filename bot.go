package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"log"
	"net/http"
	"github.com/zmb3/spotify"
	"strings"
	"github.com/bwmarrin/discordgo"
	"github.com/tkanos/gonfig"
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
	spotifyState string
	spotifyAuth = spotify.NewAuthenticator("http://localhost", spotify.ScopeUserReadPrivate)
	spotifyCh = make(chan *spotify.Client)
	DiscordSecret string
	configuration Configuration;
)
	
func init() {
	configuration = Configuration{}
	err := gonfig.GetConf("config/conf.json", &configuration)
	if err != nil {
		panic(err)
	}
	os.Setenv("SPOTIFY_ID", configuration.SpotifyClientId)
	os.Setenv("SPOTIFY_SECRET", configuration.SpotifySecret)
	DiscordSecret = configuration.DiscordSecret
}

func testConnectSpotify(){

	fmt.Println("this is just a test")

	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})
	go http.ListenAndServe(":8080", nil)

	url := spotifyAuth.AuthURL(spotifyState)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)

	// wait for auth to complete
	client := <-spotifyCh

	// use the client to make calls that require authorization
	user, err := client.CurrentUser()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("You are logged in as:", user.ID)

}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	fmt.Println("completeAuth")
}

func main() {

	// testConnectSpotify();
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

	if (checkValidLink(m.Content)) {
		s.ChannelMessageSend(m.ChannelID, "Spotify link found, adding to playlist")
		addToPlaylist(m.Content)
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
func addToPlaylist(content string) bool {
	return true
}