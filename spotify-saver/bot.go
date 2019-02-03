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
    "math/rand"
    "time"
    "log"
    "bufio"
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
	s.UpdateStatus(0, "Hey")
}

func pickRandomLine(fileName string, length int) string{
	// Pulls a random line from a text file
	processedString := "";
	file, err := os.Open(fileName)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)

    var lines []string;

    for scanner.Scan() {
    	lines = append(lines, scanner.Text())
    }

    s := rand.NewSource(time.Now().Unix())
	r := rand.New(s) // initialize local pseudorandom generator 


	for len(processedString) < length {
		message := lines[r.Intn(len(lines))]
		fmt.Println(message);

		reg, err := regexp.Compile("[^a-zA-Z0-9 !@#=$:&\\/*,.?]+")
	    if err != nil {
	        log.Fatal(err)
	    }
	    processedString = processedString + " " + reg.ReplaceAllString(message, "")
   	}
    return processedString // error return
}

func truncateString(str string, num int) string {
	bnoden := str
	if len(str) > num {
		if num > 3 {
			num -= 3
		}
		bnoden = str[0:num] + "..."
	}
	return bnoden
}
func generateSweeper(x int) string{

	// Keep it below 15 because Discord's character limit
	if x>14 {
		x=14
	}
	if x<0 {
		x=4
	}

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s) // initialize local pseudorandom generator 
	
	a := make([][]uint8, x)
	for i := range a {
	    a[i] = make([]uint8, x)
	}

	// Place the bombs
	for y := 0; y < x; y++ {
		for z := 0; z < x; z++ {
			if r.Intn(8) == 7 {
				a[y][z] = 9
			}
		}
		
	}


	// Add up the numbers. This is horrible but I just wanted to quickly hack together a minesweeper board
	for y := 0; y < x; y++ {
		for z := 0; z < x; z++ {
			if a[y][z] == 9{
				continue
			}
			if y<x-1 {
				if a[y+1][z] == 9{
				a[y][z]++
				}
				if z>0 && a[y+1][z-1] == 9{
					a[y][z]++
				
				}
				if z<x-1 && a[y+1][z+1] == 9{
					a[y][z]++
				}

			}
			if y>0 {
				if a[y-1][z] == 9{
				a[y][z]++
				}
				if  z>0 && a[y-1][z-1] == 9{
					a[y][z]++
				}
				if  z<x-1 && a[y-1][z+1] == 9{
					a[y][z]++
				}
			}
			if z<x-1 && a[y][z+1] == 9 {
				a[y][z]++
			}
			if z>0 && a[y][z-1] == 9 {
				a[y][z]++
			}
		}
		
	}

	// Format the array into a string
	returnString := " ";
	for y := 0; y < x; y++ {
		for  z := 0; z < x; z++ {
			if z!=0 {
				// returnString = returnString + " "
			}
			returnString = returnString + "|| `" + strconv.Itoa(int(a[y][z])) + "` ||"
		}
		returnString = returnString + "\n"
	}
	returnString = returnString + ""
	returnString = strings.Replace(returnString, "9", "X", -1)
	return returnString
}

// This runs every time a message is received
func messageCreate(spaghetti *discordgo.Session, m *discordgo.MessageCreate) {
		fmt.Println(m.Content)
		mentionsMe := false
		for _, us := range m.Mentions {
			if us.ID == spaghetti.State.User.ID {
				mentionsMe = true
				break
			}
		}

		// If the message has a spotify link, add it to a playlist
		if(getSongId(m.Content)!="") {
			addSongToPlaylist(configuration.SpotifyPlaylistId,getSongId(m.Content))
			spaghetti.MessageReactionAdd(m.ChannelID, m.ID, "🎵")
		}

		// Returns the playlist that the bot has been adding to
		if(m.Content=="!playlist"){
			spaghetti.ChannelMessageSend(m.ChannelID, "https://open.spotify.com/user/"+configuration.UserID+"/playlist/"+configuration.SpotifyPlaylistId)
		}

		// Stupid joke lol
		if(strings.Contains(m.Content,"ligma")&&(m.Content!="What's ligma?")){
			spaghetti.ChannelMessageSend(m.ChannelID, "What's ligma?")
		}

		// Create a minesweeper board with spoiler tags
		if(strings.HasPrefix(m.Content, "!sweep") ){
			x, _ := strconv.Atoi(strings.Split(m.Content, " ")[1])
			if (x == 1){
				spaghetti.ChannelMessageSend(m.ChannelID, "||You Lost||")
			}else{
				board := generateSweeper(x);
				spaghetti.ChannelMessageSend(m.ChannelID, board)
			}
			
		}

		// Make the bot say something
		if(strings.HasPrefix(m.Content, "!say") ){
			// say := strings.Fields(m.Content)[1]
			say := strings.Split(m.Content, " ")[1:]
			joinedSay := strings.Join(say," ")
			spaghetti.ChannelMessageSend(m.ChannelID, joinedSay)
			spaghetti.ChannelMessageDelete(m.ChannelID, m.ID)
		}

		// Add some reacts to a message to test if bot is active
		if(m.Content=="test"){
			spaghetti.MessageReactionAdd(m.ChannelID, m.ID, "🇹")
			spaghetti.MessageReactionAdd(m.ChannelID, m.ID, "🇪")
			spaghetti.MessageReactionAdd(m.ChannelID, m.ID, "🇸")
			spaghetti.MessageReactionAdd(m.ChannelID, m.ID, "🔨")
		}


		// Pick a random line from text file on mention
		if(mentionsMe){
			spaghetti.ChannelTyping(m.ChannelID)

			// Set a minimum length for the random response
			length := 30
			split :=  strings.Split(m.Content, " ")
			if(len(split) > 1){
				newLen, err := strconv.ParseInt(split[1], 10, 64)
				if err == nil {
					fmt.Println(length)
				}else{
					newLen = 30
				}
				length = int(newLen)
			}else{
				length = 30
			}
			
			finalMessage := pickRandomLine("lines.txt", int(length))

			finalMessage = truncateString(finalMessage, 1800)
			// Show the "is typing" message for some amount of time
			time.Sleep(time.Duration(len(finalMessage)*3) * time.Millisecond)
			fmt.Println(finalMessage)
			_, err2 := spaghetti.ChannelMessageSend(m.ChannelID, finalMessage)
			if err2 != nil {
				errorString := err2.Error()
				fmt.Println(errorString)
				spaghetti.ChannelMessageSend(m.ChannelID, errorString)
				}
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

