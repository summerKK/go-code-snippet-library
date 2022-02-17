stopSh() {
  PID=$(ps -ef | grep $1 | grep -v grep | awk '{print $2}' | xargs kill -9)
  echo "$PID"
}

stopServer() {
  PID=$(lsof -i tcp:$1 | grep "LISTEN" | awk '{print $2}' | xargs kill -9)
  echo "$PID"
}

stopSh "rpc/user.go"
stopSh "api/user.go"
stopSh "rpc/product.go"
stopSh "api/product.go"
stopSh "rpc/order.go"
stopSh "api/order.go"
stopSh "api/pay.go"
stopSh "rpc/pay.go"

# 关闭服务
stopServer 8000
stopServer 8001
stopServer 8002
stopServer 8003
stopServer 9000
stopServer 9001
stopServer 9002
stopServer 9003
