package main

import (
	"log"
	"os/user"
	"bytes"
	"os"
	"nagus/nagus"
	"encoding/json"
)

func GetConfigFolder() string {
	usr, err := user.Current()
	if err != nil {
		log.Panic(err)
	}
	var pathBuffer bytes.Buffer
	pathBuffer.WriteString(usr.HomeDir)
	pathBuffer.WriteRune(os.PathSeparator)
	pathBuffer.WriteString(".nagus")
	if _, err = os.Stat(pathBuffer.String()); os.IsNotExist(err) {
		os.Mkdir(pathBuffer.String(), os.ModePerm)
		imageConf := bytes.NewBufferString(pathBuffer.String())
		imageConf.WriteRune(os.PathSeparator)
		imageConf.WriteString("images")
		os.Mkdir(pathBuffer.String(), os.ModePerm)
	}
	return pathBuffer.String()
}

func main() {
	pathBuffer := bytes.NewBufferString(GetConfigFolder())
	pathBuffer.WriteRune(os.PathSeparator)
	confBuffer := bytes.NewBufferString(pathBuffer.String())
	confBuffer.WriteString("config.json")

	if _, err := os.Stat(confBuffer.String()); os.IsNotExist(err) {
		var conf *nagus.NagusConfig = &nagus.NagusConfig{
			ApiKey: "insert-api-key-here",
		}
		jsonBuffer, err := json.Marshal(conf)
		if err != nil {
			log.Println("Could not initialise config!!!")
			log.Panic(err.Error())
			return
		}

		file, err := os.Create(confBuffer.String())
		if err != nil {
			log.Println("Could not initialise config!!!")
			log.Panic(err.Error())
			return
		}
		file.WriteString(string(jsonBuffer))
		file.Close()
	}

	log.Printf("Reading config from %s", confBuffer.String())
	conf, err := nagus.ReadConfig(confBuffer.String())

	if err != nil {
		log.Panic(err)
		return
	}

	log.Printf("Using API Key %s", conf.ApiKey)

	bot, err := nagus.BuildBot(conf, pathBuffer.String())
	if err != nil {
		log.Panic(err)
		return
	}

	err = bot.Main()
	if err != nil {
		log.Panic(err)
	}
}
