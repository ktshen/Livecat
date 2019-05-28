var express = require('express');
var router = express.Router();

// redis settings
var redis = require('redis');
var client = redis.createClient(6379, "127.0.0.1"); // this creates a new client
var msg;

// check if connected to redis
client.on('connect', function() {
    console.log('Redis client connected');
});

client.on('error', function (err) {
    console.log('Something went wrong ' + err);
});


// access redis by SET
// redis.print Shows "Reply: OK" on terminal
client.set('youtube', JSON.stringify({ name: "youtube", class: "video", time: "2019-01-19 13:07:08", intro: "Have fun with my youtube" , link: "www.youtube.com"}
), redis.print);
client.get('youtube', function (error, result) {
    if (error) {
        console.log(error);
        throw error;
    }
    console.log('GET result ->' + result);
});

/* GET users listing. */
router.get('/', function(req, res, next) {
  res.send(msg+'hihihihihi');
});

module.exports = router;
