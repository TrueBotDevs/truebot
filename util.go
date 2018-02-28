package main

import(
    "github.com/bwmarrin/discordgo"
    "strings"
)

func ping(s *discordgo.Session, msg *discordgo.MessageCreate, arg string){
    s.ChannelMessageSend(msg.ChannelID, "Pong!")
}

func pong(s *discordgo.Session, msg *discordgo.MessageCreate, arg string){
    s.ChannelMessageSend(msg.ChannelID, "Ping!")
}

func copycat(s *discordgo.Session, msg *discordgo.MessageCreate, arg string){
    s.ChannelMessageSend(msg.ChannelID, msg.Content)
}

func say(s *discordgo.Session, msg *discordgo.MessageCreate, arg string){
    channel, msg := grabArg(arg)
    dgSession.ChannelMessageSend(channel, arg)
}

func isLive(s *discordgo.Session, msg *discordgo.MessageCreate, arg string){
    sender := msg.Author
    channel, _ := s.Channel(msg.ChannelID)
    guildID := channel.GuildID
    guild, _ := s.Guild(guildID)
    var vChannel *discordgo.Channel
    var vID string
    //Join the voice channel the sender is in
    for _, state := range guild.VoiceStates{
        if state.UserID == sender.ID{
            vID = state.ChannelID
            vChannel, _ = s.Channel(vID)
        }
    }
    if(strings.Contains(vChannel.Name, "🔴 ")){
        s.ChannelEdit(vID,strings.Replace(vChannel.Name, "🔴 ", "", 1))
    }else{
        s.ChannelEdit(vID, "🔴 " + vChannel.Name)
    }
	
	if(len(arg) == 1){
		del, _ := s.Channel(msg.ChannelID)
        delThis := del.LastMessageID
        s.ChannelMessageDelete(msg.ChannelID, delThis)
	}
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
    CmdList["pong"] = pong
    CmdList["say"] = say
    CmdList["copycat"] = copycat
    CmdList["live"] = isLive
}