package main

import(
    "github.com/bwmarrin/discordgo"
	"io"
	"fmt"
	"github.com/jonas747/dca"
	"github.com/rylio/ytdl"
)

var vc *discordgo.VoiceConnection
var songs = make([]string,1)
var songStream *dca.StreamingSession

func playSong(s *discordgo.Session, msg *discordgo.MessageCreate, arg string){
    //checkYT(s,msg,arg)
	sender := msg.Author //Person making the command (User)
    channel, _ := s.Channel(msg.ChannelID) //Channel sender made request in (ChannelID)
    guildID := channel.GuildID //Server ID
    guild, _ := s.Guild(guildID) //Actual server object
    var vChannel *discordgo.Channel //Voice channel
    var vID string //Voice ChannelID
	

	
    videoInfo, err := ytdl.GetVideoInfo(arg)
	if err != nil {
		s.ChannelMessageSend(msg.ChannelID,"Video not found")
		// Handle the error
	}else{
		//if songStream != nil && songStream.running{
			//songs = songs.append(arg)
		//}
		songName := videoInfo.Title
		s.UpdateStatus(0,songName)
		defer s.UpdateStatus(0,"")
	
		//Join the voice channel the sender is in (finds sender channel)
		for _, state := range guild.VoiceStates{
			if state.UserID == sender.ID{
				vID = state.ChannelID
				vChannel, _ = s.Channel(vID)
			}
		}
		//notify of voice channel
		if vChannel != nil{
			s.ChannelMessageSend(msg.ChannelID,"You are in " + vChannel.Name + ", the play command is under development")
			vc, _ = s.ChannelVoiceJoin(guildID, vID, false, false)
		}else{
			s.ChannelMessageSend(msg.ChannelID,"You are not in a voice channel, the play command is under development")
		}    
	
		//need to check for song that is given in argument
		//
		err := vc.Speaking(true)
		if err != nil{
			fmt.Println("Well shit, can't speak")
		}
		defer vc.Speaking(false)
	
	
		options := dca.StdEncodeOptions
		options.RawOutput = true
		options.Bitrate = 96
		options.Application = "lowdelay"
		options.Volume = 100


		format := videoInfo.Formats.Extremes(ytdl.FormatAudioBitrateKey, true)[0]
		downloadURL, err := videoInfo.GetDownloadURL(format)
		if err != nil {
			// Handle the error
		}

		encodingSession, err := dca.EncodeFile(downloadURL.String(), options)
		if err != nil {
			// Handle the error
		}
		defer encodingSession.Cleanup()
    
		done := make(chan error)    
		dca.NewStream(encodingSession, vc, done)
		err = <- done
		if err != nil && err != io.EOF {
			// Handle the error
		}
		vc.Disconnect()
    }
}

func stopMusic(s *discordgo.Session, msg *discordgo.MessageCreate, arg string){
	if vc != nil{
		vc.Disconnect()
	}
	s.UpdateStatus(0,"")
}



func init() {
    CmdList["play"] = playSong
    AliasList["queue"] = playSong
    CmdList["stop"] = stopMusic
    AliasList["stahp"] = stopMusic
}