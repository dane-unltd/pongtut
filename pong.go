package main

import (
	"bytes"
	"code.google.com/p/go.net/websocket"
	"encoding/binary"
	"log"
	"net/http"
	"time"
)

type PlayerId uint32
type Action uint32

type UserCommand struct {
	Actions uint32
}

type ClientConn struct {
	cmdBuf     chan UserCommand
	currentCmd UserCommand
	id         PlayerId
	ws         *websocket.Conn
	inBuf      [1500]byte
}

var newConn = make(chan *ClientConn)
var clients = make(map[PlayerId]*ClientConn)

var maxId = PlayerId(0)

func newId() PlayerId {
	maxId++
	return maxId
}

func active(id PlayerId, action Action) bool {
	if (clients[id].currentCmd.Actions & (1 << action)) > 0 {
		return true
	}
	return false
}

func wsHandler(ws *websocket.Conn) {
	cl := &ClientConn{}
	cl.ws = ws
	cl.cmdBuf = make(chan UserCommand, 5)

	cmd := UserCommand{}

	log.Println("incoming connection")

	newConn <- cl
	for {
		pkt := cl.inBuf[0:]
		n, err := ws.Read(pkt)
		pkt = pkt[0:n]
		if err != nil {
			log.Println(err)
			break
		}
		buf := bytes.NewBuffer(pkt)
		err = binary.Read(buf, binary.LittleEndian, &cmd)
		if err != nil {
			log.Println(err)
			break
		}
		cl.cmdBuf <- cmd
	}
}

func main() {
	http.Handle("/ws/", websocket.Handler(wsHandler))
	http.Handle("/www/", http.StripPrefix("/www/",
		http.FileServer(http.Dir("./www"))))
	go func() {
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	//runing at 30 FPS
	frameNS := time.Duration(int(1e9) / 30)
	clk := time.NewTicker(frameNS)

	//main loop
	for {
		select {
		case <-clk.C:
			updateInputs()
			updateSimulation()
			sendUpdates()
		case cl := <-newConn:
			id := newId()
			clients[id] = cl
			login(id)

			buf := &bytes.Buffer{}
			serialize(buf, true)
			err := websocket.Message.Send(cl.ws, buf.Bytes())
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func updateInputs() {
	for _, cl := range clients {
		for {
			select {
			case cmd := <-cl.cmdBuf:
				cl.currentCmd = cmd
			default:
				goto done
			}
		}
	done:
	}
}

var removeList = make([]PlayerId, 3)

func sendUpdates() {
	buf := &bytes.Buffer{}
	serialize(buf, false)
	removeList = removeList[0:0]
	for id, cl := range clients {
		err := websocket.Message.Send(cl.ws, buf.Bytes())
		if err != nil {
			removeList = append(removeList, id)
			log.Println(err)
		}
	}
	copyState()
	for _, id := range removeList {
		delete(clients, id)
		disconnect(id)
	}
}
