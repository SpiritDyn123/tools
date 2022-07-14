package tcp

import (
	"github.com/sirupsen/logrus"
	"net"
	"sync"
	"sync/atomic"
)

func createTcpSession(conn net.Conn, opts tcpOptions) *tcpSession {
	return &tcpSession{
		id: atomic.AddUint64(&g_conn_id, 1),
		conn: conn,
		opts: opts,
		close_chan: make(chan struct{}),
		send_chan: make(chan []interface{}, 200),
	}
}

type tcpSession struct {
	id   uint64
	conn net.Conn
	opts tcpOptions

	mutex		sync.Mutex
	b_close    bool
	close_chan chan struct{}
	send_chan  chan []interface{}
	once sync.Once
}

func (this *tcpSession) ID() uint64 {
	return this.id
}

func (this *tcpSession) Addr() net.Addr {
	return this.conn.RemoteAddr()
}

func (this *tcpSession) Send(msgs ...interface{}) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if this.b_close {
		return
	}

	select {
	case this.send_chan <- msgs:
	default:
		//todo log
		this.Close()
	}
}

func (this *tcpSession) Close() {
	this.once.Do(func(){
		this.mutex.Lock()
		defer this.mutex.Unlock()

		this.b_close = true
		this.conn.Close()
		close(this.close_chan)
		close(this.send_chan)
	})
}

func (this *tcpSession) sendLoop(codec ICodec) {
	for {
		select {
		case msg, ok := <- this.send_chan:
			if !ok {
				return
			}
			err := codec.Write(this.conn, msg, &this.opts)
			if err != nil {
				logrus.Errorf("tcpSession::sendLoop write err:%v", err)
				goto __CLOSE_END
			}
		case <-this.close_chan:
			goto __CLOSE_END
		}
	}

__CLOSE_END:
	this.Close()
}

func (this *tcpSession) readLoop(el IEventListener, codec ICodec) {
	for {
		msgs, err := codec.Read(this.conn, &this.opts)
		if err != nil {
			logrus.Errorf("tcpSession::readLoop read err:%v", err)
			el.OnDisconnect(this, err)
			return
		}

		el.OnMessage(this, msgs)
	}
}

