package dns

import (
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"
)

type checker struct {
	timeout time.Duration
}

func NewChecker(timeout time.Duration) *checker {
	return &checker{
		timeout: timeout,
	}
}

func (c *checker) CheckRecord(record Record) bool {
	// use a random replacement for wildcards to ensure we actually have a wildcard record
	dummy := fmt.Sprintf("kubermatic%d", 1000+rand.Intn(999))
	hostname := strings.Replace(record.Name, "*", dummy, -1)

	timeout := time.Now().Add(c.timeout)

	for time.Now().Before(timeout) {
		var success bool

		if record.Kind == RecordKindA {
			success = c.checkA(hostname, record.Target)
		} else {
			success = c.checkCNAME(hostname, record.Target)
		}

		if success {
			return true
		}

		time.Sleep(5 * time.Second)
	}

	return false
}

func (c *checker) checkA(hostname string, target string) bool {
	ips, _ := net.LookupIP(hostname)
	for _, ip := range ips {
		if ip.String() == target {
			return true
		}
	}

	return false
}

func (c *checker) checkCNAME(hostname string, target string) bool {
	cname, _ := net.LookupCNAME(hostname)

	return fmt.Sprintf("%s.", strings.TrimSuffix(target, ".")) == cname
}
