package main

import (
	"fmt"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gopkg.in/ini.v1"
	"net/http"
	"os"
	"os/exec"
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

	// echo back request to check if HTTP/2 works etc
	e.GET("/echo", func(c echo.Context) error {
		req := c.Request()
		format := `
<code>
    Protocol: %s<br>
    Host: %s<br>
    Remote Address: %s<br>
    Method: %s<br>
    Path: %s<br>
</code>
`
		return c.HTML(http.StatusOK, fmt.Sprintf(format, req.Proto, req.Host, req.RemoteAddr, req.Method, req.URL.Path))
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

	e.GET("/sse", serveSSE)

	if config.UI.Tls {
		e.Logger.Fatal(e.StartTLS(fmt.Sprintf("%s:%d", config.UI.Hostname, config.UI.Port), config.UI.Certificate, config.UI.Key))
	} else {
		e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%d", config.UI.Hostname, config.UI.Port)))
	}
}

// downloadTrack tries to download a given URL and saves the song to the database.
func downloadTrack(c echo.Context) error {
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
		c.Logger().Fatal(err)
	}

	return c.String(http.StatusAccepted, "")
}
