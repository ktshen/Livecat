var mongoose = require('mongoose');
var config = require('../config');

var connCrawler = mongoose.createConnection(config.db.crawler, config.db.options);
var connWeb = mongoose.createConnection(config.db.web, config.db.options);
var connUser = mongoose.createConnection(config.db.user, config.db.options);
var connKeyword = mongoose.createConnection(config.db.keyword, config.db.options);
var connHost = mongoose.createConnection(config.db.host, config.db.options);

Schema = mongoose.Schema;

//stream data had been moved to elastic
var streamSchema = new Schema({
	title:{ type:String },
	description:{ type:String },
	platform:{ type:String },
	videoid:{ type:String },
	host:{ type:String },
	status:{ type:String },
	thumbnails:{ type:String },
	published:{ type:String },
	tags:{ type:String },
	generaltag:{ type:String },
	timestamp:{ type:String },
	language:{ type:String },
	viewcount:{ type:Number, default:0 },
	viewers:{ type:Number, default:0 },
	videourl:{ type:String },
	videoembedded:{ type:String },
	chatroomembedded:{ type:String },
	channel:{ type:String }
});

var webSchema = new Schema({
	host:{ type:String },
});

var userSchema = new Schema({
	name:{ type:String },
	picture:{ type:String },
	email:{ type:String },
	uid:{ type:String },
	lastlogin:{ type: Date, default: Date.now },
	search:[ {
		keyword:{ type:String },
		time:{ type: Date, default: Date.now }
	}],
	click_through:[{
		videourl:{type:String},
		time:{ type: Date, default: Date.now }
	}],
	statistics:{
		platform:{
			YouTube:{ type:Number, default:0 },
			Twitch:{ type:Number, default:0 },
			Facebook:{ type:Number, default:0 },
			Xigua:{ type:Number, default:0 },
			Live17:{ type:Number, default:0 }
		},
		category:{
			News:{ type:Number, default:0 },
			Sports:{ type:Number, default:0 },
			Games:{ type:Number, default:0 },
			Internet_Celebrities:{ type:Number, default:0 },
			Auctions:{ type:Number, default:0 },
			Webcam:{ type:Number, default:0 },
			Music:{ type:Number, default:0 },
			Nature:{ type:Number, default:0 },
			Cartoons:{ type:Number, default:0 }
		}
	},
	history:[ {
  		videourl: { type:String },
  		starttime: { type:String },
  		endtime: { type:String },
  		duration: { type:Number }
}]
});

var keywordSchema = new Schema({
	keyword:{ type:String },
	search_through:{ type:Number, default:1 }
});

var hostSchema = new Schema({
	platform:{ type:String },
	pagelink:{ type:String },
	account:{ type:String },
	chatroom:{ type:String },
	email:{ type:String }
});


exports.liveStreams = connCrawler.model('stream', streamSchema, 'Livestreams');
exports.webPlatform = connWeb.model('web', webSchema, 'Platform');
exports.User = connUser.model('user', userSchema, 'User');
exports.Keyword = connKeyword.model('keyword', keywordSchema, 'Keyword');
exports.Host = connHost.model('host', hostSchema, 'Host');
