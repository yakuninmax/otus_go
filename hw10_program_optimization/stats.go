package hw10programoptimization

import (
	"bufio"
	"errors"
	"io"
	"strings"

	"github.com/valyala/fastjson"
)

var ErrEmptyDomain = errors.New("empty domain")

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	// Check domain for empty string.
	if len(domain) == 0 {
		return nil, ErrEmptyDomain
	}

	// Create DomainStat map.
	stats := make(DomainStat)

	// Create scanner.
	scanner := bufio.NewScanner(r)

	// Scan input data.
	for scanner.Scan() {
		// Get E-Mail field.
		email := fastjson.GetString(scanner.Bytes(), "Email")

		// Check domain suffix.
		if strings.HasSuffix(email, "."+domain) {
			stats[strings.ToLower(strings.SplitN(email, "@", 2)[1])]++
		}
	}

	return stats, nil
}
