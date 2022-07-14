package tcp

import (
	"net"
	"sync"
)

var (
	g_conn_id uint64
)

type TcpServer struct {
	ln      net.Listener
	el      IEventListener
	codec   ICodec
	opts    tcpOptions
	wg      sync.WaitGroup
	m_conns sync.Map
}

func (this *TcpServer) RunLoop(addr string, el IEventListener, codec ICodec, opts ...TcpOption) (err error ) {
	this.el = el
	this.codec = codec
	for _, opt := range opts {
		opt.f(&this.opts)
	}

	this.ln, err = net.Listen("tcp", addr)
	if err != nil {
		return
	}

	for {
		conn, err := this.ln.Accept()
		if err != nil {
			return err
		}

		this.wg.Add(1)
		go this.onAccept(conn)
	}
}

func (this *TcpServer) Stop() {
	if this.ln != nil {
		this.ln.Close()
	}

	this.m_conns.Range(func(k, v interface{}) bool{
		v.(ISession).Close()
		return true
	})
	this.wg.Wait()
}

func (this *TcpServer) onAccept(conn net.Conn) {
	defer this.wg.Done()

	session := createTcpSession(conn, this.opts)
	this.m_conns.Store(session.ID(), session)

	defer func(){
		this.m_conns.Delete(session.ID())
		session.Close()
	}()

	this.el.OnConnect(session)

	//写routine
	this.wg.Add(1)
	go func(){
		defer this.wg.Done()
		session.sendLoop(this.codec)
	}()

	//读routine
	session.readLoop(this.el, this.codec)

}
