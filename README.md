# Copilot SDK Weather Agent

A quick experiment with GitHub's Copilot SDK in Go. Built this to get hands-on with the SDK and see how it handles tool integration—in this case, fetching weather data from OpenWeatherMap.

Nothing fancy, just a working example of how the Copilot SDK can be used to create an interactive agent with custom tools.

## What It Does

Interactive command-line chat that lets you ask about weather in different cities. The Copilot SDK handles the conversation flow and automatically calls the weather tool when needed.

Example:
```
> What's the weather like in Seattle?
> How cold is it in Anchorage right now?
```

## Setup

### Prerequisites

- Go 1.21+
- OpenWeatherMap API key ([get one here](https://openweathermap.org/api))
- GitHub Copilot SDK access

### Installation

1. Clone and install dependencies:
```bash
git clone <your-repo-url>
cd go-copilot-sdk-weather-agent
go mod tidy
```

2. Create a `.env` file in the project root:
```bash
API_KEY=your_openweathermap_api_key_here
```

3. Run it:
```bash
go run main.go
```

## Usage Examples

Start the agent and ask about weather in any city:

```
└─$ go run main.go
Copilot Weather Agent (Type 'exit' to quit)
> How cold is it in Anchorage right now?

It is currently 14.1°F in Anchorage.

> What's the weather like in Miami?

It is currently 61.9°F in Miami with 61% humidity.

> 
```

Type `exit` or `quit` to stop the program.

## How It Works

The project uses:
- **GitHub Copilot SDK** for the conversational interface and tool orchestration
- **OpenWeatherMap API** for actual weather data
- A simple tool definition that the SDK can call when it needs weather info

The SDK figures out when to invoke the weather tool based on the conversation context, handles the API call, and formats the response naturally.

## Notes

This is a learning project to explore the Copilot SDK. Planning to build something more substantial next (Fail2Ban monitoring tool), but this was a good way to understand the basics of tool integration and streaming responses.

Feel free to use as a starting point for your own Copilot SDK experiments.
