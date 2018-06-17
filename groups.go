package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
	"log"
)

func joinGroup(s *discordgo.Session, msg *discordgo.MessageCreate, arg string) {
	group, _ := grabArg(arg)
	channel, _ := s.Channel(msg.ChannelID)
	guildID := channel.GuildID
	if group != ""{
		_, roleID := checkRoles(group)
    
		if roleID != ""{
			s.GuildMemberRoleAdd(guildID, msg.Author.ID, roleID)
			s.ChannelMessageSend(msg.ChannelID, "You have joined "+group+"!")
		}else{
			s.ChannelMessageSend(msg.ChannelID, "I can add you to the following groups: "+getRoleList())
		}
	}else{
		s.ChannelMessageSend(msg.ChannelID, "I can add you to the following groups: "+getRoleList())
	}
}

func leaveGroup(s *discordgo.Session, msg *discordgo.MessageCreate, arg string) {
	group, _ := grabArg(arg)
	channel, _ := s.Channel(msg.ChannelID)
	guildID := channel.GuildID
	if group != ""{
        _, roleID := checkRoles(group)
    
	    if roleID != ""{
	        s.GuildMemberRoleRemove(guildID, msg.Author.ID, roleID)
	        s.ChannelMessageSend(msg.ChannelID, "You have left "+group+"!")
	    }else{
	        s.ChannelMessageSend(msg.ChannelID, "I can remove you from the following groups: "+getRoleList())
	    }
	}else{
		s.ChannelMessageSend(msg.ChannelID, "I can add you to the following groups: "+getRoleList())
	}
}



func addGroup(s *discordgo.Session, msg *discordgo.MessageCreate, arg string) {
    group, fullName := grabArg(arg)
	//GuildRoleCreate to get a new guild, then GuildRoleEdit(guildID, roleID, name string, color int, hoist bool, perm int, mention bool) to name and stuff?
	if (fullName != "" && group != "" && (msg.Author.ID == "82987575245017088" || msg.Author.ID == "83742858800009216" || msg.Author.ID == "86255231913979904" || msg.Author.ID == "82621075220856832")){
		isRole,_ := checkRoles(group)
	    if !isRole{
	        channel, _ := s.Channel(msg.ChannelID)
	        guildID := channel.GuildID
	        newRole, err := s.GuildRoleCreate(guildID)
	
	        if err != nil {
		        fmt.Println(err)
	        }//99AAB5
	        _, err2 := s.GuildRoleEdit(guildID, newRole.ID, fullName, newRole.Color, false, 1109460032, true)
	        if err2 != nil {
		        fmt.Println(err2)
	        }
		    newItem := "INSERT INTO roles (role,title,id,uID) values (?,?,?,?)"
		    stmt, err := db.Prepare(newItem)
		    if err != nil {
		    	panic(err)
		    }
		    defer stmt.Close()
		    _, err3 := stmt.Exec(group,fullName,newRole.ID,msg.Author.ID)
            if err3 != nil {
                panic(err3)
            }
		    s.ChannelMessageSend(msg.ChannelID, "Added New Group: "+fullName+"\nYou are now free to join it using `!join "+group+"`")
	    }else{
		    s.ChannelMessageSend(msg.ChannelID, "Group already exists. Join it using !join "+group)
	    }
    }else{
	    s.ChannelMessageSend(msg.ChannelID, "Proper usage: `!addGroup role Role Title`\ni.e. `!addGroup cars Rocket Cars`\nPlease note that this command is restricted to bot admins")
	}
}

func checkRoles(arg string) (bool, string) {
	arg = strings.ToLower(arg)
    qte, err := db.Query("SELECT role, id FROM roles WHERE role = '" + arg+"'")
    if err != nil {
        log.Fatal("Query error CR:", err)
    }
    defer qte.Close()

	var role string
	var id string
	for qte.Next() {
		err = qte.Scan(&role,&id)
        if err != nil {
            log.Fatal("Parse error CR:", err)
        }
        if role != ""{
			return true, id
		}
    }
	return false, ""
}

func getRoleList() string{
	list := "```\n"
    qte, err := db.Query("SELECT title FROM roles")
    if err != nil {
        log.Fatal("Query error GRL:", err)
    }
    defer qte.Close()
	var role string
	for qte.Next(){
		err = qte.Scan(&role)
        if err != nil {
            log.Fatal("Parse error GRL:", err)
        }
        list = list+strings.Title(role)+"\n"
    }
	list = list+"```"
	return list
}


func init() {
	CmdList["addgroup"] = addGroup
	CmdList["join"] = joinGroup
	CmdList["leave"] = leaveGroup
}
