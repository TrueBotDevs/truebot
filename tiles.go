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
  sName []string
  pID []string
  sID []string
}

// TODO: This should really be a switch statement, but fuck it.  I'll change it later.
func tilesChooser(s *discordgo.Session, msg *discordgo.MessageCreate, arg string)  {
  if arg == "start" {
    tilesStart(s, msg)
  } else if arg == "join" {
    tilesJoin(s, msg)
  } else if arg == "standby" {
    tilesStandby(s, msg)
  } else if arg == "play" {
    tilesPlay(s, msg)
  } else if arg == "check" {
    tilesCheck(s, msg)
  } else if arg == "clear" {
    tilesClear(s, msg)
  } else {
    tilesHelp(s, msg)
  }
}

func tilesStart(s *discordgo.Session, msg *discordgo.MessageCreate)  {
  if q.started {
    s.ChannelMessageSend(msg.ChannelID, "<@" + msg.Author.ID + ">" + ", there is currently a queue already started.  Use `!tiles check` to see who's in the queue.")
    return
  }
  s.ChannelMessageSend(tilesID, "Hey <@&276514629700681728>! " + "<@" + msg.Author.ID + ">" + " wants to start a hanchan! Type `!tiles join` to enter the queue.")

  q.started = true
  q.pName = append(q.pName, msg.Author.Username)
  q.pID = append(q.pID, msg.Author.ID)
}

func tilesJoin(s *discordgo.Session, msg *discordgo.MessageCreate)  {
  if !q.started {
    s.ChannelMessageSend(msg.ChannelID, "<@" + msg.Author.ID + ">" + ", there is no queue right now.  Type `!tiles start` to start one!")
    return
  }
  s.ChannelMessageSend(msg.ChannelID, "<@" + msg.Author.ID + ">" + ", you've been added to the queue!  Currently we have " + strconv.Itoa(len(q.pID)) + " players in the queue, and " + strconv.Itoa(len(q.sID)) + " standby players.")
  q.pName = append(q.pName, msg.Author.Username)
  q.pID = append(q.pID, msg.Author.ID)

  if len(q.pID) + len(q.sID) >= 4 {
    s.ChannelMessageSend(tilesID, "<@" + q.owner + ">" + ", you have enough players in queue and standby to play a hanchan.  Use `!tiles play` to ping players and clear the queue.")
  }
}

func tilesStandby(s *discordgo.Session, msg *discordgo.MessageCreate)  {
  if !q.started {
    s.ChannelMessageSend(msg.ChannelID, "<@" + msg.Author.ID + ">" + ", there is no queue right now.  Type `!tiles start` to start one!")
    return
  }
  s.ChannelMessageSend(msg.ChannelID, "<@" + msg.Author.ID + ">" + ", you've been added to the standby queue!  Currently we have " + strconv.Itoa(len(q.pID)) + " players in the queue, and " + strconv.Itoa(len(q.sID)) + " standby players.")
  q.sName = append(q.sName, msg.Author.Username)
  q.sID = append(q.sID, msg.Author.ID)

  if len(q.pID) + len(q.sID) >= 4 {
    s.ChannelMessageSend(tilesID, "<@" + q.owner + ">" + ", you have enough players in queue and standby to play a hanchan.  Use `!tiles play` to ping players and clear the queue.")
  }
}

func tilesPlay(s *discordgo.Session, msg *discordgo.MessageCreate)  {
  if msg.Author.ID != q.owner {
    s.ChannelMessageSend(msg.ChannelID, "<@" + msg.Author.ID ">, you are not the owner of this queue. Ignoring.")
  }
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
    s.ChannelMessageSend(tilesID, "Hey <@" + q.pID[i] + ">!")
  }
  if len(q.pID) < 4 {
    need := 4 - len(q.pID)
    for i := 0; i <= need; i++ {
      s.ChannelMessageSend(tilesID, "Hey <@" + q.sID[i] + ">!")
    }
  }
  s.ChannelMessageSend(tilesID, "<@" + q.owner + "> has summoned you to start the hanchan!")

  q.started false
  q.pID = q.pID[:0]
  q.sID = q.sID[:0]
}

func tilesCheck(s *discordgo.Session, msg *discordgo.MessageCreate)  {
  s.ChannelMessageSend(msg.ChannelID, "<@" + msg.Author.ID + ">, there are currently " + strconv.Itoa(len(q.pID)) + " players in the queue with " + strconv.Itoa(len(q.sID)) + " standbys.")
}

func tilesClear(s *discordgo.Session, msg *discordgo.MessageCreate)  {
  if msg.Author.ID != q.owner {
    s.ChannelMessageSend(msg.ChannelID, "<@" + msg.Author.ID + ">, you are not the owner of this queue. Ignoring.")
    return
  }
  s.ChannelMessageSend(msg.ChannelID, "<@" + msg.Author.ID + ">, the queue has been cleared at your request.")
  q.pID = q.pID[:0]
  q.sID = q.sID[:0]
}

func tilesHelp(s *discordgo.Session, msg *discordgo.MessageCreate)  {
  s.ChannelMessageSend(msg.ChannelID, "<@" + msg.Author.ID + ">" + "\n```Usage: !tiles [start|join|standby|play|check|clear|help]\n\n!tiles start - Start a new hanchan queue.\n!tiles join - Join a hanchan queue.\n!tiles standby - Join a hanchan queue as a standby.\n!tiles play - Pings all queued players and clears the queue.\n!tiles check - Checks for a currently queuing hanchan.\n!tiles clear - Clears the current hanchan queue. (game owner only)\n!tiles help - This text.```")
}

func init()  {
  CmdList["tiles"] = tilesChooser
}
