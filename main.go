package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/ChrisR/pokedex/internal/pokecache"
)

type Config struct {
	next  string
	prev  string
	cache pokecache.Cache
}

type cliCommand struct {
	name        string
	description string
	callback    func(*Config) error
}

func main() {
	var interval time.Duration = 5

	commandMap := getCommands()
	scanner := bufio.NewScanner(os.Stdin)
	config := &Config{next: "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20", cache: pokecache.NewCache(interval)}

	for {
		fmt.Print("\nPokedex > ")
		scanner.Scan()

		commands := cleanInput(scanner.Text())
		userCommand, ok := commandMap[commands[0]]
		//fmt.Print("Your command was: ", commands[0])

		if ok {
			userCommand.callback(config)
		} else {
			fmt.Println("invalid command")
		}
	}
}

func getCommands() map[string]cliCommand {
	commandMap := map[string]cliCommand{
		"map": {
			name:        "map",
			description: "Displays the next 20 locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous 20 locations",
			callback:    commandMapb,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
	}

	return commandMap
}

func cleanInput(text string) []string {
	//var tokens []string

	text = strings.ToLower(text)
	r := regexp.MustCompile(`[^\s]+`)
	return r.FindAllString(text, -1)

}

func requestData(url string, config *Config) error {

	body, found := config.cache.Get(url)

	if !found {
		res, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
			return err
		}
		body, err = io.ReadAll(res.Body)
		res.Body.Close()
		if res.StatusCode > 299 {
			log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
			return err
		}
		if err != nil {
			log.Fatal(err)
			return err
		}

		config.cache.Add(url, body)
	}
	//fmt.Printf("%s", body)

	type location struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}
	type response struct {
		Count   float64    `json:"count"` // key will be "name"
		Next    string     `json:"next"`  // key will be "id"
		Prev    string     `json:"previous"`
		Results []location `json:"results"`
	}

	var resp response
	err := json.Unmarshal(body, &resp)
	if err != nil {
		fmt.Println("error:", err)
		return err
	}

	for _, v := range resp.Results {
		fmt.Println(v.Name)
	}

	config.next = resp.Next
	config.prev = resp.Prev

	return nil
}

func commandMap(config *Config) error {
	if config.next == "" {
		fmt.Println("You are on the last page")
		return nil
	}

	return requestData(config.next, config)

}

func commandMapb(config *Config) error {
	if config.prev == "" {
		fmt.Println("You are on the first page")
		return nil
	}

	return requestData(config.prev, config)

}

func commandExit(config *Config) error {
	fmt.Print("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *Config) error {
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")
	commandMap := getCommands()

	for _, v := range commandMap {
		fmt.Println(v.name, ": ", v.description)
	}
	return nil
}
