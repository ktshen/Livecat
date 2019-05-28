const asyncRedis = require("async-redis");
const redis_port = 8869;
const redis_host = "140.115.153.185"
const redis_pwd = "livestream"

const client = asyncRedis.createClient({
    port: redis_port,                   // replace with your port
    host: redis_host,                   // replace with your hostanme or IP address
    password  : redis_pwd,              // replace with your password
    });

// check if connected to redis
client.on('connect', function() {
    console.log('Redis client connected');
});
client.on("error", function (err) {
    console.log("Error " + err);
});


exports.redis_select = async (channel) => {
  client.select(channel);
  console.log("[Redis SELECT]\n " + channel);
};

exports.redis_keys = async (key, cb) => {
  var data = await client.get(key);
  console.log("[Redis KEYS]\n " + data);
  cb();
  return data;
};

exports.redis_get = async (key, cb) => {
  var data = await client.get(key);
  console.log("[Redis GET]\n " + data);
  cb();
  return data;
};

exports.redis_lrange = async (key, start, end, cb) => {
  var data = await client.lrange(key, start, end);
  console.log("[Redis LRANGE]\n " + data);
  cb();
  return data;
};

exports.redis_llen = async (key, cb) => {
  var data = await client.llen(key);
  console.log("[Redis LLEN]\n " + data);
  cb();
  return data;
};

exports.redis_smembers = async (key, cb) => {
  var data = await client.smembers(key);
  console.log("[Redis SMEMBERS]\n " + data);
  cb();
  return data;
};

exports.redis_srandmember = async (key, num, cb) => {
  var data = await client.srandmember(key, num);
  console.log("[Redis SRANDMEMBER]\n " + data);
  cb();
  return data;
};
