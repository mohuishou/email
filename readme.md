邮件推送服务
======

> 提供邮件发送服务api，后台队列发送

## 安装
- 二进制包

- 源码
```
go get github.com/mohuishou/email
```

## 配置文件

将`config.example.ini`重命名为`config.ini`

```ini
;系统配置
[system]
;并发数
worker_number = 5
;延时时间
delay = 2s 
max_queue_size = 100

;邮箱配置
[email]
;邮箱地址
address = noreplay@example.com
;smtp服务器地址
server = smtp.example.com
;端口号
port = 465
;密码
password = password

;token秘钥，用于生成应用许可
;长度为16/24/32位
[token]
key = 1234567890123456
```

## 运行
- 生成应用秘钥
```bash
./email new -name=email -ip=127.0.0.1 -author=mohuishou
```
|参数|说明|
|:----:|:----:|
|name|应用名|
|ip|允许的来源地址|
|author|应用作者|

- 运行
```bash
./email -c=config.ini -p=8080
```
|参数|说明|
|:----:|:----:|
|c|配置文件地址|
|p|端口号|

## Author
[mohuishou](github.com/mohuishou)