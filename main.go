package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"time"

	"github.com/SilverCory/ConStats/speedtest"
	"github.com/SilverCory/ConStats/sql"
)

// Configuration - The main configuration for ConStats
type Configuration struct {
	Command string
	Args    []string
	MySQL   *MySQLConfiguration
}

// MySQLConfiguration - MySQL part..
type MySQLConfiguration struct {
	Host string
}

var defaultConfiguration = &Configuration{
	Command: "./speedtest",
	Args:    []string{"--json", "--secure"},
	MySQL: &MySQLConfiguration{
		Host: "user:password@/dbname",
	},
}

func main() {

	// TODO load config
	CurrentConfig := defaultConfiguration

	if _, err := os.Stat("./config.json"); os.IsNotExist(err) {

		data, err := json.MarshalIndent(defaultConfiguration, "", "\t")

		if err != nil {
			fmt.Println("There was an error saving the default config!", err)
			os.Exit(1)
		}

		err = ioutil.WriteFile("./config.json", data, 0644)
		if err != nil {
			fmt.Println("There was an error saving the default config!", err)
			os.Exit(1)
		}

		fmt.Println("The default configuration was saved. Please edit this!")
		os.Exit(0)

	} else {

		data, err := ioutil.ReadFile("./config.json")
		if err != nil {
			fmt.Println("There was an error loading the config!", err)
			os.Exit(1)
		}

		err = json.Unmarshal(data, &CurrentConfig)
		if err != nil {
			fmt.Println("There was an error loading the config!", err)
			os.Exit(1)
		}

	}

	mysqlStorage := sql.Create()
	mysqlStorage.Host = CurrentConfig.MySQL.Host

	speed := speedtest.Create()
	speed.Args = CurrentConfig.Args
	speed.Command = CurrentConfig.Command

	doTest(speed, mysqlStorage)

	ticker := time.NewTicker(15 * time.Minute)

	// Infinite loop..
	for {
		select {
		case <-ticker.C:
			doTest(speed, mysqlStorage)
		}
	}

}

func doTest(speed *speedtest.SpeedTest, storage *sql.MySQL) {

	fmt.Println("Starting speed test...")

	result, err := speed.Test()
	if err != nil {
		fmt.Println("An error occured: ", err)
		storage.Save(nil, nil)
		return
	}

	runTime, err := time.Parse("2006-01-02T15:04:05.999999999", result.TimeStamp)
	if err != nil {
		fmt.Println("An error occured: ", err)
		storage.Save(nil, nil)
		return
	}

	fmt.Printf("Time        : %s\n", runTime.Format(time.RFC1123Z))
	fmt.Printf("Ping        : %.2f\n", result.Ping)
	fmt.Printf("Upload      : %.0f\n", result.Upload)
	fmt.Printf("Download    : %.0f\n", result.Download)

	fmt.Println("=================== Saving to SQL ===================")
	err = storage.Save(result, &runTime)
	if err != nil {
		fmt.Println("An error occured: ", err)
	} else {
		fmt.Println("Saved successfully!")
	}
}
