># <font size=4> 一、如何运行

>## <font size=3> 1.先要使用<font color="red">go get</font>安装<font color="red">gin、gorm</font>等，可以查看官方文档、B站教程、CSDN、或者chatGPT <br> 2.<font color="red">不要单独运行main.go</font>，这样routes.go就无法运行，可以同时运行两个文件。也可以在终端使用<font color="red">go build</font>,然后<font color="red">./loginTest</font>运行项目。<font color="red">每次修改完都要go build一下</font> <br> 3.<font color="red">运行时记得在`config/application.yml`改数据库配置</font>

___

># <font size=4> 二、各个文件夹的作用

>## <font size=3> 1.`common`文件夹用于存放一些通用的功能模块或工具函数，如数据库 <br> 2.`config`文件夹用于存放配置文件 <br> 3.`controller`文件夹用于存放应用程序的控制器代码。控制器是应用程序的核心部分之一，它负责接收客户端请求并作出响应 <br> 4.`middleware`文件夹用于存放中间件代码，例如日志记录、认证、CORS等 <br> 5.`model`文件夹用于存放应用程序的数据模型定义，通常使用gorm库来实现对象关系映射 <br> 6.`util`文件夹用于存放工具函数文件 <br> 7.`main.go`是Gin应用程序的主入口文件 <br> 8.`routes.go`是Gin应用程序的路由定义文件

___

># <font size=4> 三、如何创建一个需求处理函数，以注册功能为例

>## <font size=3> 1.首先在`controller`里添加对于函数，函数名大写如Register(),函数具体写法参见`controller/userController.go`的代码注释<br> 2.然后在`routes`中创建一个路由，并对接到`controller`对应的函数，如：
``` go
r.POST("/api/auth/register", controller.Register)
```