package main

import (
    "github.com/bwmarrin/discordgo"
    "math/rand"
    "strconv"
)

var (
    cmdsID = "246063490614165504"
)

func exclamationAzorae(s *discordgo.Session, msg *discordgo.MessageCreate, arg string) {
    sender := msg.Author
    channel, _ := s.Channel(msg.ChannelID)
    guildID := channel.GuildID
    guild, _ := s.Guild(guildID)
    var vChannel *discordgo.Channel
    for _, state := range guild.VoiceStates {
        if state.UserID == sender.ID {
            v := state.ChannelID
            vChannel, _ = s.Channel(v)
        }
    }
    if vChannel != nil {
        s.ChannelMessageSend(cmdsID, "<@83742858800009216> you have been pinged to "+vChannel.Name)
        del, _ := s.Channel(cmdsID)
        delThis := del.LastMessageID
        s.ChannelMessageDelete(cmdsID, delThis)
    }
    nextArg, _ := grabArg(arg)
    if nextArg == "stealth" {
        s.ChannelMessageDelete(msg.ChannelID, msg.ID)
    }
}

func imGay(s *discordgo.Session, msg *discordgo.MessageCreate, arg string) {
    s.ChannelMessageSend(msg.ChannelID, "<@"+msg.Author.ID+">, you're a "+strconv.Itoa(rand.Intn(10)-5)+" on the gay scale!")
}
func ihop(s *discordgo.Session, msg *discordgo.MessageCreate, arg string) {
    s.ChannelMessageSend(msg.ChannelID, "I'm down")
}

func dammitSlurpee(s *discordgo.Session, msg *discordgo.MessageCreate, arg string) {
    s.ChannelMessageSend(msg.ChannelID, "$damnitSlurpee")
}

func funnyBecause(s *discordgo.Session, msg *discordgo.MessageCreate, arg string) {
    s.ChannelMessageSend(msg.ChannelID, "It's funny because it's true")
}

func thatsImpossible(s *discordgo.Session, msg *discordgo.MessageCreate, arg string) {
    s.ChannelMessageSend(msg.ChannelID, "No. No! That's not true! That's impossible!")
}

func gayScale(s *discordgo.Session, msg *discordgo.MessageCreate, arg string) {
    s.ChannelMessageSend(msg.ChannelID, "The gay scale was created during AGDQ 2017 by ProGamingWithEd as a way to determine the skill of one's Mahjong play.  The scale ranges from -5 to 5 so as to give the highest levels of accuracy. The scale was later adopted by BGC, Inc. as a whole to determine the relative skill of any given action.")
}

func dotDone(s *discordgo.Session, msg *discordgo.MessageCreate, arg string) {
    index := rand.Intn(7)
    timeString := []string{"00:05.51", "2:11:18", "4:21:09", "03:12.40", "18:39", "47:35", "13:37.69"}
    s.ChannelMessageSend(msg.ChannelID, "<@"+msg.Author.ID+"> completed their run in "+timeString[index]+"!")
}

func init() {
    CmdList["azorae"] = exclamationAzorae
    CmdList["gay"] = imGay
    CmdList["gayscale"] = gayScale
    CmdList["damnit"] = dammitSlurpee
    CmdList["false"] = funnyBecause
    CmdList["true"] = thatsImpossible
    CmdList["ihop"] = ihop
    CmdList["done"] = dotDone
    AliasList["imgay"] = imGay
    AliasList["iamgay"] = imGay
    AliasList["me"] = imGay
    AliasList["greyscale"] = gayScale
    AliasList["grayscale"] = gayScale
    AliasList["dammit"] = dammitSlurpee
}
