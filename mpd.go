package main

import (
	"fmt"
	"github.com/fhs/gompd/v2/mpd"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"time"
)

// MPD API calls

func updateDb(c echo.Context) error {
	// Connect to MPD server
	conn, err := mpd.Dial("tcp", "localhost:6600")
	if err != nil {
		c.Logger().Error(err)
	}
	defer conn.Close()

	jobId, err := conn.Update("")
	if err != nil {
		c.Logger().Error(err)
	}

	return c.String(http.StatusOK, fmt.Sprintf("Database update started with job id %d", jobId))
}

func previousTrack(c echo.Context) error {
	// Connect to MPD server
	conn, err := mpd.Dial("tcp", "localhost:6600")
	if err != nil {
		c.Logger().Error(err)
	}
	defer conn.Close()

	err = conn.Previous()
	if err != nil {
		c.Logger().Error(err)
	}

	return c.String(http.StatusOK, "Playing previous track in queue")
}

func nextTrack(c echo.Context) error {
	// Connect to MPD server
	conn, err := mpd.Dial("tcp", "localhost:6600")
	if err != nil {
		c.Logger().Error(err)
	}
	defer conn.Close()

	err = conn.Next()
	if err != nil {
		c.Logger().Error(err)
	}

	return c.String(http.StatusOK, "PLaying next track in queue")
}

func stopPlayback(c echo.Context) error {
	// Connect to MPD server
	conn, err := mpd.Dial("tcp", "localhost:6600")
	if err != nil {
		c.Logger().Error(err)
	}
	defer conn.Close()

	err = conn.Stop()
	if err != nil {
		c.Logger().Error(err)
	}

	return c.String(http.StatusOK, "Playback stopped")
}

func resumePlayback(c echo.Context) error {
	// Connect to MPD server
	conn, err := mpd.Dial("tcp", "localhost:6600")
	if err != nil {
		c.Logger().Error(err)
	}
	defer conn.Close()

	status, err := conn.Status()
	if err != nil {
		c.Logger().Error(err)
	}
	if status["state"] == "stop" {
		err := conn.Play(-1)
		if err != nil {
			c.Logger().Error(err)
		}
	} else {
		err = conn.Pause(false)
		if err != nil {
			c.Logger().Error(err)
		}
	}

	return c.String(http.StatusOK, "Playback resumed")
}

func pausePlayback(c echo.Context) error {
	// Connect to MPD server
	conn, err := mpd.Dial("tcp", "localhost:6600")
	if err != nil {
		c.Logger().Error(err)
	}
	defer conn.Close()

	err = conn.Pause(true)
	if err != nil {
		c.Logger().Error(err)
	}

	return c.String(http.StatusOK, "Playback paused")
}

func seek(c echo.Context) error {
	// Connect to MPD server
	conn, err := mpd.Dial("tcp", "localhost:6600")
	if err != nil {
		c.Logger().Error(err)
	}
	defer conn.Close()

	seconds, err := strconv.Atoi(c.Param("seconds"))
	if err != nil {
		c.Logger().Error(err)
	}

	if seconds < 0 {
		return c.String(http.StatusBadRequest, "seconds must be positive integer")
	}

	// TODO: Duration type seems to be used incorrectly
	err = conn.SeekCur(time.Duration(seconds)*time.Second, false)
	if err != nil {
		c.Logger().Error(err)
	}

	return c.String(http.StatusOK, fmt.Sprintf("Seeked current track to %d seconds", seconds))
}

func toggleRepeat(c echo.Context) error {
	// Connect to MPD server
	conn, err := mpd.Dial("tcp", "localhost:6600")
	if err != nil {
		c.Logger().Error(err)
	}
	defer conn.Close()

	status, err := conn.Status()
	if err != nil {
		c.Logger().Error(err)
	}
	var msg string
	if status["repeat"] == "1" {
		err = conn.Repeat(false)
		msg = "Toggled Repeat mode to off"
	} else {
		err = conn.Repeat(true)
		msg = "Toggled Repeat mode to on"
	}
	if err != nil {
		c.Logger().Error(err)
	}

	return c.String(http.StatusOK, msg)
}

func toggleRandom(c echo.Context) error {
	// Connect to MPD server
	conn, err := mpd.Dial("tcp", "localhost:6600")
	if err != nil {
		c.Logger().Error(err)
	}
	defer conn.Close()

	status, err := conn.Status()
	if err != nil {
		c.Logger().Error(err)
	}
	var msg string
	if status["random"] == "1" {
		err = conn.Random(false)
		msg = "Toggled Random mode to off"
	} else {
		err = conn.Random(true)
		msg = "Toggled Random mode to on"
	}
	if err != nil {
		c.Logger().Error(err)
	}

	return c.String(http.StatusOK, msg)
}

func setVolume(c echo.Context) error {
	// Connect to MPD server
	conn, err := mpd.Dial("tcp", "localhost:6600")
	if err != nil {
		c.Logger().Error(err)
	}
	defer conn.Close()

	level, err := strconv.Atoi(c.Param("level"))
	if err != nil {
		c.Logger().Error(err)
	}

	if level > 100 || level < 0 {
		return c.String(http.StatusBadRequest, "Volume must be between 0 and 100")
	}

	err = conn.SetVolume(level)
	if err != nil {
		c.Logger().Error(err)
	}

	return c.String(http.StatusOK, fmt.Sprintf("Set volume to %d", level))
}

// Queue

func deleteTrackFromQueue(c echo.Context) error {
	// Connect to MPD server
	conn, err := mpd.Dial("tcp", "localhost:6600")
	if err != nil {
		c.Logger().Error(err)
	}
	defer conn.Close()

	songId, err := strconv.Atoi(c.Param("song_id"))
	if err != nil {
		c.Logger().Error(err)
	}

	err = conn.DeleteID(songId)
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.String(http.StatusOK, fmt.Sprintf("Removed song %d from queue", songId))
}

func moveTrackInQueue(c echo.Context) error {
	// Connect to MPD server
	conn, err := mpd.Dial("tcp", "localhost:6600")
	if err != nil {
		c.Logger().Error(err)
	}
	defer conn.Close()

	songId, err := strconv.Atoi(c.Param("song_id"))
	if err != nil {
		c.Logger().Error(err)
	}

	position, err := strconv.Atoi(c.Param("position"))
	if err != nil {
		c.Logger().Error(err)
	}

	err = conn.MoveID(songId, position)
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.String(http.StatusOK, fmt.Sprintf("Moved song %d to position %d", songId, position))
}
