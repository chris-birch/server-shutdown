package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"time"
)

func maxIdleTime() float64 {
	val, exists := os.LookupEnv("SERVER_IDLE_TIME")
	if !exists {
		fmt.Println("SERVER_IDLE_TIME doesn't exist, using default of 3600s")
		return 3600
	} else {
		idleTime, err := strconv.ParseFloat(val, 64)

		if err != nil {
			fmt.Println("SERVER_IDLE_TIME must be a number e.g. 3600")
			os.Exit(1)
		}

		return idleTime
	}

}

func userCount() int {
	// Execute the 'who' command on the Linux shell and count the lines
	cmd := "who | wc -l"
	out, err := exec.Command("bash", "-c", cmd).Output()

	if err != nil {
		fmt.Printf("%s", err)
	}

	// Extract just the number(s) from the result
	re := regexp.MustCompile("[0-9]+")
	output := re.FindAllString(string(out), -1)

	//convert string to int
	ucount, err := strconv.Atoi(output[0])

	if err != nil {
		fmt.Println(err)
	}

	return ucount
}

func checkIdleStatus() {
	var maxTime = maxIdleTime()
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	idleState := make(chan bool)     // Create a timer channel
	go idleCount(maxTime, idleState) // Start idle count routine

	for {
		select {
		case <-ticker.C:
			if userCount() == 0 {
				idleState <- true
			} else {
				idleState <- false
			}
		default:
			continue
		}
	}

}

func idleCount(idleTime float64, in <-chan bool) {
	for {
		state := <-in //State of idleTimer

		if state {
			fmt.Printf("No users are currently logged in. Starting %s idle timer\n", time.Duration(idleTime*float64(time.Second)))
			idleTimer := time.AfterFunc(time.Duration(idleTime)*time.Second, shutdown)

			for {
				state := <-in
				if !state {
					fmt.Println("User logged in, idle timer stopped")
					idleTimer.Stop()
					break
				}
			}

		}
	}
}

func shutdown() {
	idleTime := maxIdleTime()
	fmt.Printf("No users have been logged in for %s. Shutting down.\n", time.Duration(idleTime*float64(time.Second)))
	cmd := "shutdown -h 1 \"Server Shutdown service will now shutdown the server\""
	_, err := exec.Command("bash", "-c", cmd).Output()

	if err != nil {
		fmt.Printf("%s", err)
	}
}

func main() {
	fmt.Println("Server Shutdown Service STARTED")

	// Start checking for an idle server
	checkIdleStatus()
}
