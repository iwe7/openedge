# 自定义模块

- [文件目录](#文件目录)
  - [Docker容器模式](#docker容器模式)
  - [Native进程模式](#native进程模式)
- [模块配置](#模块配置)
- [函数计算runtime](#函数计算runtime)
- [模块退出](#模块退出)
- [模块SDK](#模块sdk)

在开发和集成自定义模块前请阅读开发编译指南，了解OpenEdge的编译环境要求。

自定义模块不限定开发语言，可运行即可。了解下面的这些约定，有利于更好更快的集成自定义模块。

## 文件目录

### Docker容器模式

容器中的工作目录推荐：/home/openedge/

容器中的配置路径是：/home/openedge/conf/conf.yml

持久化数据输出目录是：/home/openedge/var

具体的文件映射如下：

> - \<openedge_host_work_dir\>/app:/home/openedge/app
> - \<openedge_host_work_dir\>/app/\<module_mark\>:/home/openedge/conf
> - \<openedge_host_work_dir\>/var:/home/openedge/var

如果数据需要持久化在设备（宿主机）上，比如数据库和日志，必须保存在容器中的/home/openedge/var目录下，否者容器销毁数据会丢失。

### Native进程模式

模块的工作目录和OpenEdge主程序的工作目录相同。

配置路径是：\<openedge_host_work_dir\>/app/\<module_mark\>/conf.yml

## 模块配置

模块支持从文件中加载yaml格式的配置，Docker容器模式下读取conf/conf.yml，Native进程模式下读取conf/<module_mark>/conf.yml。

也支持从输入参数中获取配置，可以是json格式的字符串。比如：

    modules:
      - name: 'my_module'
        mark: 'my_module_conf_dir'
        entry: 'my_module_docker_image'
        params:
          - '-c'
          - '{"name":"my_module","address":"127.0.0.1:1234",...}'

也支持从环境变量中获取配置。比如：

    modules:
      - name: 'my_module'
        mark: 'my_module_conf_dir'
        entry: 'my_module_docker_image'
        env:
          name: my_module
          address: '127.0.0.1:1234'
          ...

## 函数计算runtime

函数计算的运行实例目前采用grpc server实现，借助grpc的多语言支持，开发者可通过[runtime.proto](https://github.com/baidu/openedge/blob/master/module/function/runtime/runtime.proto)生成自己的grpc server框架并实现消息处理的逻辑。

## 模块退出

模块退出时，OpenEdge会向模块发送SIGTERM信号，如果超过则强制退出。超时时间由app/app.yml中的grace配置项指定，默认30秒。

## 模块SDK

如果模块使用golang开发，可使用openedge中提供的部分SDK，位于github.com/baidu/openedge/module。

[mqtt dispatcher](https://github.com/baidu/openedge/blob/master/module/mqtt/dispatcher.go)可用于向mqtt hub订阅消息，支持自动重连。该mqtt dispatcher不支持消息持久化，消息持久化应由mqtt hub负责，mqtt dispatcher订阅的消息如果QoS为1，则消息处理完后需要主动回复ack，否者hub会重发消息。[remote mqtt模块参考](https://github.com/baidu/openedge/blob/master/openedge-remote-mqtt/main.go)

[function runtime server](https://github.com/baidu/openedge/blob/master/module/function/runtime/server.go)封装了grpc server和函数计算调用逻辑，方便开发者实现消息处理的函数计算runtime。[python2.7 runtime参考](https://github.com/baidu/openedge/blob/master/openedge-function-runtime-python27/openedge_function_runtime_python27.py)

[api.client](https://github.com/baidu/openedge/blob/master/api/client.go)可用于调用主程序的API来启动、重启或停止临时模块。账号为使用该api client的常驻模块名，密码使用```module.GetEnv(module.EnvOpenEdgeModuleToken)```从环境变量中获取。

[module.Load](https://github.com/baidu/openedge/blob/master/module/module.go)可用于加载模块的配置，支持从指定文件中读取yml格式的配置，也支持从启动参数中读取json格式的配置。

[module.Wait](https://github.com/baidu/openedge/blob/master/module/module.go)可用于等待模块退出的信号。

[logger](https://github.com/baidu/openedge/blob/master/module/logger/logger.go)可用于打印日志。
