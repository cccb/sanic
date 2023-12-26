document.addEventListener('DOMContentLoaded', function() {
  function fillTable(withPositionColumn, tableBodyElementName, numElements, playingIndex) {
    for (let i = 0; i < numElements; i++) {
      let trElem = "";
      let positionColumn = "";

      if (i == playingIndex)
        trElem += `<tr class="playing">`;
      else
        trElem += `<tr>`;

      if (withPositionColumn)
        positionColumn += `<td>${i + 1}</td>`

      let exampleTableEntry =
      `                                 \
      ${trElem}                         \
        ${positionColumn}               \
        <td>T.E.E.D.</td>               \
        <td>Garden (Calibre Remix)</td> \
        <td>undefined</td>              \
        <td>undefined</td>              \
        <td>06:01</td>                  \
      </tr>                             \
      `

      tableBodyElement = document.getElementById(tableBodyElementName);
      tableBodyElement.innerHTML += exampleTableEntry;
    }
  }

  fillTable(true, "queue-table-body", 8, 3);
  fillTable(false, "playlist-table-body", 15, -1);
});
