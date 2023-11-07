package main

import (
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/net/websocket"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
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

	e.GET("/ws", wsServe)

	e.Logger.Fatal(e.StartTLS(":1323", "cert.pem", "key.pem"))
	//e.Logger.Fatal(e.Start(":1323"))
}

func wsServe(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()
		for {
			// Read
			msg := ""
			err := websocket.Message.Receive(ws, &msg)
			if err != nil {
				c.Logger().Error(err)
			}
			// Forward MPD communication
			if strings.HasPrefix(strings.ToUpper(msg), "MPD#") {
				// TODO: forward request to mpd and response back to client
				err := websocket.Message.Send(ws, "MPD command received, processing... processing...")
				if err != nil {
					c.Logger().Error(err)
				}
			}
			// Download video link as audio file
			if strings.HasPrefix(strings.ToUpper(msg), "YT#") {
				// TODO: implement yt-dlp integration
				err := websocket.Message.Send(ws, "YT-DLP command received, processing... processing...")
				if err != nil {
					c.Logger().Error(err)
				}
			}
			//fmt.Printf("%s\n", msg)
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}
