var async = require('async');
var redisModel = require('../models/redis_command');
var util = require('./util');

const dataPerPage = 40;
const choicePerPage = 30;
const sortedDataPerSubject = 4;

// for index page
exports.get_home = function(req, res) {
    var streamListArr = [];
    var filterListArr = [];

    async.series([
        function(callback) { // Show data on the page user selected
            if (!req.query.pg){ // Get page number chosen by user
                dataPage = 1;
            }else{
                dataPage = req.query.pg;
            }
            start = (dataPage - 1)*dataPerPage; // which data start asking from DB
            end = dataPage*dataPerPage - 1; // end by which data

            redisModel.redis_select("0");
            streamList = redisModel.redis_lrange("All", start, end, function(){ // query data from DB
                callback(null, "one");
            });
        },
        function(callback) { // set filter
            redisModel.redis_select("0");
            filterList = redisModel.redis_smembers("Filter", function(){ // Get the filter list from DB
                callback(null, "two");
            });
        },
        function(callback) {
            redisModel.redis_select("0");
            dataNum = redisModel.redis_llen("All", function(){ // To count how many data in DB and decide how many pages can be show
      			callback(null, "three");
      		});
        },
        ], function(err) { // This function gets called after the previous tasks have called their "task callbacks"
            if (err) return next(err);
            // Because function to query DB is written in async, those finctions return promises and have to be resolved
            // set stream list
            Promise.resolve(streamList).then(function(redisMsgArr){
                redisMsgArr.forEach(
                    function(element) {
                        jsonMsg = JSON.parse(element); // parse data to json object
                        streamListArr.push(jsonMsg); // save json object into array
                });
                // set filter list
                Promise.resolve(filterList).then(function(redisMsgArr){
                    filterListArr = redisMsgArr;
                    // set page navigation
                    Promise.resolve(dataNum).then(function(redisMsg){
                        if(redisMsg%dataPerPage == 0){
                            pageNum = parseInt(redisMsg/dataPerPage); // number of pagination
                        }else{
                            pageNum = parseInt(redisMsg/dataPerPage +1); // number of pagination
                        }
                        console.log("[PAGE NUM] " + pageNum);

                        res.render( 'index', {
                            title : 'Live stream',
                            baseurl : req.path,
                            choices : "",
                            posts : streamListArr,
                            filter : filterListArr,
                            pages : util.create_array(1, pageNum),
                            currentPage : req.query.pg
                            });

                    }).catch(function(error){
                        console.error(error);
                    });
                }).catch(function(error){
                    console.error(error);
                });
            }).catch(function(error){
                console.error(error);
            });
    });
};

// for Host page
exports.get_host = function(req, res) {
        var streamListArr = [];
        var filterListArr = [];

        async.series([
            function(callback) { // Show data on the page user selected
                if (!req.query.pg){ // Get page number chosen by user
                    dataPage = 1;
                }else{
                    dataPage = req.query.pg;
                }
                start = (dataPage - 1)*dataPerPage; // which data start asking from DB
                end = dataPage*dataPerPage - 1; // end by which data

                redisModel.redis_select("0");
                streamList = redisModel.redis_lrange("All", start, end, function(){ // query data from DB
                    callback(null, "one");
                });
            },
            function(callback) { // set filter
                redisModel.redis_select("0");
                filterList = redisModel.redis_smembers("Filter", function(){ // Get the filter list from DB
                    callback(null, "two");
                });
            },
            function(callback) { // To count how many data in DB and decide how many pages can be show
                redisModel.redis_select("0");
                dataNum = redisModel.redis_llen("All", function(){
      				callback(null, "three");
      			});
            },function(callback) { // To get filted data
                redisModel.redis_select("0");
                    hostList = redisModel.redis_srandmember("Host", choicePerPage, function(){ // (filted by "Host")
                        callback(null, "four");
                    });
            },
          ], function(err) { // This function gets called after the previous tasks have called their "task callbacks"
                if (err) return next(err);
                // Because function to query DB is written in async, those finctions return promises and have to be resolved
                    // set stream list
                    Promise.resolve(streamList).then(function(redisMsgArr){
                        redisMsgArr.forEach(
                            function(element) {
                                jsonMsg = JSON.parse(element);
                                streamListArr.push(jsonMsg);
                            });
                            // set filter list
                            Promise.resolve(filterList).then(function(redisMsgArr){
                                filterListArr = redisMsgArr;
                                // set page navigation
                                Promise.resolve(dataNum).then(function(redisMsg){
                                    if(redisMsg%dataPerPage == 0){
                                        pageNum = parseInt(redisMsg/dataPerPage); // number of pagination
                                    }else{
                                        pageNum = parseInt(redisMsg/dataPerPage +1); // number of pagination
                                    }
                                    console.log("[PAGE NUM]\n " + pageNum);
                                    Promise.resolve(hostList).then(function(redisMsg){
                                        hostListArr = redisMsg;
                                        console.log("[HOST_LIST]\n " + hostListArr);

                                        res.render( 'general', {
                                            title : 'Live stream',
                                            baseurl : req.path,
                                            choices : hostListArr,
                                            posts : streamListArr,
                                            filter : filterListArr,
                                            pages : util.create_array(1, pageNum),
                                            currentPage : req.query.pg
                                            });
                                    }).catch(function(error){
                                        console.error(error);
                                    });
                                }).catch(function(error){
                                    console.error(error);
                                });
                            }).catch(function(error){
                                console.error(error);
                            });
                    }).catch(function(error){
                        console.error(error);
                    });
          });
};

// for Platform page
exports.get_platform = function(req, res) {
    var streamListArr = [];
    var filterListArr = [];

    async.waterfall([
        function(next) {
            redisModel.redis_select("0");
            allPlatformsPromise = redisModel.redis_smembers("Platform", function(){ // Get all platforms saved in DB
                return next(null, allPlatformsPromise);
  			});
        },
        function(allPlatformsPromise, next) { // Parse the platform promise into string array
            Promise.resolve(allPlatformsPromise).then(function(redisMsg){
                platformListArr = redisMsg;
                return next(null, platformListArr);
            }).catch(function(error){
                console.error(error);
            });
        },function(platformListArr, next) { // Get data from each platform
            redisModel.redis_select("3");
            platformListArr.forEach(
                function(element) {
                    async.series([
                        function(callback) {
                            tmp = redisModel.redis_lrange(element, 0, sortedDataPerSubject-1, function(){ // get 8 data per platform
                                callback(null, "one");
                            });
                        },
                      ], function(err) {
                            if (err) return next(err);
                            Promise.resolve(tmp).then(function(redisMsgArr){
                                redisMsgArr.forEach(
                                    function(element) {
                                        jsonMsg = JSON.parse(element);
                                        streamListArr.push(jsonMsg);
                                    });
                                    streamListArr.push("0"); // Save "0" to array at the end of a pltform
                                return next(null, null);
                            }).catch(function(error){
                                console.error(error);
                            });
                    });
            });
        },function(platformListArr, next) {
            redisModel.redis_select("0");
            filterList = redisModel.redis_smembers("Filter", function(){ // Get the filter list from DB
                return next(null, null);
            });
        },
      ],function(err) { // This function gets called after the previous tasks have called their "task callbacks"
            if (err) return next(err);
            Promise.resolve(filterList).then(function(redisMsgArr){
                filterListArr = redisMsgArr;

                res.render( 'sortedpage', {
                    title : 'Live stream',
                    subtitles : platformListArr,
                    stream : streamListArr,
                    baseurl : req.path,
                    filter : filterListArr,
                    });

            }).catch(function(error){
                    console.error(error);
            });
      });
};

// for Category page
exports.get_category = function(req, res) {
    var streamListArr = [];
    var filterListArr = [];

    async.series([
        function(callback)  { // Show data on the page user selected
            if (!req.query.pg){ // Get page number chosen by user
                dataPage = 1;
            }else{
                dataPage = req.query.pg;
            }
            start = (dataPage - 1)*dataPerPage; // which data start asking from DB
            end = dataPage*dataPerPage - 1; // end by which data

            redisModel.redis_select("0");
            streamList = redisModel.redis_lrange("All", start, end, function(){ // query data from DB
                callback(null, "one");
            });
        },
        function(callback) { // set filter
            redisModel.redis_select("0");
            filterList = redisModel.redis_smembers("Filter", function(){ // Get filter list from DB
                callback(null, "two");
            });
        },
        function(callback) { // To count how many data in DB and decide how many pages can be show
            redisModel.redis_select("0");
            dataNum = redisModel.redis_llen("All", function(){
  				callback(null, "three");
  			});
        },function(callback) { // To get filter object (testing)
            redisModel.redis_select("0");
            platformList = redisModel.redis_srandmember("Category", choicePerPage, function(){
                callback(null, "four");
            });
        },
      ], function(err) { // This function gets called after the previous tasks have called their "task callbacks"
            if (err) return next(err);
            // Because function to query DB is written in async, those finctions return promises and have to be resolved
                // set stream list
                Promise.resolve(streamList).then(function(redisMsgArr){
                    redisMsgArr.forEach(
                        function(element) {
                            jsonMsg = JSON.parse(element);
                            streamListArr.push(jsonMsg);
                    });
                    // set filter list
                    Promise.resolve(filterList).then(function(redisMsgArr){
                        filterListArr = redisMsgArr;
                        // set page navigation
                        Promise.resolve(dataNum).then(function(redisMsg){
                            if(redisMsg%dataPerPage == 0){
                                pageNum = parseInt(redisMsg/dataPerPage); // number of pagination
                            }else{
                                pageNum = parseInt(redisMsg/dataPerPage +1); // number of pagination
                            }
                            console.log("[PAGE NUM]\n " + pageNum);
                            Promise.resolve(platformList).then(function(redisMsg){
                                platformListArr = redisMsg;
                                console.log("[PLATFORM_LIST]\n " + platformListArr);

                                res.render( 'general', {
                                    title : 'Live stream',
                                    choices : platformListArr,
                                    baseurl : req.path,
                                    posts : streamListArr,
                                    filter : filterListArr,
                                    pages : util.create_array(1, pageNum),
                                    currentPage : req.query.pg
                                    });

                            }).catch(function(error){
                                console.error(error);
                            });
                        }).catch(function(error){
                            console.error(error);
                        });
                    }).catch(function(error){
                        console.error(error);
                    });
            }).catch(function(error){
                console.error(error);
            });
      });
};

exports.get_host_filted_list = function(req, res) {
    var streamListArr = [];
    var filterListArr = [];
    var hostListArr = [];

    async.series([
        function(callback) { // Show data on the page user selected
            if (!req.query.pg){ // Get page number chosen by user
                dataPage = 1;
            }else{
                dataPage = req.query.pg;
            }
            start = (dataPage - 1)*dataPerPage; // which data start asking from DB
            end = dataPage*dataPerPage - 1; // end by which data

            redisModel.redis_select("2");
            streamList = redisModel.redis_lrange(req.params.host, start, end, function(){
                callback(null, "one");
            });
        },
        function(callback) { // set filter
            redisModel.redis_select("0");
            filterList = redisModel.redis_smembers("Filter", function(){
                callback(null, "two");
            });
        },
        function(callback) { // To count how many data in DB and decide how many pages can be show
            redisModel.redis_select("2");
            dataNum = redisModel.redis_llen(req.params.host, function(){
                callback(null, "three");
            });
        },
      ], function(err) { // This function gets called after the previous tasks have called their "task callbacks"
            if (err) return next(err);
            // Because function to query DB is written in async, those finctions return promises and have to be resolved
            // set stream list
            Promise.resolve(streamList).then(function(redisMsgArr){
                redisMsgArr.forEach(
                    function(element) {
                        jsonMsg = JSON.parse(element);
                        streamListArr.push(jsonMsg);
                });
                // set filter list
                Promise.resolve(filterList).then(function(redisMsgArr){
                    filterListArr = redisMsgArr;
                    // set page navigation
                    Promise.resolve(dataNum).then(function(redisMsg){
                        if(redisMsg%dataPerPage == 0){
                            pageNum = parseInt(redisMsg/dataPerPage); // number of pagination
                        }else{
                            pageNum = parseInt(redisMsg/dataPerPage +1); // number of pagination
                        }
                        console.log("[PAGE NUM]\n " + pageNum);

                            res.render( 'general', {
                                title : 'Live stream',
                                choices : "",
                                baseurl : req.path,
                                posts : streamListArr,
                                filter : filterListArr,
                                pages : util.create_array(1, pageNum),
                                currentPage : req.query.pg
                                });

                    }).catch(function(error){
                        console.error(error);
                    });
                }).catch(function(error){
                    console.error(error);
                });
        }).catch(function(error){
            console.error(error);
        });
    });
};

exports.get_platform_filted_list = function(req, res) {
    var streamListArr = [];
    var filterListArr = [];
    var hostListArr = [];

    async.series([
        function(callback) { // Show data on the page user selected
            if (!req.query.pg){ // Get page number chosen by user
                dataPage = 1;
            }else{
                dataPage = req.query.pg;
            }
            start = (dataPage - 1)*dataPerPage; // which data start asking from DB
            end = dataPage*dataPerPage - 1; // end by which data

            redisModel.redis_select("3");
            streamList = redisModel.redis_lrange(req.params.platform, start, end, function(){
                callback(null, "one");
            });
        },
        function(callback) { // set filter
            redisModel.redis_select("0");
            filterList = redisModel.redis_smembers("Filter", function(){
                callback(null, "two");
            });
        },
        function(callback) { // To count how many data in DB and decide how many pages can be show
            redisModel.redis_select("3");
            dataNum = redisModel.redis_llen(req.params.platform, function(){
                callback(null, "three");
            });
        },
      ], function(err) { // This function gets called after the previous tasks have called their "task callbacks"
            if (err) return next(err);
            // Because function to query DB is written in async, those finctions return promises and have to be resolved
            // set stream list
            Promise.resolve(streamList).then(function(redisMsgArr){
                redisMsgArr.forEach(
                    function(element) {
                        jsonMsg = JSON.parse(element);
                        streamListArr.push(jsonMsg);
                });
                // set filter list
                Promise.resolve(filterList).then(function(redisMsgArr){
                    filterListArr = redisMsgArr;
                    // set page navigation
                    Promise.resolve(dataNum).then(function(redisMsg){
                        if(redisMsg%dataPerPage == 0){
                            pageNum = parseInt(redisMsg/dataPerPage); // number of pagination
                        }else{
                            pageNum = parseInt(redisMsg/dataPerPage +1); // number of pagination
                        }
                        console.log("[PAGE NUM]\n " + pageNum);  // number of pagination

                            res.render( 'general', {
                                title : 'Live stream',
                                choices : "",
                                baseurl : req.path,
                                posts : streamListArr,
                                filter : filterListArr,
                                pages : util.create_array(1, pageNum),
                                currentPage : req.query.pg
                                });

                    }).catch(function(error){
                        console.error(error);
                    });
            }).catch(function(error){
                console.error(error);
            });
        }).catch(function(error){
            console.error(error);
        });
    });
};

exports.get_category_filted_list = function(req, res) {
    var streamListArr = [];
    var filterListArr = [];
    var hostListArr = [];

    async.series([
        function(callback) { // Show data on the page user selected
            if (!req.query.pg){ // Get page number chosen by user
                dataPage = 1;
            }else{
                dataPage = req.query.pg;
            }
            start = (dataPage - 1)*dataPerPage; // which data start asking from DB
            end = dataPage*dataPerPage - 1; // end by which data

            redisModel.redis_select("1");
            streamList = redisModel.redis_lrange(req.params.category, start, end, function(){
                callback(null, "one");
            });
        },
        function(callback) { // set filter
            redisModel.redis_select("0");
            filterList = redisModel.redis_smembers("Filter", function(){
                callback(null, "two");
            });
        },
        function(callback) { // To count how many data in DB and decide how many pages can be show
            redisModel.redis_select("1");
            dataNum = redisModel.redis_llen(req.params.category, function(){
                callback(null, "three");
            });
        },
      ], function(err) { // This function gets called after the previous tasks have called their "task callbacks"
            if (err) return next(err);
            // Because function to query DB is written in async, those finctions return promises and have to be resolved
            // set stream list
            Promise.resolve(streamList).then(function(redisMsgArr){
                redisMsgArr.forEach(
                    function(element) {
                        jsonMsg = JSON.parse(element);
                        streamListArr.push(jsonMsg);
                });
                // set filter list
                Promise.resolve(filterList).then(function(redisMsgArr){
                    filterListArr = redisMsgArr;
                    // set page navigation
                    Promise.resolve(dataNum).then(function(redisMsg){
                        if(redisMsg%dataPerPage == 0){
                            pageNum = parseInt(redisMsg/dataPerPage); // number of pagination
                        }else{
                            pageNum = parseInt(redisMsg/dataPerPage +1); // number of pagination
                        }
                        console.log("[PAGE NUM]\n " + pageNum);  // number of pagination

                            res.render( 'general', {
                                title : 'Live stream',
                                choices : "",
                                baseurl : req.path,
                                posts : streamListArr,
                                filter : filterListArr,
                                pages : util.create_array(1, pageNum),
                                currentPage : req.query.pg
                                });

                    }).catch(function(error){
                        console.error(error);
                    });
                }).catch(function(error){
                    console.error(error);
                });
        }).catch(function(error){
            console.error(error);
        });
    });
};
