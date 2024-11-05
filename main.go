package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"time"
)

func maxIdleTime() float64 {
	val, exists := os.LookupEnv("FOO_USER")
	if !exists {
		log.Println("FOO_USER doesn't exist, using default of 10s")
		return 10
	} else {
		idleTime, err := strconv.ParseFloat(val, 64)

		if err != nil {
			log.Fatalln("FOO_USER must be a number e.g. 3600")
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

	idleState := make(chan bool)     // create timer channel
	go idleCount(maxTime, idleState) // Start idle count routine

	for {
		select {
		case <-ticker.C:
			if userCount() == 0 { // No users logged in
				idleState <- true
			} else {
				idleState <- false
			}
		default:
			continue
		}
	}

}

func idleCount(maxIdleTime float64, in <-chan bool) {
	for {
		state := <-in //State of idleTimer

		if state {
			log.Printf("Starting %s idle count", time.Duration(maxIdleTime*float64(time.Second)))
			idleTimer := time.AfterFunc(time.Duration(maxIdleTime)*time.Second, shutdown)

			for {
				state := <-in
				if !state {
					log.Println("Idle count stopped")
					idleTimer.Stop()
					break
				}
			}

		}
	}
}

func shutdown() {
	log.Println("We're going to shutdown now")
	os.Exit(0)
}

func main() {
	// log to custom file
	LOG_FILE := "/tmp/cost_saver.log"
	// open log file
	logFile, err := os.OpenFile(LOG_FILE, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Panic(err)
	}
	defer logFile.Close()

	// Set log output
	log.SetOutput(logFile)

	log.Println("Cost-Saver Started")

	checkIdleStatus()
}
