#!/bin/sh

tailscaled --tun=userspace-networking --socks5-server=localhost:1055 &
until tailscale up --authkey=${TPROXY_TAILSCALE_AUTH_KEY} --hostname=${TPROXY_TAILSCALE_HOST_NAME}
do
    sleep 0.1
done
echo Tailscale Started
ALL_PROXY=socks5://localhost:1055/ tproxy