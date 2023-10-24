package main

import (
	"fmt"
	"github.com/fhs/gompd/v2/mpd"
	"log"
	"time"
)

func main() {
	w, err := mpd.NewWatcher("tcp", ":6600", "")
	if err != nil {
		log.Fatalln(err)
	}
	defer w.Close()

	// Log errors.
	go func() {
		for err := range w.Error {
			log.Println("Error:", err)
		}
	}()

	// Log events.
	go func() {
		for subsystem := range w.Event {
			log.Println("Changed subsystem:", subsystem)
		}
	}()

	// Do other stuff...
	time.Sleep(3 * time.Minute)
}

func main2() {
	// Connect to MPD server
	conn, err := mpd.Dial("tcp", "localhost:6600")
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	// Loop printing the current status of MPD
	for {
		status, err := conn.Status()
		if err != nil {
			log.Fatalln(err)
		}
		song, err := conn.CurrentSong()
		if err != nil {
			log.Fatalln(err)
		}
		if status["state"] == "play" {
			fmt.Println(fmt.Sprintf("%s - %s", song["Artist"], song["Title"]))
		} else {
			fmt.Println(fmt.Sprintf("State: %s", status["state"]))
		}
		time.Sleep(1 * time.Second)
	}
}
