package moesifmiddleware

import (
	"net"
	"net/http"
	"strings"
)

func validIp(ipAddress string) bool {
	ip := net.ParseIP(ipAddress)
	return ip.To4() != nil || ip.To16() != nil
}

func getClientIpFromXForwardedFor(ipAddress string) string {

	// Split the string
	ips := strings.Split(ipAddress, ", ")

	// Sometimes IP addresses in this header can be 'unknown' (http://stackoverflow.com/a/11285650).
	// Therefore taking the left-most IP address that is not unknown
	// A Squid configuration directive can also set the value to "unknown" (http://www.squid-cache.org/Doc/config/forwarded_for/)
	for _, ip := range ips {

		// Azure Web App's also adds a port for some reason, so we'll only use the first part (the IP)
		if strings.Contains(ip, ":") {
			ip = strings.Split(ip, ":")[0]
		}

		// x-forwarded-for may return multiple IP addresses in the format:
		// "client IP, proxy 1 IP, proxy 2 IP"
		// Therefore, the right-most IP address is the IP address of the most recent proxy
		// and the left-most IP address is the IP address of the originating client.
		// source: http://docs.aws.amazon.com/elasticloadbalancing/latest/classic/x-forwarded-headers.html
		if validIp(ip) {
			return ip
		}

	}
	// Return empty String
	return ""
}

func getClientIp(request *http.Request) string {

	// Standard headers used by Amazon EC2, Heroku, and others.
	if xc, ok := request.Header["X-Client-Ip"]; ok {
		if validIp(xc[0]) {
			return xc[0]
		}
	}

	// Load-balancers (AWS ELB) or proxies.
	if lb, ok := request.Header["X-Forwarded-For"]; ok {
		xForwardedFor := getClientIpFromXForwardedFor(lb[0])
		if validIp(xForwardedFor) {
			return xForwardedFor
		}
	}

	// Cloudflare.
	// @see https://support.cloudflare.com/hc/en-us/articles/200170986-How-does-Cloudflare-handle-HTTP-Request-headers-
	// CF-Connecting-IP - applied to every request to the origin.
	if cf, ok := request.Header["Cf-Connecting-Ip"]; ok {
		if validIp(cf[0]) {
			return cf[0]
		}
	}

	// Akamai and Cloudflare: True-Client-IP.
	if tr, ok := request.Header["True-Client-Ip"]; ok {
		if validIp(tr[0]) {
			return tr[0]
		}
	}

	// Default nginx proxy/fcgi; alternative to x-forwarded-for, used by some proxies.
	if ngx, ok := request.Header["X-Real-Ip"]; ok {
		if validIp(ngx[0]) {
			return ngx[0]
		}
	}

	// (Rackspace LB and Riverbed's Stingray)
	// http://www.rackspace.com/knowledge_center/article/controlling-access-to-linux-cloud-sites-based-on-the-client-ip-address
	// https://splash.riverbed.com/docs/DOC-1926
	if rs, ok := request.Header["X-Cluster-Client-Ip"]; ok {
		if validIp(rs[0]) {
			return rs[0]
		}
	}

	if xf, ok := request.Header["X-Forwarded"]; ok {
		if validIp(xf[0]) {
			return xf[0]
		}
	}

	if ff, ok := request.Header["Forwarded-For"]; ok {
		if validIp(ff[0]) {
			return ff[0]
		}
	}

	if f, ok := request.Header["Forwarded"]; ok {
		if validIp(f[0]) {
			return f[0]
		}
	}

	// Default Address
	return request.RemoteAddr
}
