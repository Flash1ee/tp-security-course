# Домашняя работа №1 по курсу "Безопасность интернет-приложений"
## HTTP Proxy Server

### Выполнил: Варин Дмитрий


## Подготовка
Приложение работает на порту `8080`, нужно проверить, свободен ли он.  
```bash
lsof -i :8080
```
Вывод, если порт занят
```text
flashie@ubuntu http-proxy % lsof -i :8080
COMMAND     PID    USER   FD   TYPE             DEVICE SIZE/OFF NODE NAME
___go_bui 29937 flashie    3u  IPv6 0x43cb0704bad15763      0t0  TCP *:http-alt (LISTEN)
```
Освободите порт:
```text
kill -9 29937
```
Если вывод пустой, то порт свободен.

## Запуск
1. Выполнить следующие команды из корневой директории.
```bash
make docker-build && make docker-run
```
2. Выполнить запрос с помощью `curl`.
```text
flashie@ubuntu http-proxy % curl -v -x http://127.0.0.1:8080 http://mail.ru
```
Ответ:
```bash
*   Trying 127.0.0.1:8080...
* Connected to 127.0.0.1 (127.0.0.1) port 8080 (#0)
> GET http://mail.ru/ HTTP/1.1
> Host: mail.ru
> User-Agent: curl/7.77.0
> Accept: */*
> Proxy-Connection: Keep-Alive
> 
* Mark bundle as not supporting multiuse
< HTTP/1.1 301 Moved Permanently
< Connection: keep-alive
< Content-Length: 185
< Content-Type: text/html
< Date: Sat, 12 Feb 2022 12:32:13 GMT
< Location: https://mail.ru/
< Server: nginx/1.14.1
< 
<html>
<head><title>301 Moved Permanently</title></head>
<body bgcolor="white">
<center><h1>301 Moved Permanently</h1></center>
<hr><center>nginx/1.14.1</center>
</body>
</html>
* Connection #0 to host 127.0.0.1 left intact
```
Если запросить `https`, сервер ответит ошибкой `404`.
```bash
flashie@ubuntu http-proxy % curl -v -x http://127.0.0.1:8080 https://mail.ru

*   Trying 127.0.0.1:8080...
* Connected to 127.0.0.1 (127.0.0.1) port 8080 (#0)
* allocate connect buffer!
* Establish HTTP proxy tunnel to mail.ru:443
> CONNECT mail.ru:443 HTTP/1.1
> Host: mail.ru:443
> User-Agent: curl/7.77.0
> Proxy-Connection: Keep-Alive
> 
< HTTP/1.1 400 Bad Request
< Date: Sat, 12 Feb 2022 12:42:37 GMT
< Content-Length: 0
< 
* Received HTTP code 400 from proxy after CONNECT
* CONNECT phase completed!
* Closing connection 0
curl: (56) Received HTTP code 400 from proxy after CONNECT

```
## Возможные ошибки
1. Контейнер уже запущен
```text
flashie@ubuntu http-proxy % make docker-run 
docker run -p 8080:8080 --name server -t http_proxy
docker: Error response from daemon: Conflict. The container name "/server" is already in use by container "adb49c571785e209070400ae7650be642b4b642501b6eca184e3a1a7ed5806b2". You have to remove (or rename) that container to be able to reuse that name.
See 'docker run --help'.
make: *** [docker-run] Error 125
```
Решение
```bash
make docker-rm && make docker-start
```
