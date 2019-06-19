var async = require('async');
var mongoose = require('mongoose');
var streamModel = require('../models/streamModel');
var util = require('./util');
var config = require('../config');
var admin = require('firebase-admin');

const utf8 = require('utf8');

var hostData = streamModel.Host;

exports.save_hostform = function(req, res){
    console.log("[FORM BODY]:");
    console.log(req.body);

    var host = new hostData({
            platform:req.body.platform,
        	pagelink:req.body.pagelink,
        	account:req.body.account,
        	chatroom:req.body.chatroom,
        	email:req.body.email
        });
       host.save(function (err, res) {
           if (err) return console.error(err);
       });

    res.send("success!");
    return;
}
