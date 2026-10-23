package clientip

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func req(remoteAddr string, headers map[string]string) *http.Request {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.RemoteAddr = remoteAddr
	for k, v := range headers {
		r.Header.Set(k, v)
	}
	return r
}

func mustNew(t *testing.T, o Options) *Resolver {
	t.Helper()
	r, err := New(o)
	if err != nil {
		t.Fatal(err)
	}
	return r
}

func TestUntrustedPeerCannotSpoof(t *testing.T) {
	r := mustNew(t, Options{})
	// 公网直连客户端伪造转发头：一律不采信。
	got := r.ClientIP(req("203.0.113.9:40000", map[string]string{
		"X-Forwarded-For":  "8.8.8.8",
		"X-Real-IP":        "8.8.8.8",
		"CF-Connecting-IP": "8.8.8.8",
	}))
	if got != "203.0.113.9" {
		t.Fatalf("spoofed header honored: %q", got)
	}
}

func TestTrustedProxyXFFRightToLeft(t *testing.T) {
	r := mustNew(t, Options{})
	// Caddy（私网对端）追加真实客户端：XFF = "伪造victim, 真实client"。
	got := r.ClientIP(req("172.18.0.2:50000", map[string]string{
		"X-Forwarded-For": "6.6.6.6, 203.0.113.77",
	}))
	if got != "203.0.113.77" {
		t.Fatalf("want real client 203.0.113.77, got %q", got)
	}
}

func TestTrustedProxySingleValueHeader(t *testing.T) {
	r := mustNew(t, Options{})
	got := r.ClientIP(req("127.0.0.1:50000", map[string]string{"X-Real-IP": "203.0.113.5"}))
	if got != "203.0.113.5" {
		t.Fatalf("X-Real-IP: got %q", got)
	}
	// XFF 存在时优先于 X-Real-IP（默认头序）。
	got = r.ClientIP(req("127.0.0.1:50000", map[string]string{
		"X-Forwarded-For": "203.0.113.9",
		"X-Real-IP":       "203.0.113.5",
	}))
	if got != "203.0.113.9" {
		t.Fatalf("header order: got %q", got)
	}
}

func TestCDNEdgeInTrustedProxies(t *testing.T) {
	// Cloudflare 场景：把边缘网段配进可信代理后，XFF 右到左跳过边缘得到真实客户端。
	r := mustNew(t, Options{TrustedProxies: []string{"127.0.0.0/8", "173.245.48.0/20"}})
	got := r.ClientIP(req("127.0.0.1:50000", map[string]string{
		"X-Forwarded-For": "203.0.113.44, 173.245.48.7",
	}))
	if got != "203.0.113.44" {
		t.Fatalf("cdn chain: got %q", got)
	}
	// 未配置该网段时，解析到边缘 IP（安全但不精确）——提示运维必须配置。
	r2 := mustNew(t, Options{TrustedProxies: []string{"127.0.0.0/8"}})
	got = r2.ClientIP(req("127.0.0.1:50000", map[string]string{
		"X-Forwarded-For": "203.0.113.44, 173.245.48.7",
	}))
	if got != "173.245.48.7" {
		t.Fatalf("unconfigured cdn: got %q", got)
	}
}

func TestAllTrustedFallsBackToLeftmost(t *testing.T) {
	r := mustNew(t, Options{})
	got := r.ClientIP(req("192.168.1.1:50000", map[string]string{"X-Forwarded-For": "10.1.2.3, 192.168.1.9"}))
	if got != "10.1.2.3" {
		t.Fatalf("internal client: got %q", got)
	}
}

func TestTrustNone(t *testing.T) {
	r := mustNew(t, Options{TrustNone: true})
	got := r.ClientIP(req("172.18.0.2:50000", map[string]string{"X-Forwarded-For": "203.0.113.9"}))
	if got != "172.18.0.2" {
		t.Fatalf("trust-none: got %q", got)
	}
}

func TestCustomHeaders(t *testing.T) {
	r := mustNew(t, Options{Headers: []string{"True-Client-IP"}})
	got := r.ClientIP(req("10.0.0.1:1", map[string]string{
		"X-Forwarded-For": "203.0.113.9",
		"True-Client-IP":  "203.0.113.10",
	}))
	if got != "203.0.113.10" {
		t.Fatalf("custom headers: got %q", got)
	}
}

func TestGarbageHeadersFallBackToPeer(t *testing.T) {
	r := mustNew(t, Options{})
	got := r.ClientIP(req("10.0.0.1:1", map[string]string{"X-Forwarded-For": "not-an-ip, , 也不是ip"}))
	if got != "10.0.0.1" {
		t.Fatalf("garbage: got %q", got)
	}
}

func TestPeerIP(t *testing.T) {
	if got := PeerIP("1.2.3.4:5678"); got != "1.2.3.4" {
		t.Fatal(got)
	}
	if got := PeerIP("[::1]:80"); got != "::1" {
		t.Fatal(got)
	}
	if got := PeerIP("1.2.3.4"); got != "1.2.3.4" {
		t.Fatal(got)
	}
}

func TestInvalidTrustedProxies(t *testing.T) {
	if _, err := New(Options{TrustedProxies: []string{"10.0.0.0/33"}}); err == nil {
		t.Fatal("expected error for invalid CIDR")
	}
}

func TestEmptyEntriesIgnored(t *testing.T) {
	// 环境变量留空时 caarlos0/env 可能给出 [""]，不应报错或改变默认行为。
	r := mustNew(t, Options{TrustedProxies: []string{""}, Headers: []string{"", " "}})
	if len(r.headers) != len(DefaultHeaders) {
		t.Fatalf("headers should fall back to defaults, got %v", r.headers)
	}
	got := r.ClientIP(req("172.18.0.2:1", map[string]string{"X-Forwarded-For": "203.0.113.9"}))
	if got != "203.0.113.9" {
		t.Fatalf("default trusted set lost: got %q", got)
	}
}

func TestMiddlewareSetsRemoteAddr(t *testing.T) {
	r := mustNew(t, Options{})
	var seen string
	h := Middleware(r)(http.HandlerFunc(func(_ http.ResponseWriter, req *http.Request) {
		seen = req.RemoteAddr
	}))
	h.ServeHTTP(httptest.NewRecorder(), req("172.18.0.2:50000", map[string]string{"X-Forwarded-For": "203.0.113.77"}))
	if seen != "203.0.113.77" {
		t.Fatalf("middleware RemoteAddr = %q", seen)
	}
}
