package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/paked/configure"
)

var (
	conf      = configure.New()
	email     = conf.String("email", "", "Discord email address")
	password  = conf.String("password", "", "Discord password")
	googleKey = conf.String("googleKey", "", "Google API key for Youtube API")
)

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if s.State.Ready.User.Username == m.Author.Username {
		return
	}

	fmt.Printf("%20s %20s %20s > %s\n", m.ChannelID, time.Now().Format(time.Stamp), m.Author.Username, m.Content)

	if m.Content[:1] == "!" {
		channel, _ := s.Channel(m.ChannelID)
		serverID := channel.GuildID
		method := strings.Split(m.Content, " ")[0][1:]

		if method == "play" {
			youtubeLink, youtubeTitle, err := GetYoutubeURL(strings.Split(m.Content, " ")[1])
			if err != nil {
				fmt.Println(err)
				s.ChannelMessageSend(m.ChannelID, "Error: No video found")
				return
			}

			if voiceInstances[serverID] != nil {
				voiceInstances[serverID].QueueVideo(youtubeLink)
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Queued: %s", youtubeTitle))
			} else {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Playing: %s", youtubeTitle))
				go CreateVoiceInstance(youtubeLink, serverID)
			}
		} else if method == "stop" && voiceInstances[serverID] != nil {
			voiceInstances[serverID].StopVideo()
		} else if method == "skip" && voiceInstances[serverID] != nil {
			voiceInstances[serverID].SkipVideo()
		} else if method == "help" {
			s.ChannelMessageSend(m.ChannelID, `
**!play** <youtube link or query> - Search/Play Youtube link, queues up if another track is playing 
**!skip** - Skip current playing track
**!stop** - Stops tracks and clears queue`
			)
		}
	}
}

func main() {
	// Pull in configuration
	conf.Use(configure.NewFlag())
	conf.Use(configure.NewEnvironment())
	if _, err := os.Stat("config.json"); err == nil {
		conf.Use(configure.NewJSONFromFile("config.json"))
	}
	conf.Parse()

	discord, err := discordgo.New(*email, *password)
	if err != nil {
		fmt.Println("Error logging in")
		fmt.Println(err)
	}

	discord.AddHandler(messageCreate)

	// Open the websocket and begin listening.
	err = discord.Open()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Listening...")
	lock := make(chan int)
	<-lock
}
