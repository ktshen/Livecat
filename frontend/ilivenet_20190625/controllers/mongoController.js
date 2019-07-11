var async = require('async');
var mongoose = require('mongoose');
var streamModel = require('../models/streamModel');
var util = require('./util');
var config = require('../config');
var admin = require('firebase-admin');
const request = require('request');
const utf8 = require('utf8');
var fs = require('fs');

//const
const dataPerPage_unlogin = 18;
const dataPerPage_login = 30;
const platformList = ["Twitch", "Facebook", "YouTube", "17直播", "西瓜直播"]
var isLogin = false;

// mongoDB Models
var liveStreams = streamModel.liveStreams;
var webPlatformData = streamModel.webPlatform;
var userData = streamModel.User;
var keywordData = streamModel.Keyword;

// elasticsearch connection
// var elasticsearch = require('elasticsearch');
// var client = new elasticsearch.Client(config.elasticsearch.config);


// get homepage
exports.get_homepage= function(req, res) {
    var newEventListArr = [];
    var streamListArr = [];
    var user;

    async.series([
        function(callback) { // check session
            if(req.session.idToken){
                isLogin = true;// [for anonymous login] login without authentication

                // check if user exists
                userData.findOne({ uid: req.session.idToken }, function(err, result) {
                    if (err) { // if error, return error
                        return callback(err)
                    } else if (!result) { // if user not exists
                        callback(null,"one");
                    } else{
                        // able to login
                        isLogin = true;
                        user = result;
                        callback(null,"one");
                    }
                  });
              }else{
                  isLogin = false;
                    callback(null,"one");
                  }
        },
        function(callback) { // for not logged in
            if(!isLogin){
                request('http://'+config.queryAPI.server+'/cover_page', { json: true }, (err, res, body) => {
                    if (err) {
                        console.log(err);
                        callback(null, "one");
                    }else if (body){
                        body.forEach(
                                function(element) {
                                    newEventListArr.push(element._source); // save json object into array
                            });
                      callback(null, "one");
                    } else {
                        callback(null, "one");
                    }
                });
            }else{
                callback(null,"one");
            }
        },
        function(callback) { // for not logged in
            if(!isLogin){
                request('http://'+config.queryAPI.server+'/?platform=Twitch&from=0&size=3', { json: true }, (err, res, body) => {
                    if (err) {
                        console.log(err);
                        callback(null, "one");
                    }else if (body && body.hits && body.hits.hits){
                        streamList = body.hits.hits;
                        streamList.forEach(
                                function(element) {
                                    streamListArr.push(element._source); // save json object into array
                            });
                      callback(null, "one");
                    } else {
                        callback(null, "one");
                    }
                });
            }else{
                callback(null,"one");
            }
        },
        function(callback) { // for not logged in
            if(!isLogin){
                request('http://'+config.queryAPI.server+'/?platform=YouTube&from=0&size=3', { json: true }, (err, res, body) => {
                    if (err) {
                        console.log(err);
                        callback(null, "one");
                    }else if (body && body.hits && body.hits.hits){
                        streamList = body.hits.hits;
                        streamList.forEach(
                                function(element) {
                                    streamListArr.push(element._source); // save json object into array
                            });
                      callback(null, "one");
                    } else {
                        callback(null, "one");
                    }
                });
            }else{
                callback(null,"one");
            }
        },
        function(callback) { // for not logged in
            if(!isLogin){
                request('http://'+config.queryAPI.server+'/?platform=Facebook&from=0&size=2', { json: true }, (err, res, body) => {
                    if (err) {
                        console.log(err);
                        callback(null, "one");
                    }else if (body && body.hits && body.hits.hits){
                        streamList = body.hits.hits;
                        streamList.forEach(
                                function(element) {
                                    streamListArr.push(element._source); // save json object into array
                            });
                      callback(null, "one");
                    } else {
                        callback(null, "one");
                    }
                });
            }else{
                callback(null,"one");
            }
        },
        function(callback) { // for not logged in
            if(!isLogin){
                request('http://'+config.queryAPI.server+'/?platform='+utf8.encode('17直播')+'&from=0&size=2', { json: true }, (err, res, body) => {
                    if (err) {
                        console.log(err);
                        callback(null, "one");
                    }else if (body && body.hits && body.hits.hits){
                        streamList = body.hits.hits;
                        streamList.forEach(
                                function(element) {
                                    streamListArr.push(element._source); // save json object into array
                            });
                      callback(null, "one");
                    } else {
                        callback(null, "one");
                    }
                });
            }else{
                callback(null,"one");
            }
        },
        function(callback) { // for not logged in
            if(!isLogin){
                request('http://'+config.queryAPI.server+'/?platform='+utf8.encode('西瓜直播')+'&from=0&size=2', { json: true }, (err, res, body) => {
                    if (err) {
                        console.log(err);
                        callback(null, "one");
                    }else if (body && body.hits && body.hits.hits){
                        streamList = body.hits.hits;
                        streamList.forEach(
                                function(element) {
                                    streamListArr.push(element._source); // save json object into array
                            });
                      callback(null, "one");
                    } else {
                        callback(null, "one");
                    }
                });
            }else{
                callback(null,"one");
            }
        },

        // for logged in
        function(callback) {
            if(isLogin){
                    request('http://'+config.queryAPI.server+'/home_page', { json: true }, (err, res, body) => {
                      if (err) {
                          console.log(err);
                          callback(null, "one");
                      }else if (body){
                         streamListArr = body;
                        callback(null, "one");
                      } else {
                          callback(null, "one");
                      }
                    });
            }else{
                callback(null,"one");
            }
        }
        ], function(err) { // This function gets called after the previous tasks have called their "task callbacks"
            if (err) return next(err);

            if(!isLogin){
                res.render( 'index_unlogin', {
                    title : config.web.title,
                    baseurl : req.path,
                    posts : streamListArr,
                    newEventListArr : newEventListArr,
                    currentPage : req.query.pg,
                    isLogin :isLogin,
                    user:user
                });
            }else{
                res.render( 'index', {
                    title : config.web.title,
                    baseurl : req.path,
                    posts : streamListArr,
                    isLogin :isLogin,
                    user:user
                });
            }
        });
};

// get all data
exports.get_all= function(req, res) {
    var streamListArr = [];
    var dataNum, pageNum, onlineNum;
    var resultfound;
    var user;

    async.series([
        function(callback) { // check session
            if(req.session.idToken){
                isLogin = true;// [for quick login]

                // check if user exists
                userData.findOne({ uid: req.session.idToken }, function(err, result) {
                    if (err) { // if error, return error
                        return callback(err)
                    } else if (!result) { // if user not exists
                        callback(null,"one");
                    } else{
                        // able to login
                        isLogin = true;
                        user = result;
                        callback(null,"one");
                    }
                  });
              }else{
                  isLogin = false;
                    callback(null,"one");
                  }
        },
        function(callback) {
            dataPage = util.get_query_page_es(req.query.pg, dataPerPage_login)
            request('http://'+config.queryAPI.server+'/?platform=all&from='+dataPage+'&size='+dataPerPage_login, { json: true }, (err, res, body) => {
              if (err) {
                  console.log(err);
                  callback(null, "one");
              }else if (body && body.hits && body.hits.hits){
                  resultfound = body.found;
                  streamList = body.hits.hits;
                  streamList.forEach(
                          function(element) {
                              streamListArr.push(element._source); // save json object into array
                      });
                callback(null, "one");
              } else {
                  callback(null, "one");
              }
            });
        },
        function(callback) { // count data num
            request('http://'+config.queryAPI.server+'/?platform=all&from='+dataPage+'&size=500', { json: true }, (err, res, body) => {
              if (err) {
                  console.log(err);
                  callback(null, "one");
              }else if (body && body.hits && body.hits.hits){
                  dataNum = body.hits.hits.length;
                  callback(null, "one");
              } else {
                  callback(null, "one");
              }
            });
        },
        function(callback) { // count data num
            request('http://'+config.queryAPI.server+'/total_streams', { json: true }, (err, res, body) => {
              if (err) {
                  console.log(err);
                  callback(null, "one");
              }else if (body && body.count){
                  onlineNum = body.count;
                  callback(null, "one");
              } else {
                  callback(null, "one");
              }
            });
        },
        ], function(err) { // This function gets called after the previous tasks have called their "task callbacks"
            if (err) return next(err);
            pageNum = util.count_pages(dataNum, dataPerPage_login);
            if(isLogin){
                res.render( 'spa', {
                    title : config.web.title,
                    pageTitle: "All Streams",
                    baseurl : req.path,
                    posts : streamListArr,
                    onlineNum : onlineNum,
                    resultfound : resultfound,
                    pages : util.create_array(1, pageNum),
                    currentPage : req.query.pg,
                    isLogin :isLogin,
                    user:user
                });
            }else{
                res.redirect("/");
            }
        });
};

exports.search = function(req, res) {
   // const search_func = async(query, cb) => {
   //     const eres = await client.search({
   //       index: 'livestreams',
   //       size: 20,
   //       body: {
   //          query: {
   //            multi_match: {
   //                query: query,
   //                fields: ["title", "description", "platform","host"],
   //                fuzziness: "AUTO"
   //            }
   //          }
   //        }
   //     });
   //     cb();
   //     return eres.hits.hits;
   // }

    var streamListArr = [];
    var resultfound;
    var dataNum;
    var pageNum;
    var user;

    async.series([
        function(callback) { // check session
            if(req.session.idToken){
                isLogin = true;// [for quick login]

                // check if user exists
                userData.findOne({ uid: req.session.idToken }, function(err, result) {
                    if (err) { // if error, return error
                        return callback(err)
                    } else if (!result) { // if user not exists
                        callback(null,"one");
                    } else{
                        // able to login
                        isLogin = true;
                        user = result;
                        callback(null,"one");
                    }
                  });
              }else{
                  isLogin = false;
                    callback(null,"one");
                  }
        },
        function(callback){ // save user searched keyword
            if(req.session.idToken){
                userData.updateOne({uid:req.session.idToken}, { $push: { search: { keyword: req.query.search } } } , function(err, docs){
                    if(err) console.log(err);
                    console.log('save keyword');
                    callback(null,"one");
                });
            }else{
                  callback(null,"one");
                }
        },
        function(callback){ //record keyword searched frequency in mongoDB
            keywordData.findOne({ keyword: req.query.search }, function(err, result) {
                if (err) { // if error, return error
                    return callback(err)
                } else if (!result) {
                    // if keyword not exists
                    keywordData.create({ keyword: req.query.search }, function(err, docs){
                        if(err) console.log(err);
                        callback(null,"one");
                    });
                } else{
                    // keyword exists
                    keywordData.updateOne({ keyword: req.query.search }, { $inc: { search_through: 1 } } , function(err, docs){
                        if(err) console.log(err);
                        callback(null,"one");
                    });
                }
              });
        },
        function(callback) {
            if(!isLogin){
                dataPage = util.get_query_page_es(req.query.pg, dataPerPage_unlogin)
                request('http://'+config.queryAPI.server+'?q='+utf8.encode(req.query.search)+'&from='+dataPage+'&size='+dataPerPage_unlogin, { json: true }, (err, res, body) => {
                  if (err) {
                      console.log(err);
                      callback(null, "one");
                  }else if (body && body.hits && body.hits.hits){
                      resultfound = body.found; // check if results are found for the keyword
                      streamList = body.hits.hits;
                      streamList.forEach(
                              function(element) {
                                  streamListArr.push(element._source); // save json object into array
                          });
                    callback(null, "one");
                  } else {
                      callback(null, "one");
                  }
                });
            } else{
                dataPage = util.get_query_page_es(req.query.pg, dataPerPage_login)
                request('http://'+config.queryAPI.server+'?q='+utf8.encode(req.query.search)+'&from='+dataPage+'&size='+dataPerPage_login, { json: true }, (err, res, body) => {
                  if (err) {
                      console.log(err);
                      callback(null, "one");
                  }else if (body && body.hits && body.hits.hits){
                      resultfound = body.found;
                      streamList = body.hits.hits;
                      streamList.forEach(
                              function(element) {
                                  streamListArr.push(element._source); // save json object into array
                          });
                    callback(null, "one");
                  } else {
                      callback(null, "one");
                  }
                });
            }
        },
        function(callback) { // count data num
            request('http://'+config.queryAPI.server+'?q='+utf8.encode(req.query.search)+'&from=0&size=5000', { json: true }, (err, res, body) => {
              if (err) {
                  console.log(err);
                  callback(null, "one");
              }else if (body && body.hits && body.hits.hits){
                  dataNum = body.hits.hits.length;
                  callback(null, "one");
              } else {
                  callback(null, "one");
              }
            });
        },
        ], function(err) { // This function gets called after the previous tasks have called their "task callbacks"
            if (err) return next(err);

            if(!isLogin){
                pageNum = util.count_pages(dataNum, dataPerPage_unlogin);
            } else{
                pageNum = util.count_pages(dataNum, dataPerPage_login);
            }
            res.render( 'searchpage', {
                title : config.web.title,
                baseurl : req.path,
                posts : streamListArr,
                resultfound : resultfound,
                keyword : req.query.search,
                pages : util.create_array(1, pageNum),
                currentPage : req.query.pg,
                isLogin :isLogin,
                user:user
            });
        });
};

exports.get_livestream = function(req, res) {
    var stream;
    var user;

    async.series([
        function(callback) { // check session
            if(req.session.idToken){
                isLogin = true;// [for quick login]

                // check if user exists
                userData.findOne({ uid: req.session.idToken }, function(err, result) {
                    if (err) { // if error, return error
                        return callback(err)
                    } else if (!result) { // if user not exists
                        callback(null,"one");
                    } else{
                        // able to login
                        isLogin = true;
                        user = result;
                        callback(null,"one");
                    }
                  });
              }else{
                  isLogin = false;
                    callback(null,"one");
                  }
        },
        function(callback) {
            request('http://'+config.queryAPI.server+'?channel='+utf8.encode(req.params.channel), { json: true }, (err, res, body) => {
              if (err) {
                  console.log(err);
                  callback(null, "one");
              }else if (body && body.hits && body.hits.hits){
                  if (body.hits.hits.length != 0){
                      stream = body.hits.hits[0]._source;
                      callback(null, "one");
                  }
              } else {
                  callback(null, "one");
              }

            });
        }
      ],function(err) { // This function gets called after the previous tasks have called their "task callbacks"
            if (err) return next(err);

            res.render( 'livestream', {
                title : config.web.title,
                baseurl : req.path,
                stream : stream,
                isLogin :isLogin,
                user:user
            });
      });
};

exports.get_all_platform= function(req, res) {
    var streamListArr = [];
    var user;
    var twitchNum, youtubeNum, facebookNum, xiguaNum, live17Num;

    async.series([
        function(callback) { // check session
            if(req.session.idToken){
                isLogin = true;// [for quick login] login without authentication

                // check if user exists
                userData.findOne({ uid: req.session.idToken }, function(err, result) {
                    if (err) { // if error, return error
                        return callback(err)
                    } else if (!result) { // if user not exists
                        callback(null,"one");
                    } else{
                        // able to login
                        isLogin = true;
                        user = result;
                        callback(null,"one");
                    }
                  });
              }else{
                  isLogin = false;
                    callback(null,"one");
                  }
        },
        function(callback) {
            if(isLogin){
                    request('http://'+config.queryAPI.server+'/platform_page', { json: true }, (err, res, body) => {
                      if (err) {
                          console.log(err);
                          callback(null, "one");
                      }else if (body){
                         streamListArr = body;
                        callback(null, "one");
                      } else {
                          callback(null, "one");
                      }
                    });
            }else{
                callback(null,"one");
            }
        },
        function(callback) { // count data num
            request('http://'+config.queryAPI.server+'/total_streams?platform='+utf8.encode("Twitch"), { json: true }, (err, res, body) => {
              if (err) {
                  console.log(err);
                  callback(null, "one");
              }else if (body && body.count){
                  twitchNum = body.count;
                  callback(null, "one");
              } else {
                  callback(null, "one");
              }
            });
        },
        function(callback) { // count data num
            request('http://'+config.queryAPI.server+'/total_streams?platform='+utf8.encode("Facebook"), { json: true }, (err, res, body) => {
              if (err) {
                  console.log(err);
                  callback(null, "one");
              }else if (body && body.count){
                  facebookNum = body.count;
                  callback(null, "one");
              } else {
                  callback(null, "one");
              }
            });
        },
        function(callback) { // count data num
            request('http://'+config.queryAPI.server+'/total_streams?platform='+utf8.encode("YouTube"), { json: true }, (err, res, body) => {
              if (err) {
                  console.log(err);
                  callback(null, "one");
              }else if (body && body.count){
                  youtubeNum = body.count;
                  callback(null, "one");
              } else {
                  callback(null, "one");
              }
            });
        },
        function(callback) { // count data num
            request('http://'+config.queryAPI.server+'/total_streams?platform='+utf8.encode("17直播"), { json: true }, (err, res, body) => {
              if (err) {
                  console.log(err);
                  callback(null, "one");
              }else if (body && body.count){
                  live17Num = body.count;
                  callback(null, "one");
              } else {
                  callback(null, "one");
              }
            });
        },
        function(callback) { // count data num
            request('http://'+config.queryAPI.server+'/total_streams?platform='+utf8.encode("西瓜直播"), { json: true }, (err, res, body) => {
              if (err) {
                  console.log(err);
                  callback(null, "one");
              }else if (body && body.count){
                  xiguaNum = body.count;
                  callback(null, "one");
              } else {
                  callback(null, "one");
              }
            });
        }
        ], function(err) { // This function gets called after the previous tasks have called their "task callbacks"
            if (err) return next(err);

            var streamsCounter = {
                twitch: twitchNum,
                youtube: youtubeNum,
                facebook: facebookNum,
                live17: live17Num,
                xigua: xiguaNum
            }

            res.render( 'platformpage', {
                title : config.web.title,
                baseurl : req.path,
                platforms:platformList,
                posts : streamListArr,
                streamsCounter : streamsCounter,
                currentPage : req.query.pg,
                isLogin :isLogin,
                user:user
            });
        });
};

// for Platform page
exports.get_platform = function(req, res) {
    var streamListArr = [];
    var resultfound;
    var dataNum;
    var pageNum;
    var onlineNum;
    var user;

    async.series([
        function(callback) { // check session
            if(req.session.idToken){
                isLogin = true;// [for quick login]

                // check if user exists
                userData.findOne({ uid: req.session.idToken }, function(err, result) {
                    if (err) { // if error, return error
                        return callback(err)
                    } else if (!result) { // if user not exists
                        callback(null,"one");
                    } else{
                        // able to login
                        isLogin = true;
                        user = result;
                        callback(null,"one");
                    }
                  });
              }else{
                  isLogin = false;
                    callback(null,"one");
                  }
        },
        function(callback) {
            dataPage = util.get_query_page_es(req.query.pg, dataPerPage_login)
            request('http://'+config.queryAPI.server+'/?platform='+utf8.encode(req.params.platform)+'&from='+dataPage+'&size='+dataPerPage_login, { json: true }, (err, res, body) => {
              if (err) {
                  console.log(err);
                  callback(null, "one");
              }else if (body && body.hits && body.hits.hits){
                  resultfound = body.found;
                  streamList = body.hits.hits;
                  streamList.forEach(
                          function(element) {
                              streamListArr.push(element._source); // save json object into array
                      });
                callback(null, "one");
              } else {
                  callback(null, "one");
              }
            });
        },
        function(callback) { // count data num
            request('http://'+config.queryAPI.server+'/?platform='+utf8.encode(req.params.platform)+'&from=0&size=5000', { json: true }, (err, res, body) => {
              if (err) {
                  console.log(err);
                  callback(null, "one");
              }else if (body && body.hits && body.hits.hits){
                  dataNum = body.hits.hits.length;
                  callback(null, "one");
              } else {
                  callback(null, "one");
              }
            });
        },
        function(callback) { // count data num
            request('http://'+config.queryAPI.server+'/total_streams?platform='+utf8.encode(req.params.platform), { json: true }, (err, res, body) => {
              if (err) {
                  console.log(err);
                  callback(null, "one");
              }else if (body && body.count){
                  onlineNum = body.count;
                  callback(null, "one");
              } else {
                  callback(null, "one");
              }
            });
        },
        ],function(err) { // This function gets called after the previous tasks have called their "task callbacks"
            if (err) return next(err);
            pageNum = util.count_pages(dataNum, dataPerPage_login);

            if(isLogin){
                res.render( 'spa', {
                    title : config.web.title,
                    pageTitle: req.params.platform,
                    onlineNum: onlineNum,
                    baseurl : req.path,
                    posts : streamListArr,
                    resultfound : resultfound,
                    pages : util.create_array(1, pageNum),
                    currentPage : req.query.pg,
                    isLogin :isLogin,
                    user:user
                    });
            }else{
                res.redirect("/");
            }
      });
};

// for privacy page
exports.get_privacy = function(req, res) {
    var streamListArr = [];
    var dataNum;
    var pageNum;
    var user;

    async.series([
        function(callback) { // check session
            if(req.session.idToken){
                isLogin = true;// [for quick login]

                // check if user exists
                userData.findOne({ uid: req.session.idToken }, function(err, result) {
                    if (err) { // if error, return error
                        return callback(err)
                    } else if (!result) { // if user not exists
                        callback(null,"one");
                    } else{
                        // able to login
                        isLogin = true;
                        user = result;
                        callback(null,"one");
                    }
                  });
              }else{
                  isLogin = false;
                    callback(null,"one");
                  }
        },
        ],function(err) { // This function gets called after the previous tasks have called their "task callbacks"
            if (err) return next(err);

            res.render( 'privacy', {
                title : config.web.title,
                baseurl : req.path,
                isLogin :isLogin,
                user:user
            });
      });
};

exports.get_dbinfo = function(req, res) {
    var streamListArr = [];
    var user;
    var keywordResult;
    var userResult;
    var twitchNum, youtubeNum, facebookNum, xiguaNum, live17Num;

    async.series([
        function(callback) { // check session
            if(req.session.idToken){
                isLogin = true;// [for quick login]

                // check if user exists
                userData.findOne({ uid: req.session.idToken }, function(err, result) {
                    if (err) { // if error, return error
                        return callback(err)
                    } else if (!result) { // if user not exists
                        callback(null,"one");
                    } else{
                        // able to login
                        isLogin = true;
                        user = result;
                        callback(null,"one");
                    }
                  });
              }else{
                  isLogin = false;
                    callback(null,"one");
                  }
        },
        function(callback) { // get all user data
            userData.find({},function(err, result) {
                if (err) throw err;
                userResult = result;
                callback(null,"one");
            });
        },
        function(callback){
            keywordData.find({}, function(err, result) {
                if (err) { // if error, return error
                    console.log(err);
                    return callback(err)
                }
                keywordResult = result;
                callback(null,"one")
              });
        },
        function(callback) { // count data num
            request('http://'+config.queryAPI.server+'/total_streams?platform='+utf8.encode("Twitch"), { json: true }, (err, res, body) => {
              if (err) {
                  console.log(err);
                  callback(null, "one");
              }else if (body && body.count){
                  twitchNum = body.count;
                  callback(null, "one");
              } else {
                  callback(null, "one");
              }
            });
        },
        function(callback) { // count data num
            request('http://'+config.queryAPI.server+'/total_streams?platform='+utf8.encode("Facebook"), { json: true }, (err, res, body) => {
              if (err) {
                  console.log(err);
                  callback(null, "one");
              }else if (body && body.count){
                  facebookNum = body.count;
                  callback(null, "one");
              } else {
                  callback(null, "one");
              }
            });
        },
        function(callback) { // count data num
            request('http://'+config.queryAPI.server+'/total_streams?platform='+utf8.encode("YouTube"), { json: true }, (err, res, body) => {
              if (err) {
                  console.log(err);
                  callback(null, "one");
              }else if (body && body.count){
                  youtubeNum = body.count;
                  callback(null, "one");
              } else {
                  callback(null, "one");
              }
            });
        },
        function(callback) { // count data num
            request('http://'+config.queryAPI.server+'/total_streams?platform='+utf8.encode("17直播"), { json: true }, (err, res, body) => {
              if (err) {
                  console.log(err);
                  callback(null, "one");
              }else if (body && body.count){
                  live17Num = body.count;
                  callback(null, "one");
              } else {
                  callback(null, "one");
              }
            });
        },
        function(callback) { // count data num
            request('http://'+config.queryAPI.server+'/total_streams?platform='+utf8.encode("西瓜直播"), { json: true }, (err, res, body) => {
              if (err) {
                  console.log(err);
                  callback(null, "one");
              }else if (body && body.count){
                  xiguaNum = body.count;
                  callback(null, "one");
              } else {
                  callback(null, "one");
              }
            });
        }
        ],function(err) { // This function gets called after the previous tasks have called their "task callbacks"
            if (err) return next(err);

            var streamsCounter = {
                twitch: twitchNum,
                youtube: youtubeNum,
                facebook: facebookNum,
                live17: live17Num,
                xigua: xiguaNum
            }
            if (user.email == config.rootuser){
                res.render( 'dbinfo', { //render 'ejs' , sending out the data object
                    title : config.web.title,
                    baseurl : req.path,
                    keywordResult : keywordResult,
                    userResult : userResult,
                    streamsCounter : streamsCounter,
                    isLogin :isLogin,
                    user:user
                    });
            }else{
                res.send("Permission Denied");
            }

      });
};

exports.get_hostform = function(req, res) {
    var streamListArr = [];
    var user;

    async.series([
        function(callback) { // check session
            if(req.session.idToken){
                isLogin = true;// [for quick login]

                // check if user exists
                userData.findOne({ uid: req.session.idToken }, function(err, result) {
                    if (err) { // if error, return error
                        return callback(err)
                    } else if (!result) { // if user not exists
                        callback(null,"one");
                    } else{
                        // able to login
                        isLogin = true;
                        user = result;
                        callback(null,"one");
                    }
                  });
              }else{
                  isLogin = false;
                    callback(null,"one");
                  }
        },
        ],function(err) { // This function gets called after the previous tasks have called their "task callbacks"
            if (err) return next(err);

            res.render( 'formpage', {
                title : config.web.title,
                baseurl : req.path,
                isLogin :isLogin,
                user:user
                });
                // res.send(dbresult);
      });
};

exports.get_manual = function(req, res) {
    var streamListArr = [];
    var user;
    var keys = [];

    async.series([
        function(callback) { // check session
            if(req.session.idToken){
                isLogin = true;// [for quick login]

                // check if user exists
                userData.findOne({ uid: req.session.idToken }, function(err, result) {
                    if (err) { // if error, return error
                        return callback(err)
                    } else if (!result) { // if user not exists
                        callback(null,"one");
                    } else{
                        // able to login
                        isLogin = true;
                        user = result;
                        callback(null,"one");
                    }
                  });
              }else{
                  isLogin = false;
                    callback(null,"one");
                  }
        },function(callback){
            fs.readFile('../public/manual.txt', function (err, data) {
                if (err) {
                    throw err;
                    callback(null,"one");
                }
                keys = JSON.parse(data.toString()).keyword;
                console.log("[KEYWORDS]: "+keys);
                callback(null,"one");
            });
        }
        ],function(err) { // This function gets called after the previous tasks have called their "task callbacks"
            if (err) return next(err);

            if (user.email == config.rootuser){
                res.render( 'manual_manager', {
                    title : config.web.title,
                    baseurl : req.path,
                    isLogin : isLogin,
                    keys : keys,
                    user: user
                    });
            }else{
                res.send("Permission Denied");
            }

      });
};
