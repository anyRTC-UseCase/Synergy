# Synergy 智慧协同

## 代码目录

   - -------------Synergy-Web    // Web端代码              
       - --------------------Synergy-Experts     // 专家端
       - --------------------Synergy-Admin       // 管理员端
   - -------------Synergy-Android      // Android端代码
   - ------------ Synergy-Server       // 服务业务端代码

## 场景描述

智慧协同是 anyRTC 提供的协同场景化解决方案，结合 anyRTC RTC SDK 和 anyRTC 内容中心等产品，将其复杂的 API 进行模块整合，实现了功能组件化，降低了开发门槛。该方案，Android 端作为工人端进行邀请专家端进入协同，在协同中，可以邀请其他专家进行多方位指导。管理员端进入后台可以查看当前协同，协同结束后会把整个过程进行录制回放，方便后续查找问题。 

## 使用步骤

### 部署服务

先进行服务部署：Synergy-Server，查看如何[跑通服务]()

### 部署专家端和管理员端以及Android客户端

继服务部署好之后，在部署其他端，因为其他端部署或者跑通代码需要配置服务接口

## License
The MIT License (MIT).