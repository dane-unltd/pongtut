package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
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
var players = make([]PlayerId, 2)

func init() {
	ents[0].model = Paddle
	ents[0].pos = Vec{-75, 0, 0}
	ents[0].size = Vec{5, 20, 10}

	ents[1].model = Paddle
	ents[1].pos = Vec{75, 0, 0}
	ents[1].size = Vec{5, 20, 10}

	ents[2].model = Ball
	ents[2].size = Vec{20, 20, 20}
}

func copyState() {
	for i, ent := range ents {
		entsOld[i] = ent
	}
}

func updateSimulation() {
	move()
}

func move() {
	for i := range ents {
		ents[i].pos.Add(&ents[i].pos, &ents[i].vel)
	}
}

func login(id PlayerId) {
	if players[0] == 0 {
		players[0] = id
		if players[1] != 0 {
			startGame()
		}
		return
	}
	if players[1] == 0 {
		players[1] = id
		startGame()
	}
}

func startGame() {
	fmt.Println("game started")
	ents[0].pos = Vec{-75, 0, 0}
	ents[1].pos = Vec{75, 0, 0}

	ents[2].pos = Vec{0, 0, 0}
	ents[2].vel = Vec{10, 0, 0}
}

func disconnect(id PlayerId) {
	if players[0] == id {
		players[0] = 0
		stopGame()
	} else if players[1] == id {
		players[1] = 0
		stopGame()
	}
}

func stopGame() {
	ents[2].pos = Vec{0, 0, 0}
	ents[2].vel = Vec{0, 0, 0}
}

func serialize(buf io.Writer, serAll bool) {
	bitMask := make([]byte, 1)
	bufTemp := &bytes.Buffer{}
	for i, ent := range ents {
		if serAll || ent.model != entsOld[i].model {
			bitMask[0] |= 1 << uint(i)
			binary.Write(bufTemp, binary.LittleEndian, ent.model)
		}
	}
	buf.Write(bitMask)
	buf.Write(bufTemp.Bytes())

	bitMask[0] = 0
	bufTemp.Reset()
	for i, ent := range ents {
		if serAll || !ent.pos.Equals(&entsOld[i].pos) {
			bitMask[0] |= 1 << uint(i)
			binary.Write(bufTemp, binary.LittleEndian, ent.pos)
		}
	}
	buf.Write(bitMask)
	buf.Write(bufTemp.Bytes())

	bitMask[0] = 0
	bufTemp.Reset()
	for i, ent := range ents {
		if serAll || !ent.vel.Equals(&entsOld[i].vel) {
			bitMask[0] |= 1 << uint(i)
			binary.Write(bufTemp, binary.LittleEndian, ent.vel)
		}
	}
	buf.Write(bitMask)
	buf.Write(bufTemp.Bytes())

	bitMask[0] = 0
	bufTemp.Reset()
	for i, ent := range ents {
		if serAll || !ent.size.Equals(&entsOld[i].size) {
			bitMask[0] |= 1 << uint(i)
			binary.Write(bufTemp, binary.LittleEndian, ent.size)
		}
	}
	buf.Write(bitMask)
	buf.Write(bufTemp.Bytes())
}
