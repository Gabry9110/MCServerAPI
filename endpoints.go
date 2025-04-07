package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mcstatus-io/mcutil"
)

var timerActive = false
var mu sync.Mutex
var stopMonitorCh = make(chan struct{})

// noEndpoint: if no endpoint is specified in the URL, warn the user
func noEndpoint(c *gin.Context) {
	c.String(http.StatusOK, "Please specify an endpoint in my URL!\n")
}

// ping: just used to check if the API is active
func ping(c *gin.Context) {
	c.String(http.StatusOK, "Pong! :3\n")
}

// start: start the server and handle the timer
func start(c *gin.Context) {
	// systemctl doesn't do anything if the service is already running, no need to check
	err := handleService("start")
	// If the systemctl command gets an error for some reason
	if err != nil {
		fmt.Println(err.Error())
		c.String(http.StatusInternalServerError, "Server Error :(\n")
		return
	}

	// Check if a timer is already active (we don't want more than one at the same time) and start one if not
	mu.Lock()
	if timerActive {
		fmt.Println("Timer already active")
		mu.Unlock()
		c.String(http.StatusOK, "Server already on with active timer\n")
		return
	}
	timerActive = true
	stopMonitorCh = make(chan struct{})
	mu.Unlock()

	// Timeout goroutine, starts executing this in parallel and continues with the code outside
	go monitorPlayersAndShutdown()

	c.String(http.StatusOK, "Started! :D\n")
}

func stop(c *gin.Context) {
	mu.Lock()
	if timerActive {
		fmt.Println("Stopping server and monitoring...")
		close(stopMonitorCh)
	}
	mu.Unlock()

	err := handleService("stop")
	if err != nil {
		fmt.Println(err.Error())
		c.String(http.StatusInternalServerError, "Server Error :(\n")
		return
	}

	c.String(http.StatusOK, "Server stopped! :D\n")
}

func monitorPlayersAndShutdown() {
	defer resetTimerFlag()

	for {
		fmt.Println("Starting timer...")
		select {
		case <-time.After(TIMEOUT): // honk shooo mimimimi
			// Keep checking
		case <-stopMonitorCh:
			fmt.Println("Monitoraggio interrotto manualmente.")
			return
		}

		// Timer over, save current server status
		// If there was an error right earlier, assume no player is online
		var onlinePlayers int64 = 0
		response, err := mcutil.Status("localhost", SERVER_PORT)
		if err != nil {
			fmt.Println("Error while checking server status: ", err)
			// Start another timer and check again
			continue
		}
		if response.Players.Online != nil {
			onlinePlayers = *response.Players.Online
		}

		// Check if someone is online
		if onlinePlayers == 0 {
			fmt.Println("No player logged-in after timeout, shutting down the server...")
			err := handleService("stop")
			if err != nil {
				fmt.Println("Error in server shutdown: ", err)
			}
			return
		}

		fmt.Printf("Logged-in players: (%d)\n The server will continue to run.\n", onlinePlayers)
	}
}

func resetTimerFlag() {
	mu.Lock()
	defer mu.Unlock()
	timerActive = false
}
