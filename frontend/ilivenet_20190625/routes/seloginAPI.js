var express = require('express');
var seloginAPI = express.Router();
var streamModel = require('../models/streamModel');

var admin = require('firebase-admin');
var serviceAccount = require('../livecat-64d83-firebase-adminsdk-rdy5e-67f5048513');
var firebaseAdmin = admin.initializeApp({
    credential: admin.credential.cert(serviceAccount),
    databaseURL:'https://livecat-64d83.firebaseio.com'
})

var userData = streamModel.User;

seloginAPI.post('/login', function(req, res, next){
    var idToken = req.body.idToken;

    if(idToken){
        admin.auth().verifyIdToken(idToken) // verify to get user info from firebase
          .then(function(decodedToken) {
              // store user information if not exists
              userData.findOne({ uid: decodedToken.uid }, function(err, result) { // check if user exists
                  if (err) { // db query error
                    return callback(err)
                } else if (!result) { // user not exists
                      var user = new userData(
                          {    name:decodedToken.name,
                      	       picture:decodedToken.picture,
                      	       email:decodedToken.email,
                      	       uid:decodedToken.uid,
                      	 });
                         user.save(function (err, res) {
                             if (err) return console.error(err);
                         });
                  } else{ // user exists and update lastlogin time
                      userData.updateOne({uid:decodedToken.uid}, {lastlogin: Date.now()}, function(err, docs){
                          if(err) console.log(err);
                          console.log('login time updated');
                      });
                  }
                  console.log(result);
              });
              // create session ID
              req.session.idToken = decodedToken.uid;

              return res.send({"redirect":"/"});
          }).catch(function(error) {
            console.log(error);
            return res.send({"redirect":"/"});
          });
    }else{
        return res.json({ret_code: 1, ret_msg: 'login error'});
    }
});

seloginAPI.post('/quicklogin', function(req, res, next){
    var idToken = req.body.idToken;
    console.log(idToken);

    userData.findOne({ uid: idToken }, function(err, result) { // check if user exists
        if (err) { // db query error
          return callback(err)
      } else if (!result) { // user not exists
            var user = new userData(
                {    name:idToken,
                     picture:"/images/anonymous.png",
                     email:idToken,
                     uid:idToken,
               });
               user.save(function (err, res) {
                   if (err) return console.error(err);
               });
        } else{ // user exists and update lastlogin time
            userData.updateOne({uid:idToken.uid}, {lastlogin: Date.now()}, function(err, docs){
                if(err) console.log(err);
                console.log('login time updated');
            });
        }
        console.log(result);
    });

    if(idToken){
        req.session.idToken = idToken;
        return res.send({"redirect":"/"});
    }else{
        return res.json({ret_code: 1, ret_msg: 'login error'});
    }
});

seloginAPI.get('/signout', function(req, res) {
    req.session.destroy();
    console.log("delete");
    return res.send({"redirect":"/"});
});


module.exports = seloginAPI;
