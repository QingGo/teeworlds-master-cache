@echo off
findstr /v master1.teeworlds.com C:\Windows\System32\drivers\etc\hosts > tmphost.txt
move tmphost.txt C:\Windows\System32\drivers\etc\hosts
echo    hosts�ļ��ѻָ�
ipconfig /flushdns
pause