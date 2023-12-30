const API_URL = `${document.location.protocol}//${document.location.host}/api`;
const VOLUME_STEP = 5;

// Get control elements

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

// Add API calls to controls

control_update_db.addEventListener("click", e => {
  console.log("Issuing database update")
  fetch(`${API_URL}/update_db`).then(async r => {
    if (r.status === 200) {
      const job_id = await r.text();
      console.log(`Update started (Job ID: ${job_id})`)
    } else {
      console.error(`API returned ${r.status}: ${r.statusText}`)
    }
  });
});
control_previous.addEventListener("click", e => {
  fetch(`${API_URL}/previous_track`);
});
control_play_pause.addEventListener("click", e => {
  if (e.target.innerHTML === "&#x23F5;&#xFE0E;") {  // TODO: check is never true
    fetch(`${API_URL}/pause`);
  } else {  // Pause
    fetch(`${API_URL}/play`);
  }
});
control_stop.addEventListener("click", e => {
  fetch(`${API_URL}/stop`);
});
control_next.addEventListener("click", e => {
  fetch(`${API_URL}/next_track`);
});
control_progress.addEventListener("change", e => {
  fetch(`${API_URL}/seek/${e.target.value}`)
});
control_repeat.addEventListener("click", e => {
  if (e.target.innerHTML === "&#x1F534; repeat") {  // TODO: check is never true
    e.target.innerHTML = "&#x1F518; repeat";
  } else {
    e.target.innerHTML = "&#x1F534; repeat";
  }
  fetch(`${API_URL}/repeat`);
});
control_shuffle.addEventListener("click", e => {
  if (e.target.innerHTML === "&#x1F534; shuffle") {  // TODO: check is never true
    e.target.innerHTML = "&#x1F518; shuffle";
  } else {
    e.target.innerHTML = "&#x1F534; shuffle";
  }
  fetch(`${API_URL}/random`);
});
control_xfade_minus.addEventListener("click", e => {
  // TODO: not yet implemented
  fetch(`${API_URL}/xfade`);
});
control_xfade_plus.addEventListener("click", e => {
  // TODO: not yet implemented
  fetch(`${API_URL}/xfade`);
});
control_volume_up.addEventListener("click", e => {
  const v = Math.min(parseInt(control_volume.value) + VOLUME_STEP, 100);
  fetch(`${API_URL}/volume/${v}`);
  control_volume.value = v;

});
control_volume_down.addEventListener("click", e => {
  const v = Math.max(parseInt(control_volume.value) - VOLUME_STEP, 0);
  fetch(`${API_URL}/volume/${v}`);
  control_volume.value = v;
});
control_volume.addEventListener("change", e => {
  fetch(`${API_URL}/volume/${e.target.value}`);
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
      if ("state" in msg.mpd_status && msg.mpd_status.state === "play") {
        control_play_pause.innerHTML = "&#x23F8;&#xFE0E;";  // Pause
      } else {
        control_play_pause.innerHTML = "&#x23F5;&#xFE0E;";  // Play
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
        } else {
          control_repeat.innerHTML = "&#x1F518; repeat"; // 🔘 Radio Button
        }
      }

      // update shuffle state
      if ("random" in msg.mpd_status) {
        if (msg.mpd_status.random === "1") {
          control_shuffle.innerHTML = "&#x1F534; shuffle"; // 🔴 Red Circle
        } else {
          control_shuffle.innerHTML = "&#x1F518; shuffle"; // 🔘 Radio Button
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
    if ("Artist" in msg.mpd_current_song && "Title" in msg.mpd_current_song) {
      control_track.value = `${msg.mpd_current_song.Artist} - ${msg.mpd_current_song.Title}`
    } else {
      control_track.value = msg.mpd_current_song.file;
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
