function switchTab(event, tabname) {
  for (let tab in document.getElementsByClassName("tab")) {
    tab.style.display = "none";
  }
  for (let button in document.getElementsByClassName("tablink")) {
    button.className = button.className.replace(" active", "");
  }


  document.getElementById(tabname).style.display = "block";
  event.currentTarget.className += " active";
}
