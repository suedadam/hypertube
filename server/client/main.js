let socket = io();

$enterBtn = document.getElementById("enterBtn");
$searchInput = document.getElementById("searchInput");

$enterBtn.onclick = function(){
    socket.emit("imdb query string", $searchInput.value);
    $searchInput.value = "";
}