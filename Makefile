mpd:
	mkdir -p /tmp/sanic/{music,playlists}
	touch /tmp/sanic/mpd_db
	mpd --no-daemon ./mpd.conf
