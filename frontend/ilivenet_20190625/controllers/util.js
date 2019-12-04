var async = require('async');
var crypto = require('crypto');

// Create array sequence(for page navigation)
exports.create_array = function(start, end) {
    var array = [];
    for (var i=start; i<=end; i++) {
        array.push(i);
    }
    return array;
}

exports.count_pages = function(dataNum, dataPerPage){
    if(dataNum%dataPerPage == 0){
        pageNum = dataNum/dataPerPage; // number of pagination
    }else{
        pageNum = dataNum/dataPerPage +1; // number of pagination
    }
    return pageNum;
}

exports.get_query_page = function(query_page){
    if (!query_page){ // Get page number chosen by user
        dataPage = 1;
    }else{
        dataPage = query_page;
    }
    return dataPage;
}

exports.get_query_page_es = function(query_page, dataPerPage){
    return (query_page-1)*dataPerPage;
}

function getRandomSalt(){
    return Math.random().toString().slice(2, 5);
}

exports.cryptToken = function(token) {
    // 密码“加盐”
    salt = getRandomSalt()
    var saltToken = token + ':' + salt;

    // 加盐密码的md5值
    var md5 = crypto.createHash('md5');
    var result = md5.update(saltToken).digest('hex');
    return result;
}
