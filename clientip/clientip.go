// Package clientip 解析真实客户端 IP（可信代理感知的转发头解析）。
//
// 规则：
//   - 只有当直连对端（TCP RemoteAddr）落在可信代理网段内时，才采信转发头；
//     否则一律使用对端 IP，公网客户端无法用伪造头绕过（限流/日志防欺骗）。
//   - 逗号列表头（如 X-Forwarded-For）从右向左跳过可信代理网段，
//     返回第一个不可信 IP——与 Express trust proxy / Symfony 同款算法，
//     天然抵御客户端在列表左侧塞伪造条目。
//   - 头名称与顺序、可信代理网段均可配置；默认头序兼顾 nginx/Caddy/Cloudflare/Akamai。
//
// 部署提示：位于 CDN/云 LB（公网 IP 的边缘）之后时，必须把边缘网段
// （如 Cloudflare IP 列表）配进 TrustedProxies，否则解析出的是边缘 IP。
package clientip

import (
	"net"
	"net/http"
	"strings"
)

// DefaultHeaders 是默认的转发头查询顺序。
// X-Forwarded-For 排最前：合规代理对其做"追加"语义，配合右到左解析最抗伪造；
// 单值头（X-Real-IP/CF-Connecting-IP/True-Client-IP）依赖边缘覆写/清洗，仅作回退。
var DefaultHeaders = []string{"X-Forwarded-For", "X-Real-IP", "CF-Connecting-IP", "True-Client-IP"}

// defaultTrustedCIDRs 是未显式配置时默认可信的代理网段：
// 环回 + RFC1918 + 链路本地 + IPv6 ULA/链路本地。
// 覆盖"反代与应用同机/同 docker 网络"的常见部署；公网直连客户端不可能以这些地址出现，
// 因此伪造头不会被采信。多租户集群（k8s 等）请显式配置更小的网段。
var defaultTrustedCIDRs = []string{
	"127.0.0.0/8", "::1/128",
	"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16", "169.254.0.0/16",
	"fc00::/7", "fe80::/10",
}

// Options 是 Resolver 配置。
type Options struct {
	// Headers 是转发头查询顺序；空则用 DefaultHeaders。
	Headers []string
	// TrustedProxies 是可信代理网段（CIDR 或单 IP）；空则用默认私网集合。
	TrustedProxies []string
	// TrustNone 为 true 时完全不采信转发头（应用直接暴露公网、无任何代理）。
	TrustNone bool
}

// Resolver 按配置解析客户端 IP，并发安全。
type Resolver struct {
	headers   []string
	trusted   []*net.IPNet
	trustNone bool
}

// New 构造 Resolver；TrustedProxies 中的非法条目返回错误（空白条目忽略，
// 兼容"环境变量留空"产生的空串）。
func New(opts Options) (*Resolver, error) {
	r := &Resolver{trustNone: opts.TrustNone}
	for _, h := range opts.Headers {
		if h = strings.TrimSpace(h); h != "" {
			r.headers = append(r.headers, h)
		}
	}
	if len(r.headers) == 0 {
		r.headers = DefaultHeaders
	}
	cidrs := defaultTrustedCIDRs
	if hasNonEmpty(opts.TrustedProxies) {
		cidrs = opts.TrustedProxies
	}
	for _, c := range cidrs {
		if strings.TrimSpace(c) == "" {
			continue
		}
		n, err := parseCIDROrIP(strings.TrimSpace(c))
		if err != nil {
			return nil, err
		}
		r.trusted = append(r.trusted, n)
	}
	return r, nil
}

func hasNonEmpty(ss []string) bool {
	for _, s := range ss {
		if strings.TrimSpace(s) != "" {
			return true
		}
	}
	return false
}

func parseCIDROrIP(s string) (*net.IPNet, error) {
	if _, n, err := net.ParseCIDR(s); err == nil {
		return n, nil
	}
	ip := net.ParseIP(s)
	if ip == nil {
		return nil, &net.ParseError{Type: "CIDR or IP", Text: s}
	}
	bits := 32
	if ip.To4() == nil {
		bits = 128
	}
	return &net.IPNet{IP: ip, Mask: net.CIDRMask(bits, bits)}, nil
}

// ClientIP 返回可信的客户端 IP（不含端口）；无法解析时返回原样字符串。
func (r *Resolver) ClientIP(req *http.Request) string {
	peer := PeerIP(req.RemoteAddr)
	if r.trustNone {
		return peer
	}
	if ip := net.ParseIP(peer); ip == nil || !r.isTrusted(ip) {
		return peer
	}
	for _, h := range r.headers {
		v := req.Header.Get(h)
		if v == "" {
			continue
		}
		if ip := r.fromHeader(v); ip != "" {
			return ip
		}
	}
	return peer
}

// fromHeader 解析头的值：右到左跳过可信代理，返回第一个不可信 IP；
// 全部可信时返回最左（最初客户端声明位）。
func (r *Resolver) fromHeader(v string) string {
	parts := strings.Split(v, ",")
	for i := len(parts) - 1; i >= 0; i-- {
		ip := net.ParseIP(strings.TrimSpace(parts[i]))
		if ip == nil {
			continue
		}
		if !r.isTrusted(ip) {
			return ip.String()
		}
	}
	for _, p := range parts {
		if ip := net.ParseIP(strings.TrimSpace(p)); ip != nil {
			return ip.String()
		}
	}
	return ""
}

func (r *Resolver) isTrusted(ip net.IP) bool {
	for _, n := range r.trusted {
		if n.Contains(ip) {
			return true
		}
	}
	return false
}

// PeerIP 从 RemoteAddr 提取对端 IP（去端口）。
func PeerIP(remoteAddr string) string {
	if host, _, err := net.SplitHostPort(remoteAddr); err == nil {
		return host
	}
	return remoteAddr
}

// Middleware 解析客户端 IP 并写回 r.RemoteAddr（替代 chi middleware.RealIP，
// 后者无条件信任转发头，可被伪造）。下游的限流/日志读到的即为可信 IP。
func Middleware(r *Resolver) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			req.RemoteAddr = r.ClientIP(req)
			next.ServeHTTP(w, req)
		})
	}
}
