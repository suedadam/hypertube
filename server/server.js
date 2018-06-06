/**
 * Created by Alnim on 05/29/2018.
 */

var express = require('express');
var app = express();
var http = require('http').Server(app);
var io = require('socket.io')(http);
var socketSearch = require('./socketSearch');
var trie        = require('./trie');
var mongoose = require('mongoose');
// import GoogleImageSearch from 'free-google-image-search'

// GoogleImageSearch.searchImage("Forrest Gump Movie Cover")
// .then((res) => {
//     console.log(res[0]); // This will return array of image URLs
// })
mongoose.connect("mongodb://localhost/IMDB");
// Mongoose Schema (kUsers)
var movieSchema = new mongoose.Schema({
    mId: {
        type: String,
        unique: true
    },
    title: String,
    genres: String,
    rating: Number
}, { collection: "movies" });


//Mongoose Modal (kUsers)
var movies = mongoose.model('Movies', movieSchema, "movies");
module.exports.movies = movies;
trie.create();

app.set('port', (process.env.PORT || 8080));
app.get('/', function (req, res) {
    res.sendFile(__dirname + "/client/index.html");
    uIP = req.ip;
});

app.use(express.static(__dirname + '/'));
io.on('connection', function (socket) {
    console.log("A view has been collected " + socket.id);
    console.log(socketSearch);
    socket.on("imdb query string", function (query) { socketSearch.onIMDBQuery(query, socket) });
});
http.listen(app.get('port'), '10.142.0.3', function () {
    console.log(' +' + app.get('port'));
});