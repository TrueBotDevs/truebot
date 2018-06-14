package main

import (
	"fmt"
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
	case "bombs":
		s.GuildMemberRoleAdd(guildID, msg.Author.ID, "424309460467449864")
		s.ChannelMessageSend(msg.ChannelID, "You have joined Bombs!")
	case "civ":
		s.GuildMemberRoleAdd(guildID, msg.Author.ID, "349279584694566934")
		s.ChannelMessageSend(msg.ChannelID, "You have joined Civ!")
	default:
		s.ChannelMessageSend(msg.ChannelID, "I can add you to the following groups: ```Overwatch\nTabletop\nMinors\nTiles\nCars\nPUBG\nDeceit\nBombs\nCiv```")
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
	case "bombs":
		s.GuildMemberRoleRemove(guildID, msg.Author.ID, "424309460467449864")
		s.ChannelMessageSend(msg.ChannelID, "You have left Bombs!")
	case "civ":
		s.GuildMemberRoleRemove(guildID, msg.Author.ID, "349279584694566934")
		s.ChannelMessageSend(msg.ChannelID, "You have left Civ!")
	default:
		s.ChannelMessageSend(msg.ChannelID, "I can remove you from the following groups: ```Overwatch\nTabletop\nMinors\nTiles\nCars\nPUBG\nDeceit\nBombs\nCiv```")
	}
}

func addGroup(s *discordgo.Session, msg *discordgo.MessageCreate, arg string) {
	//GuildRoleCreate to get a new guild, then GuildRoleEdit(guildID, roleID, name string, color int, hoist bool, perm int, mention bool) to name and stuff?
	channel, _ := s.Channel(msg.ChannelID)
	guildID := channel.GuildID
	newRole, err := s.GuildRoleCreate(guildID)

	if err != nil {
		fmt.Println(err)
	}
	_, err2 := s.GuildRoleEdit(guildID, newRole.ID, arg, 0x99AAB5, false, 0, true)
	if err2 != nil {
		fmt.Println(err2)
	}
}

func init() {
	CmdList["addgroup"] = addGroup
	CmdList["join"] = joinGroup
	CmdList["leave"] = leaveGroup
}
