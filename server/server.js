/**
 * Created by Alnim on 05/29/2018.
 */

var express = require('express');
var app = express();
var http = require('http').Server(app);
var io = require('socket.io')(http);
var mongoose = require('mongoose');
mongoose.connect("127.0.0.1", "IMDB");

app.set('port', (process.env.PORT || 8080));
app.get('/', function (req, res) {
    res.sendFile(__dirname + '../client/index.html');
    uIP = req.ip;
});

app.use(express.static(__dirname + '/'));

http.listen(app.get('port'), '0.0.0.0', function () {
    console.log(' +' + app.get('port'));
});