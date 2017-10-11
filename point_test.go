package martiangeo

import (
	"testing"
)

func TestTransfer(t *testing.T) {
	p := Point{128.543, 37.065}
	gcj := GCJ(p)
	t.Log(gcj.ToWGS())
	t.Log(gcj.ToBD())

	wgs := WGS(p)
	t.Log(wgs.ToGCJ())
	t.Log(wgs.ToBD())

	bd := BD(p)
	t.Log(bd.ToGCJ())
	t.Log(bd.ToWGS())

}
