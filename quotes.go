package main

import(
    "fmt"
    "log"
    "math/rand"
    "time"
    "strings"
    "github.com/bwmarrin/discordgo"
) 

func getQuote(s *discordgo.Session, msg *discordgo.MessageCreate, comp string){
    qte, err := db.Query("SELECT quote FROM quotes WHERE quote LIKE \"%"+comp+"%\"")
    if err != nil {
		log.Fatal("Query error:", err)
	}
    defer qte.Close()
    
    var quote string
    var quotes [255]string
    var index = 0
    var newIndex = 1
    for qte.Next(){
        err = qte.Scan(&quote)
        if err != nil {
            log.Fatal("Parse error:", err)
        }
        quotes[index] = quote
        index++
    }
    s1 := rand.NewSource(time.Now().UnixNano())
    r1 := rand.New(s1)
    if index == 0{
        getQuote(s,msg," ")
        return
    }
    newIndex = r1.Intn(index)
    //fmt.Println(quotes)
    //fmt.Println(quotes[newIndex])
    s.ChannelMessageSend(msg.ChannelID,quotes[newIndex])
    fmt.Println(quotes[newIndex])
}

func addQuote(quote string) string{
    if strings.Contains(quote, "<@"){
        return "Fuck you, don't @ people in quotes"
    }else{
        newItem := "INSERT INTO quotes (quote) values (?)"
        stmt, err := db.Prepare(newItem)
        if err != nil { panic(err) }
        defer stmt.Close()

        _, err2 := stmt.Exec(quote)
        if err2 != nil { panic(err2) }
        return quote + " Added to the database"
    }
}

func init() {
    CmdList["quote"] = getQuote
    CmdList["addquote"] = addQuote
}