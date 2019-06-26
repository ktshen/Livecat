var async = require('async');
var mongoose = require('mongoose');
var streamModel = require('../models/streamModel');
var util = require('./util');
var config = require('../config');
var admin = require('firebase-admin');
var fs = require('fs');

const utf8 = require('utf8');

var hostData = streamModel.Host;

exports.save_hostform = function(req, res){
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

       res.render( 'success', {
           title : config.web.title,
           baseurl : req.path,
           isLogin : false,
       });
    return;
}

exports.save_manual = function(req, res){
    console.log(req.body);
    console.log(req.body.keyword);
    fs.writeFile('../public/manual.txt', JSON.stringify(req.body), function (err) {
        if (err)
            console.log(err);
        else
            console.log('Write operation complete.');
    });
    return;
}
