package main

import(
    "github.com/bwmarrin/discordgo"
    "strings"
)

func ping(s *discordgo.Session, msg *discordgo.MessageCreate, arg string){
    s.ChannelMessageSend(msg.ChannelID, "Pong!")
}

func copycat(s *discordgo.Session, msg *discordgo.MessageCreate, arg string){
    s.ChannelMessageSend(msg.ChannelID, msg.Content)
}

//This might want to go in the main file
func grabArg(s string) (string,string){
    arg := strings.Split(s, " ")[0]
    remainder := " "
    if len(s) > len(arg)+1{
            remainder = strings.Replace(s, arg + " ", "", 1)
    }
    return arg,remainder
}

func init() {
    CmdList["ping"] = ping
    CmdList["copycat"] = copycat
}