# features

## frontend

- Ribbon menu
  - [x] Show mpd connection state
  - [ ] Display config dialog (is this even needed?)
  - [x] `Update DB` button
    - [x] Disable if running update is detected
  - [ ] `Previous Track` button
  - [ ] `Next Track` button
  - [x] `Stop` button
  - [x] `Play` button
    - [ ] `Pause` button
  - [x] Track seeker
  - [ ] `Repeat` toggle
  - [ ] `Shuffle` toggle
  - [ ] xfade
    - [ ] decrease
    - [ ] increase
  - [x] Volume 
    - [x] increase
    - [x] decrease
    - [x] set with bar
  - [x] `Now playing`
    - [x] shows current track
    - [x] marquee effect
  - [x] `Time`
- Queue
  - [x] Show queue
  - [ ] Highlight current track
  - [ ] Move track up
  - [ ] Move track down
  - [ ] Remove track
  - [ ] `Clear queue` button?
- File browser
  - [ ] List all directories
  - [ ] Open folders have different icon (📂 vs 📁)
  - [ ] Folders with subfolders have a ➕ sign
  - [ ] Clicked folders contents are displayed in the results
  - [ ] Select tracks in results
  - [ ] `Add` selected tracks to queue button
- Search
  - [ ] Search files results
  - [ ] Select tracks in results
  - [ ] `Add` selected tracks to queue button
- Playlist browser
  - [ ] Show current playlists
  - [ ] `Replace` current queue with playlist button
  - [ ] `Attach` playlist to current queue button
  - [ ] `Save` current queue as playlist button
    - [x] Show dialog
  - [ ] `Delete` selected playlist button
    
## backend

- Websocket
  - [ ] `#status` requests mpd infos:
    - `status` 
    - `currentsong`
    - `playlistinfo`
  - [ ] `#download` requests download of URL (`yt-dlp`)
    - *TBA*
- API endpoints
  - [x] `/api/update_db`
  - [x] `/api/previous_track`
  - [x] `/api/next_track`
  - [x] `/api/stop`
  - [x] `/api/play`
  - [ ] `/api/pause`
  - [x] `/api/seek/:seconds`
  - [ ] `/api/repeat`
  - [ ] `/api/random`
  - [x] `/api/volume/:level`
  - [ ] `/api/xfade/:seconds`
  - [ ] `/api/queue_clear`
  - [ ] `/api/queue_add/:songid`
  - [ ] `/api/queue_del/:songid`
  - [ ] `/api/queue_move/:songid/:pos`
  - [ ] `/api/list_database/:path`
  - [ ] `/api/list_playlists`
  - [ ] `/api/save_playlist`
  - [ ] `/api/delete_playlist`


# foo

- client: connect websocket
- server: on_connect: send full state
- server: subscribe to changes
- server: on_change: send to client

