# HeartBeatMonitor

HeartbeatMonitor — это простая система мониторинга серверов и приложений, которая собирает и отображает метрики с удалённых агентов.
Изначально разрабатывалось для мониторинга серверов web-приложения BPMSoft и Creatio.

Реализованы сбор следующих метрик:

- Занятое место на диске;
- Размер конкретного каталога;
- Сбор данных, полученных путем выполнения SQL-запросов в PostgreSQL;
- Серверное время.
- Состояние Docker-контейнеров 

## Установка приложения

``` bash
sudo apt update
sudo apt install musl musl-dev musl-tools

wget ftp://185.246.118.224/HeartBeatMonitor-linux-amd64.zip -O HeartBeatMonitor-linux-amd64.zip
sudo unzip -o HeartBeatMonitor-linux-amd64.zip -d /srv/HBM \
sudo chmod +x /srv/HBM/HeartBeatMonitor-linux-amd64
```

## Обновление приложения

Актуальная версия приложения размещена на FTP-сервере с анонимным доступом.
Для обновления достаточно скачать актуальную версию приложения и разархивировать в каталог приложения.

``` bash
wget ftp://185.246.118.224/HeartBeatMonitor-linux-amd64.zip -O HeartBeatMonitor-linux-amd64.zip \
&& sudo unzip -o HeartBeatMonitor-linux-amd64.zip -d /srv/HBM -x appsettings.json \
&& sudo chmod +x /srv/HBM/HeartBeatMonitor-linux-amd64 \
&& sudo systemctl restart hbm.service
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


## Docker-контейнер [!PAUSED]

> [!CAUTION]
> Работы по запуску приложения в Докер не ведутся. 
> Работ в ближайшее время не ожидается. 
> Данные в текущем разделе представлены в информационных целях.

Докеризация практически не поддерживается.
Не реализована поддержка подмены `appsettings.json`

Работоспособность плагинов

| Плагин                                           | Статус                                      |
| ------------------------------------------------ | ------------------------------------------- |
| Плагин отображения версии запущенного приложения | Работает                                    |
| HardDriveFreeSpacePlugin                         | Не работает                                 |
| PostgreSqlQueryPlugin                            | Работает                                    |
| FolderSizePlugin                                 | Не работает                                 |
| ServerTimePlugin                                 | Не работает, но отображает время контейнера |
| ServiceStatusPlugin                              | Не работает                                 |
| DockerContainerStatusPlugin                      | Не работает                                 |

Единственное возможное использование в текущей версии: 
- агенты запускаются на удаленных серверах нативно
- веб-интерфейс может запускаться как Docker-контейнер

### Сборка Docker-image

``` shell
docker build \
    --build-arg COMMIT_HASH="App_in_docker" \
    --build-arg VERSION_DATE_PART=$(date +%Y%m%d.%H%M) \
    -t heartbeatmonitor:latest .
```

### Запуск Docker-container

``` shell
docker run -d \
    -p 80:80 \
    --name heartbeatmonitor \
    heartbeatmonitor:latest
```
