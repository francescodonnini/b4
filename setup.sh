delay=600
tc qdisc add dev eth0 root netem delay $((delay))ms
exec /b4