package main

import(
    "github.com/bwmarrin/discordgo"
)

var vc *discordgo.VoiceConnection

func playSong(s *discordgo.Session, msg *discordgo.MessageCreate, arg string){
    sender := msg.Author
    channel, _ := s.Channel(msg.ChannelID)
    guildID := channel.GuildID
    guild, _ := s.Guild(guildID)
    var vChannel *discordgo.Channel
    var vID string
    songName := "you for a fool"
    //Join the voice channel the sender is in
    for _, state := range guild.VoiceStates{
        if state.UserID == sender.ID{
            vID = state.ChannelID
            vChannel, _ = s.Channel(vID)
        }
    }
    if vChannel != nil{
        s.ChannelMessageSend(msg.ChannelID,"You are in " + vChannel.Name + ", the play command is under development")
        vc, _ = s.ChannelVoiceJoin(guildID, vID, false, false)
    }else{
        s.ChannelMessageSend(msg.ChannelID,"You are not in a voice channel, the play command is under development")
    }    
    s.UpdateStatus(0,songName)
}

func stopMusic(s *discordgo.Session, msg *discordgo.MessageCreate, arg string){
    vc.Disconnect()
}

func init() {
    CmdList["play"] = playSong
    AliasList["queue"] = playSong
    CmdList["stop"] = stopMusic
    AliasList["stahp"] = stopMusic
}