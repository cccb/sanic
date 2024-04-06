package main

import (
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

	return c.String(http.StatusOK, strconv.Itoa(jobId))
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

	return c.String(http.StatusOK, "")
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

	return c.String(http.StatusOK, "")
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

	return c.String(http.StatusOK, "")
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

	return c.String(http.StatusOK, "")
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

	return c.String(http.StatusOK, "")
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

	err = conn.SeekCur(time.Duration(seconds)*time.Second, false)
	if err != nil {
		c.Logger().Error(err)
	}

	return c.String(http.StatusOK, "")
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
	if status["repeat"] == "1" {
		err = conn.Repeat(false)
	} else {
		err = conn.Repeat(true)
	}
	if err != nil {
		c.Logger().Error(err)
	}

	return c.String(http.StatusOK, "")
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
	if status["random"] == "1" {
		err = conn.Random(false)
	} else {
		err = conn.Random(true)
	}
	if err != nil {
		c.Logger().Error(err)
	}

	return c.String(http.StatusOK, "")
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

	return c.String(http.StatusOK, "")
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

	return c.String(http.StatusOK, "")
}
