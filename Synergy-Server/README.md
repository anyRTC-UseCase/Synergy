## 一、部署环境

1. golang版本:1.16以上
2. mysql版本:5.5以上
3. redis版本:3.2.12以上
4. 智慧协同服务需部署在公网环境

### 二、数据库

`部署MySQL数据库`

```
[root@localhost ~]# yum -y install mariadb mariadb-server
[root@localhost ~]# systemctl enable mariadb
[root@localhost ~]# systemctl start mariadb
```

`导入sql脚本`

```
# 脚本位置 database目录下
[root@localhost database]# mysql -uroot
MariaDB [(none)]> source /usr/local/src/ARTeamViewService/database/ar_teamview.sql
```

`部署redis数据库`

```
[root@localhost ~]# yum -y install redis
[root@localhost ~]# systemctl enable redis
[root@localhost ~]# systemctl start redis
```

### 三、启动项目

`配置文件详情`

> [root@localhost ARTeamViewService]#  vim conf/config.go

`服务配置文件`

> [root@localhost ARTeamViewService]# vim conf/config.toml  (Linux环境下)

```
# 服务端口号（TCP）
port = 12681
# 进程pid目录(注意windows环境为空)
pidDir = "/run"

# 控制台右上角RESTful API密钥
customerId = ""
# 控制台右上角RESTful API密码
customerCert = ""
# 用户控制台创建项目 APP ID
appId = ""
# 用户控制台创建项目启用权限密钥
appToken = ""
# 开发者通知服务用于签名的secret(第一次启动项目时填空,在配置回调时联系客服获取)
devpSecret = ""

# http接口 token有效期（单位：秒）
tokenExpire = 172800
# rtc token有效期（单位：秒）
rtcTokenExpire = 2700
# rtm token有效期（单位：秒）
rtmTokenExpire = 172800

# 录像服务接口前缀(如果是私有云换自己的录像接口前缀)
httpVodPrefix = "https://api.agrtc.cn"
# 录像文件播放下载地址前缀
# 如果是第三方存储供应商,录像文件播放下载地址为"http://host:port/directory1/directory2/xxx.mp4",则此时前缀为:
# httpVodFilePrefix = "http://host:port/directory1/directory2/"
# 如果是私有oss,并且配置信息中fileNamePrefix配置了文件存储位置(参看下面 第三方云存储的配置信息),此时前缀为:
httpVodFilePrefix = "http://host:port/directory3/directory4/"

# 如果是私有oss，请在服务器使用crontab 设置定时任务删除录像文件，可忽略ossEndpoint参数
# 录像文件保存天数
vodValidDays = 30
# Bucket对应的Endpoint，以华东1（杭州）为例，填写为https://oss-cn-hangzhou.aliyuncs.com
ossEndpoint = ""

# 第三方云存储的配置信息
# 如果是阿里云oss，参见https://docs.anyrtc.io/cn/Recording/restful/cloud_recording_api_rest
# 如果是私有oss,参见https://github.com/anyRTC/Tools/tree/master/RecodingStoreServer
# 以下配置为私有oss
[dConfig.storageConfig]
# 文件传输接口URL
bucket= "http://host:port/arapi/v1/fdfs/file/uploadOssFile"
# vendor：100，100表示自定义私有上传文件系统，不上传云厂商OSS。
vendor= 100
# (选填）JSONArray 类型，由多个字符串组成的数组，指定录制文件在第三方云存储中的存储位置。 
fileNamePrefix= ["directory3", "directory4"]

# redis信息(根据自己的redis信息配置)
[dConfig.redis]
port = 6379
host = "127.0.0.1"
family = 4
password = ""
db = 0
# mysql信息 (根据自己的mysql信息配置)
[dConfig.mysql]
user = "root"
host = "127.0.0.1"
port = 3306
password = ""
database = "ar_teamview"
# 日志信息
[dConfig.logger]
# 日志保存天数
logExpire = 3
# 日志文件目录(windows环境为"./logs")
logDir = "/var/log/arteamviewservice"
```

`初始化swag接口文档`

> 不需要文档可跳过此步骤

1. 修改项目根目录main.go文件

   ```
   #47行
   // @host 192.168.1.140:12681
   192.168.1.140修改为你自己的服务器ip
   
   #251行
   URL: "http://192.168.1.140:12681/swagger/doc.json", //The url pointing to API definition
   192.168.1.140修改为你自己的服务器ip
   ```

2. 安装swag包

   ```
   go install github.com/swaggo/swag/cmd/swag@latest
   ```

3. 验证

   ```
   swag -version
   ```

4. 初始化

   ```
   swag init
   ```

5. 项目启动后访问下面url即可得到swag文档

   ```
   http://192.168.1.140:12681/swagger/index.html
   192.168.1.140修改为你自己的服务器ip
   ```

`编译启动项目`

`linux环境`

```
## 项目根目录下执行
[root@localhost ARTeamViewService]#  go build -o arteamviewservice -ldflags -w ./main.go

## 得到arteamviewservice可执行文件,运行即可
[root@localhost ARTeamViewService]# ./arteamviewservice    			# 前台启动
[root@localhost ARTeamViewService]# nohup ./arteamviewservice &     # 后台启动
```

`windows环境`

```
## 项目根目录下执行 (CMD中执行)
F:\ARTeamViewService> go build -o arteamviewservice.exe -ldflags -w ./main.go

## 得到arteamviewservice.exe,点击运行即可
```

### 四、配置回调

1. 事件回调地址(host换成自己的ip):

   ```
   http://host:12681/arapi/v1/teamview/teamViewRtcNotify  
   ```

2. 录像回调地址(host换成自己的ip):

   ```
   http://host:12681/arapi/v1/teamview/teamViewVodNotify
   ```

3. 联系客服配置回调地址,并获取开发者通知服务用于签名的secret

4. 修改配置文件devpSecret并重启项目

### 五、项目主要功能

- signIn:注册登录为一体的登录
- insertRoom:创建房间
- joinRoom:加入房间
- leaveRoom:离开房间
- getRoomList:获取房间列表
- getSpecialist:获取专家列表
- getUserInfo:获取用户信息
- teamViewRtcNotify:事件回调
- teamViewVodNotify:录像回调
- insertUserOnlineInfo:记录用户的在线信息

### 六、更新摘要

- 2.1更新
  1. 更新golang到1.18
  2. 添加接口swag文档
  3. 录像添加SubscribeAudioUids和SubscribeVideoUids字段,音频录取全部,视频只录取移动端
  4. 录像stop和31的回调fileList为空,组装录像文件url
  5. 房间和录像文件只保存固定天数
  6. 房间状态添加通话结束,四个状态的逻辑处理
  7. 数据库更新（见项目/database/ar_teamview.sql中的2.0数据库更新）
