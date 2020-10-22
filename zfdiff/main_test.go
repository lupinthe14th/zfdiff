package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestZfDiff(t *testing.T) {
	t.Parallel()
	type in struct {
		a, b string
	}
	tests := []struct {
		in      in
		want    []string
		wantErr bool
	}{
		{in: in{a: "apple", b: "banana"}, want: []string{}, wantErr: true},
		{in: in{a: "./testdata/t1", b: "./testdata/t2"}, want: []string{"reply.sources.example.com.	60	IN	MX	1 sources.example.com."}, wantErr: false},
		{in: in{a: "./testdata/t2", b: "./testdata/t1"}, want: []string{"reply.sources.example.com.	60	IN	MX	1 sources.example.com."}, wantErr: false},
	}
	for i, tt := range tests {
		i, tt := i, tt
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()
			got, err := zfDiff(tt.in.a, tt.in.b)
			if tt.wantErr == true && err == nil {
				t.Fatalf("in: %v got: %#v err: %v want: %#v wantErr: %v", tt.in, got, err, tt.want, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("in: %v got: %#v err: %v want: %#v wantErr: %v", tt.in, got, err, tt.want, tt.wantErr)
			}
		})
	}
}

func TestRrList(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in      string
		want    []string
		wantErr bool
	}{
		{in: "banana", want: []string{}, wantErr: true},
		{in: "./testdata/t1", want: []string{
			"example.com.	60	AWS	ALIAS	A dns-name.elb.amazonaws.com. ABCDEFABCDE false",
			"sources.example.com.	60	IN	A	192.0.2.0",
			"reply.sources.example.com.	60	IN	MX	1 sources.example.com.",
		}, wantErr: false},
		{in: "./testdata/t2", want: []string{
			"example.com.	60	AWS	ALIAS	A dns-name.elb.amazonaws.com. ABCDEFABCDE false",
			"sources.example.com.	60	IN	A	192.0.2.0",
		}, wantErr: false},
	}
	for i, tt := range tests {
		i, tt := i, tt
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()
			got, err := rrList(tt.in)
			if tt.wantErr == true && err == nil {
				t.Fatalf("in: %v got: %#v err: %v want: %#v wantErr: %v", tt.in, got, err, tt.want, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("in: %v got: %#v err: %v want: %#v wantErr: %v", tt.in, got, err, tt.want, tt.wantErr)
			}
		})
	}
}

func TestParseRR(t *testing.T) {
	t.Parallel()
	type in struct {
		s, zone, ttl string
	}
	tests := []struct {
		in   in
		want string
	}{
		{in: in{s: "reply.sources   60      IN          MX    1 sources.example.com.", zone: "example.com.", ttl: "300"}, want: "reply.sources.example.com.	60	IN	MX	1 sources.example.com."},
		{in: in{s: "*        3600        IN        TXT        \"v=spf1 include:spf.example.com -all\"", zone: "example.com.", ttl: "300"}, want: "*.example.com.	3600	IN	TXT	\"v=spf1 include:spf.example.com -all\""},
		{in: in{s: "www 300 IN A 192.0.2.0", zone: "example.com.", ttl: "300"}, want: "www.example.com.	300	IN	A	192.0.2.0"},
		{in: in{s: "www IN A 192.0.2.0", zone: "example.com.", ttl: "300"}, want: "www.example.com.	300	IN	A	192.0.2.0"},
		{in: in{s: "www IN A 192.0.2.0", zone: "example.com.", ttl: "60"}, want: "www.example.com.	60	IN	A	192.0.2.0"},
		{in: in{s: "www 300 AWS ALIAS A dns-name.elb.amazonaws.com. ABCDEFABCDE false", zone: "example.com.", ttl: "60"}, want: "www.example.com.	300	AWS	ALIAS	A dns-name.elb.amazonaws.com. ABCDEFABCDE false"},
		{in: in{s: "www AWS ALIAS A dns-name.elb.amazonaws.com. ABCDEFABCDE false", zone: "example.com.", ttl: "60"}, want: "www.example.com.	60	AWS	ALIAS	A dns-name.elb.amazonaws.com. ABCDEFABCDE false"},
	}
	for i, tt := range tests {
		i, tt := i, tt
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()
			got := parseRR(tt.in.s, tt.in.zone, tt.in.ttl)
			if got != tt.want {
				t.Fatalf("in: %v got: %#v want: %#v ", tt.in, got, tt.want)
			}
		})
	}
}
