### Make your own Teeworlds master cache server
```
# DEBUG
go run .\cmd\teeworlds-master-cache.go
# Production Env
go build .\cmd\teeworlds-master-cache.go
# linux
GIN_MODE=release ./teeworlds-master-cache
```
### How to Use Your Teeworlds master cache server
* replace '127.0.0.1' in file 'add_to_host.bat' to your own cache server ip
* run 'add_to_host.bat' with **Administrator privileges**
* If your don't want to use custom master cache server any more, run 'remove_from_host.bat' with **Administrator privileges**

### 使用说明
* 把'add_to_host.bat'里的'127.0.0.1'改成缓存master服务器的ip
* 想使用缓存master服务器时，使用管理员权限运行'add_to_host.bat'
* 不想使用缓存master服务器时，使用管理员权限运行'remove_from_host.bat'