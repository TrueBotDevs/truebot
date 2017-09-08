package main

import(
    "fmt"
    "log"
    "time"
    "strings"
    "strconv"
	"github.com/bwmarrin/discordgo"
)

var(
    //serverId = "82683153025601536"
    channelId = "355504765305618432"
)

func parseDate(date string) (string, time.Duration){
    compString := "weeks days hours minutes seconds"
    lookingForDates := true
    dateArgs := strings.Split(date," ")
    dateIndex := 0
    var parsedDuration time.Duration
    if strings.HasPrefix(date, "in "){
        date = strings.Replace(date, "in ", "", 1)
    }
    for lookingForDates {
        if dateIndex >= len(dateArgs)-1{
            lookingForDates = false
            break
        }
        timeInt := strings.Split(date," ")[dateIndex:dateIndex+1][0]
        timeStr := strings.Split(date," ")[dateIndex+1:dateIndex+2][0]
        convertedInt, err := strconv.ParseInt(timeInt,10,32); 
        if err != nil{
            lookingForDates = false
            break
        }
        if strings.Contains(compString,timeStr) == false{
            lookingForDates = false
            break
        }
        fmt.Println(timeInt + " " + timeStr)
        dateIndex += 2
        if strings.Contains("seconds",timeStr){
            parsedDuration += time.Duration(convertedInt)*time.Second
        }else if strings.Contains("days",timeStr){
            parsedDuration += time.Duration(convertedInt*24)*time.Hour
        }else if strings.Contains("hours",timeStr){
            parsedDuration += time.Duration(convertedInt)*time.Hour
        }else if strings.Contains("minutes",timeStr){
            parsedDuration += time.Duration(convertedInt)*time.Minute
        }else if strings.Contains("weeks",timeStr){
            parsedDuration += time.Duration(convertedInt*24*7)*time.Hour
        }else if strings.Contains("years",timeStr){
            parsedDuration += time.Duration(convertedInt*24*365)*time.Hour
        }
    }
    if dateIndex < len(dateArgs) {
        return strings.Join(strings.Split(date," ")[dateIndex:]," "), parsedDuration
    } else {
        return " ", parsedDuration
    }    
}
func addReminder(s *discordgo.Session, msg *discordgo.MessageCreate, arg string){
    remainderMsg, timeToWait := parseDate(arg)
    
    newItem := "INSERT INTO reminders (reminder,date,userId) values (?,?,?)"
    stmt, err := db.Prepare(newItem)
    if err != nil { panic(err) }
    defer stmt.Close()

    remindDate := time.Now().Add(timeToWait)
    _, err2 := stmt.Exec(remainderMsg,remindDate.Unix(),msg.Author.ID)
    if err2 != nil { panic(err2) }
    
	s.ChannelMessageSend(msg.ChannelID, "Ok, <@" + msg.Author.ID + ">, I will remind you in " + timeToWait.String() + "```" + remainderMsg + "```")
}
func doRemind(){
    for true{
        if hasSession{               
            currentTime := strconv.FormatInt(time.Now().Unix(),10)
            query := "SELECT userId, reminder, reminderId FROM reminders WHERE date <= " + currentTime + " AND isDone = 0;"
            qte, err := db.Query(query)
            if err != nil {
                dgSession.ChannelMessageSend(channelId, err.Error())
                log.Fatal("Query error:", err)
            }
            defer qte.Close()
            
            var rem string
            var uId int
            var reminderId int
            var reminders [10000]string
            var reminderIds [10000]int
            var index = 0
            for qte.Next(){
                err = qte.Scan(&uId, &rem, &reminderId)
                if err != nil {
                    dgSession.ChannelMessageSend(channelId, "Shit's really fucked")
                    log.Fatal("Parse error:", err)
                }
                dgSession.ChannelMessageSend(channelId, "Hey <@" + strconv.Itoa(uId) + ">```" + rem + "```")
                fmt.Println(rem)
                reminders[index] = rem
                reminderIds[index] = reminderId
                index++
            }
            for i := 0; i < index; i++{
                deleteCmd := "UPDATE reminders SET isDone = 1 WHERE reminderId = ?"
                stmt, err2 := db.Prepare(deleteCmd)
                if err2 != nil { 
                    dgSession.ChannelMessageSend(channelId, "Shit's kinda fucked")
                    panic(err2) }
                defer stmt.Close()

                _, err3 := stmt.Exec(reminderIds[i])
                if err3 != nil { panic(err3)
                    dgSession.ChannelMessageSend(channelId, "Shit's super fucked")
                }
                time.Sleep(1000 * time.Millisecond)
            }
        }
    }
}
func init() {
    fmt.Println("Don't forget to register on site for SGDQ 2017!")
	CmdList["remindme"] = addReminder
    go doRemind()
}