package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/armon/go-socks5"
	"golang.org/x/net/context"
)

// DNSResolver resolves any name to localhost
type DNSResolver struct {
	resolvingMap map[string]net.IP
}

// NewDNSResolver returns new resolver
func NewDNSResolver(resolvingMap string) *DNSResolver {
	r := &DNSResolver{}
	r.resolvingMap = make(map[string]net.IP)

	err := r.loadMapFromString(resolvingMap)
	if err != nil {
		fmt.Println("ERROR", err)
		os.Exit(1)
	}

	return r
}

// Resolve implements custom name resolution
func (r *DNSResolver) Resolve(ctx context.Context, name string) (context.Context, net.IP, error) {
	ip, err := r.resolve(name)
	log.Printf("Resolved \"%v\": %v", name, ip)
	return ctx, ip, err
}

func (r *DNSResolver) resolve(name string) (ip net.IP, err error) {
	for d, ip := range r.resolvingMap {
		if d == name || d == "*" {
			return ip, nil
		}
	}

	addr, err := net.ResolveIPAddr("ip", name)
	if err != nil {
		return nil, err
	}
	return addr.IP, err
}

// LoadMap adds domain/ip mapping from comma separated string
func (r *DNSResolver) loadMapFromString(config string) error {
	// Load Domain:IP mapping(s) from comma separated config string
	for _, domainIP := range strings.Split(config, ",") {
		domainIPTuple := strings.Split(domainIP, ":")
		switch len(domainIPTuple) {
		case 2:
			domain := strings.TrimSpace(domainIPTuple[0])
			if len(domain) == 0 {
				domain = WildcardDomain
			}

			ip := net.ParseIP(domainIPTuple[1])
			if ip == nil {
				return fmt.Errorf("invalid address: %s", domainIP)
			}

			if _, ok := r.resolvingMap[domain]; ok {
				if r.resolvingMap[domain].String() != ip.String() {
					return fmt.Errorf("already have domain \"%s\" resolving to \"%s\" cannot resolve also to \"%s\"", domain, r.resolvingMap[domain], ip)
				}
			}

			r.resolvingMap[domain] = ip
		default:
			return fmt.Errorf("invalid mapping: %s", domainIP)
		}
	}

	return nil
}

// WildcardDomain matches any domain name
const WildcardDomain = "*"

func main() {
	var resolvingMap, listenAddr string

	flag.StringVar(&listenAddr, "l", ":1080", "Listen address")
	flag.StringVar(&resolvingMap, "r", ":127.0.0.1", "Comma separated list of \"domain:IP\" for DNS resolving, \""+WildcardDomain+"\" or empty matches any name.\n\tDomains not in list/wildcard resolve through regular system DNS.")
	flag.Parse()

	// Create a SOCKS5 server with custom resolver
	resolver := NewDNSResolver(resolvingMap)
	conf := &socks5.Config{Resolver: resolver}
	server, err := socks5.New(conf)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Listening at %s, domain mappings: %s\n", listenAddr, resolver.resolvingMap)

	// Create SOCKS5 proxy
	if err := server.ListenAndServe("tcp", listenAddr); err != nil {
		panic(err)
	}
}
