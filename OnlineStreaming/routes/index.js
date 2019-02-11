var express = require('express');
var router = express.Router();
var redisControl = require('../controllers/redisController')

// GET home page.
router.get('/', redisControl.get_home);

// GET host page.
router.get('/host', redisControl.get_host);

// GET category page.
router.get('/category', redisControl.get_category);

// GET views page.
router.get('/views', redisControl.get_home);

// GET platform page.
router.get('/platform', redisControl.get_platform);

// GET host filted page
router.get('/host/:host', redisControl.get_host_filted_list);

// GET platform filted page.
router.get('/platform/:platform', redisControl.get_platform_filted_list);

// GET category filted page.
router.get('/category/:category', redisControl.get_category_filted_list);



module.exports = router;
