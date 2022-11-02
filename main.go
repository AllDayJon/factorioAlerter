package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	// Get the environment variables
	getEnv()
	// Change the working directory
	changeDir()
	// Check if the log file exists
	checkForLog()
	forever()
}

func getEnv() {
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		fmt.Println(pair[0], pair[1])
	}
	if _, ok := os.LookupEnv("DISCORD_WEBHOOK"); ok {
		fmt.Println("DISCORD_WEBHOOK exists.. Continuing")
	} else {
		fmt.Println("DISCORD_WEBHOOK does not exist.. Exiting")
		os.Exit(1)
	}
}

// Get the environment variable
func getDiscordWebhook() string {
	return os.Getenv("DISCORD_WEBHOOK")
}

// Change the working directory to /app/factorio/
func changeDir() {
	// Check if the directory exists
	if _, err := os.Stat("factorio"); os.IsNotExist(err) {
		fmt.Println("Directory does not exist.. Exiting")
		os.Exit(1)
	} else {
		fmt.Println("Directory exists.. Setting new CWD")
		err := os.Chdir("factorio")
		if err != nil {
			return
		}
	}
}

func checkForLog() string {
	// Check if the file exists
	if _, err := os.Stat("factorio-current.log"); os.IsNotExist(err) {
		fmt.Println("factorio-current.log does not exist.. Exiting")
		os.Exit(1)
	} else {
		fmt.Println("factorio-current.log exists.. Continuing")

		// Open the file
		file, err := os.Open("factorio-current.log")
		if err != nil {
			return "Error opening file"
		} else {
			fmt.Println("factorio-current.log opened.. Continuing")
			scanner := bufio.NewScanner(file)
			var joined string
			for scanner.Scan() {
				line := scanner.Text()
				if strings.Contains(scanner.Text(), "joined the game") {
					if strings.Contains(scanner.Text(), time.Now().Format("15:04:05")) {
						joined += strings.Split(line, " ")[3] + " has joined the game." + "\\n"
						fmt.Println("Time difference: ", time.Now().Sub(time.Now()))
						fmt.Println(joined)

					}

				}
			}
			//close the file
			err := file.Close()
			if err != nil {
				return "Error closing file"
			}
			if joined == "" {
				fmt.Println("No new players have joined the game.")
			} else {
				return joined
			}

		}

	}

	return "Error"
}

func sendMessageToDiscord(message string) {
	if message == "" || message == "Error" {
		fmt.Println("Message is empty.. Not sending to discord")
	} else {
		fmt.Println("Message is not empty.. Continuing")
		req, err := http.NewRequest("POST", getDiscordWebhook(), strings.NewReader(`{"content": "`+message+`"}`))
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				return
			}
		}(resp.Body)

		// Print the response to console
		fmt.Println("response Status:", resp.Status)
		fmt.Println("response Headers:", resp.Header)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println()
		fmt.Println("response Body:", string(body))
		fmt.Println()
		fmt.Println("request Body:", message)
	}
}

func forever() {
	for {
		sendMessageToDiscord(checkForLog())
		fmt.Println("Sleeping for 10 seconds..")
		time.Sleep(10 * time.Second)

	}
}
