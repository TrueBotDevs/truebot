package main

import (
    "github.com/bwmarrin/discordgo"
    "strings"
)

func joinGroup(s *discordgo.Session, msg *discordgo.MessageCreate, arg string) {
    group, _ := grabArg(arg)
    channel, _ := s.Channel(msg.ChannelID)
    guildID := channel.GuildID
    switch strings.ToLower(group) {
    case "overwatch":
        s.GuildMemberRoleAdd(guildID, msg.Author.ID, "190633106842058754")
        s.ChannelMessageSend(msg.ChannelID, "You have joined Overwatch!")
    case "tabletop":
        s.GuildMemberRoleAdd(guildID, msg.Author.ID, "270691313911857165")
        s.ChannelMessageSend(msg.ChannelID, "You have joined Tabletop Simulator!")
    case "minors":
        s.GuildMemberRoleAdd(guildID, msg.Author.ID, "250769687997186048")
        s.ChannelMessageSend(msg.ChannelID, "You have joined the Minors!")
    case "tiles":
        s.GuildMemberRoleAdd(guildID, msg.Author.ID, "276514629700681728")
        s.ChannelMessageSend(msg.ChannelID, "You have joined Meme Tiles!")
    case "cars":
        s.GuildMemberRoleAdd(guildID, msg.Author.ID, "277545381993381889")
        s.ChannelMessageSend(msg.ChannelID, "You have joined Rocket Cars!")
    case "pubg":
        s.GuildMemberRoleAdd(guildID, msg.Author.ID, "304049346846916608")
        s.ChannelMessageSend(msg.ChannelID, "You have joined PUBG!")
    case "deceit":
        s.GuildMemberRoleAdd(guildID, msg.Author.ID, "376550070029385749")
        s.ChannelMessageSend(msg.ChannelID, "You have joined Deceit!")
    default:
        s.ChannelMessageSend(msg.ChannelID, "I can add you to the following groups: ```Overwatch\nTabletop\nMinors\nTiles\nCars\nPUBG\nDeceit```")
    }
}

func leaveGroup(s *discordgo.Session, msg *discordgo.MessageCreate, arg string) {
    group, _ := grabArg(arg)
    channel, _ := s.Channel(msg.ChannelID)
    guildID := channel.GuildID
    switch strings.ToLower(group) {
    case "overwatch":
        s.GuildMemberRoleRemove(guildID, msg.Author.ID, "190633106842058754")
        s.ChannelMessageSend(msg.ChannelID, "You have left Overwatch!")
    case "tabletop":
        s.GuildMemberRoleRemove(guildID, msg.Author.ID, "270691313911857165")
        s.ChannelMessageSend(msg.ChannelID, "You have left Tabletop Simulator!")
    case "minors":
        s.GuildMemberRoleRemove(guildID, msg.Author.ID, "250769687997186048")
        s.ChannelMessageSend(msg.ChannelID, "You have left the Minors!")
    case "tiles":
        s.GuildMemberRoleRemove(guildID, msg.Author.ID, "276514629700681728")
        s.ChannelMessageSend(msg.ChannelID, "You have left Meme Tiles!")
    case "cars":
        s.GuildMemberRoleRemove(guildID, msg.Author.ID, "277545381993381889")
        s.ChannelMessageSend(msg.ChannelID, "You have left Rocket Cars!")
    case "pubg":
        s.GuildMemberRoleRemove(guildID, msg.Author.ID, "304049346846916608")
        s.ChannelMessageSend(msg.ChannelID, "You have left PUBG!")
    case "deceit":
        s.GuildMemberRoleRemove(guildID, msg.Author.ID, "376550070029385749")
        s.ChannelMessageSend(msg.ChannelID, "You have left Deceit!")
    default:
        s.ChannelMessageSend(msg.ChannelID, "I can remove you from the following groups: ```Overwatch\nTabletop\nMinors\nTiles\nCars\nPUBG\nDeceit```")
    }
}

func init() {
    CmdList["join"] = joinGroup
    CmdList["leave"] = leaveGroup
}
