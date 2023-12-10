package main

import (
	"github.com/fhs/gompd/v2/mpd"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"strconv"
	"time"
)

// MPD API calls

func previousTrack(c echo.Context) error {
	// Connect to MPD server
	conn, err := mpd.Dial("tcp", "localhost:6600")
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	err = conn.Previous()
	if err != nil {
		log.Fatalln(err)
	}

	return c.String(http.StatusOK, "")
}

func nextTrack(c echo.Context) error {
	// Connect to MPD server
	conn, err := mpd.Dial("tcp", "localhost:6600")
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	err = conn.Next()
	if err != nil {
		log.Fatalln(err)
	}

	return c.String(http.StatusOK, "")
}

func stopPlayback(c echo.Context) error {
	// Connect to MPD server
	conn, err := mpd.Dial("tcp", "localhost:6600")
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	err = conn.Stop()
	if err != nil {
		log.Fatalln(err)
	}

	return c.String(http.StatusOK, "")
}

func resumePlayback(c echo.Context) error {
	// Connect to MPD server
	conn, err := mpd.Dial("tcp", "localhost:6600")
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	err = conn.Pause(false)
	if err != nil {
		log.Fatalln(err)
	}

	return c.String(http.StatusOK, "")
}

func pausePlayback(c echo.Context) error {
	// Connect to MPD server
	conn, err := mpd.Dial("tcp", "localhost:6600")
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	err = conn.Pause(true)
	if err != nil {
		log.Fatalln(err)
	}

	return c.String(http.StatusOK, "")
}

func seek(c echo.Context) error {
	// Connect to MPD server
	conn, err := mpd.Dial("tcp", "localhost:6600")
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	seconds, err := strconv.Atoi(c.Param("seconds"))
	if err != nil {
		log.Fatalln(err)
	}

	if seconds < 0 {
		return c.String(http.StatusBadRequest, "seconds must be positive integer")
	}

	err = conn.SeekCur(time.Duration(seconds)*time.Second, false)
	if err != nil {
		log.Fatalln(err)
	}

	return c.String(http.StatusOK, "")
}

func toggleRepeat(c echo.Context) error {
	// Connect to MPD server
	conn, err := mpd.Dial("tcp", "localhost:6600")
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	status, err := conn.Status()
	if err != nil {
		log.Fatalln(err)
	}
	if status["repeat"] == "1" {
		err = conn.Repeat(false)
	} else {
		err = conn.Repeat(true)
	}
	if err != nil {
		log.Fatalln(err)
	}

	return c.String(http.StatusOK, "")
}

func toggleRandom(c echo.Context) error {
	// Connect to MPD server
	conn, err := mpd.Dial("tcp", "localhost:6600")
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	status, err := conn.Status()
	if err != nil {
		log.Fatalln(err)
	}
	if status["toggleRandom"] == "1" {
		err = conn.Random(false)
	} else {
		err = conn.Random(true)
	}
	if err != nil {
		log.Fatalln(err)
	}

	return c.String(http.StatusOK, "")
}

func setVolume(c echo.Context) error {
	// Connect to MPD server
	conn, err := mpd.Dial("tcp", "localhost:6600")
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	level, err := strconv.Atoi(c.Param("level"))
	if err != nil {
		log.Fatalln(err)
	}

	if level > 100 || level < 0 {
		return c.String(http.StatusBadRequest, "Volume must be between 0 and 100")
	}

	err = conn.SetVolume(level)
	if err != nil {
		log.Fatalln(err)
	}

	return c.String(http.StatusOK, "")
}
