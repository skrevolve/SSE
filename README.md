# SSE
Server-Sent Event (HTTP)

- 양방향 통신이 필요하지 않고 Amazon API Gateway, WebSocket 비용보다 좀더 저렴하게 사용하고 싶은 경우 사용
- 서버의 데이터를 지속적으로 Event Streaming
- server -> client 단방향 통신
- 재접속의 저수준 처리 자동지원

![image](https://github.com/skrevolve/SSE/assets/41939976/f824459a-cabd-4a28-a660-8b661c1b6819)
![image](https://github.com/skrevolve/SSE/assets/41939976/7a8710e1-cc34-4664-a5cf-0f61154c8ec0)


### client
```sh
npm install node eventsource
node main
```

### server
```sh
go run main.go
```