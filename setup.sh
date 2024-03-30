min=500
max=1000
stdmin=10
stdmax=80
delay=$((RANDOM%(max-min+1)+min))
var=$((RANDOM%(stdmax-stdmin+1)+stdmin))
echo "delay=$((delay)); var=$((var))"
tc qdisc add dev eth0 root netem delay $((delay))ms $((var))ms