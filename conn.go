package main

import (
	"github.com/pebbe/zmq4"
	"encoding/json"
	"os"
	"fmt"
	"io"
	"time"
)

type RawMessage []interface{}

func parseToMessage(msg []byte) RawMessage {
	var r RawMessage
	json.Unmarshal(msg, &r)
	return r
}
func (r RawMessage) toCommand() Command {
	var m Command

	kind := r[0].(string)

	if kind == "Subscribe" {
		m = NewSubscribeCommand(r)
	} else if kind == "PatchAttrs" {
		m = NewPatchAttrsCommand(r)
	} else if kind == "PostElem" {
		m = NewPostElemCommand(r)
	} else if kind == "PutElem" {
		m = NewPutElemCommand(r)
	} else if kind == "DeleteElem" {
		m = NewDeleteElemCommand(r)
	} else if kind == "PatchStyles" {
		m = NewPatchStylesCommand(r)
	} else if kind == "Close" {
		m = CloseCommand{}
	} else {
		m = NewErrCommand(r)
	}
	return m
}


type Conn interface {
	Start()
	Send(string)
	Recv() Command
	Stop()
}

type ZMQConn struct {
	addr string
	sock *zmq4.Socket
}

func NewZMQConn(addr string) ZMQConn {
	
	fmt.Printf(`{"guise":"%s"}`, addr)
	os.Stdout.Close() // hmm, this seems like a bad global effect...
	sock, _ := zmq4.NewSocket(zmq4.PAIR)
	return ZMQConn{addr, sock}
}

func (z ZMQConn) Start() {
	z.sock.Bind(z.addr)
}

func (z ZMQConn) Send(s string) {
	z.sock.Send(s, 0)
}

func (z ZMQConn) Recv() Command {
	someBytes, _ := z.sock.RecvBytes(0)
	os.Stderr.WriteString("guise: " + string(someBytes) + "\n")
	r := parseToMessage(someBytes)
	return r.toCommand()
}

func (z ZMQConn) Stop() {
	z.sock.Close()
}

type StreamConn struct {
	in io.Reader
	decoder *json.Decoder
	out io.Writer
}

func StdioConn() StreamConn {
	d := json.NewDecoder(os.Stdin)
	return StreamConn{os.Stdin, d, os.Stdout}
}

func FileConn(in string, out string) StreamConn {
	f, _ := os.Open(in)
	g, _ := os.OpenFile(out, os.O_WRONLY, 777)
	d := json.NewDecoder(f)
	return StreamConn{f, d, g}
}

func (f StreamConn) Start() {}
func (f StreamConn) Send(s string) {
	f.out.Write([]byte(s + "\n\x00"))
}
func (f StreamConn) Recv() Command {
	var r RawMessage
	for ! f.decoder.More() {
		time.Sleep(time.Millisecond * 10)
	}
	f.decoder.Decode(&r)
	return r.toCommand()
}
func (f StreamConn) Stop() {}