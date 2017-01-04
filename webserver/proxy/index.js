var path = require("path");
var fs = require("fs");
var express = require('express');
var morgan = require("morgan");
var app      = express();
var httpProxy = require('http-proxy');
var apiProxy = httpProxy.createProxyServer();
var profileService = 'http://profile:5000',
    authenticationService = 'http://authentication:5001',
    photoService = 'http://photo:5002',
    voteService = 'http://vote:5003',
    commentService = 'http://comment:5004';


var accessLogStream = fs.createWriteStream(__dirname + '/access.log', {flags: 'a'});
app.use(morgan('combined', {stream: accessLogStream}));

console.log(path.join(__dirname, '../webapp'));

app.use('/', express.static(path.join(__dirname, '../webapp')));

app.all("/image*", function(req, res) {
    console.log('redirecting to PhotoService');
    apiProxy.web(req, res, {target: photoService});
});

app.all("/token-auth*", function(req, res) {
    console.log('redirecting to authenticationService');
    apiProxy.web(req, res, {target: authenticationService});
});
app.all("/users*", function(req, res) {
    console.log('redirecting to profileService');
    apiProxy.web(req, res, {target: profileService});
});

app.all("/votes*", function(req, res) {
    console.log('redirecting to voteService');
    apiProxy.web(req, res, {target: voteService});
});

app.all("/comments*", function(req, res) {
    console.log('redirecting to commentService');
    apiProxy.web(req, res, {target: commentService});
});

app.listen(4999, function () {
    console.log('Example app listening on port 4999!')
});