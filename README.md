# 如何运行

## 首先要使用go get安装gin、gorm等，可以查看官方文档、B站教程、CSDN、或者chatGPT

## 不要单独运行main.go，这样routes.go就无法运行，可以同时运行两个文件。也可以在终端使用go build,然后./loginTest运行项目。每次修改完都要go build一下

# 各个文件夹的作用

## common文件夹用于存放一些通用的功能模块或工具函数，如数据库等

## config文件夹用于存放配置文件

## controller文件夹用于存放应用程序的控制器代码。控制器是应用程序的核心部分之一，它负责接收客户端请求并作出响应

## middleware文件夹用于存放中间件代码，例如日志记录、认证、CORS等

## model文件夹用于存放应用程序的数据模型定义，通常使用gorm库来实现对象关系映射。

## util文件夹用于存放工具函数文件

## main.go是Gin应用程序的主入口文件

## routes.go是Gin应用程序的路由定义文件


# 如何创建一个需求处理函数，以注册功能为例

## 首先在controller里添加对于函数，函数名大写如Register(),函数具体写法参见'controller/userController.go'的代码注释

## 然后在routes中创建一个路由，并对接到controller对于的函数，如：r.POST("/api/auth/register", controller.Register)