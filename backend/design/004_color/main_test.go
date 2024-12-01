package main

import "testing"

func TestConvColorCode2RGB(t *testing.T) {
	r, g, b := ConvColorCode2RGB("5B5A56")
	t.Logf("resp:%v,%v,%v", r, g, b)
}
