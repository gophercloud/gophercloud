package websso

import (
	"net"
	"testing"
)

func TestIsLoopbackAddr(t *testing.T) {
	tests := []struct {
		name string
		addr net.Addr
		want bool
	}{
		{name: "IPv4 loopback", addr: &net.TCPAddr{IP: net.ParseIP("127.0.0.1")}, want: true},
		{name: "IPv6 loopback", addr: &net.TCPAddr{IP: net.ParseIP("::1")}, want: true},
		{name: "non-loopback", addr: &net.TCPAddr{IP: net.ParseIP("192.0.2.1")}},
		{name: "non-TCP", addr: testAddr("127.0.0.1:9990")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if actual := isLoopbackAddr(test.addr); actual != test.want {
				t.Fatalf("isLoopbackAddr() = %t, want %t", actual, test.want)
			}
		})
	}
}

type testAddr string

func (testAddr) Network() string  { return "test" }
func (a testAddr) String() string { return string(a) }
