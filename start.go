package main

import (
	"fmt"
	"time"

	"github.com/codegangsta/cli"
)

// StartCommand starts local klient. Requires sudo.
func StartCommand(c *cli.Context) int {
	s, err := newService()
	if err != nil {
		fmt.Printf("Error starting %s: '%s'\n", KlientName, err)
		return 1
	}

	if err := s.Start(); err != nil {
		fmt.Printf("Error starting %s: '%s'\n", KlientName, err)
		return 1
	}

	fmt.Println("Waiting until started...")
	if err := WaitUntilStarted(KlientAddress, 5, 1*time.Second); err != nil {
		fmt.Printf("Error: Klient was unable to start properly.")
		return 1
	}

	fmt.Printf("Successfully started %s\n", KlientName)
	return 0
}

func WaitUntilStarted(address string, attempts int, pauseIntv time.Duration) error {
	var err error
	// Try multiple times to connect to Klient, and return the final error
	// if needed.
	for i := 0; i < 5; i++ {
		time.Sleep(pauseIntv)

		if err = HealthCheck(address); err == nil {
			break
		}
	}
	return err
}
