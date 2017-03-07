package main

import(
    "github.com/bwmarrin/discordgo"
)
var(
    cmdsID = "246063490614165504"
)

func exclamationAzorae(s *discordgo.Session, msg *discordgo.MessageCreate, arg string){
    sender := msg.Author
    channel, _ := s.Channel(msg.ChannelID)
    guildID := channel.GuildID
    guild, _ := s.Guild(guildID)
    var vChannel *discordgo.Channel
    for _, state := range guild.VoiceStates{
        if state.UserID == sender.ID{
            v := state.ChannelID
            vChannel, _ = s.Channel(v)
        }
    }
    if vChannel != nil{
        s.ChannelMessageSend(cmdsID,"<@83742858800009216> you have been pinged to " + vChannel.Name)
        del, _ := s.Channel(cmdsID)
        delThis := del.LastMessageID
        s.ChannelMessageDelete(cmdsID, delThis)
    }
    nextArg, _ := grabArg(arg)
    if nextArg == "stealth"{
        s.ChannelMessageDelete(msg.ChannelID,msg.ID)
    }
}

func init() {
    CmdList["azorae"] = exclamationAzorae
}