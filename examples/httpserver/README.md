# HTTP Server Example
This is a simple example of HTTP server that uses the ratelimiter module as a middleware.

## Getting started
1. Update the `.env` content following YAML format. Set up the server config according to your liking.
2. Run the http server using the following command
```
make start
```
3. Test healthcheck endpoint
```
curl http://localhost:8080/ping
```
4. In the default `.env`, the rate limit is set at 5 rpm (requests per minute). You may test this by visiting `http://localhost:8080/test` on your browser. If you refresh the page 5 times, you'll get the following response at the end.
```
Response
HTTP/1.1 429 Too Many Requests
Ratelimit-Limit: 5
Ratelimit-Reset-After: 27.454376s
Ratelimit-Retry-After: 27.454376s
Date: Tue, 30 Mar 2021 00:00:00 GMT
Content-Length: 31
Content-Type: text/plain; charset=utf-8

request has exceeded rate limit
```