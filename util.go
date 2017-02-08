package main

import(
    "github.com/bwmarrin/discordgo"
)

func ping(s *discordgo.Session, msg *discordgo.MessageCreate, arg string){
    s.ChannelMessageSend(msg.ChannelID, "Pong!")
}

func copycat(s *discordgo.Session, msg *discordgo.MessageCreate, arg string){
    s.ChannelMessageSend(msg.ChannelID, msg.Content)
}

func init() {
    CmdList["ping"] = ping
    CmdList["copycat"] = copycat
}