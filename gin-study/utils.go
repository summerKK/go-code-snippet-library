package gin

import (
	"fmt"
	"net"
	"path"
	"reflect"
	"runtime"
	"strings"
)

func joinGroupPath(elems ...string) string {
	joined := path.Join(elems...)
	lastComponent := elems[len(elems)-1]
	if len(lastComponent) > 0 && lastComponent[len(lastComponent)-1] == '/' && joined[len(joined)-1] != '/' {
		joined += "/"
	}

	return joined
}

func filterFlags(content string) string {
	for i, c := range content {
		if c == ' ' || c == ';' {
			return content[:i]
		}
	}

	return content
}

func nameOfFunction(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

func chooseData(custom, wildcard interface{}) interface{} {
	if custom == nil {
		if wildcard == nil {
			panic("negotiation config is invalid")
		}
		return wildcard
	}
	return custom
}

func parseAccept(accept string) []string {
	parts := strings.Split(accept, ",")

	for i, part := range parts {
		index := strings.IndexByte(part, ';')
		if index >= 0 {
			part = accept[0:index]
		}
		part = strings.TrimSpace(part)
		parts[i] = part
	}

	return parts
}

func debugPrint(format string, values ...interface{}) {
	if IsDebugging() {
		fmt.Printf("[GIN-debug] "+format, values...)
	}
}

func lastChar(str string) uint8 {
	size := len(str)
	if size == 0 {
		panic("The length of the string can't be 0")
	}

	return str[size-1]
}

func ipInMasks(ip net.IP, masks []interface{}) bool {
	for _, proxy := range masks {
		var mask *net.IPNet
		var err error

		switch t := proxy.(type) {
		case string:
			if _, mask, err = net.ParseCIDR(t); err != nil {
				panic(err)
			}
		case net.IP:
			mask = &net.IPNet{IP: t, Mask: net.CIDRMask(len(t)*8, len(t)*8)}
		case net.IPNet:
			mask = &t
		}

		if mask.Contains(ip) {
			return true
		}
	}

	return false
}

// the ForwardedFor middleware unwraps the X-Forwarded-For headers, be careful to only use this
// middleware if you've got servers in front of this server. The list with (known) proxies and
// local ips are being filtered out of the forwarded for list, giving the last not local ip being
// the real client ip.
func ForwardedFor(proxies ...interface{}) HandlerFunc {
	if len(proxies) == 0 {
		// default to local ips
		var reservedLocalIps = []string{"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16"}

		proxies = make([]interface{}, len(reservedLocalIps))

		for i, v := range reservedLocalIps {
			proxies[i] = v
		}
	}

	return func(c *Context) {
		// the X-Forwarded-For header contains an array with left most the client ip, then
		// comma separated, all proxies the request passed. The last proxy appears
		// as the remote address of the request. Returning the client
		// ip to comply with default RemoteAddr response.

		// check if remoteaddr is local ip or in list of defined proxies
		remoteIp := net.ParseIP(strings.Split(c.Request.RemoteAddr, ":")[0])

		if !ipInMasks(remoteIp, proxies) {
			c.Next()
			return
		}

		if forwardedFor := c.Request.Header.Get("X-Forwarded-For"); forwardedFor != "" {
			parts := strings.Split(forwardedFor, ",")

			for i := len(parts) - 1; i >= 0; i-- {
				part := parts[i]

				ip := net.ParseIP(strings.TrimSpace(part))

				if ipInMasks(ip, proxies) {
					continue
				}

				// returning remote addr conform the original remote addr format
				c.Request.RemoteAddr = ip.String() + ":0"

				// remove forwarded for address
				c.Request.Header.Set("X-Forwarded-For", "")
				break
			}
		}

		c.Next()
	}
}
