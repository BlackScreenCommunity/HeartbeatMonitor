# HeartBeatMonitor

## Обновление версии

``` bash
wget ftp://185.246.118.224/HeartBeatMonitor-linux-amd64.zip
sudo unzip -o HeartBeatMonitor-linux-amd64.zip -d /srv/HBM -x appsettings.json
sudo chmod +x /srv/HBM/HeartBeatMonitor-linux-amd64
sudo systemctl restart hbm.service
```

## Пример файла *.service

``` bash
[Unit]
Description=HeartBeatMonitor

[Service]
WorkingDirectory=/srv/www/HBM
ExecStart=/srv/www/HBM/HeartBeatMonitor-linux-amd64
Restart=always
SyslogIdentifier=HeartBeatMonitor
PrivateTmp=true

[Install]
WantedBy=multi-user.target
```
