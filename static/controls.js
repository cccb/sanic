const API_URL = `${document.location.protocol}://${document.location.host}/api`;

// Get control elements

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
const queue_table = document.querySelector("#queue tbody");

// Add API calls to controls

control_update_db.addEventListener("click", e => {
  fetch(`${API_URL}/update_db`);
});
control_previous.addEventListener("click", e => {
  fetch(`${API_URL}/previous_track`);
});
control_play_pause.addEventListener("click", e => {
  if (e.target.innerText === "&#x23F5;&#xFE0E;") {  // Play
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
  fetch(`${API_URL}/repeat`);
});
control_shuffle.addEventListener("click", e => {
  fetch(`${API_URL}/shuffle`);
});
control_xfade_minus.addEventListener("click", e => {
  // TODO: not yet implemented
  fetch(`${API_URL}/xfade`);
});
control_xfade_plus.addEventListener("click", e => {
  // TODO: not yet implemented
  fetch(`${API_URL}/xfade`);
});

// Create WebSocket connection.
const socket = new WebSocket(`${document.location.protocol === "https" ? "wss" : "ws"}://${document.location.host}/ws`);

// Connection opened
socket.addEventListener("open", (e) => {
  socket.send("Hello Server!");
});

// Listen for messages
socket.addEventListener("message", (e) => {
  console.log("Message from server");
  const msg = JSON.parse(e.data);
  if ("error" in msg.mpd_status) {
    console.error(msg.mpd_status.error);
  }

  if ("updating_db" in msg.mpd_status) {
    control_update_db.disable();
  } else {
    control_update_db.enable();
  }
  if ("state" in msg.mpd_status && msg.mpd_status.state === "play") {
    control_play_pause.innerText = "&#x23F8;&#xFE0E;";  // Pause
  } else {
    control_play_pause.innerText = "&#x23F5;&#xFE0E;";  // Play
  }
  if ("elapsed" in msg.mpd_status) {
    control_progress.value = msg.mpd_status.elapsed;
  }
  if ("duration" in msg.mpd_status) {
    control_progress.max = msg.mpd_status.duration;
  }
  if ("repeat" in msg.mpd_status) {
    control_repeat.checked = msg.mpd_status.repeat;
  }
  if ("random" in msg.mpd_status) {
    control_shuffle.checked = msg.mpd_status.random;
  }
  if ("xfade" in msg.mpd_status) {
    control_xfade.value = msg.mpd_status.xfade;
  }
});
