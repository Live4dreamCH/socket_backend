# 前后端网络接口

## 简介

网络协议使用HTTP和WebSocket。HTTP用于一般的“请求-响应式”的业务逻辑，且主要是控制；WebSocket用于“服务器推消息”、消息传输和文件传输。

协议内容主要使用JSON格式，除了文件传输部分。

IP地址114.116.234.101；
HTTP端口43851；
ws端口43852，消息路由为"/msg"，文件路由为"/file"；

## 功能列表

### 登录

登录分两步, http登录和ws登录
先进行http登录, 向服务器发送用户名和密码, 得到sid
然后向服务器的消息端口和文件端口建立ws连接, 发送sid以认证: 正确则保留连接, 错误/超时5s则断开连接

#### HTTP部分

POST /login

```json
"请求": {
    // 用户ID, 数字
    "id": 10086,
    // 密码, 字符串
    "psw": "password"
},
"回复": {
    // 请求字段缺失或数值类型错误，状态码400
    "请求格式错误": {
        "res": "NO",
        "reason": "json bind error"
    },
    "登陆失败": {
        "res": "NO",
        "reason": "wrong password"
    },
    "登录成功": {
        "res": "OK",
        // sid是sesionID, 一段base64, URL编码的随机字符串
        "sid": "random_string",
        // 用户昵称
        "name": "Tom"
    }
}
```

#### ws部分

ws是长连接，一个连接需要应用于多种操作。因此需要有一个“op”字段，表示操作类型。

ws连接是全双工的，相当于两条并行、反向的单向连接。在发送方连续发送报文时，发送方不知道接收方的“回复”，是回复的哪条报文，因此要从建立连接起，对“我发送的报文”进行计数，并填入“seq”字段中。
seq从0开始计数，一般第0个报文就是ws的登录报文，见下文。对方回复时，在“ack”字段中复述seq。收到此ack后就可以释放相关资源/展示回复结果。

（相比TCP，无需进行超时重传，因为基于TCP，它已经做了；只有连接断开时，才需要进行处理：重新连接，然后重新登录）

对方不一定需要回复，只有业务需要，才会回复。因而，如果不期望回复，可以不发送seq。

```json
// 目前设计的ws连接有2个, 一个文件, 一个消息. 两个连接 登录过程一致
"请求": {
    "op": "login",
    "seq": 0,
    "sid": "random_string"
},
"回复": {
    // 若失败/超时, 返回这条json后服务器就主动断开连接
    "登陆失败/超时": {
        "op": "login",
        "ack": 0,
        "res": "NO"
    },
    "登录成功": {
        "op": "login",
        "ack": 0,
        "res": "OK"
    }
}
```

### 加好友

A加B为好友, A先向服务器发送http请求, 服务器检查请求格式/检查是否已经为好友后http回复A

随后在B上线后, 服务器通过ws通知B(ws两个连接中的消息连接)

B在用户操作(同意/拒绝)后, 也向服务器发送http请求, 服务器检查请求格式/检查是否已经为好友后http回复B
若B同意, 则服务器在数据库中为AB两人建立一个会话, 使用上一句中的http向B返回会话号
若B拒绝, 则在数据库中删除A->B的好友关系

随后在A上线后, 服务器通过ws通知A加好友结果/会话号

#### HTTP部分

```json
"请求": {
    // POST /addfriend
    "A的好友请求": {
        "sid": "random_string",
        // B的uid
        "frid": 10086
    },
    // POST /resfriend
    "B接受好友的http请求": {
        "sid": "random_string",
        // A的uid
        "frid": 10086,
        "ans": "accept"
    },
    // POST /resfriend
    "B拒绝好友的http请求": {
        "sid": "random_string",
        // A的uid
        "frid": 10086,
        "ans": "refuse"
    }
},
"回复": {
    // 请求字段缺失或数值类型错误，状态码400
    "请求格式错误": {
        "res": "NO",
        "reason": "json bind error"
    },
    "sid错误": {
        "res": "NO",
        "reason": "wrong sid"
    },
    "已是好友/已发送单方向好友申请": {
        "res": "NO",
        "reason": "already friend"
    },
    // 数据库错误详情会保存在服务器端日志, 调试有问题找后端
    "数据库中添加/删除好友失败": {
        "res": "NO",
        "reason": "db err"
    },
    "服务器成功收到好友请求": {
        "res": "OK",
        // 若B接受好友申请, 则返回AB间的会话号
        "conv_id": 123
    }
}
```

#### ws部分

```json
"服务器把来自A的好友申请告诉B": {
    "op": "friend request",
    "seq": 1, // 不期望回复，非必需
    // A的uid
    "frid": 10086,
    // A的昵称
    "name": "Tom"
},
"服务器把B的答复告诉A": {
    "op": "friend answer",
    "seq": 1, // 不期望回复，非必需
    // B的uid
    "frid": 10086,
    // B的昵称
    "name": "Tom",
    // B的答复，可能为accept或refuse
    "ans": "accept",
    // 若"ans"="accept", 则返回AB间的会话号
    "conv_id": 123
}
```

### 收发消息

收发消息完全采用ws的消息连接。
为了兼容C-S-C和P2P两种传输方式，C->S的发消息与S->C的收消息（服务器推消息）格式一致。

#### 客户端收发消息

不管是C-S-C还是P2P，客户端收发的报文格式是一致的。

```json
{
    "op": "msg",
    "seq": 1,
    // 会话号
    "conv_id": 10086,
    // 发送消息的用户的uid
    "sender": 233,
    // 发送消息的时间
    "time": "2020-05-07 10:42:13",
    // 消息类型
    "type": "text",
    // 消息内容
    "content": "钟离nb！"
}
```

#### 服务器回复

服务器暂不需要回复，视后续功能需要再添加。

<!-- 若通过C-S-C方式，客户端向服务器用上述格式发送了一条消息，则服务器应当返回此消息的编号。此编号用于在客户端数据库中存放

```json
{
    "op": "msg response",
    "ack": 1,
    "msg id": 1
}
``` -->

### p2p通信

服务端需要根据json信息进行sdp转发，并进行异常处理。使用消息ws通道进行信息通信和信息接受。

#### 发送sdp

```json
{
    // 操作符为connect/connect response
    "op":"connect",
    // 服务器有回复,需要seq
    "seq":13,
    // 发送者的uid
    "from":10086,
    // 转发对象的uid
    "to":233,
    // sdp内容
    "sdp":"balabala"
}
```

#### 服务器回复

当对方不在线时, 服务器会回复:

```json
{
    "op": "conncet error", 
    "seq": 42, 
    "ack": 13, 
    "reason": "offline"
}
```

### 文件传输

文件传输单独使用一个ws连接，并且不允许并发传输两个文件（一个文件传完了再传另一个）。

由于ws分包机制的不确定性以及断点续传的要求，应用层需要把文件拆成一个一个的小段来传输。

形式上，文件传输总是由一个文字型的ws包开始，它内部用json格式说明了待传文件的信息，把它叫做协商包；其后紧随若干个二进制型的ws包，它们就是文件的内容，把它叫做内容包，每个包长度1kB，最后一个可以短一些。

#### 开始包

```json
{
    "op":"start",
    "conv_id":2,
    "name": "yuanshen.exe",
    "len": 2905792,
    "start":0,
    "end":2838
}
```

#### 暂停包

```json
{
    "op":"pause",
    "conv_id":2,
    "name": "yuanshen.exe",
}
```

#### 继续包

```json
{
    "op":"start",
    "conv_id":2,
    "name": "yuanshen.exe",
    "len": 2905792,
    "start":1000,
    "end":2838
}
```

### 用户昵称/好友列表/会话列表/会话成员查询

#### 用户昵称

请求: HTTP POST /name

```json
{
    "sid": "random_string",
    // 目标用户ID, 数字
    "id": 1
}
```

回复:

```json
{
    // 请求字段缺失或数值类型错误，状态码400
    "请求格式错误": {
        "res": "NO",
        "reason": "json bind error"
    },
    "sid错误": {
        "res": "NO",
        "reason": "wrong sid"
    },
    "无此用户": {
        "res": "NO",
        "reason": "wrong uid"
    },
    "返回昵称": {
        "res": "OK",
        "name": "mhq"
    }
}
```

#### 好友列表

请求: HTTP POST /friendlist

```json
{
    "sid": "random_string"
}
```

回复:

```json
{
    // 请求字段缺失或数值类型错误，状态码400
    "请求格式错误": {
        "res": "NO",
        "reason": "json bind error"
    },
    "sid错误": {
        "res": "NO",
        "reason": "wrong sid"
    },
    "返回好友列表": {
        "res": "OK",
        "friendlist": [
            {"id":1, "name":"mhq"},
            {"id":2, "name":"lch"}
        ]
    }
}
```

#### 会话查询

请求: HTTP POST /convlist

```json
{
    "sid": "random_string"
}
```

回复会话号和会话名, 若会话是两好友间的会话, 返回对方昵称作为会话名:

```json
{
    // 请求字段缺失或数值类型错误，状态码400
    "请求格式错误": {
        "res": "NO",
        "reason": "json bind error"
    },
    "sid错误": {
        "res": "NO",
        "reason": "wrong sid"
    },
    "返回会话列表": {
        "res": "OK",
        "convlist": [
            {"conv_id":1, "name":"软工小组"},
            {"conv_id":2, "name":"mhq"}
        ]
    }
}
```

#### 会话成员查询

请求: HTTP POST /convmemlist

```json
{
    "sid": "random_string",
    // 待查询的会话号
    "conv_id":2
}
```

回复成员uid和成员昵称:

```json
{
    // 请求字段缺失或数值类型错误，状态码400
    "请求格式错误": {
        "res": "NO",
        "reason": "json bind error"
    },
    "sid错误": {
        "res": "NO",
        "reason": "wrong sid"
    },
    "conv_id错误": {
        "res": "NO",
        "reason": "wrong conv_id"
    },
    "返回会话成员列表": {
        "res": "OK",
        "convmemlist": [
            {"id":1, "name":"mhq"},
            {"id":2, "name":"lch"}
        ]
    }
}
```
