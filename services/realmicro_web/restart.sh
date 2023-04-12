#!/bin/bash

echo "Usage: "
echo "  [realmicro-web script]"
echo ""
echo "Available Commands:"
echo "  restart"
echo "  stop"
echo ""
echo "Input Commands: "


RestartServer() {
    git pull
    go build -o realmicro-web
    supervisorctl restart realmicro-web
}

StopServer() {
    supervisorctl stop realmicro-web
}

while :
do
  read cmd
  case $cmd in
    restart) RestartServer
      echo "重启服务完成;"
      break
    ;;
    stop) StopServer
      echo "关闭服务完成;"
      break
    ;;
    *) echo "cmd not found Re-enter;"
      continue
    ;;
  esac
done