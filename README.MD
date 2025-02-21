# HeartBeatMonitor

HeartbeatMonitor — это простая система мониторинга серверов и приложений, которая собирает и отображает метрики с удалённых агентов.
Изначально разрабатывалось для мониторинга серверов web-приложения BPMSoft и Creatio.

Реализованы сбор следующих метрик:

- Занятое место на диске;
- Размер конкретного каталога;
- Сбор данных, полученных путем выполнения SQL-запросов в PostgreSQL;
- Серверное время.

## Обновление приложения

Актуальная версия приложения размещена на FTP-сервере с анонимным доступом.
Для обновления достаточно скачать актуальную версию приложения и разархивировать в каталог приложения.

``` bash
wget ftp://185.246.118.224/HeartBeatMonitor-linux-amd64.zip -O HeartBeatMonitor-linux-amd64.zip
sudo unzip -o HeartBeatMonitor-linux-amd64.zip -d /srv/HBM -x appsettings.json
sudo chmod +x /srv/HBM/HeartBeatMonitor-linux-amd64
sudo systemctl restart hbm.service
```

## Пример файла hbm.service

``` bash
[Unit]
Description=HeartBeatMonitor

[Service]
WorkingDirectory=/srv/HBM
ExecStart=/srv/HBM/HeartBeatMonitor-linux-amd64
Restart=always
SyslogIdentifier=HeartBeatMonitor
PrivateTmp=true

[Install]
WantedBy=multi-user.target
```
