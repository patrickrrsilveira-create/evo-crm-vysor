package spintax

import (
	"math/rand"
	"strings"
)

func Process(text string) string {
	for {
		start := strings.LastIndex(text, "{")
		if start == -1 {
			break
		}
		end := strings.Index(text[start:], "}")
		if end == -1 {
			break
		}
		end += start

		options := strings.Split(text[start+1:end], "|")
		chosen := options[rand.Intn(len(options))]
		text = text[:start] + chosen + text[end+1:]
	}
	return text
}
