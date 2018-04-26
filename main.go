package main

import (
	"log"
	"os/user"
	"bytes"
	"os"
	"nagus/nagus"
	"encoding/json"
)

func main() {
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
	}
	pathBuffer.WriteRune(os.PathSeparator)
	pathBuffer.WriteString("config.json")

	if _, err = os.Stat(pathBuffer.String()); os.IsNotExist(err) {
		var conf *nagus.NagusConfig = &nagus.NagusConfig{
			ApiKey: "insert-api-key-here",
		}
		jsonBuffer, err := json.Marshal(conf)
		if err != nil {
			log.Println("Could not initialise config!!!")
			log.Panic(err.Error())
			return
		}

		file, err := os.Create(pathBuffer.String())
		if err != nil {
			log.Println("Could not initialise config!!!")
			log.Panic(err.Error())
			return
		}
		file.WriteString(string(jsonBuffer))
		file.Close()
	}

	log.Printf("Reading config from %s", pathBuffer.String())
	conf, err := nagus.ReadConfig(pathBuffer.String())

	if err != nil {
		log.Panic(err)
		return
	}

	log.Printf("Using API Key %s", conf.ApiKey)

	bot, err := nagus.BuildBot(conf)
	if err != nil {
		log.Panic(err)
		return
	}

	err = bot.Main()
	if err != nil {
		log.Panic(err)
	}
}
