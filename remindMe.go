package main

import(
    "fmt"
    "time"
    "strings"
    "strconv"
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
        }
    }
    if dateIndex < len(dateArgs) {
        return strings.Join(strings.Split(date," ")[dateIndex:]," "), parsedDuration
    } else {
        return " ", parsedDuration
    }    
}
func addReminder(){
    
}
func init() {
    fmt.Println("Don't forget to register on site for SGDQ 2017!")
}