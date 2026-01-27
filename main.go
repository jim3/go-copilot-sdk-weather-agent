package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	copilot "github.com/github/copilot-sdk/go"
	"github.com/joho/godotenv"
)

var baseURL string = "https://api.openweathermap.org/data/2.5/weather?q="
var units string = "imperial"

// Define the parameter type for the tool
type WeatherParams struct {
	City string `json:"city" jsonschema:"The city name"`
}

type Main struct {
	Temp      float64 `json:"temp"`
	FeelsLike float64 `json:"feels_like"`
	TempMin   float64 `json:"temp_min"`
	TempMax   float64 `json:"temp_max"`
	Pressure  int     `json:"pressure"`
	Humidity  int     `json:"humidity"`
	SeaLevel  int     `json:"sea_level"`
	GrndLevel int     `json:"grnd_level"`
}

type WeatherResult struct {
	Main Main `json:"main"`
}

// Define the return type for the tool
type ToolWeatherResult struct {
	City        string `json:"city"`
	Temperature string `json:"temperature"`
	Humidity    string `json:"humidity"`
}

func createWeatherTool() copilot.Tool {
	return copilot.DefineTool(
		"get_weather",
		"Get the current weather for a city",
		func(params WeatherParams, inv copilot.ToolInvocation) (ToolWeatherResult, error) {
			err := godotenv.Load()
			if err != nil {
				log.Fatal("Error loading .env file")
			}

			APIKEY := os.Getenv("API_KEY")
			URL := fmt.Sprintf("%s%s&units=%s&appid=%s", baseURL, params.City, units, APIKEY)

			req, err := http.NewRequest(http.MethodGet, URL, nil)
			if err != nil {
				fmt.Printf("client request failed: %s\n", err)
				os.Exit(1)
			}

			res, err := http.DefaultClient.Do(req)
			if err != nil {
				fmt.Printf("client response failed: %s\n", err)
				os.Exit(1)
			}
			fmt.Printf("STATUS CODE:%d\n\n", res.StatusCode)

			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				fmt.Printf("%s\n", err)
			}

			var w WeatherResult
			err = json.Unmarshal(resBody, &w)
			if err != nil {
				fmt.Printf("%s\n", err)
			}

			return ToolWeatherResult{
				City:        params.City,
				Temperature: fmt.Sprintf("%.1f°F", w.Main.Temp),
				Humidity:    fmt.Sprintf("%d", w.Main.Humidity),
			}, nil
		},
	)
}

func main() {
	// Initialize the Copilot Client
	getWeather := createWeatherTool()
	client := copilot.NewClient(nil)
	if err := client.Start(); err != nil {
		log.Fatal(err)
	}
	defer client.Stop()

	// Create a session

	session, err := client.CreateSession(&copilot.SessionConfig{
		Model:     "gpt-4.1",
		Streaming: true,
		Tools:     []copilot.Tool{getWeather},
	})
	if err != nil {
		log.Fatal(err)
	}

	// Handle streaming events

	session.On(func(event copilot.SessionEvent) {
		if event.Type == "assistant.message_delta" {
			fmt.Print(*event.Data.DeltaContent)
		}
		if event.Type == "session.idle" {
			fmt.Println()
		}
	})

	// ------------------ The REPL Loop ------------------

	// Create a scanner to read user input from the command line
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Copilot Weather Agent (Type 'exit' to quit)")
	fmt.Print("> ")

	// Read user input in a loop
	for scanner.Scan() {
		userInput := scanner.Text()

		if userInput == "exit" || userInput == "quit" {
			break
		}

		// Send the user input to the Copilot session
		_, err = session.SendAndWait(copilot.MessageOptions{
			Prompt: userInput, // What is the the current temperature in Anchorage?
		}, 0)
		if err != nil {
			log.Fatal(err)
		}

		// Print the prompt symbol again for the next question
		fmt.Print("\n> ")
	}

	os.Exit(0)
}
