package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fhs/gompd/v2/mpd"
	"github.com/labstack/echo/v4"
	"io"
	"time"
)

// Event represents Server-Sent Event.
// SSE explanation: https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events#event_stream_format
type Event struct {
	// ID is used to set the EventSource object's last event ID value.
	ID []byte
	// Data field is for the message. When the EventSource receives multiple consecutive lines
	// that begin with data:, it concatenates them, inserting a newline character between each one.
	// Trailing newlines are removed.
	Data []byte
	// Event is a string identifying the type of event described. If this is specified, an event
	// will be dispatched on the browser to the listener for the specified event name; the website
	// source code should use addEventListener() to listen for named events. The onmessage handler
	// is called if no event name is specified for a message.
	Event []byte
	// Retry is the reconnection time. If the connection to the server is lost, the browser will
	// wait for the specified time before attempting to reconnect. This must be an integer, specifying
	// the reconnection time in milliseconds. If a non-integer value is specified, the field is ignored.
	Retry []byte
	// Comment line can be used to prevent connections from timing out; a server can send a comment
	// periodically to keep the connection alive.
	Comment []byte
}

// MarshalTo marshals Event to given Writer
func (ev *Event) MarshalTo(w io.Writer) error {
	// Marshalling part is taken from: https://github.com/r3labs/sse/blob/c6d5381ee3ca63828b321c16baa008fd6c0b4564/http.go#L16
	if len(ev.Data) == 0 && len(ev.Comment) == 0 {
		return nil
	}

	if len(ev.Data) > 0 {
		if _, err := fmt.Fprintf(w, "id: %s\n", ev.ID); err != nil {
			return err
		}

		sd := bytes.Split(ev.Data, []byte("\n"))
		for i := range sd {
			if _, err := fmt.Fprintf(w, "data: %s\n", sd[i]); err != nil {
				return err
			}
		}

		if len(ev.Event) > 0 {
			if _, err := fmt.Fprintf(w, "event: %s\n", ev.Event); err != nil {
				return err
			}
		}

		if len(ev.Retry) > 0 {
			if _, err := fmt.Fprintf(w, "retry: %s\n", ev.Retry); err != nil {
				return err
			}
		}
	}

	if len(ev.Comment) > 0 {
		if _, err := fmt.Fprintf(w, ": %s\n", ev.Comment); err != nil {
			return err
		}
	}

	if _, err := fmt.Fprint(w, "\n"); err != nil {
		return err
	}

	return nil
}

// serveSSE handles sending Server-Sent-Events.
func serveSSE(c echo.Context) error {
	// TODO: figure out how to retrieve IP from Forwarded header behind proxy: https://echo.labstack.com/docs/ip-address
	c.Logger().Printf("SSE client connected, ip: %v", c.RealIP())

	w := c.Response()
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Connect to MPD server
	mpdConn, err := mpd.Dial("tcp", "localhost:6600")
	if err != nil {
		c.Logger().Error(err)
		event := Event{
			Event: []byte("mpd"),
			Data:  []byte(fmt.Sprintf("connection error: %s", err.Error())),
		}
		if err := event.MarshalTo(w); err != nil {
			return err
		}
		w.Flush()
		return nil
	}
	defer mpdConn.Close()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	var lastJsonStatus []byte
	var lastJsonCurrentSong []byte
	var lastJsonQueue []byte
	for {
		select {
		case <-c.Request().Context().Done():
			c.Logger().Printf("SSE client disconnected, ip: %v", c.RealIP())
			return nil
		case <-ticker.C:
			c.Logger().Printf("Getting MPD status for %v", c.RealIP())

			status, err := mpdConn.Status()
			if err != nil {
				c.Logger().Error(err)
			}
			jsonStatus, err := json.Marshal(status)
			if err != nil {
				c.Logger().Error(err)
			}
			//c.Logger().Print("status " + string(jsonStatus))
			// Only send new event if different from last time
			if !bytes.Equal(jsonStatus, lastJsonStatus) {
				statusEvent := Event{
					Event: []byte("status"),
					Data:  []byte(string(jsonStatus)),
				}
				if err := statusEvent.MarshalTo(w); err != nil {
					return err
				}
				lastJsonStatus = jsonStatus
			}

			currentsong, err := mpdConn.CurrentSong()
			if err != nil {
				c.Logger().Error(err)
			}
			jsonCurrentSong, err := json.Marshal(currentsong)
			if err != nil {
				c.Logger().Error(err)
			}
			//c.Logger().Print("current_song " + string(jsonCurrentSong))
			// Only send new event if different from last time
			if !bytes.Equal(jsonCurrentSong, lastJsonCurrentSong) {
				currentSongEvent := Event{
					Event: []byte("currentsong"),
					Data:  []byte(string(jsonCurrentSong)),
				}
				if err := currentSongEvent.MarshalTo(w); err != nil {
					return err
				}
				lastJsonCurrentSong = jsonCurrentSong
			}

			queue, err := mpdConn.PlaylistInfo(-1, -1)
			if err != nil {
				c.Logger().Error(err)
			}
			jsonQueue, err := json.Marshal(queue)
			if err != nil {
				c.Logger().Error(err)
			}
			//c.Logger().Print("queue " + string(jsonQueue))
			// Only send new event if different from last time
			if !bytes.Equal(jsonQueue, lastJsonQueue) {
				queueEvent := Event{
					Event: []byte("queue"),
					Data:  []byte(string(jsonQueue)),
				}
				if err := queueEvent.MarshalTo(w); err != nil {
					return err
				}
				lastJsonQueue = jsonQueue
			}

			// Ping to prevent timeout
			pingEvent := Event{
				Comment: []byte("ping"),
			}
			if err := pingEvent.MarshalTo(w); err != nil {
				return err
			}

			w.Flush()
		}
	}
}
