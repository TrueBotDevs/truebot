package main

import(
    "github.com/bwmarrin/discordgo"
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
        s.ChannelMessageSend(msg.ChannelID,"<@83742858800009216> you have been pinged to " + vChannel.Name)
        del, _ := s.Channel(msg.ChannelID)
        delThis := del.LastMessageID
        s.ChannelMessageDelete(msg.ChannelID, delThis)
    }
}

func init() {
    CmdList["azorae"] = exclamationAzorae
}