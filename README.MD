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

## Запуск приложения

### Входные параметры

Актуальные параметры для запуска можно посмотреть при помощи команды:

``` shell
./HeartBeatMonitor-linux-amd64 -h
```

**configFilePath**. Строковые параметр. Определяет путь до файла с конфигурацией приложения.
    Если аргумент не передан, то приложение пытается открыть файл `appsettings.json`, в корневом каталоге приложения

### Примеры запуска

``` shell 
# Запуск без аргументов
./HeartBeatMonitor-linux-amd64

# Указываем путь до конфигурационного файла
./HeartBeatMonitor-linux-amd64 -configFilePath=/home/user/.config/hbm/appsettings.json
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



| Плагин                                           | Статус                                   | Комментарий                                                                                                                                   |
| ------------------------------------------------ | ---------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------- |
| Плагин отображения версии запущенного приложения | Работает                                 |                                                                                                                                               |
| Плагин отслеживания свободного места на диске    | Не работает                              | [WIP] Работает при монтировании корневого каталога на чтение                                                                                  |
| Плагин взаимодействия с PostgreSQL               | Работает                                 |                                                                                                                                               |
| Плагин отслеживания размера каталога             | Не работает                              | [WIP] Работает при монтировании корневого каталога на чтение                                                                                  |
| GПлагин отображения серверного времени            | Работает, но отображает время контейнера | [Research]                                                                                                                                    |
| Плагин отображения сервисов systemd              | Работает, с дополнительными условиями    | Работает при монтировании D-Bus-сокета  --security-opt apparmor=unconfined \, -v /run/dbus/system_bus_socket:/run/dbus/system_bus_socket:ro \ |
| Плагин отображения состояния Docker-контейнеров  | Работает, с дополнительными условиями    | Работает при монтировании Docker-сокет  -v /var/run/docker.sock:/var/run/docker.sock \                                                        |


Единственное возможное использование в текущей версии: 
- агенты запускаются на удаленных серверах нативно
- веб-интерфейс может запускаться как Docker-контейнер

### Сборка Docker-image

``` shell
docker build \
    --no-cache \
    --build-arg COMMIT_HASH=$(git rev-parse --short HEAD) \
    --build-arg VERSION_DATE_PART=$(date +%Y%m%d.%H%M) \
    -t heartbeatmonitor:latest .
```

### Запуск Docker-container

``` shell
docker run -d \
    -p 8089:8088 \
    --name heartbeatmonitor \
    --security-opt apparmor=unconfined \
    -v /run/dbus/system_bus_socket:/run/dbus/system_bus_socket:ro \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -v $(pwd)/appsettings.json:/config/appsettings.json \
    heartbeatmonitor:latest
```
