## 恩尼格玛密码机的 Go 语言实现 🔐
> 若以今日之智慧，破解昨日之密码，我辈也当为图灵。

这是一个用 Go 语言实现的恩尼格玛密码机，拥有三个转子和十组飞线。

## 这是什么？ 🤔
恩尼格玛（Enigma）是二战期间德国军方使用的加密设备。这台机器的密码系统最终被以*艾伦·图灵*为首的英国布莱切利园密码破译小组成功破解，成为了密码学历史上的重要里程碑。

想深入了解的可以看看：
- [维基百科：恩尼格玛密码机](https://zh.wikipedia.org/wiki/%E6%81%A9%E5%B0%BC%E6%A0%BC%E7%8E%9B%E5%AF%86%E7%A0%81%E6%9C%BA) 📚 - 这里的内容比我说的靠谱
- [Enigma的工作原理](https://www.youtube.com/watch?v=J46hu4RMB5I) ⚙️ - 看不懂英文的请自行脑补
- [Enigma的缺陷](https://www.youtube.com/watch?v=Sqpe5vZoKTo) 🔍 - 哦吼？666

如果你觉得看文档太枯燥，可以去看《模仿游戏》这部电影。 🎬

## 如何使用 🚀
首先，你得有个 Go 语言环境。

克隆下来，运行：
```bash
git clone https://github.com/kaolengmian7/enigma-machine.git
go run main.go
```

### 加密消息 🔒
就像德国人给总理发电报一样简单：
```curl
curl -X POST http://localhost:8080/api/encrypt \
-H "Content-Type: application/json" \
-d '{
    "message": "HELLO WORLD",
    "plugboard": ["AB", "CD", "EF", "GH", "IJ"],
    "positions": [0, 0, 0]
}'
```
返回：`{"result":"YFNDZ AAEZV"}`

### 解密消息 🔓
跟加密一样，只是接口不同：
```curl
curl -X POST http://localhost:8080/api/decrypt \
-H "Content-Type: application/json" \
-d '{
    "message": "YFNDZ AAEZV",
    "plugboard": ["AB", "CD", "EF", "GH", "IJ"],
    "positions": [0, 0, 0] 
}'
```
返回：`{"result":"HELLO WORLD"}`

## 一个小请求 ⭐
如果你觉得这个项目还不错，不妨给个 star。这年头，star 比德军的密码还难得到🐶。

## 破译
如果有大佬实现了破译代码，欢迎提交 mr 🎉
