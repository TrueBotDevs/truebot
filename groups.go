package main

import(
    "github.com/bwmarrin/discordgo"
    "strings"
)

func joinGroup(s *discordgo.Session, msg *discordgo.MessageCreate, arg string){
    group,_ := grabArg(arg)
    channel, _ := s.Channel(msg.ChannelID)
    guildID := channel.GuildID
    switch strings.ToLower(group){
        case "overwatch":
            s.GuildMemberRoleAdd(guildID,msg.Author.ID,"190633106842058754")
            s.ChannelMessageSend(msg.ChannelID, "You have joined Overwatch!")
        case "tabletop":
            s.GuildMemberRoleAdd(guildID,msg.Author.ID,"270691313911857165")
            s.ChannelMessageSend(msg.ChannelID, "You have joined Tabletop Simulator!")
        case "minors":
            s.GuildMemberRoleAdd(guildID,msg.Author.ID,"250769687997186048")
            s.ChannelMessageSend(msg.ChannelID, "You have joined the Minors!")
        case "tiles":
            s.GuildMemberRoleAdd(guildID,msg.Author.ID,"276514629700681728")
            s.ChannelMessageSend(msg.ChannelID, "You have joined Meme Tiles!")
        case "cars":
            s.GuildMemberRoleAdd(guildID,msg.Author.ID,"277545381993381889")
            s.ChannelMessageSend(msg.ChannelID, "You have joined Rocket Cars!")
        default:
            s.ChannelMessageSend(msg.ChannelID, "I can only add you to: Overwatch, Tabletop, Minors, Tiles, and Cars.  Ask Trooble or Slurpee for more groups")
    }
}

func leaveGroup(s *discordgo.Session, msg *discordgo.MessageCreate, arg string){
    group,_ := grabArg(arg)
    channel, _ := s.Channel(msg.ChannelID)
    guildID := channel.GuildID
    switch group{
        case "overwatch":
            s.GuildMemberRoleRemove(guildID,msg.Author.ID,"190633106842058754")
            s.ChannelMessageSend(msg.ChannelID, "You have left Overwatch!")
        case "tabletop":
            s.GuildMemberRoleRemove(guildID,msg.Author.ID,"270691313911857165")
            s.ChannelMessageSend(msg.ChannelID, "You have left Tabletop Simulator!")
        case "minors":
            s.GuildMemberRoleRemove(guildID,msg.Author.ID,"250769687997186048")
            s.ChannelMessageSend(msg.ChannelID, "You have left the Minors!")
        case "tiles":
            s.GuildMemberRoleRemove(guildID,msg.Author.ID,"276514629700681728")
            s.ChannelMessageSend(msg.ChannelID, "You have left Meme Tiles!")
        case "cars":
            s.GuildMemberRoleRemove(guildID,msg.Author.ID,"277545381993381889")
            s.ChannelMessageSend(msg.ChannelID, "You have left Rocket Cars!")
        default:
            s.ChannelMessageSend(msg.ChannelID, "I can only remove you from: Overwatch, Tabletop, Minors, Tiles, and Cars.  Ask Trooble or Slurpee for more groups")
    }
}

func init() {
    CmdList["join"] = joinGroup
    CmdList["leave"] = leaveGroup
}