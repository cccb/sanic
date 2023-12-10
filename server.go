package main

import (
	"encoding/json"
	"fmt"
	"github.com/fhs/gompd/v2/mpd"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func main() {
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
	g.GET("/previous_track", previousTrack)
	g.GET("/next_track", nextTrack)
	g.GET("/stop", stopPlayback)
	g.GET("/play", resumePlayback)
	g.GET("/pause", pausePlayback)
	g.GET("/seek/:seconds", seek)
	g.GET("/repeat", toggleRepeat)
	g.GET("/random", toggleRandom)
	g.GET("/volume/:level", setVolume)

	g.GET("/download", downloadTrack)

	e.GET("/ws", wsServe)

	//e.Logger.Fatal(e.StartTLS(":1323", "cert.pem", "key.pem"))
	e.Logger.Fatal(e.Start(":1323"))
}

func wsServe(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()

		// Connect to MPD server
		mpdConn, err := mpd.Dial("tcp", "localhost:6600")
		if err != nil {
			log.Fatalln(err)
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
				log.Println(msg)
				if strings.ToLower(msg) == "#status" {
					// TODO: Get current MPD status and return it
					status, err := mpdConn.Status()
					if err != nil {
						log.Fatalln(err)
					}
					jsonData, err := json.Marshal(status)
					if err != nil {
						log.Fatalln(err)
					}
					err = websocket.Message.Send(ws, fmt.Sprintf("{\"mpd_status\":%s}", string(jsonData)))
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
