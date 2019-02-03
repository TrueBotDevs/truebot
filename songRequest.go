package main

import (
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/go-ini/ini"
	"github.com/jonas747/dca"
	"github.com/rylio/ytdl"
	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var vc *discordgo.VoiceConnection
var songs = make([]Song, 0)
var songsFinished = true

var maxResults = flag.Int64("max-results", 25, "Max YouTube results")

var youtubeKey string
var botCommandsChannel string

//Song struct
type Song struct {
	Message *discordgo.MessageCreate
	Session *discordgo.Session
	Arg     string
}

func checkSong(s *discordgo.Session, msg *discordgo.MessageCreate, arg string) {
	arg = strings.TrimSpace(arg)
	var vid string
	var vid2 string

	if arg != "" {
		vid = ytSearch(s, msg, arg)
		vid2 = "https://www.youtube.com/watch?v=" + vid
		fmt.Println("ID= " + vid2)

	} else {
		vid = ""
	}

	_, err := ytdl.GetVideoInfo(vid2)
	if err != nil || vid == "" {
		s.ChannelMessageSend(msg.ChannelID, "Video not found")
	} else {
		song := Song{msg, s, vid2}
		songs = append(songs, song)
		//fmt.Println(len(songs))

		if len(songs) > 1 {
			s.ChannelMessageSend(botCommandsChannel, "```Your song is currently number "+strconv.Itoa(len(songs))+" in the queue.```")
		}
		if songsFinished == true {
			songsFinished = false
			go playSong()
		}
	}

}

func playSong() {
	for true {
		if hasSession {
			song := songs[0]

			s := song.Session
			msg := song.Message
			arg := song.Arg

			sender := msg.Author                   //Person making the command (User)
			channel, _ := s.Channel(msg.ChannelID) //Channel sender made request in (ChannelID)
			guildID := channel.GuildID             //Server ID
			guild, _ := s.Guild(guildID)           //Actual server object
			var vChannel *discordgo.Channel        //Voice channel
			var vID string                         //Voice ChannelID

			videoInfo, err := ytdl.GetVideoInfo(arg)
			if err != nil {
				s.ChannelMessageSend(msg.ChannelID, "Video not found")
				// Handle the error
			} else {
				songName := videoInfo.Title
				s.UpdateStatus(0, songName)
				defer s.UpdateStatus(0, "")

				//Join the voice channel the sender is in (finds sender channel)
				for _, state := range guild.VoiceStates {
					if state.UserID == sender.ID {
						vID = state.ChannelID
						vChannel, _ = s.Channel(vID)
					}
				}
			}
			//notify of voice channel
			if vChannel != nil {
				s.ChannelMessageSend(botCommandsChannel, "```Playing in: "+vChannel.Name+"\nSong: "+videoInfo.Title+"\nLength: "+videoInfo.Duration.String()+"\nRequested by: "+msg.Author.Username+"```")
				//s.ChannelMessageSend(msg.ChannelID,"You are in " + vChannel.Name + ", the play command is under development")
				vc, _ = s.ChannelVoiceJoin(guildID, vID, false, false)
			} else {
				s.ChannelMessageSend(msg.ChannelID, "You are not in a voice channel, the play command is under development")
			}

			//need to check for song that is given in argument
			err = vc.Speaking(true)
			if err != nil {
				fmt.Println("Well shit, can't speak")
			}
			defer vc.Speaking(false)

			options := dca.StdEncodeOptions
			options.RawOutput = true
			options.Bitrate = 96
			options.Application = "lowdelay"
			options.Volume = 100

			fmt.Println(arg)
			format := videoInfo.Formats.Extremes(ytdl.FormatAudioBitrateKey, true)[0]
			downloadURL, err := videoInfo.GetDownloadURL(format)
			if err != nil {
				fmt.Println(err)
			}

			encodingSession, err := dca.EncodeFile(downloadURL.String(), options)

			if err != nil {
				fmt.Println(err)
			}
			defer encodingSession.Cleanup()

			done := make(chan error)

			dca.NewStream(encodingSession, vc, done)

			err = <-done

			if err != nil && err != io.EOF {
				fmt.Println(err)
			}

			if len(songs) > 1 {
				songs = songs[1:len(songs)]
			} else {
				songs = nil
			}

			if len(songs) < 1 {
				vc.Disconnect()
				songsFinished = true
				break
			}
		}
	}
}

func skipSong(s *discordgo.Session, msg *discordgo.MessageCreate, arg string) {
	sender := msg.Author                   //Person making the command (User)
	channel, _ := s.Channel(msg.ChannelID) //Channel sender made request in (ChannelID)
	guildID := channel.GuildID             //Server ID
	guild, _ := s.Guild(guildID)           //Actual server object
	var vID string                         //Voice ChannelID

	for _, state := range guild.VoiceStates {
		if state.UserID == sender.ID {
			vID = state.ChannelID
		}
	}

	if vc != nil {
		vc.Disconnect()
		if len(songs) != 0 {
			vc, _ = s.ChannelVoiceJoin(guildID, vID, false, false)
		}
	}
}

func stopMusic(s *discordgo.Session, msg *discordgo.MessageCreate, arg string) {
	if vc != nil {
		vc.Disconnect()
	}

	songs = nil

	s.UpdateStatus(0, "")
}

func songInfo(s *discordgo.Session, msg *discordgo.MessageCreate, arg string) {
	if vc != nil {
		song := songs[0]
		msg1 := song.Message
		arg1 := song.Arg

		videoInfo, err := ytdl.GetVideoInfo(arg1)
		if err != nil {
			s.ChannelMessageSend(msg.ChannelID, "Video not found")
			// Handle the error
		} else {
			songName := videoInfo.Title
			s.ChannelMessageSend(msg.ChannelID, "```Song: "+songName+"\nLength: "+videoInfo.Duration.String()+"\nRequested by: "+msg1.Author.Username+"```")
		}
	} else {
		s.ChannelMessageSend(msg.ChannelID, "```Not Currently Playing```")
	}
}

func ytSearch(s *discordgo.Session, msg *discordgo.MessageCreate, arg string) string {
	var vids = make([]string, 0)

	flag.Parse()

	client := &http.Client{
		Transport: &transport.APIKey{Key: youtubeKey},
	}

	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Error creating new YouTube client: %v", err)
	}

	// Make the API call to YouTube.
	call := service.Search.List("id,snippet").
		Q(arg).
		MaxResults(*maxResults)
	response, err := call.Do()
	if err != nil {
		log.Fatalf("Error making search API call: %v", err)
	}

	// Iterate through each item and add it to the correct list.
	for _, item := range response.Items {
		switch item.Id.Kind {
		case "youtube#video":
			vids = append(vids, item.Id.VideoId)
			//videos[item.Id.VideoId] = item.Snippet.Title
		}
	}

	//fmt.Println(strconv.Itoa(len(vids)) + " " + vids[0])

	if len(vids) > 0 {
		return vids[0]
	}
	return ""
}

func init() {
	CmdList["songinfo"] = songInfo
	AliasList["current"] = songInfo
	CmdList["play"] = checkSong
	AliasList["queue"] = checkSong
	AliasList["songrequest"] = checkSong
	//AliasList["searchtest"] = checkSong
	//CmdList["searchtest"] = ytSearch
	CmdList["skip"] = skipSong
	CmdList["stop"] = stopMusic
	AliasList["stahp"] = stopMusic
	cfg, err := ini.Load("./config/truebot.ini")
	if err != nil {
		fmt.Println("Was not able to load YouTube API Key - ", err)
	}
	youtubeKey = cfg.Section("api-keys").Key("youtube").String()
	botCommandsChannel = cfg.Section("channels").Key("bot-commands").String()
}
