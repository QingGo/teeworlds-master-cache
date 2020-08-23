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