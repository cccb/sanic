mpd:
	mkdir -p /tmp/sanic/{music,playlists}
	mpd --no-daemon ./mpd.conf
