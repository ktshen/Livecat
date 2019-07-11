var async = require('async');
var mongoose = require('mongoose');
var streamModel = require('../models/streamModel');
var util = require('./util');
var config = require('../config');
const request = require('request');
const utf8 = require('utf8');
const utils = require('util')

//const

// mongoDB Models
var liveStreams = streamModel.liveStreams;
var webPlatformData = streamModel.webPlatform;
var userData = streamModel.User;


exports.save_history = function(req, res){
    // save user video info to mongoDB
    if(req.session.idToken){
        userData.updateOne({uid:req.session.idToken}, { $push: { history: req.body } } , function(err, docs){
            if(err) console.log(err);
            console.log('history time updated');
        });
    }
    return;
}

exports.save_click = function(req, res){
    console.log("[BODY]: "+req);
    // elastic DB video click through ++
    request('http://120.126.16.88:17777/update_click_through?videourl='+req.body.videourl, { json: true }, (err, res, body) => {
      if (err) {
          console.log(err);
      }
    });

    // save user clicked info
    if(req.session.idToken){
        // user clicked video info
        userData.updateOne({uid:req.session.idToken}, { $push: { click_through: {videourl:req.body.videourl,time:new Date().toISOString()} } }, function(err, docs){
            if(err) console.log(err);
        });

        // increase user clicked platform
        switch(req.body.platform) {
            case "西瓜直播":
                userData.updateOne({uid:req.session.idToken}, { $inc: { "statistics.platform.Xigua":1}  }, function(err, docs){
                    if(err) console.log(err);
                });
                break;
            case "Facebook":
                userData.updateOne({uid:req.session.idToken}, { $inc: { "statistics.platform.Facebook":1}  }, function(err, docs){
                    if(err) console.log(err);
                });
                break;
            case "17直播":
                userData.updateOne({uid:req.session.idToken}, { $inc: { "statistics.platform.Live17":1}  }, function(err, docs){
                    if(err) console.log(err);
                });
                break;
            case "YouTube":
                userData.updateOne({uid:req.session.idToken}, { $inc: { "statistics.platform.YouTube":1}  }, function(err, docs){
                    if(err) console.log(err);
                });
                break;
            case "Twitch":
                userData.updateOne({uid:req.session.idToken}, { $inc: { "statistics.platform.Twitch":1}  }, function(err, docs){
                    if(err) console.log(err);
                });
                break;
            default:
                break;
        }
    }
    res.send({"result":"save done"});
    return;
}
