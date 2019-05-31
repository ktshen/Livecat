# 啟動eserver
```
sudo gunicorn -b 120.126.16.88:17777 -k gevent -w 32 eserver:app
```

# 啟動 dbmanager
```
python3 dbmanager.py
```

# 啟動elastic docker
sudo docker-compose -f elastic-docker.yml up

