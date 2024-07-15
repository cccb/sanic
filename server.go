package main

import (
	"encoding/json"
	"fmt"
	"github.com/fhs/gompd/v2/mpd"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/net/websocket"
	"gopkg.in/ini.v1"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

// Config holds the configuration for the mpd connection and for the web server.
type Config struct {
	MPD struct {
		Hostname string `ini:"hostname"`
		Port     int    `ini:"port"`
		Username string `ini:"username"`
		Password string `ini:"password"`
	} `ini:"mpd"`
	UI struct {
		Hostname    string `ini:"hostname"`
		Port        int    `ini:"port"`
		Tls         bool   `ini:"tls"`
		Certificate string `ini:"cert"`
		Key         string `ini:"key"`
	} `ini:"ui"`
}

func main() {
	iniData, err := ini.Load("config.ini")
	if err != nil {
		fmt.Printf("Fail to read configuration file: %v", err)
		os.Exit(1)
	}

	var config Config

	err = iniData.MapTo(&config)
	if err != nil {
		fmt.Printf("Fail to parse configuration file: %v", err)
		os.Exit(1)
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:  "static",
		HTML5: true, // SPA mode; not-found will be redirected to root
	}))
	e.Use(echoprometheus.NewMiddleware("sanic")) // adds middleware to gather metrics

	e.GET("/metrics", echoprometheus.NewHandler())

	e.GET("/", func(c echo.Context) (err error) {
		// HTTP/2 Server Push
		pusher, ok := c.Response().Writer.(http.Pusher)
		if ok {
			if err = pusher.Push("/style.css", nil); err != nil {
				return
			}
			if err = pusher.Push("/index.js", nil); err != nil {
				return
			}
			if err = pusher.Push("/favicon.ico", nil); err != nil {
				return
			}
		}
		return c.File("index.html")
	})

	g := e.Group("/api")
	g.GET("/update_db", updateDb)
	g.GET("/previous_track", previousTrack)
	g.GET("/next_track", nextTrack)
	g.GET("/stop", stopPlayback)
	g.GET("/play", resumePlayback)
	g.GET("/pause", pausePlayback)
	g.GET("/seek/:seconds", seek)
	g.GET("/repeat", toggleRepeat)
	g.GET("/random", toggleRandom)
	g.GET("/volume/:level", setVolume)

	g.GET("/queue/:song_id/delete", deleteTrackFromQueue)
	g.GET("/queue/:song_id/move/:position", moveTrackInQueue)
	g.GET("/queue/replace/:playlist_name", replaceQueue)
	g.GET("/queue/attach/:playlist_name", attachPlaylist)

	g.GET("/playlists", listPlaylists)
	g.POST("/playlists/:name", savePlaylist)
	g.GET("/playlists/:name", listPlaylist)
	g.DELETE("/playlists/:name", deletePlaylist)

	g.GET("/database/:pattern", searchDatabase)

	g.GET("/download", downloadTrack)

	e.GET("/ws", wsServe)

	if config.UI.Tls {
		e.Logger.Fatal(e.StartTLS(fmt.Sprintf("%s:%d", config.UI.Hostname, config.UI.Port), config.UI.Certificate, config.UI.Key))
	} else {
		e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%d", config.UI.Hostname, config.UI.Port)))
	}
}

// wsServe handles websocket connections.
func wsServe(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()

		// Connect to MPD server
		mpdConn, err := mpd.Dial("tcp", "localhost:6600")
		if err != nil {
			//log.Fatalln(err)
			c.Logger().Error(err)
			err = websocket.Message.Send(ws, fmt.Sprintf("{\"mpd_error\":\"%s\"}", err.Error()))
			if err != nil {
				c.Logger().Error(err)
			}
		}
		defer mpdConn.Close()

		for {
			// Read
			msg := ""
			err := websocket.Message.Receive(ws, &msg)
			if err != nil {
				c.Logger().Error(err)
				break
			} else {
				// log.Print(msg)
				if strings.ToLower(msg) == "#status" {
					status, err := mpdConn.Status()
					if err != nil {
						c.Logger().Error(err)
					}
					currentsong, err := mpdConn.CurrentSong()
					if err != nil {
						c.Logger().Error(err)
					}
					queue, err := mpdConn.PlaylistInfo(-1, -1)
					if err != nil {
						c.Logger().Error(err)
					}
					jsonStatus, err := json.Marshal(status)
					if err != nil {
						c.Logger().Error(err)
					}
					jsonCurrentSong, err := json.Marshal(currentsong)
					if err != nil {
						c.Logger().Error(err)
					}
					jsonQueue, err := json.Marshal(queue)
					if err != nil {
						c.Logger().Error(err)
					}
					err = websocket.Message.Send(ws, fmt.Sprintf("{\"mpd_status\":%s,\"mpd_current_song\":%s,\"mpd_queue\":%s}", string(jsonStatus), string(jsonCurrentSong), string(jsonQueue)))
					if err != nil {
						c.Logger().Error(err)
					}

				} else if strings.HasPrefix(strings.ToLower(msg), "#download ") {
					// Download video link as audio file
					uri := strings.SplitN(msg, " ", 2)[1]
					// TODO: implement yt-dlp integration
					err := websocket.Message.Send(ws, fmt.Sprintf("Downloading %s", uri))
					if err != nil {
						c.Logger().Error(err)
					}
				}
			}
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}

// downloadTrack tries to download a given URL and saves the song to the database.
func downloadTrack(c echo.Context) error {
	// yt-dlp \
	// --no-wait-for-video \
	// --no-playlist \
	// --windows-filenames \
	// --newline \
	// --extract-audio \
	// --audio-format mp3 \
	// --audio-quality 0 \
	// -f bestaudio/best \
	// ${video_url}

	cmd := exec.Command(
		"yt-dlp",
		"--no-wait-for-video",
		"--no-playlist",
		"--windows-filenames",
		"--newline",
		"--extract-audio",
		"--audio-format", "mp3",
		"--audio-quality", "0",
		"--format", "bestaudio/best",
		c.Param("url"),
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalln(err)
	}

	return c.String(http.StatusAccepted, "")
}
