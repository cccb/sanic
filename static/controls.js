// Configuration

const API_URL = `${document.location.protocol}//${document.location.host}/api`;
const VOLUME_STEP = 5;

// Get control elements

const dialog_save_playlist = document.getElementById("save-playlist");
const control_playlist_name = document.getElementById("control-playlist-name");
const dialog_save_playlist_submit = document.querySelector("#save-playlist button");
const dialog_save_playlist_close = document.querySelector("#save-playlist .close");

const connection_state = document.getElementById("connection-state");
const control_update_db = document.getElementById("control-update-db");
const control_previous = document.getElementById("control-previous");
const control_play_pause = document.getElementById("control-play-pause");
const control_stop = document.getElementById("control-stop");
const control_next = document.getElementById("control-next");
const control_progress = document.getElementById("control-progress");
const control_repeat = document.getElementById("control-repeat");
const control_shuffle = document.getElementById("control-shuffle");
const control_xfade = document.getElementById("control-xfade");
const control_xfade_minus = document.getElementById("control-xfade-minus");
const control_xfade_plus = document.getElementById("control-xfade-plus");
const control_volume = document.getElementById("control-volume");
const control_volume_up = document.getElementById("control-volume-up");
const control_volume_down = document.getElementById("control-volume-down");
const queue_table = document.querySelector("#queue tbody");
const control_track = document.getElementById("control-track");
const control_time = document.getElementById("control-time");
const tabs = document.getElementById("tabs");
const tab_browser = document.getElementById("tab-browser");
const tab_search = document.getElementById("tab-search");
const tab_playlists = document.getElementById("tab-playlists");
const control_playlist_list = document.getElementById("control-playlist-list");
const control_replace_playlist = document.getElementById("control-replace-playlist");
const control_attach_playlist = document.getElementById("control-attach-playlist");
const control_save_playlist = document.getElementById("control-save-playlist");
const control_delete_playlist = document.getElementById("control-delete-playlist");
const result_table = document.querySelector("#result tbody");

// UI controls

control_replace_playlist.addEventListener("click", e => {
  fetch(`${API_URL}/`).then(async r => {
    if (r.status !== 200) {
      console.error(`API returned ${r.status}: ${r.statusText}`);
    }
  });
});

control_attach_playlist.addEventListener("click", e => {
  fetch(`${API_URL}/`).then(async r => {
    if (r.status !== 200) {
      console.error(`API returned ${r.status}: ${r.statusText}`);
    }
  });
});

control_save_playlist.addEventListener("click", e => {
  dialog_save_playlist.showModal()
});

dialog_save_playlist_close.addEventListener("click", e => {
  dialog_save_playlist.close()
});

dialog_save_playlist_submit.addEventListener("click", e => {
  fetch(`${API_URL}/playlists`, {method: "PUT"}).then(async r => {
    if (r.status === 201) {
      console.log(`Playlist "${control_playlist_name.value}" saved`)
    }
  });
});

control_delete_playlist.addEventListener("click", e => {
  const playlist_id = control_playlist_list.value;
  fetch(`${API_URL}/playlists/${playlist_id}`, {method: "DELETE"}).then(r => {
    if (r.status === 204) {
      console.log(`Playlist ${playlist_id} successfully deleted.`);
    } else {
      console.error(`API returned ${r.status}: ${r.statusText}`);
    }
  });
});

tab_browser.addEventListener("click", e => {
  if (!tab_browser.classList.contains("active")) {
    tab_browser.classList.add("active");
    tab_search.classList.remove("active")
    tab_playlists.classList.remove("active")
    document.getElementById("file-browser").style.display = "block";
    document.getElementById("search").style.display = "none";
    document.getElementById("playlist-browser").style.display = "none";
  }
});

tab_search.addEventListener("click", e => {
  if (!tab_search.classList.contains("active")) {
    tab_browser.classList.remove("active");
    tab_search.classList.add("active")
    tab_playlists.classList.remove("active")
    document.getElementById("file-browser").style.display = "none";
    document.getElementById("search").style.display = "block";
    document.getElementById("playlist-browser").style.display = "none";
  }
});

tab_playlists.addEventListener("click", e => {
  fetch(`${API_URL}/playlists`).then(async r => {
    if (r.status === 200) {
      const playlists = await r.json();
      control_playlist_list.options.length = 0;  // clear playlists
      playlists.forEach(p => {
        const option = document.createElement("option")
        option.innerText = p["playlist"];
        option.value = p["playlist"];
        option.addEventListener("click", () => {
          fetch(`${API_URL}/playlists/${p["playlist"]}`).then(async r => {
            if (r.status === 200) {
              const songs = await r.json();
              console.log(songs)
              result_table.innerHTML = "";
              songs.forEach(song => {
                const tr = document.createElement("tr");
                const artist = document.createElement("td");
                artist.innerText = song["Artist"];
                const title = document.createElement("td");
                title.innerText = song["Title"];
                const time = document.createElement("td");
                const seconds = parseInt(song["Time"]);
                const time_hours = Math.floor(seconds / 3600);
                const time_minutes = Math.floor((seconds - time_hours * 3600) / 60);
                const time_seconds = Math.floor(seconds - time_hours * 3600 - time_minutes * 60);
                time.innerText = `${time_hours}:${time_minutes.toString().padStart(2, '0')}:${time_seconds.toString().padStart(2, '0')}`
                tr.appendChild(artist);
                tr.appendChild(title);
                tr.appendChild(document.createElement("td")); // album
                tr.appendChild(document.createElement("td")); // genre
                tr.appendChild(time);
                result_table.appendChild(tr);
              });
            }
          })
        });
        control_playlist_list.appendChild(option)
      });
    }
  });
  if (!tab_playlists.classList.contains("active")) {
    tab_browser.classList.remove("active");
    tab_search.classList.remove("active")
    tab_playlists.classList.add("active")
    document.getElementById("file-browser").style.display = "none";
    document.getElementById("search").style.display = "none";
    document.getElementById("playlist-browser").style.display = "block";
  }
});

// Add API calls to controls

control_update_db.addEventListener("click", e => {
  console.log("Issuing database update")
  fetch(`${API_URL}/update_db`).then(async r => {
    if (r.status === 200) {
      const job_id = await r.text();
      console.log(`Update started (Job ID: ${job_id})`);
      e.target.disabled = true;
    } else {
      console.error(`API returned ${r.status}: ${r.statusText}`);
    }
  });
});
control_previous.addEventListener("click", e => {
  fetch(`${API_URL}/previous_track`).then(async r => {
    if (r.status >= 400) {
      console.error(`API returned ${r.status}: ${r.statusText}`);
    }
  });
});
control_play_pause.addEventListener("click", e => {
  if (e.target.innerHTML === "⏸︎") {
    fetch(`${API_URL}/pause`).then(async r => {
      if (r.status >= 400) {
        console.error(`API returned ${r.status}: ${r.statusText}`);
      }
    });
  } else {  // Pause
    fetch(`${API_URL}/play`).then(async r => {
      if (r.status >= 400) {
        console.error(`API returned ${r.status}: ${r.statusText}`);
      }
    });
  }
});
control_stop.addEventListener("click", e => {
  fetch(`${API_URL}/stop`).then(async r => {
    if (r.status >= 400) {
      console.error(`API returned ${r.status}: ${r.statusText}`);
    }
  });
});
control_next.addEventListener("click", e => {
  fetch(`${API_URL}/next_track`).then(async r => {
    if (r.status >= 400) {
      console.error(`API returned ${r.status}: ${r.statusText}`);
    }
  });
});
control_progress.addEventListener("change", e => {
  fetch(`${API_URL}/seek/${e.target.value}`).then(async r => {
    if (r.status >= 400) {
      console.error(`API returned ${r.status}: ${r.statusText}`);
    }
  });
});
control_repeat.addEventListener("click", e => {
  if (e.target.dataset.state === "on") {  // TODO: check is never true
    e.target.innerHTML = "&#x1F518; repeat";
    e.target.dataset.state = "off";
  } else {
    e.target.innerHTML = "&#x1F534; repeat";
    e.target.dataset.state = "on";
  }
  fetch(`${API_URL}/repeat`).then(async r => {
    if (r.status >= 400) {
      console.error(`API returned ${r.status}: ${r.statusText}`);
    }
  });
});
control_shuffle.addEventListener("click", e => {
  if (e.target.dataset.state === "on") {  // TODO: check is never true
    e.target.innerHTML = "&#x1F518; shuffle";
    e.target.dataset.state = "off";
  } else {
    e.target.innerHTML = "&#x1F534; shuffle";
    e.target.dataset.state = "on";
  }
  fetch(`${API_URL}/random`).then(async r => {
    if (r.status >= 400) {
      console.error(`API returned ${r.status}: ${r.statusText}`);
    }
  });
});
control_xfade_minus.addEventListener("click", e => {
  // TODO: not yet implemented
  fetch(`${API_URL}/xfade`).then(async r => {
    if (r.status >= 400) {
      console.error(`API returned ${r.status}: ${r.statusText}`);
    }
  });
});
control_xfade_plus.addEventListener("click", e => {
  // TODO: not yet implemented
  fetch(`${API_URL}/xfade`).then(async r => {
    if (r.status >= 400) {
      console.error(`API returned ${r.status}: ${r.statusText}`);
    }
  });
});
control_volume_up.addEventListener("click", e => {
  const v = Math.min(parseInt(control_volume.value) + VOLUME_STEP, 100);
  fetch(`${API_URL}/volume/${v}`).then(async r => {
    if (r.status >= 400) {
      console.error(`API returned ${r.status}: ${r.statusText}`);
    }
  });
  control_volume.value = v;

});
control_volume_down.addEventListener("click", e => {
  const v = Math.max(parseInt(control_volume.value) - VOLUME_STEP, 0);
  fetch(`${API_URL}/volume/${v}`).then(async r => {
    if (r.status >= 400) {
      console.error(`API returned ${r.status}: ${r.statusText}`);
    }
  });
  control_volume.value = v;
});
control_volume.addEventListener("change", e => {
  fetch(`${API_URL}/volume/${e.target.value}`).then(async r => {
    if (r.status >= 400) {
      console.error(`API returned ${r.status}: ${r.statusText}`);
    }
  });
});

// Create WebSocket connection.
const socket = new WebSocket(`${document.location.protocol === "https:" ? "wss" : "ws"}://${document.location.host}/ws`);

// Connection opened
socket.addEventListener("open", (e) => {
  socket.send("Hello Server!");
});

// Listen for messages and update UI state
socket.addEventListener("message", (e) => {
  // Print out mpd response
  console.log(`DEBUG: ${e.data}`);  // DEBUG

  const msg = JSON.parse(e.data);

  if ("mpd_status" in msg) {
    if (msg.mpd_status == null) {
      connection_state.innerHTML = "&#x274C; Disconnected";  // ✅ Check Mark Button
    } else {
      // print error if present
      if ("error" in msg.mpd_status) {
        console.error(msg.mpd_status.error);
      }

      // update "Update DB" button
      if ("updating_db" in msg.mpd_status) {
        control_update_db.disabled = true;
      } else {
        if (control_update_db.disabled) {
          console.log("Database update done.")
        }
        control_update_db.disabled = false;
      }

      // update play/pause button
      if ("state" in msg.mpd_status && msg.mpd_status.state !== "play") {  // TODO: only update DOM if necessary
        control_play_pause.innerHTML = "&#x23F5;&#xFE0E;";  // Play
      } else {
        control_play_pause.innerHTML = "&#x23F8;&#xFE0E;";  // Pause
      }

      // update playback time
      if ("elapsed" in msg.mpd_status && "duration" in msg.mpd_status) {
        const elapsed_hours = Math.floor(msg.mpd_status.elapsed / 3600);
        const elapsed_minutes = Math.floor((msg.mpd_status.elapsed - elapsed_hours * 3600) / 60);
        const elapsed_seconds = Math.floor(msg.mpd_status.elapsed - elapsed_hours * 3600 - elapsed_minutes * 60);
        const duration_hours = Math.floor(msg.mpd_status.duration / 3600);
        const duration_minutes = Math.floor((msg.mpd_status.duration - duration_hours * 3600) / 60);
        const duration_seconds = Math.floor(msg.mpd_status.duration - duration_hours * 3600 - duration_minutes * 60);
        control_time.value = `${elapsed_hours}:${elapsed_minutes.toString().padStart(2, '0')}:${elapsed_seconds.toString().padStart(2, '0')}/${duration_hours}:${duration_minutes.toString().padStart(2, '0')}:${duration_seconds.toString().padStart(2, '0')}`;
      }
      if ("elapsed" in msg.mpd_status) {
        control_progress.value = msg.mpd_status.elapsed;
      }
      if ("duration" in msg.mpd_status) {
        control_progress.max = msg.mpd_status.duration;
      }

      // update repeat state
      if ("repeat" in msg.mpd_status) {
        if (msg.mpd_status.repeat === "1") {
          control_repeat.innerHTML = "&#x1F534; repeat"; // 🔴 Red Circle
          control_repeat.dataset.state = "on";
        } else {
          control_repeat.innerHTML = "&#x1F518; repeat"; // 🔘 Radio Button
          control_repeat.dataset.state = "off";
        }
      }

      // update shuffle state
      if ("random" in msg.mpd_status) {
        if (msg.mpd_status.random === "1") {
          control_shuffle.innerHTML = "&#x1F534; shuffle"; // 🔴 Red Circle
          control_shuffle.dataset.state = "on";
        } else {
          control_shuffle.innerHTML = "&#x1F518; shuffle"; // 🔘 Radio Button
          control_shuffle.dataset.state = "off";
        }
      }

      // update crossfade state
      if ("xfade" in msg.mpd_status) {
        control_xfade.value = msg.mpd_status.xfade;
      }

      // update volume
      if ("volume" in msg.mpd_status) {
        control_volume.value = msg.mpd_status.volume;
      }
    }
  }

  // update song info
  if ("mpd_current_song" in msg && msg.mpd_current_song != null) {
    let track;
    if ("Artist" in msg.mpd_current_song && "Title" in msg.mpd_current_song) {
      track = `<span>${msg.mpd_current_song.Artist} - ${msg.mpd_current_song.Title}</span>`
    } else {
      track = `<span>${msg.mpd_current_song.file}</span>`;
    }
    if (control_track.innerHTML.toString() !== track) {
      control_track.innerHTML = track;
    }
  }

  // update queue
  if ("mpd_queue" in msg && msg.mpd_queue != null) {
    const tbody = document.createElement("tbody");
    msg.mpd_queue.forEach(elem => {
      const tr = document.createElement("tr");
      tr.dataset.song_id = elem.Id;
      if ("songid" in msg.mpd_status && msg.mpd_status.songid === elem.Id) {
        tr.classList.add("playing");
      } else {
        tr.classList.remove("playing");
      }
      // TODO: check if current row is currently playing track
      const pos = document.createElement("td");
      pos.innerText = elem.Pos;
      const artist = document.createElement("td");
      if ("Artist" in elem) {
        artist.innerText = elem.Artist;
      }
      const track = document.createElement("td");
      if ("Title" in elem) {
        track.innerText = elem.Title;
      } else {
        track.innerText = elem.file;
      }
      const album = document.createElement("td");
      // album.innerText = "";
      const length = document.createElement("td");
      const duration_hours = Math.floor(elem.duration / 3600);
      const duration_minutes = Math.floor((elem.duration - duration_hours * 3600) / 60);
      const duration_seconds = Math.floor(elem.duration - duration_hours * 3600 - duration_minutes * 60);
      length.innerText = `${duration_hours}:${duration_minutes.toString().padStart(2, '0')}:${duration_seconds.toString().padStart(2, '0')}`;
      const actions = document.createElement("td");
      // TODO: maybe use a instead of button?
      const moveUp = document.createElement("button");
      moveUp.classList.add("borderless");
      moveUp.innerHTML = "&#x1F53A;"; // 🔺 Red Triangle Pointed Down
      moveUp.addEventListener("click", event => {
        console.log(`DEBUG: move song ${elem.Pos} up`);
        // fetch(`${API_URL}/queue_del/${elem.Pos}`).then(r => {
        //   console.log(r.text());
        // });
      });
      // TODO: maybe use a instead of button?
      const moveDown = document.createElement("button");
      moveDown.classList.add("borderless");
      moveDown.innerHTML = "&#x1F53B;"; // 🔻 Red Triangle Pointed Up
      moveDown.addEventListener("click", event => {
        console.log(`DEBUG: move song ${elem.Pos} down`);
        // fetch(`${API_URL}/queue_del/${elem.Pos}`).then(r => {
        //   console.log(r.text());
        // });
      });
      // TODO: maybe use a instead of button?
      const remove = document.createElement("button");
      remove.classList.add("borderless");
      remove.innerHTML = "&#x274C;"; // ❌ Cross mark
      remove.addEventListener("click", event => {
      console.log(`DEBUG: remove song id ${elem.Id} from queue`);
        fetch(`${API_URL}/queue/delete/${elem.Id}`).then(r => {
          console.log(r.text());
        });
      });
      actions.appendChild(moveUp);
      actions.appendChild(moveDown);
      actions.appendChild(remove);
      tr.appendChild(pos);
      tr.appendChild(artist);
      tr.appendChild(track);
      tr.appendChild(album);
      tr.appendChild(length);
      tr.appendChild(actions);
      tbody.appendChild(tr);
    });
    const currentQueue = document.querySelector("#queue tbody")
    if (currentQueue.innerHTML !== tbody.innerHTML) {
      console.log("Updating queue")
      currentQueue.outerHTML = tbody.outerHTML;
    }
  }

  if ("mpd_error" in msg) {
    console.error(`MPD Error: ${msg.mpd_error}`)
  }
});

// Request MPD status every second
window.setInterval(() => {
  if (socket.readyState === socket.OPEN) {
    socket.send("#status");
    connection_state.innerHTML = "&#x2705; Connected";  // ❌ Cross Mark
  } else {
    connection_state.innerHTML = "&#x274C; Disconnected";  // ✅ Check Mark Button
  }
}, 1000);
