# HeartBeatMonitor

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
