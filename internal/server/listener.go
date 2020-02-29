package server

import (
	"errors"
	"fmt"
	"gola/internal/bootstrap"
	"gola/internal/logger"
	"math"
	"net"
	"os"
	"sync"
)

// DozListener 監聽
type DozListener struct {
	sync.Once
	net.Listener
	buf   chan struct{}
	out   chan struct{}
	errCh chan error
	debug bool
}

// DozConn 連線
type DozConn struct {
	net.Conn
	onClosed func()
}

// Close 關閉連線
func (conn *DozConn) Close() error {
	// log.Println("連線關閉...", conn.Conn.RemoteAddr().String())
	err := conn.Conn.Close()
	conn.onClosed()
	return err
}

func (l *DozListener) buffOut() {
	select {
	case <-l.errCh:
		l.out <- <-l.buf
	default:
		<-l.buf
	}
}

func (l *DozListener) onlyDoWithBuffer(do func() error) error {
	if l.buf != nil {
		return do()
	}

	return nil
}

// Accept 接收連線
func (l *DozListener) Accept() (net.Conn, error) {

	err := l.onlyDoWithBuffer(func() error {
		// 如果現在Buffer滿了或關閉了，不接收連線
		select {
		case <-l.errCh:
			return errors.New("DozListener Closed")
		case l.buf <- struct{}{}:
			// log.Println("等待連線...")
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	conn, err := l.Listener.Accept()
	if err != nil {
		doErr := l.onlyDoWithBuffer(func() error {
			l.buffOut()
			select {
			case <-l.errCh:
				return errors.New("DozListener Closed")
			default:
			}
			return nil
		})

		if doErr != nil {
			err = doErr
		}

		if l.debug {
			logger.Warn("接收連線Error... " + err.Error())
		}
		return nil, err
	}
	// log.Println("接收連線...", conn.RemoteAddr().String())

	return &DozConn{
		Conn:     conn,
		onClosed: l.buffOut,
	}, nil
}

// Close 關閉監聽
func (l *DozListener) Close() error {
	var err error
	l.Do(func() {
		l.onlyDoWithBuffer(func() error {
			l.out = make(chan struct{}, cap(l.buf))
			// close(l.buf)
			close(l.errCh)
			return nil
		})
		err = l.Listener.Close()
		text := "關閉監聽..."
		if err != nil {
			text += " " + err.Error()
		}
		if l.debug {
			logger.Warn(text)
		}
	})

	return err
}

// Wait 等待連線關閉
func (l *DozListener) Wait() error {
	err := l.onlyDoWithBuffer(func() error {
		sig := bootstrap.WaitOnceSignal()

		var notOver bool
		var count int

		count = len(l.buf)
		notOver = count > 0

		for notOver {
			if l.debug {
				logger.Warn(fmt.Sprintf("還有%d條連線，等待關閉...", count))
			}
			select {
			case <-l.out:
				count = len(l.buf)
				if count == 0 {
					notOver = false
					break
				}
			case <-sig:
				os.Exit(127)
			}
		}
		if l.debug {
			logger.Success("連線已經清空")
		}

		return nil
	})

	return err
}

// NewDozListner 建立新的監聽
func NewDozListner(l net.Listener, poolSize int, debug bool) *DozListener {
	dl := &DozListener{
		Listener: l,
		errCh:    make(chan error),
		debug:    debug,
	}

	if poolSize == 0 {
		poolSize = math.MaxInt64
	}
	dl.buf = make(chan struct{}, poolSize)
	return dl
}
