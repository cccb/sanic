function addSongToQueue(song) {
  const table = document.querySelector("#queue tbody");
  const tr = document.createElement("tr");
  const pos = document.createElement("td");
  pos.innerText = "?";  // TODO: figure out queue length +1
  const artist = document.createElement("td");
  artist.innerText = song.artist;
  const track = document.createElement("td");
  track.innerText = song.track;
  const album = document.createElement("td");
  album.innerText = song.album;
  const length = document.createElement("td");
  length.innerText = song.length;
  const actions = document.createElement("td");
  actions.classList.add("actions");
  // TODO: maybe use a instead of button?
  const moveUp = document.createElement("button");
  moveUp.classList.add("borderless");
  moveUp.innerHTML = "&#x1F53A;"; // 🔺 Red Triangle Pointed Down
  moveUp.addEventListener("click", event => {
    console.log(`DEBUG: move song ${song.id} up`);
  });
  // TODO: maybe use a instead of button?
  const moveDown = document.createElement("button");
  moveDown.classList.add("borderless");
  moveDown.innerHTML = "$#x1F53B;"; // 🔻 Red Triangle Pointed Up
  moveDown.addEventListener("click", event => {
    console.log(`DEBUG: move song ${song.id} down`);
  });
  // TODO: maybe use a instead of button?
  const remove = document.createElement("button");
  remove.classList.add("borderless");
  remove.innerHTML = "$#x274C;"; // ❌ Cross mark; 🗑️ Wastebasket
  remove.addEventListener("click", event => {
    console.log(`DEBUG: remove song ${song.id} from queue`);
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
  table.appendChild(tr);
}
