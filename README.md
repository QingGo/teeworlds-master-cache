### Make your own Teeworlds master cache server
```
# DEBUG
go run ./cmd/teeworlds-master-cache.go
# Build
go build ./cmd/teeworlds-master-cache.go
# run
nohup ./teeworlds-master-cache -PostToken <yourtoken> >> teeworlds-master-cache.log 2>&1 &
```
### How to Use Your Teeworlds master cache server
* replace '127.0.0.1' in file 'add_to_host.bat' to your own cache server ip
* run 'add_to_host.bat' with **Administrator privileges**
* If your don't want to use custom master cache server any more, run 'remove_from_host.bat' with **Administrator privileges**

### 使用说明
* 把'add_to_host.bat'里的'127.0.0.1'改成缓存master服务器的ip
* 想使用缓存master服务器时，使用管理员权限运行'add_to_host.bat'
* 不想使用缓存master服务器时，使用管理员权限运行'remove_from_host.bat'

为了方便增减服务器记录的操作，用React开发了前端页面并放在github page。
* 前端的项目地址是：https://github.com/QingGo/teeworlds-master-cache-frontend 
* page的页面地址是：https://qinggo.github.io/teeworlds-master-cache-frontend/

为了配置https使用了haproxy，相关配置如下：
``` bash
sudo mkdir /etc/haproxy/sslforfree/
sudo mkdir /etc/ssl/xip.io
sudo openssl genrsa -out /etc/ssl/xip.io/xip.io.key 1024
sudo openssl req -new -key /etc/ssl/xip.io/xip.io.key -out /etc/ssl/xip.io/xip.io.csr
sudo openssl x509 -req -days 365 -in /etc/ssl/xip.io/xip.io.csr -signkey /etc/ssl/xip.io/xip.io.key -out /etc/ssl/xip.io/xip.io.crt
sudo cat /etc/ssl/xip.io/xip.io.crt /etc/ssl/xip.io/xip.io.key | sudo tee /etc/haproxy/sslforfree/ssl.pem

# 手动启动haproxy命令：
sudo haproxy -f /etc/haproxy/haproxy.cfg
# 手动关闭haproxy命令：
sudo killall haproxy
# 重载haproxy配置：
sudo haproxy -f /etc/haproxy/haproxy.cfg -sf
# 重载haproxy配置（不中断服务）：
sudo haproxy -f /etc/haproxy/haproxy.cfg -p /var/run/haproxy.pid -sf $(cat /var/run/haproxy.pid)

```
但是自己生成的https证书还是会有警告，ip又无法申请https证书，最后决定把api应用部署到heroku。主要是需要把端口配置从环境变量获取，新增Procfile配置文件，以及修改一下生成token的逻辑保障安全。在网页端可以关联github库，每次推master分支时自动部署。
heroku cli是通过git分支heroku的信息和项目关联在一起的
``` bash
# 登录前先设置代理，和网页用的要一致
heroku login
git remote get-url heroku
git remote set-url heroku https://git.heroku.com/teeworld-master-cache.git
heroku ps
heroku logs --tail
heroku open
```

在heroku上部署有两个问题。一是免费账号30分钟无人访问会自动睡眠，内存里存的服务器列表数据丢失。二是没找到提供udp端口给客户端访问的方式，这是致命的一点。

因此打算把heroku上的应用作为反向代理，只利用其提供https的功能。但是heroku访问到国内的服务器会慢，服务器部署在海外的话，游戏客户端获取列表会和其它请求都会变慢。下一步考虑找找国内有没有类似的PaaS免费提供商。或者有没有别的给ip地址加一层https的方案。