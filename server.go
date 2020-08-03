package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
)

// BEGIN Structs
type User struct {
	UserName      string
	Discriminator string
	Token         string
	ImageURL      string
	Guilds        []UserGuild
}

type UserGuild struct {
	Name     string
	ID       string
	Enabled  bool
	Prefix   string
	Delay    int
	Channels []Channel
}

// ENDOF Structs

// BEGIN Handlers
var user *User

/* binExecute receives a set of parameters from the http request following the standard URL query
bin?key=value. It then parses the key and if it corresponds to any known key, will execute and return to
index page with a status of 302 (See https://en.wikipedia.org/wiki/List_of_HTTP_status_codes)
*/
func binExecute(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	if request.Form["command"] != nil {
		if request.FormValue("command") == "refresh" {
			refresh(client)
			updateConfigYaml()
			fmt.Println("Sucessfully refreshed")
		} else if request.Form["command"][0] == "changeGuildState" && request.Form["id"] != nil {
			id := request.FormValue("id")
			var changedGuildSlice []Guild
			for _, guild := range conf.Session.Guilds {
				if guild.ID == id {
					guild.Enabled = !guild.Enabled
					changedGuildSlice = append(changedGuildSlice, guild)
					fmt.Println("[GUILDS] Server with ID: " + id + "'s state has been changed sucessfully to: " + strconv.FormatBool(guild.Enabled))
				} else {
					changedGuildSlice = append(changedGuildSlice, guild)
				}
			}
			conf.Session.Guilds = changedGuildSlice
			updateConfigYaml()
		} else if request.FormValue("command") == "changeChannelState" && request.Form["guildid"] != nil && request.Form["channelid"] != nil {
			guildid := request.FormValue("guildid")
			channelid := request.FormValue("channelid")
			var changedGuildSlice []Guild
			for _, parent := range conf.Session.Guilds {

				if parent.ID == guildid {
					var changedChannelSlice []Channel
					for _, channel := range parent.Channels {
						if channel.ID == channelid {
							channel.Enabled = !channel.Enabled
							changedChannelSlice = append(changedChannelSlice, channel)
							fmt.Println("[CHANNELS] Channel with ID: " + channel.ID + "'s state was changed to: " + strconv.FormatBool(channel.Enabled))
						} else {
							changedChannelSlice = append(changedChannelSlice, channel)
						}
					}
					parent.Channels = changedChannelSlice
				}
				changedGuildSlice = append(changedGuildSlice, parent)
			}
			conf.Session.Guilds = changedGuildSlice
			updateConfigYaml()
		} else if request.Form["command"][0] == "changeGuildDelay" && request.Form["guildid"] != nil && request.Form["delay"] != nil {
			guildid := request.FormValue("guildid")
			delay, _ := strconv.Atoi(request.FormValue("delay"))
			var changedGuildSlice []Guild
			for _, guild := range conf.Session.Guilds {
				if guild.ID == guildid {
					guild.Delay = delay
					changedGuildSlice = append(changedGuildSlice, guild)
					fmt.Println("[GUILDS] Guild with ID: " + guildid + "'s delay was changed to: " + strconv.Itoa(delay) + " s")
				} else {
					changedGuildSlice = append(changedGuildSlice, guild)
				}
			}
			conf.Session.Guilds = changedGuildSlice
			updateConfigYaml()
		}
	}
	http.Redirect(writer, request, "/settings", 302)
}

func (user *User) updateGuilds() {
	var userGuilds []UserGuild
	for _, guild := range conf.Session.Guilds {
		var appendGuild UserGuild
		appendGuild.Name = guild.Name
		appendGuild.Enabled = guild.Enabled
		appendGuild.ID = guild.ID
		appendGuild.Prefix = guild.Prefix
		appendGuild.Delay = guild.Delay
		appendGuild.Channels = guild.Channels
		userGuilds = append(userGuilds, appendGuild)
	}
	user.Guilds = userGuilds
}

func (user *User) settingsHandler(writer http.ResponseWriter, request *http.Request) {
	user.updateGuilds()
	baseTemplate, _ := template.ParseFiles("./static/index.html")
	baseTemplate.Execute(writer, user)
}

func errorHandler(writer http.ResponseWriter, request *http.Request, status int) {
	writer.WriteHeader(status)
	fmt.Println("404")
	if status == http.StatusNotFound {
		fmt.Fprint(writer, "test")
	}
}

// Start works calling the ListenAndServe as a GoRoutine (non-blocking call).
func Start() {
	userInfo, _ := client.User("@me")
	var userGuilds []UserGuild
	for _, guild := range conf.Session.Guilds {
		var appendGuild UserGuild
		appendGuild.Name = guild.Name
		appendGuild.Enabled = guild.Enabled
		appendGuild.ID = guild.ID
		appendGuild.Prefix = guild.Prefix
		appendGuild.Delay = guild.Delay
		appendGuild.Channels = guild.Channels
		userGuilds = append(userGuilds, appendGuild)
	}
	user = &User{UserName: userInfo.Username, Discriminator: userInfo.Discriminator, Token: conf.Token, ImageURL: userInfo.AvatarURL(""), Guilds: userGuilds}
	http.Handle("/", http.FileServer(http.Dir("./static/")))
	http.HandleFunc("/settings", user.settingsHandler)
	http.HandleFunc("/bin", binExecute)
	err := http.ListenAndServe(":9898", nil)
	if err != nil {
		log.Fatal("Could not serve, error: ", err)
	}
}
