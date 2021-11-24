# nginx_exporter

#### 介绍
从ingress-nginx官方代码中的expoter迁移出来 用来监控虚拟机上的nginx的expoter

#### 软件架构
nginx 通过lua模块monitor.lua 将nginx log 以json格式发到 `/tmp/prometheus-nginx.socket`, nginx_exporter 通过这个socket获得数据并组装成metrics。 


#### 依赖
1） nginx 必须要编译有lua模块
2） lua 必须要有 cjson 模块


#### 安装luajit 和 cjson
```shell
yum install gcc -y

cd /usr/local/src/
wget --no-check-certificate https://luajit.org/download/LuaJIT-2.0.5.zip 
unzip LuaJIT-2.0.5.zip
cd LuaJIT-2.0.5/
make install PREFIX=/usr/local/luajit 

cd /usr/local/src/
wget --no-check-certificate https://kyne.com.au/~mark/software/download/lua-cjson-2.1.0.zip
unzip lua-cjson-2.1.0.zip 
cd lua-cjson-2.1.0/
# 这里要修改makefile文件，不然编译报错
sed -i 's#^LUA_INCLUDE_DIR = .*#LUA_INCLUDE_DIR =   /usr/local/src/LuaJIT-2.0.5/src#' Makefile
make && make install 
```



#### 使用说明

1.  xxxx
2.  xxxx
3.  xxxx

#### 参与贡献

1.  Fork 本仓库
2.  新建 Feat_xxx 分支
3.  提交代码
4.  新建 Pull Request


#### 特技

1.  使用 Readme\_XXX.md 来支持不同的语言，例如 Readme\_en.md, Readme\_zh.md
2.  Gitee 官方博客 [blog.gitee.com](https://blog.gitee.com)
3.  你可以 [https://gitee.com/explore](https://gitee.com/explore) 这个地址来了解 Gitee 上的优秀开源项目
4.  [GVP](https://gitee.com/gvp) 全称是 Gitee 最有价值开源项目，是综合评定出的优秀开源项目
5.  Gitee 官方提供的使用手册 [https://gitee.com/help](https://gitee.com/help)
6.  Gitee 封面人物是一档用来展示 Gitee 会员风采的栏目 [https://gitee.com/gitee-stars/](https://gitee.com/gitee-stars/)
