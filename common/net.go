package common

import (
	"context"
	"crypto/tls"
	"net"
	"sync"

	log "github.com/liudanking/goutil/logutil"

	quic "github.com/lucas-clemente/quic-go"
)

const (
	KQuicProxy = "quic-proxy"
)

type QuicListener struct {
	quic.Listener
	chAcceptConn chan *AcceptConn
}

type AcceptConn struct {
	conn net.Conn
	err  error
}

func NewQuicListener(l quic.Listener) *QuicListener {
	ql := &QuicListener{
		Listener:     l,
		chAcceptConn: make(chan *AcceptConn, 1),
	}
	go ql.doAccept()
	return ql
}

func (ql *QuicListener) doAccept() {
	for {
		sess, err := ql.Listener.Accept(context.TODO())
		if err != nil {
			log.Error("accept session failed:%v", err)
			continue
		}
		log.Info("accept a session")

		go func(sess quic.Session) {
			for {
				stream, err := sess.AcceptStream(context.TODO())
				if err != nil {
					log.Notice("accept stream failed:%v", err)
					sess.Close()
					return
				}
				log.Info("accept stream %v", stream.StreamID())
				ql.chAcceptConn <- &AcceptConn{
					conn: &QuicStream{sess: sess, Stream: stream},
					err:  nil,
				}
			}
		}(sess)
	}
}

func (ql *QuicListener) Accept() (net.Conn, error) {
	ac := <-ql.chAcceptConn
	return ac.conn, ac.err
}

type QuicStream struct {
	sess quic.Session
	quic.Stream
}

func (qs *QuicStream) LocalAddr() net.Addr {
	return qs.sess.LocalAddr()
}

func (qs *QuicStream) RemoteAddr() net.Addr {
	return qs.sess.RemoteAddr()
}

type QuicDialer struct {
	skipCertVerify bool
	sess           quic.Session
	sync.Mutex
}

func NewQuicDialer(skipCertVerify bool) *QuicDialer {
	return &QuicDialer{
		skipCertVerify: skipCertVerify,
	}
}

func (qd *QuicDialer) Dial(network, addr string) (net.Conn, error) {
	qd.Lock()
	defer qd.Unlock()

	if qd.sess == nil {
		sess, err := quic.DialAddr(addr, &tls.Config{
			InsecureSkipVerify: qd.skipCertVerify,
			NextProtos:         []string{KQuicProxy},
		}, nil)
		if err != nil {
			log.Error("dial session failed:%v", err)
			return nil, err
		}
		qd.sess = sess
	}

	stream, err := qd.sess.OpenStreamSync(context.TODO())
	if err != nil {
		log.Info("[1/2] open stream from session no success:%v, try to open new session", err)
		qd.sess.Close()
		sess, err := quic.DialAddr(addr, &tls.Config{
			InsecureSkipVerify: true,
			NextProtos:         []string{KQuicProxy},
		}, nil)
		if err != nil {
			log.Error("[2/2] dial new session failed:%v", err)
			return nil, err
		}
		qd.sess = sess

		stream, err = qd.sess.OpenStreamSync(context.TODO())
		if err != nil {
			log.Error("[2/2] open stream from new session failed:%v", err)
			return nil, err
		}
		log.Info("[2/2] open stream from new session OK")
	}

	log.Info("addr:%s, stream_id:%v", addr, stream.StreamID())
	return &QuicStream{sess: qd.sess, Stream: stream}, nil
}
