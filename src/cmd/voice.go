package nipplebot

import (
	"encoding/binary"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
	"github.com/oleiade/lane"
)

var (
	run            *exec.Cmd
	voiceInstances = map[string]*VoiceInstance{}
)

const (
	channels  int = 2
	frameRate int = 48000
	frameSize int = 960
)

// VoiceInstance is created for each connected server
type VoiceInstance struct {
	discord      *discordgo.Session
	queue        *lane.Queue
	pcmChannel   chan []int16
	serverID     string
	skip         bool
	stop         bool
	trackPlaying bool
}

func (vi *VoiceInstance) playVideo(url string) {
	vi.trackPlaying = true

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Http.Get\nerror: %s\ntarget: %s\n", err, url)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("reading answer: non 200 status code received: '%s'", err)
	}

	run = exec.Command("ffmpeg", "-i", "-", "-f", "s16le", "-ar", strconv.Itoa(frameRate), "-ac", strconv.Itoa(channels), "pipe:1")
	run.Stdin = resp.Body
	stdout, err := run.StdoutPipe()
	if err != nil {
		fmt.Println("StdoutPipe Error:", err)
		return
	}

	err = run.Start()
	if err != nil {
		fmt.Println("RunStart Error:", err)
		return
	}

	// buffer used during loop below
	audiobuf := make([]int16, frameSize*channels)

	//vi.discord.Voice.Speaking(true)
	//defer vi.discord.Voice.Speaking(false)

	for {
		// read data from ffmpeg stdout
		err = binary.Read(stdout, binary.LittleEndian, &audiobuf)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			break
		}
		if err != nil {
			fmt.Println("error reading from ffmpeg stdout :", err)
			break
		}
		if vi.stop == true || vi.skip == true {
			err = run.Process.Kill()
			break
		}
		vi.pcmChannel <- audiobuf
	}

	vi.trackPlaying = false
}

// StopVideo marks to stop all tracks and clears queue on the next binary read.
func (vi *VoiceInstance) StopVideo() {
	vi.stop = true
}

// SkipVideo skips the current playing track
func (vi *VoiceInstance) SkipVideo() {
	vi.skip = true
}

func (vi *VoiceInstance) connectVoice() {
	vi.discord, _ = discordgo.New(email, password)

	// Open the websocket and begin listening.
	err := vi.discord.Open()
	if err != nil {
		fmt.Println(err)
	}

	channels, err := vi.discord.GuildChannels(vi.serverID)

	var voiceChannel string
	var voiceChannels []string
	for _, channel := range channels {
		if channel.Type == discordgo.ChannelTypeGuildVoice {
			voiceChannels = append(voiceChannels, channel.ID)
			if strings.Contains(strings.ToLower(channel.Name), "music") && voiceChannel == "" {
				voiceChannel = channel.ID
			}
		}
	}

	if voiceChannel == "" {
		fmt.Println("Selecting first channel")
		voiceChannel = voiceChannels[0]
	}

	err = vi.discord.ChannelVoiceJoin(vi.serverID, voiceChannel, false, true)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Hacky loop to prevent returning when voice isn't ready
	// TODO: Find a better way.
	for vi.discord.Voice.Ready == false {
		runtime.Gosched()
	}
}

// QueueVideo places a Youtube link in a queue
func (vi *VoiceInstance) QueueVideo(youtubeLink string) {
	fmt.Println("Queuing video")
	vi.queue.Enqueue(youtubeLink)
}

func (vi *VoiceInstance) processQueue() {
	if vi.trackPlaying == false {
		for {
			vi.skip = false
			link := vi.queue.Dequeue()
			if link == nil || vi.stop == true {
				break
			}
			vi.playVideo(link.(string))
		}

		// No more tracks in queue? Cleanup.
		fmt.Println("Closing connections")
		close(vi.pcmChannel)
		vi.discord.Voice.Close()
		vi.discord.Close()
		delete(voiceInstances, vi.serverID)
		fmt.Println("Done")
	}
}

// CreateVoiceInstance accepts both a youtube query and a server id, boots up the voice connection, and plays the track.
func CreateVoiceInstance(youtubeLink string, serverID string) {
	vi := new(VoiceInstance)
	voiceInstances[serverID] = vi

	fmt.Println("Connecting Voice...")
	vi.serverID = serverID
	vi.queue = lane.NewQueue()
	vi.connectVoice()

	vi.pcmChannel = make(chan []int16, 2)
	go SendPCM(vi.discord.Voice, vi.pcmChannel)

	vi.QueueVideo(youtubeLink)
	vi.processQueue()
}
