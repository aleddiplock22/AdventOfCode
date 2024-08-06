package main

import "fmt"

/* ~coordinate system for hexagons~

Going to essentially treat this like cube coords

	-r
+s / \+q
  |   |
 -q\ / -s
	+r

So we track a vector (q,r,s) to represent the point in space much like a
(z,y,x) coord for a cube system.


I GOT THE COORDINATE SYSTEM WRONG FIRST TIME IGNORE DIAGRAM??? BUT THIS IS THE IDEA:
https://www.redblobgames.com/grids/hexagons/
*/

type HexCoord struct {
	r int // N/S
	s int // NW/SE axis
	q int // NE/SW axis
}

func (c *HexCoord) Add(c2 HexCoord) {
	c.r += c2.r
	c.s += c2.s
	c.q += c2.q
}

func (c *HexCoord) Sub(c2 HexCoord) {
	c.r -= c2.r
	c.s -= c2.s
	c.q -= c2.q
}

func (c HexCoord) Print() {
	fmt.Printf("(r=%v,s=%v,q=%v)\n", c.r, c.s, c.q)
}

var NORTH_EAST = HexCoord{
	r: 1,
	s: 0,
	q: -1,
}

var SOUTH_WEST = HexCoord{
	r: -1,
	s: 0,
	q: 1,
}

var NORTH_WEST = HexCoord{
	r: 0,
	s: 1,
	q: -1,
}

var SOUTH_EAST = HexCoord{
	r: 0,
	s: -1,
	q: 1,
}

var EAST = HexCoord{
	r: 1,
	s: -1,
	q: 0,
}

var WEST = HexCoord{
	r: -1,
	s: 1,
	q: 0,
}

func AddHexs(hex1 HexCoord, hex2 HexCoord) HexCoord {
	return HexCoord{
		hex1.r + hex2.r,
		hex1.s + hex2.s,
		hex1.q + hex2.q,
	}
}

func (c HexCoord) GetNeighbourCoords() []HexCoord {
	var neighbours []HexCoord
	for _, dir := range []HexCoord{NORTH_EAST, NORTH_WEST, SOUTH_EAST, SOUTH_WEST, EAST, WEST} {
		neighbours = append(neighbours, AddHexs(c, dir))
	}
	return neighbours
}
