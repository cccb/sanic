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

// updateDb Updates the music database: find new files, remove deleted files, update modified files.
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

// previousTrack Plays previous song in the queue.
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

// nextTrack Plays next song in the queue.
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

// stopPlayback Stops playing.
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

// resumePlayback Begins playing the playlist or if paused resumes playback.
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

// pausePlayback Pauses playback.
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

// seek Seeks to the position defined by seconds within the current song.
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

// toggleRepeat Toggles repeat state between 1 or 0.
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

// toggleRandom Toggles random state between 1 or 0.
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

// setVolume Sets volume to level, the range of volume is 0-100.
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

// deleteTrackFromQueue removed track with song_id from queue
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

// moveTrackInQueue moves song with song_id to the new place position in the queue.
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

// attachPlaylist adds the playlist with the name playlist_name to the queue.
func attachPlaylist(c echo.Context) error {
	// Connect to MPD server
	conn, err := mpd.Dial("tcp", "localhost:6600")
	if err != nil {
		c.Logger().Error(err)
	}
	defer conn.Close()

	name := c.Param("playlist_name")

	err = conn.PlaylistLoad(name, -1, -1)
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, "")
}

// replaceQueue replaces the current queue with  the playlist with the name playlist_name.
func replaceQueue(c echo.Context) error {
	// Connect to MPD server
	conn, err := mpd.Dial("tcp", "localhost:6600")
	if err != nil {
		c.Logger().Error(err)
	}
	defer conn.Close()

	name := c.Param("playlist_name")

	err = conn.Clear()
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusInternalServerError, err.Error())
	}
	err = conn.PlaylistLoad(name, -1, -1)
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, "")
}

// Playlists

// listPlaylists return a list of all stored playlists.
func listPlaylists(c echo.Context) error {
	// Connect to MPD server
	conn, err := mpd.Dial("tcp", "localhost:6600")
	if err != nil {
		c.Logger().Error(err)
	}
	defer conn.Close()

	playlists, err := conn.ListPlaylists()
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, playlists)
}

// listPlaylist returns the contents of the playlist defined by name.
func listPlaylist(c echo.Context) error {
	// Connect to MPD server
	conn, err := mpd.Dial("tcp", "localhost:6600")
	if err != nil {
		c.Logger().Error(err)
	}
	defer conn.Close()

	name := c.Param("name")

	playlist, err := conn.PlaylistContents(name)
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, playlist)
}

// deletePlaylist deletes the playlist defined by name.
func deletePlaylist(c echo.Context) error {
	// Connect to MPD server
	conn, err := mpd.Dial("tcp", "localhost:6600")
	if err != nil {
		c.Logger().Error(err)
	}
	defer conn.Close()

	name := c.Param("name")

	err = conn.PlaylistRemove(name)
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.String(http.StatusNoContent, "")
}

// savePlaylist saves the current queue to a playlist with the given name.
func savePlaylist(c echo.Context) error {
	// Connect to MPD server
	conn, err := mpd.Dial("tcp", "localhost:6600")
	if err != nil {
		c.Logger().Error(err)
	}
	defer conn.Close()

	name := c.Param("name")

	err = conn.PlaylistSave(name)
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.String(http.StatusCreated, "")
}

// searchDatabase search the database path given by pattern and returns all entries that contain the pattern either in their artist, album or title.
func searchDatabase(c echo.Context) error {
	// Connect to MPD server
	conn, err := mpd.Dial("tcp", "localhost:6600")
	if err != nil {
		c.Logger().Error(err)
	}
	defer conn.Close()

	pattern := c.Param("pattern")

	artistResult, err := conn.Search("artist", pattern)
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusBadRequest, err.Error())
	}
	albumResult, err := conn.Search("album", pattern)
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusBadRequest, err.Error())
	}
	titleResult, err := conn.Search("title", pattern)
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusBadRequest, err.Error())
	}

	songs := append(append(artistResult, albumResult...), titleResult...)

	// make list unique
	uniqueList := make([]mpd.Attrs, 0, len(songs))
	keep := make(map[string]bool)
	for _, song := range songs {
		if _, ok := keep[song["file"]]; !ok {
			keep[song["file"]] = true
			uniqueList = append(uniqueList, song)
		}
	}

	return c.JSON(http.StatusOK, uniqueList)
}
