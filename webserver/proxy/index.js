var path = require("path");
var fs = require("fs");
var env = require('node-env-file');

// Does not overwrite existing environment variables.
if (fs.existsSync(__dirname + '/.env')) {
    env(__dirname + '/.env');
    console.log(env.data)
}
var express = require('express');
var morgan = require("morgan");
var app      = express();
var httpProxy = require('http-proxy');
var apiProxy = httpProxy.createProxyServer();
var profileService =  process.env.PROFILE_URL,
    authenticationService =  process.env.AUTHENTICATION_URL,
    photoService =  process.env.PHOTO_URL,
    voteService =  process.env.VOTE_URL,
    commentService =  process.env.COMMENT_URL;

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
app.listen(process.env.PORT, function () {
    console.log('Proxy listening on port %s!', process.env.PORT)
});