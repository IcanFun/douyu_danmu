package main

import (
	"regexp"
	"testing"
)

func TestReg(t *testing.T) {
	reg := regexp.MustCompile("nn@=([^n]*)/txt@")
	matchs := reg.FindAllStringSubmatch("20c8c9/level@=25/sahf@=0/cst@=1541143735020/bnn@=女流/bl@=9/brid@=156277/hc@=043410607b0ea218395c3/nn@=aa/txt@", -1)
	if matchs != nil && len(matchs) > 0 {
		t.Logf(matchs[len(matchs)-1][1])
	}
}
