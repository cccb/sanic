package main

import (
	"fmt"
	"github.com/fhs/gompd/v2/mpd"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
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

	e.GET("/ws", wsServe)

	e.Logger.Fatal(e.StartTLS(":1323", "cert.pem", "key.pem"))
	//e.Logger.Fatal(e.Start(":1323"))
}

func wsServe(c echo.Context) error {
	fmt.Println("wsServe")
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()
		fmt.Println("handler")
		for {
			// Read
			msg := ""
			err := websocket.Message.Receive(ws, &msg)
			if err != nil {
				c.Logger().Error(err)
				break
			} else {
				if strings.HasPrefix(strings.ToUpper(msg), "MPD#") {
					// Forward MPD communication
					// TODO: forward request to mpd and response back to client
					err := websocket.Message.Send(ws, "MPD command received, processing... processing...")
					if err != nil {
						c.Logger().Error(err)
					}

				} else if strings.HasPrefix(strings.ToUpper(msg), "YT#") {
					// Download video link as audio file
					// TODO: implement yt-dlp integration
					err := websocket.Message.Send(ws, "YT-DLP command received, processing... processing...")
					if err != nil {
						c.Logger().Error(err)
					}
				}
			}
			//fmt.Println(msg)
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}

// API calls

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

	return c.String(http.StatusNoContent, "")
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

	return c.String(http.StatusNoContent, "")
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

	return c.String(http.StatusNoContent, "")
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

	return c.String(http.StatusNoContent, "")
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

	return c.String(http.StatusNoContent, "")
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

	return c.String(http.StatusNoContent, "")
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

	return c.String(http.StatusNoContent, "")
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

	return c.String(http.StatusNoContent, "")
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

	return c.String(http.StatusNoContent, "")
}
