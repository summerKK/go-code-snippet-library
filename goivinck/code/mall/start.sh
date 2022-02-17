sh -c "/usr/local/go-1.17/go1.17/bin/go run /Users/summer/Docker/goproject/goivinck/code/mall/service/user/rpc/user.go -f /Users/summer/Docker/goproject/goivinck/code/mall/service/user/rpc/etc/user.yaml &"

wait

sh -c "/usr/local/go-1.17/go1.17/bin/go run /Users/summer/Docker/goproject/goivinck/code/mall/service/user/api/user.go -f /Users/summer/Docker/goproject/goivinck/code/mall/service/user/api/etc/user.yaml &"

wait

sh -c "/usr/local/go-1.17/go1.17/bin/go run /Users/summer/Docker/goproject/goivinck/code/mall/service/product/rpc/product.go -f /Users/summer/Docker/goproject/goivinck/code/mall/service/product/rpc/etc/product.yaml &"

wait

sh -c "/usr/local/go-1.17/go1.17/bin/go run /Users/summer/Docker/goproject/goivinck/code/mall/service/product/api/product.go -f /Users/summer/Docker/goproject/goivinck/code/mall/service/product/api/etc/product.yaml &"

wait

sh -c "/usr/local/go-1.17/go1.17/bin/go run /Users/summer/Docker/goproject/goivinck/code/mall/service/order/rpc/order.go -f /Users/summer/Docker/goproject/goivinck/code/mall/service/order/rpc/etc/order.yaml &"

wait

sh -c "/usr/local/go-1.17/go1.17/bin/go run /Users/summer/Docker/goproject/goivinck/code/mall/service/order/api/order.go -f /Users/summer/Docker/goproject/goivinck/code/mall/service/order/api/etc/order.yaml &"

wait

sh -c "/usr/local/go-1.17/go1.17/bin/go run /Users/summer/Docker/goproject/goivinck/code/mall/service/pay/rpc/pay.go -f /Users/summer/Docker/goproject/goivinck/code/mall/service/pay/rpc/etc/pay.yaml &"

wait

sh -c "/usr/local/go-1.17/go1.17/bin/go run /Users/summer/Docker/goproject/goivinck/code/mall/service/pay/api/pay.go -f /Users/summer/Docker/goproject/goivinck/code/mall/service/pay/api/etc/pay.yaml &"
