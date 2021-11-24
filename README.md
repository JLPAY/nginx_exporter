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

1.  将Lua脚本copy到 /data/nginx/lua 目录（这个目录可以自己定义，和nginx配置文件一致就行）;

2.  修改nginx的http模块配置，新增如下配置
```nginx
http {

    # lua脚本的目录路径
    lua_package_path "/data/nginx/lua/?.lua;;";

    init_by_lua_block {
        collectgarbage("collect")

        -- init modules
        local ok, res

        ok, res = pcall(require, "monitor")
        if not ok then
                error("require failed: " .. tostring(res))
        else
                monitor = res
        end

        ok, res = pcall(require, "plugins")
        if not ok then
                error("require failed: " .. tostring(res))
        else
                plugins = res
        end
        -- load all plugins that'll be used here
        plugins.init({  })
    }

    init_worker_by_lua_block {
        monitor.init_worker(10000)
        plugins.run()
    }

    log_by_lua_block {
        monitor.call()
        plugins.run()
    }
    
......
}

```

3.  启动nginx_exproter
```shell
# 编译ngx_exporter
git clone https://gitee.com/xianglinzeng/nginx_exporter.git
cd nginx_exporte
go mod tidy
go build -o nginx_exporter

# nginx_exporter 参数
# -port 指定启动端口
# -v    指定日志级别  1 2 3 4 5 越高日志越详细，默认是2，不指定也行，调试使用5
./nginx_exporter -port=9999 -v=5
```




