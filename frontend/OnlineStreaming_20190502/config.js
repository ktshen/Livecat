const mongo = "120.126.16.96:27017"

const config = {
     web: {
         title: 'iLiveNet'
     },

    db: {
        crawler: 'mongodb://'+mongo+'/Crawler',
        web: 'mongodb://'+mongo+'/Web',
        user: 'mongodb://'+mongo+'/User',
        keyword: 'mongodb://'+mongo+'/Keyword',
        host: 'mongodb://'+mongo+'/Host',
        session: 'mongodb://'+mongo+'/sessiondb',
        options:{
            server: {
                auto_reconnect: true,
                poolSize: 10
            },
            useNewUrlParser: true
        }
    },

    sessionStorage: {
        url: 'mongodb://'+mongo+'/sessiondb'
    },

    elasticsearch: {
        config:{
            host: '120.126.16.88:9200',
            log: 'trace'
        }
    },

    queryAPI: {
        server: '120.126.16.88:17777'
    }

};

module.exports = config;
