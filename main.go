package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/barnybug/cli53"
	"github.com/miekg/dns"
)

// These variables are set in build step
var (
	Version  = "unset"
	Revision = "unset"
)

func zfDiff(a, b string) ([]string, error) {
	out := make([]string, 0)

	as, err := rrList(a)
	if err != nil {
		return out, err
	}
	bs, err := rrList(b)
	if err != nil {
		return out, err
	}

	sort.Strings(as)
	sort.Strings(bs)
	memo := make(map[string]int)
	for i := range as {
		memo[as[i]]++
	}
	for i := range bs {
		memo[bs[i]]++
	}

	for k, v := range memo {
		if v == 1 {
			out = append(out, k)
		}
	}
	sort.Strings(out)
	return out, nil
}

// ZoneファイルからRRのリストを作成
func rrList(f string) ([]string, error) {
	out := []string{}
	fd, err := os.Open(f)
	if err != nil {
		return out, err
	}
	defer func() {
		fd.Close()
	}()

	z := dns.ParseZone(fd, "", "")
	for k := range z {
		if k.Error != nil {
			log.Printf("error: %v", k.Error)
		}
		r := parseComment(k.RR, k.Comment)
		out = append(out, r.String())
	}
	return out, nil
}

// See: https://github.com/barnybug/cli53/blob/1fe271a0d2b14217aaa0d4e1e546e8b59401b9ca/bind.go#L17-L39
func parseComment(rr dns.RR, comment string) dns.RR {
	if strings.HasPrefix(comment, "; AWS ") {
		kvs, err := cli53.ParseKeyValues(comment[6:])
		if err == nil {
			routing := kvs.GetString("routing")
			if fn, ok := cli53.RoutingTypes[routing]; ok {
				route := fn()
				route.Parse(kvs)
				rr = &cli53.AWSRR{
					rr,
					route,
					kvs.GetOptString("healthCheckId"),
					kvs.GetString("identifier"),
				}
			} else {
				fmt.Printf("Warning: parse AWS extension - routing=\"%s\" not understood\n", routing)
			}
		} else {
			fmt.Printf("Warning: parse AWS extension: %s", err)
		}
	}
	return rr
}

func main() {
	if len(os.Args) != 3 {
		log.Printf("argument count error: %v", len(os.Args))
		log.Printf("args: %v", os.Args)
		os.Exit(1)
	}
	out, err := zfDiff(os.Args[1], os.Args[2])
	if err != nil {
		log.Printf("error: %v", err)
		os.Exit(1)
	}
	for _, v := range out {
		fmt.Println(v)
	}
}
