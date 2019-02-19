package dns

import (
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"
)

type validator struct {
	timeout time.Duration
}

func NewValidator(timeout time.Duration) *validator {
	return &validator{
		timeout: timeout,
	}
}

func (c *validator) ValidateRecord(record Record) bool {
	// use a random replacement for wildcards to ensure we actually have a wildcard record
	dummy := fmt.Sprintf("kubermatic%d", 1000+rand.Intn(999))
	hostname := strings.Replace(record.Name, "*", dummy, -1)

	timeout := time.Now().Add(c.timeout)

	for time.Now().Before(timeout) {
		var success bool

		if record.Kind == RecordKindA {
			success = c.validateA(hostname, record.Target)
		} else {
			success = c.validateCNAME(hostname, record.Target)
		}

		if success {
			return true
		}

		time.Sleep(5 * time.Second)
	}

	return false
}

func (c *validator) validateA(hostname string, target string) bool {
	ips, _ := net.LookupIP(hostname)
	for _, ip := range ips {
		if ip.String() == target {
			return true
		}
	}

	return false
}

func (c *validator) validateCNAME(hostname string, target string) bool {
	cname, _ := net.LookupCNAME(hostname)

	return fmt.Sprintf("%s.", strings.TrimSuffix(target, ".")) == cname
}
