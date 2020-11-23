package gnet

import (
	"fmt"
	"net"
	"testing"
)

func TestIsIPv4(t *testing.T) {
	type args struct {
		ip net.IP
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "127.0.0.1 is ipv4",
			args: args{ip: net.ParseIP("127.0.0.1")},
			want: true,
		},
		{
			name: "fe80::1 is not ipv4",
			args: args{ip: net.ParseIP("fe80::1")},
			want: false,
		},
		{
			name: "nil IP is not ipv4",
			args: args{ip: net.ParseIP("")},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsIPv4(tt.args.ip); got != tt.want {
				t.Errorf("IsIPv4() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsIPv6(t *testing.T) {
	type args struct {
		ip net.IP
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "127.0.0.1 is not ipv6",
			args: args{ip: net.ParseIP("127.0.0.1")},
			want: false,
		},
		{
			name: "fe80::1 is ipv6",
			args: args{ip: net.ParseIP("fe80::1")},
			want: true,
		},
		{
			name: "nil IP is not ipv6",
			args: args{ip: net.ParseIP("")},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsIPv6(tt.args.ip); got != tt.want {
				t.Errorf("IsIPv6() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ExampleIPV4OrV6(t *testing.T) {
	fmt.Println(IsIPv6(net.ParseIP("fe80::1")))
	fmt.Println(IsIPv4(net.ParseIP("fe80::1")))
	fmt.Println(IsIPv6(net.ParseIP("127.0.0.1")))
	fmt.Println(IsIPv4(net.ParseIP("127.0.0.1")))

	// Output:
	// true
	// false
}
