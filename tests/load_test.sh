echo "Load testing hash password call"
for i in `seq 1 2000`; do curl --data "password=angryMonkey" http://localhost:8080/hash; done
echo "Load test complete"