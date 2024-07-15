// Configuration

const API_URL = `${document.location.protocol}//${document.location.host}/api`;
const VOLUME_STEP = 5;

// Get control elements

const dialog_save_playlist = document.getElementById("save-playlist");
const control_playlist_name = document.getElementById("control-playlist-name");
const dialog_save_playlist_submit = document.querySelector("#save-playlist form button");
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
const control_refresh_playlists = document.getElementById("control-refresh-playlists");
const control_replace_playlist = document.getElementById("control-replace-playlist");
const control_attach_playlist = document.getElementById("control-attach-playlist");
const control_save_playlist = document.getElementById("control-save-playlist");
const control_delete_playlist = document.getElementById("control-delete-playlist");
const result_table = document.querySelector("#result tbody");
const control_search_pattern = document.getElementById("control-search-pattern");
const control_search_submit = document.getElementById("control-search-submit");

// Utility functions

secondsToTrackTime = (t) => {
  const hours = Math.floor(t / 3600);
  const minutes = Math.floor((t - hours * 3600) / 60);
  const seconds = Math.floor(t - hours * 3600 - minutes * 60);

  return `${hours}:${minutes.toString().padStart(2, "0")}:${seconds.toString().padStart(2, "0")}`;
}

removeTrackFromQueue = (event) => {
  const song_id = event.target.parentElement.parentElement.dataset.song_id;

  console.log(`DEBUG: remove song id ${song_id} from queue`);
  fetch(`${API_URL}/queue/${song_id}/delete`).then(r => {
    console.log(r.text());
  });
}

moveTrackInQueue = (event, direction) => {
  const song_id = event.target.parentElement.parentElement.dataset.song_id;
  // TODO: figure out position in queue by counting HTML elements?
  const position = parseInt(event.target.parentElement.parentElement.firstChild.innerText);

  console.log(`DEBUG: move song ${song_id} down in queue to position ${position + direction}`);
  fetch(`${API_URL}/queue/${song_id}/move/${position + direction}`).then(r => {
    console.log(r.text());
  });
}

refreshPlaylists = () => {
  console.log("Refreshing playlists from MPD server")
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
            if (r.status === 200) fillResultTable(await r.json());
          })
        });
        control_playlist_list.appendChild(option)
      });
    }
  });
}

fillResultTable = (songs) => {
  result_table.innerHTML = "";
  songs.forEach(song => {
    const tr = document.createElement("tr");
    const artist = document.createElement("td");
    artist.innerText = song["Artist"];
    const title = document.createElement("td");
    title.innerText = song["Title"];
    const time = document.createElement("td");
    time.innerText = secondsToTrackTime(parseInt(song["Time"]))
    tr.appendChild(artist);
    tr.appendChild(title);
    tr.appendChild(document.createElement("td")); // album
    tr.appendChild(document.createElement("td")); // genre
    tr.appendChild(time);
    result_table.appendChild(tr);
  });
}

// UI controls

tab_browser.addEventListener("click", () => {
  if (!tab_browser.classList.contains("active")) {
    tab_browser.classList.add("active");
    tab_search.classList.remove("active")
    tab_playlists.classList.remove("active")
    document.getElementById("file-browser").style.display = "block";
    document.getElementById("search").style.display = "none";
    document.getElementById("playlist-browser").style.display = "none";
  }
});

tab_search.addEventListener("click", () => {
  if (!tab_search.classList.contains("active")) {
    tab_browser.classList.remove("active");
    tab_search.classList.add("active")
    tab_playlists.classList.remove("active")
    document.getElementById("file-browser").style.display = "none";
    document.getElementById("search").style.display = "block";
    document.getElementById("playlist-browser").style.display = "none";
  }
});

tab_playlists.addEventListener("click", () => {
  if (!tab_playlists.classList.contains("active")) {
    tab_browser.classList.remove("active");
    tab_search.classList.remove("active")
    tab_playlists.classList.add("active")
    document.getElementById("file-browser").style.display = "none";
    document.getElementById("search").style.display = "none";
    document.getElementById("playlist-browser").style.display = "block";
  }
});

// Show "Save playlist" modal
control_save_playlist.addEventListener("click", () => {
  dialog_save_playlist.showModal()
});

// Close "Save playlist" modal
dialog_save_playlist_close.addEventListener("click", () => {
  dialog_save_playlist.close()
});


// Add API calls to controls

control_search_submit.addEventListener("click", event => {
  event.preventDefault()
  fetch(`${API_URL}/database/${control_search_pattern.value}`).then(async r => {
    if (r.status === 200) {
      fillResultTable([...new Set(await r.json())]);
    } else {
      console.error(`API returned ${r.status}: ${r.statusText}`);
    }
  })
});

control_refresh_playlists.addEventListener("click", () => {
  refreshPlaylists();
});

control_replace_playlist.addEventListener("click", () => {
  fetch(`${API_URL}/queue/replace/${control_playlist_list.value}`).then(async r => {
    if (r.status !== 200) {
      console.error(`API returned ${r.status}: ${r.statusText}`);
    }
  });
});

control_attach_playlist.addEventListener("click", () => {
  fetch(`${API_URL}/queue/attach/${control_playlist_list.value}`).then(async r => {
    if (r.status !== 200) {
      console.error(`API returned ${r.status}: ${r.statusText}`);
    }
  });
});

// Save current queue as new playlist and refresh playlist list
dialog_save_playlist_submit.addEventListener("click", () => {
  fetch(`${API_URL}/playlists/${control_playlist_name.value}`, {method: "POST"}).then(async r => {
    if (r.status === 201) {
      console.log(`Playlist "${control_playlist_name.value}" saved`)
      refreshPlaylists()
    } else {
      console.error(`API returned ${r.status}: ${r.statusText}`);
    }
  });
});

control_delete_playlist.addEventListener("click", () => {
  const playlist_name = control_playlist_list.value;
  fetch(`${API_URL}/playlists/${control_playlist_list.value}`, {method: "DELETE"}).then(r => {
    if (r.status === 204) {
      console.log(`Playlist "${playlist_name}" successfully deleted.`);
      refreshPlaylists();
    } else {
      console.error(`API returned ${r.status}: ${r.statusText}`);
    }
  });
});

tab_browser.addEventListener("click", () => {
  if (!tab_browser.classList.contains("active")) {
    tab_browser.classList.add("active");
    tab_search.classList.remove("active")
    tab_playlists.classList.remove("active")
    document.getElementById("file-browser").style.display = "block";
    document.getElementById("search").style.display = "none";
    document.getElementById("playlist-browser").style.display = "none";
  }
});

tab_search.addEventListener("click", () => {
  if (!tab_search.classList.contains("active")) {
    tab_browser.classList.remove("active");
    tab_search.classList.add("active")
    tab_playlists.classList.remove("active")
    document.getElementById("file-browser").style.display = "none";
    document.getElementById("search").style.display = "block";
    document.getElementById("playlist-browser").style.display = "none";
  }
});

tab_playlists.addEventListener("click", () => {
  refreshPlaylists();
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

control_update_db.addEventListener("click", () => {
  console.log("Issuing database update")
  fetch(`${API_URL}/update_db`).then(async r => {
    if (r.status === 200) {
      console.log(await r.text());
      event.target.disabled = true;
    } else {
      console.error(`API returned ${r.status}: ${r.statusText}`);
    }
  });
});

control_previous.addEventListener("click", () => {
  fetch(`${API_URL}/previous_track`).then(async r => {
    if (r.status >= 400) {
      console.error(`API returned ${r.status}: ${r.statusText}`);
    }
  });
});

control_play_pause.addEventListener("click", event => {
  if (event.target.innerHTML === "⏸︎") {  // Resume playback
    fetch(`${API_URL}/pause`).then(async r => {
      if (r.status >= 400) {
        console.error(`API returned ${r.status}: ${r.statusText}`);
      }
    });
  } else {  // Pause playback
    fetch(`${API_URL}/play`).then(async r => {
      if (r.status >= 400) {
        console.error(`API returned ${r.status}: ${r.statusText}`);
      }
    });
  }
});

control_stop.addEventListener("click", () => {
  fetch(`${API_URL}/stop`).then(async r => {
    if (r.status >= 400) {
      console.error(`API returned ${r.status}: ${r.statusText}`);
    }
  });
});

control_next.addEventListener("click", () => {
  fetch(`${API_URL}/next_track`).then(async r => {
    if (r.status >= 400) {
      console.error(`API returned ${r.status}: ${r.statusText}`);
    }
  });
});

control_progress.addEventListener("change", event => {
  fetch(`${API_URL}/seek/${event.target.value}`).then(async r => {
    if (r.status >= 400) {
      console.error(`API returned ${r.status}: ${r.statusText}`);
    }
  });
});

control_progress.addEventListener("input", event => {
  control_time.value = `${secondsToTrackTime(event.target.value)}/${secondsToTrackTime(event.target.max)}`;
});

control_repeat.addEventListener("click", event => {
  if (event.target.dataset.state === "on") {  // TODO: check is never true
    event.target.innerHTML = "&#x1F518; repeat";
    event.target.dataset.state = "off";
  } else {
    event.target.innerHTML = "&#x1F534; repeat";
    event.target.dataset.state = "on";
  }
  fetch(`${API_URL}/repeat`).then(async r => {
    if (r.status >= 400) {
      console.error(`API returned ${r.status}: ${r.statusText}`);
    }
  });
});

control_shuffle.addEventListener("click", event => {
  if (event.target.dataset.state === "on") {  // TODO: check is never true
    event.target.innerHTML = "&#x1F518; shuffle";
    event.target.dataset.state = "off";
  } else {
    event.target.innerHTML = "&#x1F534; shuffle";
    event.target.dataset.state = "on";
  }
  fetch(`${API_URL}/random`).then(async r => {
    if (r.status >= 400) {
      console.error(`API returned ${r.status}: ${r.statusText}`);
    }
  });
});

control_xfade_minus.addEventListener("click", () => {
  // TODO: not yet implemented
  fetch(`${API_URL}/xfade`).then(async r => {
    if (r.status >= 400) {
      console.error(`API returned ${r.status}: ${r.statusText}`);
    }
  });
});

control_xfade_plus.addEventListener("click", () => {
  // TODO: not yet implemented
  fetch(`${API_URL}/xfade`).then(async r => {
    if (r.status >= 400) {
      console.error(`API returned ${r.status}: ${r.statusText}`);
    }
  });
});

control_volume_up.addEventListener("click", () => {
  const volume = Math.min(parseInt(control_volume.value) + VOLUME_STEP, 100);
  fetch(`${API_URL}/volume/${volume}`).then(async r => {
    if (r.status >= 400) {
      console.error(`API returned ${r.status}: ${r.statusText}`);
    }
  });
  control_volume.value = volume;

});

control_volume_down.addEventListener("click", () => {
  const volume = Math.max(parseInt(control_volume.value) - VOLUME_STEP, 0);
  fetch(`${API_URL}/volume/${volume}`).then(async r => {
    if (r.status >= 400) {
      console.error(`API returned ${r.status}: ${r.statusText}`);
    }
  });
  control_volume.value = volume;
});

control_volume.addEventListener("change", event => {
  fetch(`${API_URL}/volume/${event.target.value}`).then(async r => {
    if (r.status >= 400) {
      console.error(`API returned ${r.status}: ${r.statusText}`);
    }
  });
});

// Websocket logic

// Create WebSocket connection.
const socket = new WebSocket(`${document.location.protocol === "https:" ? "wss" : "ws"}://${document.location.host}/ws`);

// Connection opened
socket.addEventListener("open", () => {
  socket.send("Hello Server!");
});

// Listen for messages and update UI state
socket.addEventListener("message", event => {
  // Print out mpd response
  console.log(`DEBUG: ${event.data}`);  // DEBUG

  const msg = JSON.parse(event.data);

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
      if ("time" in msg.mpd_status) {
        const [elapsed, duration] = msg.mpd_status.time.split(":", 2)
        control_progress.value = elapsed;
        control_progress.max = duration;
        // triggers the update of control_time element
        const e = new Event("input");
        control_progress.dispatchEvent(e);
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
      track = `${msg.mpd_current_song.Artist} - ${msg.mpd_current_song.Title}`
    } else {
      track = msg.mpd_current_song.file;
    }
    if (control_track.innerHTML !== `<span>${track}</span>`) {
      control_track.innerHTML = `<span>${track}</span>`;
    }
  }

  // update queue
  if ("mpd_queue" in msg && msg.mpd_queue != null) {
    const tbody = document.createElement("tbody");
    msg.mpd_queue.forEach(song => {
      const tr = document.createElement("tr");
      tr.dataset.song_id = song.Id;
      if ("songid" in msg.mpd_status && msg.mpd_status.songid === song.Id) {
        tr.classList.add("playing");
      }
      const pos = document.createElement("td");
      pos.innerText = song.Pos;
      const artist = document.createElement("td");
      if ("Artist" in song) {
        artist.innerText = song.Artist;
      }
      const track = document.createElement("td");
      if ("Title" in song) {
        track.innerText = song.Title;
      } else {
        track.innerText = song.file;
      }
      const album = document.createElement("td");
      // TODO: Do songs have album info attached to them?
      album.innerText = "";
      const length = document.createElement("td");
      length.innerText = secondsToTrackTime(song.duration);
      const actions = document.createElement("td");
      const moveUp = document.createElement("button");
      moveUp.classList.add("borderless");
      if (parseInt(song.Pos) !== 0) {
        moveUp.innerHTML = "&#x1F53A;"; // 🔺 Red Triangle Pointed Down
        moveUp.addEventListener("click", event => { moveTrackInQueue(event, -1) });
      } else {
        moveUp.innerHTML = "&emsp;";
      }
      const moveDown = document.createElement("button");
      moveDown.classList.add("borderless");
      if (parseInt(song.Pos) !== msg.mpd_queue.length - 1) {
        moveDown.innerHTML = "&#x1F53B;"; // 🔻 Red Triangle Pointed Up
        moveDown.addEventListener("click", event => {moveTrackInQueue(event, 1)});
      } else {
        moveDown.innerHTML = "&emsp;";
      }
      const remove = document.createElement("button");
      remove.classList.add("borderless");
      remove.innerHTML = "&#x274C;"; // ❌ Cross mark
      remove.addEventListener("click", removeTrackFromQueue);
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
    connection_state.innerHTML = "&#x2705; Connected";  // ✅ Check Mark Button
  } else {
    connection_state.innerHTML = "&#x274C; Disconnected";  // ❌ Cross Mark
  }
}, 1000);

