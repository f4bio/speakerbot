package nipplebot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
	"time"
)

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) error {
	if s.State.Ready.User.Username == m.Author.Username {
		return nil
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
				_, err := s.ChannelMessageSend(m.ChannelID, "Error: No video found")
				return err
			}

			if voiceInstances[serverID] != nil {
				voiceInstances[serverID].QueueVideo(youtubeLink)
				_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Queued: %s", youtubeTitle))
				return err

			} else {
				_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Playing: %s", youtubeTitle))
				go CreateVoiceInstance(youtubeLink, serverID)
				return err
			}
		} else if method == "stop" && voiceInstances[serverID] != nil {
			voiceInstances[serverID].StopVideo()
		} else if method == "skip" && voiceInstances[serverID] != nil {
			voiceInstances[serverID].SkipVideo()
		} else if method == "help" {
			msg := fmt.Sprintf("%s\\n%s\\%s",
				"**!play** <youtube link or query> - Search/Play Youtube link, queues up if another track is playing",
				"**!skip** - Skip current playing track",
				"**!stop** - Stops tracks and clears queue")
			_, err := s.ChannelMessageSend(m.ChannelID, msg)
			return err
		}
	}
}
