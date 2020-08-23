#!/bin/bash

# for mac
sed -i "" /master1\.teeworlds\.com/d /etc/hosts
# for linux
sed -i /master1\.teeworlds\.com/d /etc/hosts
echo "移除master1.teeworlds.com在/etc/hosts的dns记录"