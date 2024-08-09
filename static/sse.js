// Server-Sent-Events

if (typeof (EventSource) !== "undefined") {
  const sse = new EventSource("/sse");
  sse.addEventListener("status", handleStatus);
  sse.addEventListener("currentsong", handleCurrentSong);
  sse.addEventListener("queue", handleQueue);
  sse.onmessage = (event) => {
    console.log("sse message: " + event.data);
  };
  sse.onerror = (err) => {
    console.error("EventSource failed:", err);
    connection_state.innerHTML = "&#x274C; Disconnected";  // ❌ Cross Mark
  };
  sse.onopen = () => {
    console.log("EventSource connected");
    connection_state.innerHTML = "&#x2705; Connected";  // ✅ Check Mark Button
  };
} else {
  console.error("Sorry, your browser does not support server-sent events...");
}

function handleStatus(event) {
  const status = JSON.parse(event.data);

  // print error if present
  if ("error" in status) {
    console.error(status.error);
  }

  // update "Update DB" button
  if ("updating_db" in status) {
    control_update_db.disabled = true;
  } else {
    if (control_update_db.disabled) {
      console.log("Database update done.")
    }
    control_update_db.disabled = false;
  }

  // update play/pause button
  // TODO: only update DOM if necessary
  if ("state" in status && status.state !== "play") {
    control_play_pause.innerHTML = "&#x23F5;&#xFE0E;";  // Play
  } else {
    control_play_pause.innerHTML = "&#x23F8;&#xFE0E;";  // Pause
  }

  if ("songid" in status) {
    control_track.dataset.songid = status.songid;
  }

  // update playback time
  if ("time" in status) {
    const [elapsed, duration] = status.time.split(":", 2)
    control_progress.value = elapsed;
    control_progress.max = duration;
    // triggers the update of control_time element
    const e = new Event("input");
    control_progress.dispatchEvent(e);
  }

  // update repeat state
  if ("repeat" in status) {
    if (status.repeat === "1") {
      control_repeat.innerHTML = "&#x1F534; repeat"; // 🔴 Red Circle
      control_repeat.dataset.state = "on";
    } else {
      control_repeat.innerHTML = "&#x1F518; repeat"; // 🔘 Radio Button
      control_repeat.dataset.state = "off";
    }
  }

  // update shuffle state
  if ("random" in status) {
    if (status.random === "1") {
      control_shuffle.innerHTML = "&#x1F534; shuffle"; // 🔴 Red Circle
      control_shuffle.dataset.state = "on";
    } else {
      control_shuffle.innerHTML = "&#x1F518; shuffle"; // 🔘 Radio Button
      control_shuffle.dataset.state = "off";
    }
  }

  // update crossfade state
  if ("xfade" in status) {
    control_xfade.value = status.xfade;
  }

  // update volume
  if ("volume" in status) {
    control_volume.value = status.volume;
  }
}

function handleCurrentSong(event) {
  const current_song = JSON.parse(event.data);

  let track;
  if ("Artist" in current_song && "Title" in current_song) {
    track = `${current_song.Artist} - ${current_song.Title}`
  } else {
    track = current_song.file;
  }
  // Only replace if necessary to not interrupt the animation
  if (control_track.innerHTML !== `<span>${track}</span>`) {
    control_track.innerHTML = `<span>${track}</span>`;
  }
}

function handleQueue(event) {
  const queue = JSON.parse(event.data);

  console.log(queue);

  const tbody = document.createElement("tbody");
  queue.forEach(song => {
    const tr = document.createElement("tr");
    tr.dataset.song_id = song.Id;
    if (control_track.dataset.songid === song.Id) {
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
    if (parseInt(song.Pos) !== queue.length - 1) {
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
  // only update queue if necessary to not interrupt user interaction
  if (currentQueue.innerHTML !== tbody.innerHTML) {
    console.log("Updating queue")
    currentQueue.outerHTML = tbody.outerHTML;
  }
}
