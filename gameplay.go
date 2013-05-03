package main

import (
	"bytes"
	"encoding/binary"
	"io"
	"math"
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
	score          uint32
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
	processInput()
	collisionCheck()
	move()
}

func move() {
	for i := range ents {
		ents[i].pos.Add(&ents[i].pos, &ents[i].vel)
	}
}

func processInput() {
	if players[0] == 0 || players[1] == 0 {
		return
	}

	newVel := 0.0
	if active(players[0], Up) {
		newVel += 5
	}
	if active(players[0], Down) {
		newVel -= 5
	}
	ents[0].vel[1] = newVel

	newVel = 0.0
	if active(players[1], Up) {
		newVel += 5
	}
	if active(players[1], Down) {
		newVel -= 5
	}
	ents[1].vel[1] = newVel
}

const FieldHeight = 120

func collisionCheck() {
	for i := range ents {
		if ents[i].pos[1] > FieldHeight/2-ents[i].size[1]/2 {
			ents[i].pos[1] = FieldHeight/2 - ents[i].size[1]/2
			if ents[i].vel[1] > 0 {
				ents[i].vel[1] = -ents[i].vel[1]
			}
		}
		if ents[i].pos[1] < -FieldHeight/2+ents[i].size[1]/2 {
			ents[i].pos[1] = -FieldHeight/2 + ents[i].size[1]/2
			if ents[i].vel[1] < 0 {
				ents[i].vel[1] = -ents[i].vel[1]
			}
		}
	}

	rSq := ents[2].size[0] / 2
	rSq *= rSq
	for i := 0; i < 2; i++ {
		//v points from the center of the paddel to the point on the
		//border of the paddel which is closest to the sphere
		v := Vec{}
		v.Sub(&ents[2].pos, &ents[i].pos)
		v.Clamp(&ents[i].size)

		//d is the vector between the closest points on the paddle and
		//the sphere
		d := Vec{}
		d.Sub(&ents[2].pos, &ents[i].pos)
		d.Sub(&d, &v)

		distSq := d.Nrm2Sq()
		if distSq < rSq {
			//move the sphere in direction of d to remove the
			//penetration
			dPos := Vec{}
			dPos.Scale(math.Sqrt(rSq/distSq)-1, &d)
			ents[2].pos.Add(&ents[2].pos, &dPos)

			dotPr := Dot(&ents[2].vel, &d)
			if dotPr < 0 {
				d.Scale(-2*dotPr/distSq, &d)
				ents[2].vel.Add(&ents[2].vel, &d)
			}
		}
	}

	if ents[2].pos[0] < -100 {
		ents[2].pos = Vec{0, 0, 0}
		ents[2].vel = Vec{2, 3, 0}
		ents[1].score++
	} else if ents[2].pos[0] > 100 {
		ents[2].pos = Vec{0, 0, 0}
		ents[2].vel = Vec{-2, 3, 0}
		ents[0].score++
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
	ents[0].score = 0
	ents[1].score = 0
	ents[2].vel = Vec{2, 3, 0}
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
	ents[0].pos = Vec{-75, 0, 0}
	ents[1].pos = Vec{75, 0, 0}
	ents[2].pos = Vec{0, 0, 0}
	ents[2].vel = Vec{0, 0, 0}
}

func serialize(buf io.Writer, serAll bool) {
	bitMask := make([]byte, 1)
	bufTemp := &bytes.Buffer{}
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

	bitMask[0] = 0
	bufTemp.Reset()
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
		if serAll || ent.score != entsOld[i].score {
			bitMask[0] |= 1 << uint(i)
			binary.Write(bufTemp, binary.LittleEndian, ent.score)
		}
	}
	buf.Write(bitMask)
	buf.Write(bufTemp.Bytes())
}
