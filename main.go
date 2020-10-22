package main

import (
	"bufio"
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
	out := make([]string, 0)
	fd, err := os.Open(f)
	if err != nil {
		return out, err
	}
	defer func() {
		fd.Close()
	}()

	l := 0
	var zone, ttl string
	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		switch l {
		case 0:
			zone = strings.Fields(scanner.Text())[1]
		case 1:
			ttl = strings.Fields(scanner.Text())[1]
		default:
			rr := parseRR(scanner.Text(), zone, ttl)
			out = append(out, rr)
		}
		l++
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

// SeeAlso: https://github.com/barnybug/cli53/blob/1fe271a0d2b14217aaa0d4e1e546e8b59401b9ca/commands.go#L574-L583
func parseRR(s, zone, ttl string) string {
	origin := fmt.Sprintf("$ORIGIN %s\n", zone)
	defaultTTL := fmt.Sprintf("$TTL %s\n", ttl)
	record, err := dns.NewRR(origin + defaultTTL + s)
	if awsrr, ok := record.(*cli53.AWSRR); ok {
		record = awsrr.RR
	}
	if err != nil {
		log.Printf("error: %v", err)
	}
	return record.String()
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
