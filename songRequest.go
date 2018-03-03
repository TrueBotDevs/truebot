package main

import(
    "github.com/bwmarrin/discordgo"
	"io"
	"fmt"
	//"time"
	"strconv"
	"github.com/jonas747/dca"
	"github.com/rylio/ytdl"
)

var vc *discordgo.VoiceConnection
var songs = make([]Song,0)
var songsFinished = true

type Song struct{
	
	Message *discordgo.MessageCreate
	Session *discordgo.Session
	Arg string
}

func checkSong(s *discordgo.Session, msg *discordgo.MessageCreate, arg string){
	_, err := ytdl.GetVideoInfo(arg)
	if err != nil {
		s.ChannelMessageSend(msg.ChannelID,"Video not found")
	}else{
		song := Song{msg,s,arg,}
		songs = append(songs, song)
		//fmt.Println(len(songs))
		
		if len(songs) > 1{
			s.ChannelMessageSend(strconv.Itoa(246063490614165504),"```Your song is currently number " + strconv.Itoa(len(songs)) + " in the queue.```")
		}
	}
	if(songsFinished == true){
		songsFinished = false;
		go playSong();
	}
}

func playSong(){
	for true{
        if hasSession{
			song := songs[0]

			s := song.Session
			msg := song.Message
			arg := song.Arg
			
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
			}
			//notify of voice channel
			if vChannel != nil{
				s.ChannelMessageSend(strconv.Itoa(246063490614165504),"```Playing in: " + vChannel.Name + "\nSong: " + videoInfo.Title + "\nRequested by: " + msg.Author.Username + "```")
				//s.ChannelMessageSend(msg.ChannelID,"You are in " + vChannel.Name + ", the play command is under development")
				vc, _ = s.ChannelVoiceJoin(guildID, vID, false, false)
			}else{
				s.ChannelMessageSend(msg.ChannelID,"You are not in a voice channel, the play command is under development")
			}    
	
			//need to check for song that is given in argument
			err = vc.Speaking(true)
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
			
			if(len(songs) > 1){
				songs = songs[1:len(songs)]
			}else{
				songs = nil
			}
			
			if(len(songs) < 1){
				vc.Disconnect()
				songsFinished = true
				break
			}
		}
	}		
}
		
func skipSong(s *discordgo.Session, msg *discordgo.MessageCreate, arg string){
	sender := msg.Author //Person making the command (User)
	channel, _ := s.Channel(msg.ChannelID) //Channel sender made request in (ChannelID)
	guildID := channel.GuildID //Server ID
	guild, _ := s.Guild(guildID) //Actual server object
	var vID string //Voice ChannelID
	
	for _, state := range guild.VoiceStates{
		if state.UserID == sender.ID{
			vID = state.ChannelID
		}
	}
	
	if vc != nil{
		vc.Disconnect()
		if len(songs) != 0{
			vc, _ = s.ChannelVoiceJoin(guildID, vID, false, false)
		}
	}
}
		
func stopMusic(s *discordgo.Session, msg *discordgo.MessageCreate, arg string){
	if vc != nil{
		vc.Disconnect()
	}
	
	songs = nil
		
	s.UpdateStatus(0,"")
}

func songInfo(s *discordgo.Session, msg *discordgo.MessageCreate, arg string){
	song := songs[0]
	msg1 := song.Message
	arg1 := song.Arg
	
	videoInfo, err := ytdl.GetVideoInfo(arg1)
	if err != nil {
		s.ChannelMessageSend(msg.ChannelID,"Video not found")
		// Handle the error
	}else{
		songName := videoInfo.Title
		s.ChannelMessageSend(msg.ChannelID,"```Song: " + songName + "\nRequested by: " + msg1.Author.Username + "```")
	}
	
}
func init() {
	CmdList["songinfo"] = songInfo
	AliasList["current"] = songInfo
    CmdList["play"] = checkSong
    AliasList["queue"] = checkSong
	AliasList["songrequest"] = checkSong
	AliasList["songtest"] = checkSong
	CmdList["skip"] = skipSong
    CmdList["stop"] = stopMusic
    AliasList["stahp"] = stopMusic
}