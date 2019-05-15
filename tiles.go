package main

import (
    "github.com/bwmarrin/discordgo"
    "strconv"
)

var (
    tilesID  = "270331769033588737"
    vTilesID = "231272263389806602"
    gTilesID = "276514629700681728"
    q        = Queue{}
)

//Queue for holding a list of players
type Queue struct {
    owner   string
    started bool
    pName   []string
    pID     []string
}

func tilesChooser(s *discordgo.Session, msg *discordgo.MessageCreate, arg string) {

    switch arg {
    case "start":
        tilesStart(s, msg)
    case "join":
        tilesJoin(s, msg)
    case "leave":
        tilesLeave(s, msg)
    case "play":
        tilesPlay(s, msg)
    case "check":
        tilesCheck(s, msg)
    case "clear":
        tilesClear(s, msg)
    case "1pin":
        s.ChannelMessageSend(msg.ChannelID, "私の一筒!")
    case "[start|join|leave|play|check|clear|help]":
        s.ChannelMessageSend(msg.ChannelID, "<@"+msg.Author.ID+">, actively go fuck yourself. Especially if you're Church.")
    case "debug":
        tilesDebug(s, msg)
    default:
        tilesHelp(s, msg)
    }
    tilesChannelUpdate(s)
}

func tilesStart(s *discordgo.Session, msg *discordgo.MessageCreate) {
    if q.started {
        s.ChannelMessageSend(msg.ChannelID, "<@"+msg.Author.ID+">, there's already a queue in progress!  Use `!tiles check` to see who's in the queue.")
        return
    }
    s.ChannelMessageSend(tilesID, "Hey <@&"+gTilesID+">! "+"<@"+msg.Author.ID+"> wants to start a hanchan! Type `!tiles join` to enter the queue.")

    q.started = true
    q.pName = append(q.pName, msg.Author.Username)
    q.pID = append(q.pID, msg.Author.ID)
    q.owner = msg.Author.ID
}

func tilesJoin(s *discordgo.Session, msg *discordgo.MessageCreate) {
    if !q.started {
        s.ChannelMessageSend(msg.ChannelID, "<@"+msg.Author.ID+">, there is no queue to join!  Type `!tiles start` to start one!")
        return
    }
    if q.owner == msg.Author.ID {
        s.ChannelMessageSend(msg.ChannelID, "<@"+msg.Author.ID+">, you made the queue so you don't need to `!tiles join`, you're already in the queue.")
        return
    }
    for i := 1; i < len(q.pID); i++ {
        if msg.Author.ID == q.pID[i] {
            s.ChannelMessageSend(msg.ChannelID, "<@"+msg.Author.ID+">, you are already in the queue!")
            return
        }
    }
    if len(q.pID) == 4 {
        s.ChannelMessageSend(msg.ChannelID, "<@"+msg.Author.ID+">, the play queue is currently maxed out right now.  Not adding player.")
        return
    }
    q.pName = append(q.pName, msg.Author.Username)
    q.pID = append(q.pID, msg.Author.ID)
    s.ChannelMessageSend(msg.ChannelID, "<@"+msg.Author.ID+">, you've been added to the queue!  There are "+strconv.Itoa(len(q.pID))+" players in the queue.")

    if len(q.pID) == 4 {
        s.ChannelMessageSend(tilesID, "<@"+q.owner+">, you have enough players in queue to play a hanchan.  Use `!tiles play` to ping players and clear the queue.")
    }
}

func tilesLeave(s *discordgo.Session, msg *discordgo.MessageCreate) {
    if !q.started {
        s.ChannelMessageSend(msg.ChannelID, "<@"+msg.Author.ID+">, there is no queue to leave! Type `!tiles start` to start one!")
        return
    }
    if q.owner == msg.Author.ID {
        // Sanity check to make sure that the owner is actually the first person in the queue's slice.  If they aren't, debug and clear the queue.
        if msg.Author.ID != q.pID[0] {
            s.ChannelMessageSend(msg.ChannelID, "Something went wrong.  This should never happen. Clearing queue.")
            tilesDebug(s, msg)
            q.pID = q.pID[:0]
            q.pName = q.pName[:0]
            q.started = false
            q.owner = ""
            s.ChannelMessageSend(msg.ChannelID, "Queue cleared.")
            return
        }
        // Next, check to see if the queue has only one player.  If it has only one, just clear the queue.
        if len(q.pID) == 1 {
            s.ChannelMessageSend(msg.ChannelID, "<@"+msg.Author.ID+">,  you are the sole player in the queue.  Clearing the queue.")
            tilesClear(s, msg)
            return
        }
        // If neither of the above, remove the owner from the queue.
        for i := 0; i < len(q.pID); i++ {
            q.pID = append(q.pID[:i], q.pID[i+1:]...)
            q.pName = append(q.pName[:i], q.pName[i+1:]...)
        }
        q.owner = q.pID[0]
        s.ChannelMessageSend(msg.ChannelID, "<@"+msg.Author.ID+">, you have been removed from the queue.  The new owner is <@"+q.owner+">.")
    }
    for i := 1; i < len(q.pID); i++ {
        if msg.Author.ID == q.pID[i] {
            q.pID = append(q.pID[:i], q.pID[i+1:]...)
            q.pName = append(q.pName[:i], q.pName[i+1:]...)
            s.ChannelMessageSend(msg.ChannelID, "<@"+msg.Author.ID+">, you have been removed from the queue.")
        }
    }
}

func tilesPlay(s *discordgo.Session, msg *discordgo.MessageCreate) {
    if !q.started {
        s.ChannelMessageSend(msg.ChannelID, "<@"+msg.Author.ID+">, there is no queue to start playing!  Type `!tiles start` to start one!")
        return
    }
    if msg.Author.ID != q.owner {
        s.ChannelMessageSend(msg.ChannelID, "<@"+msg.Author.ID+">, you are not the owner of this queue. Ignoring.")
        return
    }
    if len(q.pID) < 4 {
        s.ChannelMessageSend(msg.ChannelID, "<@"+msg.Author.ID+">, there aren't enough people in the queue to start. Get more people!")
        return
    }

    // TODO: Make this a ready check style goroutine, maybe?
    for i := 1; i < len(q.pID); i++ {
        s.ChannelMessageSend(tilesID, "Hey <@"+q.pID[i]+">!")
    }
    s.ChannelMessageSend(tilesID, "<@"+q.owner+"> has summoned you to start the hanchan!")

    q.started = false
    q.pID = q.pID[:0]
    q.pName = q.pName[:0]
    q.owner = ""
}

func tilesCheck(s *discordgo.Session, msg *discordgo.MessageCreate) {
    if !q.started {
        s.ChannelMessageSend(msg.ChannelID, "<@"+msg.Author.ID+">, there is no queue to check!  Type `!tiles start` to start one!")
        return
    }
    s.ChannelMessageSend(msg.ChannelID, "<@"+msg.Author.ID+">, there are "+strconv.Itoa(len(q.pID))+" players in the queue.")
    s.ChannelMessageSend(msg.ChannelID, "The owner of this queue is: "+q.pName[0])
    s.ChannelMessageSend(msg.ChannelID, "The following people are in the queue:")
    for i := 1; i < len(q.pID); i++ {
        s.ChannelMessageSend(msg.ChannelID, q.pName[i])
    }
}

func tilesClear(s *discordgo.Session, msg *discordgo.MessageCreate) {
    if !q.started {
        s.ChannelMessageSend(msg.ChannelID, "<@"+msg.Author.ID+">, there is no queue to clear!  Type `!tiles start` to start one!")
        return
    }
    s.ChannelMessageSend(msg.ChannelID, "<@"+msg.Author.ID+">, the queue has been cleared at your request.")
    q.pID = q.pID[:0]
    q.pName = q.pName[:0]
    q.started = false
    q.owner = ""
}

func tilesDebug(s *discordgo.Session, msg *discordgo.MessageCreate) {
    s.ChannelMessageSend(msg.ChannelID, "Logging current queue.")
    s.ChannelMessageSend(msg.ChannelID, "Player IDs queued:")
    for i := 0; i < len(q.pID); i++ {
        s.ChannelMessageSend(msg.ChannelID, q.pID[i])
    }
    s.ChannelMessageSend(msg.ChannelID, "Player usernames queued:")
    for i := 0; i < len(q.pName); i++ {
        s.ChannelMessageSend(msg.ChannelID, q.pName[i])
    }
    s.ChannelMessageSend(msg.ChannelID, "Status of queue: "+strconv.FormatBool(q.started)+". Owner: "+q.owner)
    s.ChannelMessageSend(msg.ChannelID, "Logging done.")
}

func tilesHelp(s *discordgo.Session, msg *discordgo.MessageCreate) {
    s.ChannelMessageSend(msg.ChannelID, "<@"+msg.Author.ID+">"+"\n```Usage: !tiles [start|join|leave|play|check|clear|help]\n\n!tiles start - Start a new hanchan queue.\n!tiles join - Join the current queue.\n!tiles leave - Leave the current queue.\n!tiles play - Pings all queued players and clears the queue.\n!tiles check - Checks for a currently queuing hanchan.\n!tiles clear - Clears the current hanchan queue.\n!tiles help - This text.```")
}

func tilesChannelUpdate(s *discordgo.Session) {
    return //This feature is being removed as the tiles nerds said they didn't need it and that voice channel goes largely unused
    switch len(q.pID) {
    case 0:
        s.ChannelEdit(vTilesID, "Meme Tiles")
    case 1:
        s.ChannelEdit(vTilesID, "🀄 Meme Tiles")
    case 2:
        s.ChannelEdit(vTilesID, "🀄🀄 Meme Tiles")
    case 3:
        s.ChannelEdit(vTilesID, "🀄🀄🀄 Meme Tiles")
    default:
        s.ChannelEdit(vTilesID, "🀄🀄🀄🀄 Meme Tiles")
    }
}

func init() {
    CmdList["tiles"] = tilesChooser
    //CmdList["Tiles"] = tilesChooser
}
