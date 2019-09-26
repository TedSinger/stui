package main

import (
	"sync"
	"github.com/pebbe/zmq4"
	"encoding/json"
	"os"
	"io"
)

type RawMessage []interface{}

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

type Duplex struct {
	In chan Command
	Out chan string
}

type Conn interface {
	Start(*sync.WaitGroup)
	GetDuplex() *Duplex
}

type ZMQConn struct {
	duplex *Duplex
	addr string
	sock *zmq4.Socket
}

func NewZMQConn(duplex *Duplex, addr string) ZMQConn {
	sock, _ := zmq4.NewSocket(zmq4.PAIR)
	return ZMQConn{duplex, addr, sock}
}

func (z ZMQConn) GetDuplex() *Duplex {return z.duplex}

func (z ZMQConn) Start(wg *sync.WaitGroup) {
	z.sock.Bind(z.addr)
	go z.Recv()
	for msg := range z.duplex.Out {
		z.sock.Send(msg, 0)
	}
	z.sock.Close()
	wg.Done()
}

func (z ZMQConn) Recv() {
	for {
		someBytes, err := z.sock.RecvBytes(0)
		if err != nil {
			close(z.duplex.In)
			break
		} else {
			var r RawMessage
			json.Unmarshal(someBytes, &r)
			z.duplex.In <- r.toCommand()
		}
	}
}

type StreamConn struct {
	duplex *Duplex
	in io.Reader
	decoder *json.Decoder
	out io.Writer
}

func StdioConn(duplex *Duplex) StreamConn {
	in := os.Stdin
	d := json.NewDecoder(in)
	return StreamConn{duplex, in, d, os.Stdout}
}

func FileConn(duplex *Duplex, in string, out string) StreamConn {
	f, _ := os.Open(in)
	g, _ := os.OpenFile(out, os.O_WRONLY, 777)
	d := json.NewDecoder(f)
	return StreamConn{duplex, f, d, g}
}

func (f StreamConn) Start(wg *sync.WaitGroup) {
	go f.Recv()
	for msg := range f.duplex.Out {
		f.out.Write([]byte(msg + "\n"))
	}
	wg.Done()
}

func (f StreamConn) Recv() {
	for {
		var r RawMessage
		err := f.decoder.Decode(&r)
		if err == io.EOF {
			close(f.duplex.In)
			break
		} else {
			f.duplex.In <- r.toCommand()
		}	
	}
}

func (f StreamConn) GetDuplex() *Duplex {return f.duplex}