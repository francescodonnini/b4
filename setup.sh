min=500
max=1000
tc qdisc add dev eth0 root netem delay $((RANDOM%(max-min+1)+min))ms