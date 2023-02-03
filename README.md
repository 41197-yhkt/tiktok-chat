# Chat Example

## RUN
```shell
go run main.go client.go message.go hub.go
go run user_client1.go
go run user_client2.go
```

## Server

### hub

hub中添加chooseCLient,根据信息的uid选择用户

### Client

client在第一次登陆时查询httpHeader中的uid字段来在Hub中注册用户

### 问题

Postman无法进行socketIO测试, 显示Error: connect ECONNREFUSED.
