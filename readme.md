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

将`config.example.yaml`重命名为`config.yaml`

```yaml
# 配置文件
# 系统配置
system:
  worker_number: 2 #单个账号的并发数
  delay: 2s #每个协程队列的延时
  max_queue_size: 100 #每个协程队列的最大值
  key: 1234567890qwerft #AES秘钥key 长度为16/24/32位

# 邮箱配置
emails:
- address:  example@example.com #邮箱地址
  server: smtp.exmail.example.com #smtp服务器地址
  port: 465 #端口号
  password: example #密码
- address:  example@example.com
  server: smtp.exmail.example.com
  port: 465
  password: example
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

## API

- 发送
```
uri: /
method: post
header: {
  token: "rdytfugihnljvcdhrjtugkyhk32fdr7tifuyvh"
}
params: {
  title: "Hello World!"
  address: "1@lailin.xyz"
  content: "<h1>Hello World</h1>"
}
```
|参数|说明|备注|
|:----:|:----:|:----:|
|token|签发的应用token|必须|
|title|邮件主题|必须|
|address|邮件地址|必须|
|content|邮件内容|必须|

- 返回
```json
{
    "status": 0,
    "msg": "地址 必须！",
    "data": ""
}
```
|参数|说明|备注|
|:----:|:----:|:----:|
|status|状态|0:失败,1:成功|
|msg|提示信息||
|data|数据信息|暂未使用|


## Author
[mohuishou](github.com/mohuishou)