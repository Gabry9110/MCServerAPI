package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	SERVER_PORT         uint16
	SERVER_SERVICE_NAME string
	TIMEOUT             time.Duration
)

func main() {
	// Check if running on Linux
	if runtime.GOOS != "linux" {
		fmt.Println("This software only works on Linux systems.")
		os.Exit(1)
	}

	// Check if systemd is being used, as it depends on it for runtime
	_, err := exec.Command("systemctl", "--version").Output()
	if err != nil {
		fmt.Println("This software requires systemd as init system.")
	}

	// Command line arguments for custom names
	SERVER_PORT = uint16(*flag.Uint("port", uint(25565), "Minecraft server port (default: 25565)"))
	SERVER_SERVICE_NAME = *flag.String("service-name", "minecraft.service", "Systemd service name for the Minecraft server (default: minecraft.service)")
	TIMEOUT, err = time.ParseDuration(*flag.String("timeout", "5m", "Timeout for auto-stop (default: 5m)"))
	if err != nil {
		fmt.Println("Use a valid timeout format (such as 5m or 2m30s)")
	}

	// Start the API
	router := gin.Default()
	router.GET("/", noEndpoint)
	router.GET("/ping", ping)
	router.GET("/start", start)
	router.GET("/stop", stop)

	router.Run("localhost:8080")
}
