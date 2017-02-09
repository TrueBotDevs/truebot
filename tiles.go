package main

import (
  "fmt"
  "github.com/bwmarrin/discordgo"
)

var (
  tilesID := "270331769033588737"
)

// Create the queue type!
type queue struct {
  owner int
  started bool
  pName := make([]string, 0)
  sName := make([]string, 0)
  pID := make([]string, 0)
  sID := make([]string, 0)
}

// TODO: This shoud really be a switch statement, but fuck it.  I'll change it later.
func tilesChooser(s *discordgo.Session, msg *discordgo.MessageCreate, arg string)  {
  if arg == "start" {
    tilesStart(s *discordgo.Session, msg *discordgo.MessageCreate)
  } else if arg == "join" {
    tilesJoin(s *discordgo.Session, msg *discordgo.MessageCreate)
  } else if arg == "standby" {
    tilesStandby(s *discordgo.Session, msg *discordgo.MessageCreate)
  } else if arg == "play" {
    tilesPlay(s *discordgo.Session, msg *discordgo.MessageCreate)
  } else if arg == "check" {
    tilesCheck(s *discordgo.Session, msg *discordgo.MessageCreate)
  } else if arg == "clear" {
    tilesClear(s *discordgo.Session, msg *discordgo.MessageCreate)
  } else {
    tilesHelp(s *discordgo.Session, msg *discordgo.MessageCreate)
  }
}

func tilesStart(s *discordgo.Session, msg *discordgo.MessageCreate)  {
  if q.started {
    s.ChannelMessageSend(msg.ChannelID, "<@" + msg.Author.ID + ">" + ", there is currently a queue already started.  Use `!tiles check` to see who's in the queue.")
    return
  }
  s.ChannelMessageSend(tilesID, "Hey <@276514629700681728>!" + "<@" + msg.Author.ID + ">" + " wants to start a hanchan! Type `!tiles join` to enter the queue.")

  q := queue{owner: msg.Author.id, started: true}
  q.pName := append(msg.Author.Username)
  q.pID := append(msg.Author.ID)
}

func tilesJoin(s *discordgo.Session, msg *discordgo.MessageCreate)  {
  if !q.started {
    s.ChannelMessageSend(msg.ChannelID, "<@" + msg.Author.ID + ">" + ", there is no queue right now.  Type `!tiles start` to start one!")
    return
  }
  s.ChannelMessageSend(msg.ChannelID, "<@" + msg.Author.ID + ">" + ", you've been added to the queue!  Currently we have " + len(q.pID) + " players in the queue, and " + len(q.sID) + " standby players.")
  q.pName := append(msg.Author.Username)
  q.pID := append(msg.Author.ID)

  if len(q.pID) + len(q.sID) >= 4 {
    s.ChannelMessageSend(tilesID, "<@" + q.owner + ">" + ", you have enough players in queue and standby to play a hanchan.  Use `!tiles play` to ping players and clear the queue.")
  }
}

func tilesStandby(s *discordgo.Session, msg *discordgo.MessageCreate)  {
  if !q.started {
    s.ChannelMessageSend(msg.ChannelID, "<@" + msg.Author.ID + ">" + ", there is no queue right now.  Type `!tiles start` to start one!")
    return
  }
  s.ChannelMessageSend(msg.ChannelID, "<@" + msg.Author.ID + ">" + ", you've been added to the standby queue!  Currently we have " + len(q.pID) + " players in the queue, and " + len(q.sID) + " standby players.")
  q.sName := append(msg.Author.Username)
  q.sID := append(msg.Author.ID)

  if len(q.pID) + len(q.sID) >= 4 {
    s.ChannelMessageSend(tilesID, "<@" + q.owner + ">" + ", you have enough players in queue and standby to play a hanchan.  Use `!tiles play` to ping players and clear the queue.")
}

func tilesplay()  {
  if !q.started {
    s.ChannelMessageSend(msg.ChannelID, "<@" + msg.Author.ID + ">" + ", there is no queue right now.  Type `!tiles start` to start one!")
    return
  }
  if len(q.pID) + len(q.sID) < 4 {
    s.ChannelMessageSend(msg.ChannelID, "<@" + msg.Author.ID + ">" +", there aren't enough people in the queue to start. Get more people!")
    return
  }

  // TODO: Make this a ready check style goroutine.
  for i := 1; i <= len(q.pID); i++ {
    s.ChannelMessageSend(tilesID, "Hey <@" + [i]q.pID + ">!")
  }
  if len(q.pID) < 4 {
    need := 4 - len(q.pID)
    for i := 0; i <= need; i++ {
      s.ChannelMessageSend(tilesID, "Hey <@" + [i]q.sID + ">!")
    }
  }
  s.ChannelMessageSend(tilesID, "<@" + q.owner + "> has summoned you to start the hanchan!")
}

func tilesHelp(s *discordgo.Session)  {
  s.ChannelMessageSend(msg.ChannelID, "<@" + m.author.ID + ">" + "\n```Usage: !tiles [start|join|standby|play|check|clear|help]\n\n!tiles start - Start a new hanchan queue.\n!tiles join - Join a hanchan queue.\n!tiles standby - Join a hanchan queue as a standby.\n!tiles play - Pings all queued players and clears the queue.\n!tiles check - Checks for a currently queuing hanchan.\n!tiles clear - Clears the current hanchan queue. (game owner only)\n!tiles help - This text.```")
}

func init()  {
  CmdList["tiles"] = tilesChooser
}
