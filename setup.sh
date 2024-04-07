min=500
max=1000
delay=$((RANDOM%(max-min+1)+min))
delay=340
tc qdisc add dev eth0 root netem delay $((delay))ms
exec /b4