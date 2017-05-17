package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/SilverCory/constats/speedtest"
	"github.com/SilverCory/constats/sql"
	"github.com/SilverCory/constats/web"
	"log"
)

// Configuration - The main configuration for ConStats
type Configuration struct {
	MyTable          string
	Command          string
	Args             []string
	IntervalMinuites int
	MySQL            *MySQLConfiguration
	WebData          *WebDataConfiguration
}

type WebDataConfiguration struct {
	CreateWebData  bool
	RunServer      bool
	Host           string
	FetchAllTables bool
}

// MySQLConfiguration - MySQL part..
type MySQLConfiguration struct {
	Host string
}

var defaultConfiguration = &Configuration{
	MyTable:          "my-computer-home",
	Command:          "./speedtest",
	Args:             []string{"--json", "--secure"},
	IntervalMinuites: 15,
	MySQL: &MySQLConfiguration{
		Host: "user:password@/dbname",
	},
	WebData: &WebDataConfiguration{
		CreateWebData:  true,
		RunServer:      true,
		Host:           ":8080",
		FetchAllTables: false,
	},
}

func main() {

	CurrentConfig := defaultConfiguration
	checkConfig(CurrentConfig)

	mysqlStorage := sql.Create()
	mysqlStorage.Host = CurrentConfig.MySQL.Host

	speed := speedtest.Create()
	speed.Args = CurrentConfig.Args
	speed.Command = CurrentConfig.Command

	if CurrentConfig.WebData.RunServer {
		go web.RunWebserver(CurrentConfig.WebData.Host)
	}

	doTest(speed, mysqlStorage, CurrentConfig.MyTable)
	doData(CurrentConfig.WebData.CreateWebData, mysqlStorage, CurrentConfig.MyTable, CurrentConfig.WebData.FetchAllTables)

	ticker := time.NewTicker(time.Duration(CurrentConfig.IntervalMinuites) * time.Minute)

	// Infinite loop..
	for {
		select {
		case <-ticker.C:
			doTest(speed, mysqlStorage, CurrentConfig.MyTable)
			doData(CurrentConfig.WebData.CreateWebData, mysqlStorage, CurrentConfig.MyTable, CurrentConfig.WebData.FetchAllTables) // Yes, this makes the data every {interval} mins.. 3000 each go. Deal with it.
		}
	}

}

func checkConfig(currentConfig *Configuration) {
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

		err = json.Unmarshal(data, &currentConfig)
		if err != nil {
			fmt.Println("There was an error loading the config!", err)
			os.Exit(1)
		}

	}

	currentConfig.MyTable = "constats_" + currentConfig.MyTable

}

func doData(createData bool, storage *sql.MySQL, table string, fetchAll bool) {

	if !createData {
		return
	}

	var (
		createTables []string
		err          error
	)

	if fetchAll {
		createTables, err = storage.FindTables()
		if err != nil {
			fmt.Println("Unable to find tables!", err)
			return
		}
	} else {
		createTables = make([]string, 1)
		createTables[0] = table
	}

	for _, table = range createTables {

		start := time.Now()

		data, err := web.GenerateData(storage, table)
		if err != nil {
			fmt.Println("There was an error generating webdata for table "+table+"!", err)
			continue
		}

		fileData, err := json.MarshalIndent(data, "", "\t")
		if err != nil {
			fmt.Println("There was an error generating webdata for table "+table+"!", err)
			continue
		}

		err = ioutil.WriteFile("./data/connectionData_"+table+".json", fileData, 0644)
		if err != nil {
			fmt.Println("There was an error generating webdata for table "+table+"!", err)
			continue
		}

		elapsed := time.Since(start)
		log.Printf("Data generation took %s for %q", elapsed, table)

	}
}

func doTest(speed *speedtest.SpeedTest, storage *sql.MySQL, table string) {

	fmt.Println("Starting speed test...")

	result, err := speed.Test()
	if err != nil {
		fmt.Println("An error occured: ", err)
		storage.Save(nil, nil, table)
		return
	}

	runTime, err := time.Parse(time.RFC3339, result.TimeStamp)
	if err != nil {
		fmt.Println("An error occured: ", err)
		storage.Save(nil, nil, table)
		return
	}

	fmt.Printf("Time        : %s\n", runTime.Format(time.RFC1123Z))
	fmt.Printf("Ping        : %.2f\n", result.Ping)
	fmt.Printf("Upload      : %.0f\n", result.Upload)
	fmt.Printf("Download    : %.0f\n", result.Download)

	fmt.Println("=================== Saving to SQL ===================")
	err = storage.Save(result, &runTime, table)
	if err != nil {
		fmt.Println("An error occured: ", err)
	} else {
		fmt.Println("Saved successfully!")
	}

}
