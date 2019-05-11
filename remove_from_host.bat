@echo off
findstr /v master1.teeworlds.com C:\Windows\System32\drivers\etc\hosts > tmphost.txt
move tmphost.txt C:\Windows\System32\drivers\etc\hosts
echo   hosts文件恢复完成
ipconfig /flushdns
pause