qiniutest manual
=========================

[![Build Status](https://travis-ci.org/qiniu/qiniutest.svg?branch=develop)](https://travis-ci.org/qiniu/qiniutest) [![GoDoc](https://godoc.org/qiniupkg.com/qiniutest?status.svg)](https://godoc.org/qiniupkg.com/qiniutest)

[![Qiniu Logo](http://open.qiniudn.com/logo.png)](http://www.qiniu.com/)

# 下载

```
go get -u qiniupkg.com/qiniutest
```

# 基础原理

* [httptest.v1/README.md](https://github.com/qiniu/httptest.v1)
 
# 命令行

```bash
qiniutest -v <QiniutestFile.qtf>
```

由 qiniutest 程序解释并执行 QTF 文件描述的测试脚本。指定 `-v` 参数表示 verbose，会获得更多的调试信息输出。

如果我们在 QTF 文件开头加上这样一行：

```bash
#!/usr/bin/env qiniutest
```

并把 `<QiniutestFile.qtf>` 文件设置为可执行：

```bash
chmod +x <QiniutestFile.qtf>
```

如此就可以直接运行它：

```bash
./<QiniutestFile.qtf>
```

# QTF 语言手册

## 命令

### match

匹配两个object。语法：

```bash
match <expected-object> <source-object>
```

关于 match 的详细解释，参见：

* [httptest.v1/README.md](https://github.com/qiniu/httptest.v1)

### host

样例1：使用环境变量来选择 stage 还是 product

```bash
match $(testenv) `env QiniuTestEnv`
match $(env) `envdecode QiniuTestEnv_$(testenv)`

host rs.qiniu.com $(env.RSHost)
```

这样，测试人员只需要配置环境：

```bash
export QiniuTestEnv_stage='{
	"RSHost": "192.168.1.10:9999",
	"AK": "...",
	"SK": "..."
}'

export QiniuTestEnv_product='{
	"RSHost": "rs.qiniu.com",
	"AK": "...",
	"SK": "..."
}'
```

然后：

```bash
QiniuTestEnv=stage qiniutest <QiniutestFile.qtf> # 测试stage环境
QiniuTestEnv=product qiniutest <QiniutestFile.qtf> # 测试product环境
```

### auth

定义 auth 信息：

```bash
auth <auth-name> <auth-interface>
```

在某次请求时引用 auth：

```bash
# 这里的 <auth-information> 可以是之前已经定义好的 <auth-name>，也可以直接是某个 <auth-interface>
auth <auth-information>
```

auth 信息通常用 AccessKey/SecretKey，或者 Username/Password，都是很敏感的信息，一般通过 env 传入，避免随着脚本入库。

样例1：

```bash
match $(testenv) `env QiniuTestEnv`
match $(env) `envdecode QiniuTestEnv_$(testenv)`

host rs.qiniu.com $(env.RSHost)
auth qboxtest `qbox $(env.AK) $(env.SK)`

post http://rs.qiniu.com/stat/`base64 testqiniu:ecug-2014-place.png`
auth qboxtest
ret 200
echo $(resp)
```

它等价于：

```bash
match $(testenv) `env QiniuTestEnv`
match $(env) `envdecode QiniuTestEnv_$(testenv)`

host rs.qiniu.com $(env.RSHost)

post http://rs.qiniu.com/stat/`base64 testqiniu:ecug-2014-place.png`
auth `qbox $(env.AK) $(env.SK)`
ret 200
echo $(resp)
```

### echo/println

echo/println功能相同，都用于调试、打印信息。语法：

```bash
echo <object1> <object2> ...
```

### req/post/get/delete/put

req 发起一个请求：

```bash
req <http-method> <url>
```

而 post/get/delete/put 是 `http-method` 分别为 POST/GET/DELETE/PUT 时的简写。如：

```bash
post <url>
```

### header

用于指定请求包或返回包的某个头部信息。语法：

```bash
header <key> <value1> <value2> ...
```

需要注意的是，在返回包匹配时，语句：

```bash
header Content-Type $(mime)
```

如果 $(resp.header.Content-Type) 是 ["application/json"]，那么得到的 $(mime) 并不是 ["application/json"]，而是 "application/json"。如果希望是 ["application/json"]，则应该这样写：

```bash
match $(mime) $(resp.header.Content-Type)
```

### body/json/form/text/binary

body 用于指定请求包或返回包的正文内容。语法：

```bash
body <content-type> <content-data>
```

其中 `<content-type>` 可以是 "application/json"、"application/text" 这样的全称，也可以简写为 "json"、"form"、"text"。

而 json/form 指令是 `<content-type>` 为 json/form 时的简写。如：

```bash
json <json-data>
```

等价于：

```bash
body json <json-data>
```

而 binary 指令是 `<content-type>` 为 "application/octet-stream" 的简写。如：

```bash
binary <binary-data>
```

等价于：

```bash
body application/octet-stream <binary-data>
```

### ret

ret 用来获得返回包。语法：

```bash
ret [<status-code>]
```

在指定了 `<status-code>` 时，会要求返回的 $(resp.code) 与该 status code 相符，否则测试失败。语句：

```bash
ret <status-code>
```

等价于：

```bash
ret
match <status-code> $(resp.code)
```

### clear

clear 用来清除已经绑定的变量。语法：

```bash
clear <var-name1> <var-name2> ...
```

### let

let 用于变量赋值，和主流命令式编程语言的 `=` 最为接近。例如：

```bash
let $(var-name) <expression>
```

等价于：

```bash
clear var-name
match $(var-name) <expression>
```

### equal/equalSet

equal 要求两个 object 的内容精确相等：

```bash
equal <object1> <object2>
```

equalSet 要求两个 object 都是 array，并且把两个 array 看作集合，要求两个集合精确相等：

```bash
equalSet <array-object1> <array-object2>
```

也就是两个 array 的元素排序后应该精确相同。

## 子命令

### base64

```bash
base64 -d -std <text>
```

对一段文本进行 base64 编码（encode）。如果指定了 `-d` 参数则为解码（decode）。如果指定了 `-std` 则使用 base64.StdEncoding（默认使用的是 UrlSafe 的 base64.URLEncoding）进行编解码。

### env

```bash
env <var-name>
```

取得环境变量对应的文本。

### decode

```bash
decode <text>
```

对一段 json 文本进行解码（decode）。

### envdecode

```bash
envdecode <var-name>
```

取得一个环境变量对应的文本，并且进行解码（json decode）。等价于：

```bash
match $(__auto_var1) `env <var-name>`
decode $(__auto_var1)
```

### qbox

```bash
qbox <AccessKey> <SecretKey>
```

返回由七牛云存储 qbox 规范定义的 `auth interface`，可被 `auth` 命令使用。如：

```bash
auth `qbox <AccessKey> <SecretKey>`
```

### authstub

```bash
authstub -uid <Uid> -utype <Utype>
```

返回由七牛内部定义的 mock authorization 授权的 `auth interface`，可被 `auth` 命令使用。如：

```bash
auth `authstub -uid 1 -utype 4`
```

这样就模拟了一个 uid 为 1 的标准用户。

