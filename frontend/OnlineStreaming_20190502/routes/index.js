var express = require('express');
var router = express.Router();
var mongoControl = require('../controllers/mongoController');
var userControl = require('../controllers/userController');
var hostControl = require('../controllers/hostController');
var admin = require('firebase-admin');
// var serviceAccount = require('../livecat-64d83-firebase-adminsdk-rdy5e-67f5048513');
// var firebaseAdmin = admin.initializeApp({
//     credential: admin.credential.cert(serviceAccount),
//     databaseURL:'https://livecat-64d83.firebaseio.com'
// })

// GET home page.
router.get('/', mongoControl.get_homepage);

// GET all stream page.
router.get('/all', mongoControl.get_all);

// GET search page.
router.get('/search', mongoControl.search);

// GET all platform page
router.get('/platform', mongoControl.get_all_platform);

// GET platform page
router.get('/platform/:platform', mongoControl.get_platform);

// GET video page
router.get('/livestream/:channel', mongoControl.get_livestream);

// GET privacy page.
router.get('/privacy', mongoControl.get_privacy);

// GET mongoDB info
router.get('/list_all', mongoControl.get_dbinfo);

// POST save history.
router.post('/savehistory', userControl.save_history);

// POST save user click action.
router.post('/saveclick', userControl.save_click);

// GET host form page
router.get('/host_form', mongoControl.get_hostform);

// POST save user click action.
router.post('/host_form_submit', hostControl.save_hostform);

module.exports = router;
