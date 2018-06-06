let socket = io();

$searchInput = document.getElementById("searchInput");
$resultsDiv = document.getElementById("resultsDiv");

$searchInput.onkeyup = function (e) {
    if (e.which <= 90 && e.which >= 48 || e.which == 8) {
        if ($searchInput.value.length > 0)
            socket.emit("imdb query string", $searchInput.value);
        else
            $resultsDiv.innerHTML = "";
    }
}

function showQueryResults(docsStr) {
    console.log("got");
    let docs = JSON.parse(docsStr);
    if (docs.length > 0) {
        // console.log("docs: ", docs);
        $resultsDiv.innerHTML = "";
        for (let i of docs)
            $resultsDiv.innerHTML += (i.title) + " | " + + (i.year) + "<br>";
    }
}
socket.on("server: search result", showQueryResults);