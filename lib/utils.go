package lib

import (
	"fmt"
	"regexp"
	"strings"
)

func ReadCertificate(cert string) (string, error) {
	certRegExp, err := regexp.Compile(`(-+[\w ]+-+[a-zA-Z0-9\s/+=]*?-+[\w ]+-+)`)
	if err != nil {
		return "", err
	}
	certHeaderFooterRegExp, err := regexp.Compile(`-+[\w ]+-+`)
	if err != nil {
		return "", err
	}

	certClean := certRegExp.FindAllString(cert, -1)[0]
	if len(certClean) == 0 {
		return "", fmt.Errorf("Invalid PEM format")
	}
	certMiddle := certHeaderFooterRegExp.Split(certClean, -1)[1]
	certHeaderFooter := certHeaderFooterRegExp.FindAllString(certClean, -1)

	certStitched := make([]string, 3)
	certStitched[0] = certHeaderFooter[0]
	certStitched[1] = certMiddle
	certStitched[2] = certHeaderFooter[1]

	fmt.Println(certStitched)

	certPartsClean := map2(certStitched, func(s string) string {
		return strings.TrimSpace(s)
	})
	fmt.Println(certPartsClean)

	return strings.Join(certPartsClean, "\n"), nil
}

func map2(data []string, f func(string) string) []string {
	mapped := make([]string, len(data))

	for i, e := range data {
		mapped[i] = f(e)
	}

	return mapped
}
