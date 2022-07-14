package tcp

import (
	"github.com/sirupsen/logrus"
	"net"
	"sync"
	"time"
)

type TcpClient struct {
	reconnect bool
	session   ISession
	el        IEventListener
	codec     ICodec

	b_stop				bool
	mutex  				sync.Mutex
	wg 					sync.WaitGroup

	opts tcpOptions
}

func (this *TcpClient) Dial(addr string, reconnect bool, el IEventListener, codec ICodec, opts ...TcpOption) (err error) {
	this.reconnect = reconnect
	this.el = el
	this.codec = codec
	for _, opt := range opts {
		opt.f(&this.opts)
	}

	var conn net.Conn
	for {
		//建立连接
		conn, err = net.DialTimeout("tcp", addr, time.Second * 3)
		if err != nil {
			logrus.Errorf("TcpClient::Dial add:%s err:%v", addr, err)
			if !this.reconnect {
				return
			}

			this.mutex.Lock()
			if this.b_stop {
				this.mutex.Unlock()
				return
			}
			this.mutex.Unlock()

			continue
		}

		//创建sesssion
		this.mutex.Lock()
		if this.b_stop {
			conn.Close()
			this.mutex.Unlock()
			return
		}
		session := createTcpSession(conn, this.opts)
		this.session = session
		this.mutex.Unlock()

		//working
		this.onConnect(session)

		//连接断开
		this.mutex.Lock()
		if this.b_stop {
			this.mutex.Unlock()
			return
		}
		this.mutex.Unlock()
	}

	return
}

func (this *TcpClient) Stop() {
	this.mutex.Lock()
	if this.b_stop {
		this.mutex.Unlock()
		return
	}
	this.b_stop = true
	if this.session != nil {
		this.session.Close()
	}
	this.mutex.Unlock()

	this.wg.Wait()
}

func (this *TcpClient) onConnect(session *tcpSession) {
	defer session.Close()

	this.el.OnConnect(session)

	//写routine
	this.wg.Add(1)
	go func() {
		this.wg.Done()
		session.sendLoop(this.codec)
	}()

	//读routine
	session.readLoop(this.el, this.codec)
}

