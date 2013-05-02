package main

import (
	"bytes"
	"encoding/binary"
	"io"
)

type Model uint32

const (
	Paddle Model = 1
	Ball   Model = 2
)

const (
	Up   Action = 0
	Down Action = 1
)

type Entity struct {
	pos, vel, size Vec
	model          Model
}

var ents = make([]Entity, 3)
var entsOld = make([]Entity, 3)

func init() {
	ents[0].model = Paddle
	ents[0].pos = [...]float64{0, -100, 0}
	ents[0].size = [...]float64{20, 5, 10}

	ents[1].model = Paddle
	ents[1].pos = [...]float64{0, 100, 0}
	ents[1].size = [...]float64{20, 5, 10}

	ents[2].model = Ball
}

func updateSimulation() {
	for i, ent := range ents {
		entsOld[i] = ent
	}
}

func serialize(buf io.Writer, serAll bool) {
	for i, ent := range ents {
		bitMask := make([]byte, 1)
		bufTemp := &bytes.Buffer{}
		if serAll || !ent.pos.Equals(&entsOld[i].pos) {
			bitMask[0] |= 1 << uint(i)
			binary.Write(buf, binary.LittleEndian, ent.pos)
		}
		buf.Write(bitMask)
		buf.Write(bufTemp.Bytes())
	}
	for i, ent := range ents {
		bitMask := make([]byte, 1)
		bufTemp := &bytes.Buffer{}
		if serAll || !ent.vel.Equals(&entsOld[i].vel) {
			bitMask[0] |= 1 << uint(i)
			binary.Write(buf, binary.LittleEndian, ent.vel)
		}
		buf.Write(bitMask)
		buf.Write(bufTemp.Bytes())
	}
	for i, ent := range ents {
		bitMask := make([]byte, 1)
		bufTemp := &bytes.Buffer{}
		if serAll || !ent.size.Equals(&entsOld[i].size) {
			bitMask[0] |= 1 << uint(i)
			binary.Write(buf, binary.LittleEndian, ent.size)
		}
		buf.Write(bitMask)
		buf.Write(bufTemp.Bytes())
	}
	for i, ent := range ents {
		bitMask := make([]byte, 1)
		bufTemp := &bytes.Buffer{}
		if serAll || ent.model != entsOld[i].model {
			bitMask[0] |= 1 << uint(i)
			binary.Write(buf, binary.LittleEndian, ent.model)
		}
		buf.Write(bitMask)
		buf.Write(bufTemp.Bytes())
	}
}
