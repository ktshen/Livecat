var async = require('async');
var redis_model = require('../models/redis_command');

// Create array sequence(for page navigation)
exports.create_array = function(start, end) {
    var array = [];
    for (var i=start; i<=end; i++) {
        array.push(i);
    }
    return array;
}

exports.count_pages = function(target, cb){
	var page_num;
	var data_per_page = 30;

	async.series([
		function(callback) {
			page_num = redis_model.redis_llen(target, function(){
				callback(null, "one");
			});
		},
		function(callback) {
			Promise.resolve(page_num).then(function(num){
				page_num = parseInt(num/data_per_page)+1;
				console.log("2: "+page_num);
				callback(null, "two");
			});
		}
	  ], function(err) { // This function gets called after the previous tasks have called their "task callbacks"
			if (err) return next(err);
			// console.log(page_num);
			// return page_num;
			cb();

	  });
	  return page_num;
}
