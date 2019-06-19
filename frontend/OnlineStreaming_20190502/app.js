var createError = require('http-errors');
var express = require('express');
var path = require('path');
// var cookieParser = require('cookie-parser');
var session = require('express-session');
const MongoStore = require('connect-mongo')(session);

var logger = require('morgan');
var partials = require('express-partials');


var indexRouter = require('./routes/index');
var routerSession = require('./routes/seloginAPI');
var config = require('./config')

var app = express();

// view engine setup
app.set('views', path.join(__dirname, 'views'));
app.set('view engine', 'ejs');
app.use(partials())

app.use(logger('dev'));
app.use(express.json());
app.use(express.urlencoded({ extended: false }));
// app.use(cookieParser());
app.use(express.static(path.join(__dirname, 'public')));

app.use(session({
  secret: 'thisislivecatonlinestreaming',
  store:new MongoStore(config.sessionStorage),
  resave: false,
  saveUninitialized: false,
  cookie: { maxAge: 100 * 24 * 60 * 60 * 1000 } //100 day life
}));


app.use('/', indexRouter);
app.use('/session', routerSession);


// catch 404 and forward to error handler
app.use(function(req, res, next) {
  next(createError(404));
});

// error handler
app.use(function(err, req, res, next) {
  // set locals, only providing error in development
  res.locals.message = err.message;
  res.locals.error = req.app.get('env') === 'development' ? err : {};

  // render the error page
  res.status(err.status || 500);
  console.log(err.message);
  res.render('error-404');
});

module.exports = app;
