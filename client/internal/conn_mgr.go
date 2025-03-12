package internal

import (
	"context"
	"os"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/netbirdio/netbird/client/internal/lazyconn"
	lazyConnManager "github.com/netbirdio/netbird/client/internal/lazyconn/manager"
	"github.com/netbirdio/netbird/client/internal/peer"
	"github.com/netbirdio/netbird/client/internal/peerstore"
)

const (
	envDisableLazyConn = "NB_LAZY_CONN_DISABLE"
)

// ConnMgr coordinates both lazy connections (established on-demand) and permanent peer connections.
//
// The connection manager is responsible for:
// - Managing lazy connections via the lazyConnManager
// - Maintaining a list of excluded peers that should always have permanent connections
// - Handling connection establishment based on peer signaling
type ConnMgr struct {
	peerStore   *peerstore.Store
	lazyConnMgr *lazyConnManager.Manager

	connStateListener *peer.ConnectionListener

	wg        sync.WaitGroup
	ctxCancel context.CancelFunc
}

func NewConnMgr(peerStore *peerstore.Store, iface lazyconn.WGIface, dispatcher *peer.ConnectionDispatcher) *ConnMgr {
	var lazyConnMgr *lazyConnManager.Manager
	if os.Getenv(envDisableLazyConn) != "true" {
		lazyConnMgr = lazyConnManager.NewManager(iface, dispatcher)
	}

	e := &ConnMgr{
		peerStore:   peerStore,
		lazyConnMgr: lazyConnMgr,
	}
	return e
}

func (e *ConnMgr) Start(parentCtx context.Context) {
	if e.lazyConnMgr == nil {
		log.Infof("lazy connection manager is disabled")
		return
	}

	ctx, cancel := context.WithCancel(parentCtx)
	e.ctxCancel = cancel

	e.wg.Add(1)
	go e.receiveLazyEvents(ctx)
}

func (e *ConnMgr) AddExcludeFromLazyConnection(peerID string) {
	e.lazyConnMgr.ExcludePeer(peerID)
}

func (e *ConnMgr) AddPeerConn(peerKey string, conn *peer.Conn) (exists bool) {
	if success := e.peerStore.AddPeerConn(peerKey, conn); !success {
		return true
	}

	if !e.isStarted() {
		conn.Open()
		return
	}

	lazyPeerCfg := lazyconn.PeerConfig{
		PublicKey:  peerKey,
		AllowedIPs: conn.WgConfig().AllowedIps,
		PeerConnID: conn.ConnID(),
	}
	excluded, err := e.lazyConnMgr.AddPeer(lazyPeerCfg)
	if err != nil {
		conn.Log.Errorf("failed to add peer to lazyconn manager: %v", err)
		conn.Open()
		return
	}

	if excluded {
		conn.Log.Infof("peer is on lazy conn manager exclude list, opening connection")
		conn.Open()
		return
	}

	conn.Log.Infof("peer added to lazy conn manager")
	return
}

func (e *ConnMgr) OnSignalMsg(peerKey string) (*peer.Conn, bool) {
	conn, ok := e.peerStore.PeerConn(peerKey)
	if !ok {
		return nil, false
	}

	if !e.isStarted() {
		return conn, true
	}

	if ok := e.lazyConnMgr.RemovePeer(peerKey); ok {
		conn.Log.Infof("removed peer from lazy conn manager")
		conn.Open()
	}
	return conn, true
}

func (e *ConnMgr) RemovePeerConn(peerKey string) {
	conn, ok := e.peerStore.Remove(peerKey)
	if !ok {
		return
	}
	defer conn.Close()

	if !e.isStarted() {
		return
	}

	e.lazyConnMgr.RemovePeer(peerKey)
	conn.Log.Infof("removed peer from lazy conn manager")
}

func (e *ConnMgr) Close() {
	if !e.isStarted() {
		return
	}

	e.ctxCancel()
	e.lazyConnMgr.Close()
	e.wg.Wait()
	e.lazyConnMgr = nil
}

func (e *ConnMgr) receiveLazyEvents(ctx context.Context) {
	defer e.wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case peerID := <-e.lazyConnMgr.OnDemand:
			e.peerStore.PeerConnOpen(peerID)
		case peerID := <-e.lazyConnMgr.Idle:
			e.peerStore.PeerConnClose(peerID)
		}
	}
}

func (e *ConnMgr) isStarted() bool {
	return e.lazyConnMgr != nil && e.ctxCancel != nil
}
