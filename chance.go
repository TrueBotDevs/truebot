package main

import (
    "github.com/bwmarrin/discordgo"
    "math/rand"
    "strconv"
    "time"
)

func rollDie(sides int) int {
    s1 := rand.NewSource(time.Now().UnixNano())
    r1 := rand.New(s1)
    result := r1.Intn(sides) + 1
    return result
}

func rollDice(numDice int, sides int) int {
    sum := 0
    for i := 0; i < numDice; i++ {
        sum += rollDie(sides)
    }
    return sum
}

func rollFunction(s *discordgo.Session, msg *discordgo.MessageCreate, arg string) {
    numDice := 1
    sides := 20
    total := rollDice(numDice, sides)
    s.ChannelMessageSend(msg.ChannelID, "Rolling "+strconv.Itoa(numDice)+"d"+strconv.Itoa(sides)+": `"+strconv.Itoa(total)+"`")
}

func init() {
    CmdList["roll"] = rollFunction
}
