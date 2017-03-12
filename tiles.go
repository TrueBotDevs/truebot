package main

import (
  "github.com/bwmarrin/discordgo"
  "strconv"
)

var (
  tilesID = "270331769033588737"
  q = Queue{}
)

// Create the queue type!
type Queue struct {
  owner string
  started bool
  pName []string
  pID []string
}

// TODO: This should really be a switch statement, but fuck it.  I'll change it later.
func tilesChooser(s *discordgo.Session, msg *discordgo.MessageCreate, arg string)  {
  if arg == "start" {
    tilesStart(s, msg)
  } else if arg == "join" {
    tilesJoin(s, msg)
  } else if arg == "leave" {
    tilesLeave(s, msg)
  } else if arg == "play" {
    tilesPlay(s, msg)
  } else if arg == "check" {
    tilesCheck(s, msg)
  } else if arg == "clear" {
    tilesClear(s, msg)
  } else if arg == "1pin" {
    s.ChannelMessageSend(msg.ChannelID, "私の一筒!")
    return
  } else if arg == "[start|join|leave|play|check|clear|help]"{
    s.ChannelMessageSend(msg.ChannelID, "<@" + msg.Author.ID + ">, actively go fuck yourself. Especially if you're Church.")
    return
  } else {
    tilesHelp(s, msg)
  }
}

func tilesStart(s *discordgo.Session, msg *discordgo.MessageCreate)  {
  if q.started {
    s.ChannelMessageSend(msg.ChannelID, "<@" + msg.Author.ID + ">, there's already a queue in progress!  Use `!tiles check` to see who's in the queue.")
    return
  }
  s.ChannelMessageSend(tilesID, "Hey <@&276514629700681728>! " + "<@" + msg.Author.ID + "> wants to start a hanchan! Type `!tiles join` to enter the queue.")

  q.started = true
  q.pName = append(q.pName, msg.Author.Username)
  q.pID = append(q.pID, msg.Author.ID)
  q.owner = msg.Author.ID
}

func tilesJoin(s *discordgo.Session, msg *discordgo.MessageCreate)  {
  if !q.started {
    s.ChannelMessageSend(msg.ChannelID, "<@" + msg.Author.ID + ">, there is no queue to join!  Type `!tiles start` to start one!")
    return
  }
  if q.owner == msg.Author.ID {
    s.ChannelMessageSend(msg.ChannelID, "<@" + msg.Author.ID + ">, you are the owner of this queue, you can't join again!")
    return
  }
  for i := 1; i < len(q.pID); i++ {
    if msg.Author.ID == q.pID[i] {
      s.ChannelMessageSend(msg.ChannelID, "<@" + msg.Author.ID + ">, you are already in the queue!")
      return
    }
  }
  if len(q.pID) == 4 {
    s.ChannelMessageSend(msg.ChannelID, "<@" + msg.Author.ID + ">, the play queue is currently maxed out right now.  Not adding player.")
    return
  }
  q.pName = append(q.pName, msg.Author.Username)
  q.pID = append(q.pID, msg.Author.ID)
  s.ChannelMessageSend(msg.ChannelID, "<@" + msg.Author.ID + ">, you've been added to the queue!  There are " + strconv.Itoa(len(q.pID)) + " players in the queue.")

  if len(q.pID) == 4 {
    s.ChannelMessageSend(tilesID, "<@" + q.owner + ">, you have enough players in queue and standby to play a hanchan.  Use `!tiles play` to ping players and clear the queue.")
  }
}

func tilesLeave(s *discordgo.Session, msg *discordgo.MessageCreate)  {
  if !q.started {
    s.ChannelMessageSend(msg.ChannelID, "<@" + msg.Author.ID + ">, there is no queue to leave! Type `!tiles start` to start one!")
    return
  }
  if q.owner == msg.Author.ID {
    s.ChannelMessageSend(msg.ChannelID, "<@" + msg.Author.ID + ">, logic for owner leaving is in the works.  Please clear the queue instead for the time being.")
    return
  }
  for i := 1; i < len(q.pID); i++ {
    if msg.Author.ID == q.pID[i] {
      q.pID = append(q.pID[:i], q.pID[i+1:]...)
      q.pName = append(q.pName[:i], q.pName[i+1:]...)
      s.ChannelMessageSend(msg.ChannelID, "<@" + msg.Author.ID + ">, you have been removed from the queue.")
    }
  }
}

func tilesPlay(s *discordgo.Session, msg *discordgo.MessageCreate)  {
  if !q.started {
    s.ChannelMessageSend(msg.ChannelID, "<@" + msg.Author.ID + ">, there is no queue to start playing!  Type `!tiles start` to start one!")
    return
  }
  if msg.Author.ID != q.owner {
    s.ChannelMessageSend(msg.ChannelID, "<@" + msg.Author.ID + ">, you are not the owner of this queue. Ignoring.")
    return
  }
  if len(q.pID) < 4 {
    s.ChannelMessageSend(msg.ChannelID, "<@" + msg.Author.ID + ">, there aren't enough people in the queue to start. Get more people!")
    return
  }

  // TODO: Make this a ready check style goroutine, maybe?
  for i := 1; i < len(q.pID); i++ {
    s.ChannelMessageSend(tilesID, "Hey <@" + q.pID[i] + ">!")
  }
  s.ChannelMessageSend(tilesID, "<@" + q.owner + "> has summoned you to start the hanchan!")

  q.started = false
  q.pID = q.pID[:0]
  q.pName = q.pName[:0]
  q.owner = ""
}

func tilesCheck(s *discordgo.Session, msg *discordgo.MessageCreate)  {
  if !q.started {
    s.ChannelMessageSend(msg.ChannelID, "<@" + msg.Author.ID + ">, there is no queue to check!  Type `!tiles start` to start one!")
    return
  }
  s.ChannelMessageSend(msg.ChannelID, "<@" + msg.Author.ID + ">, there are " + strconv.Itoa(len(q.pID)) + " players in the queue.")
  s.ChannelMessageSend(msg.ChannelID, "The owner of this queue is: " + q.pName[0])
  s.ChannelMessageSend(msg.ChannelID, "The following people are in the queue:")
  for i := 1; i < len(q.pID); i++ {
    s.ChannelMessageSend(msg.ChannelID, q.pName[i])
  }
}

func tilesClear(s *discordgo.Session, msg *discordgo.MessageCreate)  {
  if !q.started {
    s.ChannelMessageSend(msg.ChannelID, "<@" + msg.Author.ID + ">, there is no queue to clear!  Type `!tiles start` to start one!")
    return
  }
  s.ChannelMessageSend(msg.ChannelID, "<@" + msg.Author.ID + ">, the queue has been cleared at your request.")
  q.pID = q.pID[:0]
  q.pName = q.pName[:0]
  q.started = false
  q.owner = ""
}

func tilesHelp(s *discordgo.Session, msg *discordgo.MessageCreate)  {
  s.ChannelMessageSend(msg.ChannelID, "<@" + msg.Author.ID + ">" + "\n```Usage: !tiles [start|join|leave|play|check|clear|help]\n\n!tiles start - Start a new hanchan queue.\n!tiles join - Join the current queue.\n!tiles leave - Leave the current queue.\n!tiles play - Pings all queued players and clears the queue.\n!tiles check - Checks for a currently queuing hanchan.\n!tiles clear - Clears the current hanchan queue.\n!tiles help - This text.```")
}

func init()  {
  CmdList["tiles"] = tilesChooser
}
