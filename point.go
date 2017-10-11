package martiangeo

import (
	"math"
)

const (
	a   = 6378245.0              // 半长轴
	ee  = 0.00669342162296594323 // 扁率
	mpi = 3000 * math.Pi / 180   // baidu magic number
)

type Point struct {
	Longitute, Latitude float64
}

func (p Point) offset() Point {
	lon := p.Longitute - 105
	lat := p.Latitude - 35
	a := 40. / 3. * (math.Sin(6*math.Pi*lon) + math.Sin(2*math.Pi*lon))

	rlon := 300. + lon + 2*lat + 0.1*lon*lon + 0.1*lon*lat + 0.1*math.Sqrt(math.Abs(lon)) + a
	rlon += 40. / 3. * (math.Sin(math.Pi*lon) + 2*math.Sin(math.Pi/3*lon) + 7.5*math.Sin(math.Pi/12*lon) + 15*math.Sin(math.Pi/30*lon))

	rlat := -100. + 2*lon + 3*lat + 0.2*lat*lat + 0.1*lon*lat + 0.2*math.Sqrt(math.Abs(lon)) + a
	rlat += 40. / 3. * (math.Sin(math.Pi*lat) + 2*math.Sin(math.Pi/3*lat) + 8*math.Sin(math.Pi/12*lat) + 16*math.Sin(math.Pi/30*lat))
	return Point{rlon, rlat}
}

// WGS84: World Geodetic System / raw GPS coordinate
// https://en.wikipedia.org/wiki/World_Geodetic_System
type WGS Point

func (s *WGS) ToGCJ() GCJ {
	if (*Point)(s).outOfChina() {
		return GCJ{s.Longitute, s.Latitude}
	}

	os := (*Point)(s).offset()
	radlat := deg2rad(s.Latitude)
	magic := math.Sin(radlat)
	magic = 1 - ee*magic*magic
	sqrtMagic := math.Sqrt(magic)
	os.Longitute = rad2deg(os.Longitute) / a / math.Cos(radlat) * sqrtMagic
	os.Latitude = rad2deg(os.Latitude) / a / (1 - ee) * magic * sqrtMagic

	return GCJ{s.Longitute + os.Longitute, s.Latitude + os.Latitude}
}

func (s WGS) ToBD() BD {
	return s.ToGCJ().ToBD()
}

// GCJ02: 国家测绘局座标 / 火星座标系 / 高德
// https://zh.wikipedia.org/wiki/%E4%B8%AD%E5%8D%8E%E4%BA%BA%E6%B0%91%E5%85%B1%E5%92%8C%E5%9B%BD%E6%B5%8B%E7%BB%98%E9%99%90%E5%88%B6
type GCJ Point

func (g GCJ) ToWGS() WGS {
	if (Point)(g).outOfChina() {
		return WGS{g.Longitute, g.Latitude}
	}

	os := (Point)(g).offset()
	radlat := deg2rad(g.Latitude)
	magic := math.Sin(radlat)
	magic = 1 - ee*magic*magic
	sqrtMagic := math.Sqrt(magic)
	os.Longitute = rad2deg(os.Longitute) / a / math.Cos(radlat) * sqrtMagic
	os.Latitude = rad2deg(os.Latitude) / a / (1 - ee) * magic * sqrtMagic

	return WGS{g.Longitute - os.Longitute, g.Latitude - os.Latitude}
}

func (g GCJ) ToBD() BD {
	lon := g.Longitute
	lat := g.Latitude
	z := math.Sqrt(lon*lon+lat*lat) + 0.00002*math.Sin(lat*mpi)
	theta := math.Atan2(lat, lon) + 0.000003*math.Cos(lon*mpi)
	lon = z*math.Cos(theta) + 0.0065
	lat = z*math.Sin(theta) + 0.006
	return BD{lon, lat}
}

// BD09: 百度座标
// http://lbsyun.baidu.com/index.php?title=coordinate
type BD Point

func (b BD) ToGCJ() GCJ {
	lon := b.Longitute - 0.0065
	lat := b.Latitude - 0.006
	z := math.Sqrt(lon*lon+lat*lat) - 0.00002*math.Sin(lat*mpi)
	theta := math.Atan2(lat, lon) - 0.000003*math.Cos(lon*mpi)
	return GCJ{z * math.Cos(theta), z * math.Sin(theta)}
}

func (b BD) ToWGS() WGS {
	return b.ToGCJ().ToWGS()
}

func (p Point) outOfChina() bool {
	return !(p.Longitute > 73.66 && p.Longitute < 135.05 && p.Latitude > 3.86 && p.Latitude < 53.55)
}

func deg2rad(deg float64) float64 {
	return math.Pi / 180 * deg
}

func rad2deg(rad float64) float64 {
	return 180 / math.Pi * rad
}
