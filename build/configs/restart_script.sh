#!/bin/sh
if (systemctl -q is-active minecraft.service)
then
    echo "Miinecraft server 已執行重啟!"
    systemctl restart minecraft.service
else
    echo "Miinecraft server 已開始執行!"
    systemctl start minecraft.service
fi