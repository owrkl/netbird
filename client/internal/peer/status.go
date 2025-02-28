package peer

import (
	"errors"
	"net/netip"
	"slices"
	"sync"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/maps"
	"google.golang.org/grpc/codes"
	gstatus "google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/netbirdio/netbird/client/iface/configurer"
	"github.com/netbirdio/netbird/client/internal/relay"
	"github.com/netbirdio/netbird/client/proto"
	"github.com/netbirdio/netbird/management/domain"
	relayClient "github.com/netbirdio/netbird/relay/client"
)

const eventQueueSize = 10

type ResolvedDomainInfo struct {
	Prefixes     []netip.Prefix
	ParentDomain domain.Domain
}

type EventListener interface {
	OnEvent(event *proto.SystemEvent)
}

// State contains the latest state of a peer
type State struct {
	Mux                        *sync.RWMutex
	IP                         string
	PubKey                     string
	FQDN                       string
	ConnStatus                 ConnStatus
	ConnStatusUpdate           time.Time
	Relayed                    bool
	LocalIceCandidateType      string
	RemoteIceCandidateType     string
	LocalIceCandidateEndpoint  string
	RemoteIceCandidateEndpoint string
	RelayServerAddress         string
	LastWireguardHandshake     time.Time
	BytesTx                    int64
	BytesRx                    int64
	Latency                    time.Duration
	RosenpassEnabled           bool
	routes                     map[string]struct{}
}

// AddRoute add a single route to routes map
func (s *State) AddRoute(network string) {
	s.Mux.Lock()
	defer s.Mux.Unlock()
	if s.routes == nil {
		s.routes = make(map[string]struct{})
	}
	s.routes[network] = struct{}{}
}

// SetRoutes set state routes
func (s *State) SetRoutes(routes map[string]struct{}) {
	s.Mux.Lock()
	defer s.Mux.Unlock()
	s.routes = routes
}

// DeleteRoute removes a route from the network amp
func (s *State) DeleteRoute(network string) {
	s.Mux.Lock()
	defer s.Mux.Unlock()
	delete(s.routes, network)
}

// GetRoutes return routes map
func (s *State) GetRoutes() map[string]struct{} {
	s.Mux.RLock()
	defer s.Mux.RUnlock()
	return maps.Clone(s.routes)
}

// LocalPeerState contains the latest state of the local peer
type LocalPeerState struct {
	IP              string
	PubKey          string
	KernelInterface bool
	FQDN            string
	Routes          map[string]struct{}
}

// Clone returns a copy of the LocalPeerState
func (l LocalPeerState) Clone() LocalPeerState {
	l.Routes = maps.Clone(l.Routes)
	return l
}

// SignalState contains the latest state of a signal connection
type SignalState struct {
	URL       string
	Connected bool
	Error     error
}

// ManagementState contains the latest state of a management connection
type ManagementState struct {
	URL       string
	Connected bool
	Error     error
}

// RosenpassState contains the latest state of the Rosenpass configuration
type RosenpassState struct {
	Enabled    bool
	Permissive bool
}

// NSGroupState represents the status of a DNS server group, including associated domains,
// whether it's enabled, and the last error message encountered during probing.
type NSGroupState struct {
	ID      string
	Servers []string
	Domains []string
	Enabled bool
	Error   error
}

// FullStatus contains the full state held by the Status instance
type FullStatus struct {
	Peers           []State
	ManagementState ManagementState
	SignalState     SignalState
	LocalPeerState  LocalPeerState
	RosenpassState  RosenpassState
	Relays          []relay.ProbeResult
	NSGroupStates   []NSGroupState
}

// Status holds a state of peers, signal, management connections and relays
type Status struct {
	mux                   sync.Mutex
	peers                 map[string]State
	changeNotify          map[string]chan struct{}
	signalState           bool
	signalError           error
	managementState       bool
	managementError       error
	relayStates           []relay.ProbeResult
	localPeer             LocalPeerState
	offlinePeers          []State
	mgmAddress            string
	signalAddress         string
	notifier              *notifier
	rosenpassEnabled      bool
	rosenpassPermissive   bool
	nsGroupStates         []NSGroupState
	resolvedDomainsStates map[domain.Domain]ResolvedDomainInfo

	// To reduce the number of notification invocation this bool will be true when need to call the notification
	// Some Peer actions mostly used by in a batch when the network map has been synchronized. In these type of events
	// set to true this variable and at the end of the processing we will reset it by the FinishPeerListModifications()
	peerListChangedForNotification bool

	relayMgr *relayClient.Manager

	eventMux     sync.RWMutex
	eventStreams map[string]chan *proto.SystemEvent
	eventQueue   *EventQueue
}

// NewRecorder returns a new Status instance
func NewRecorder(mgmAddress string) *Status {
	return &Status{
		peers:                 make(map[string]State),
		changeNotify:          make(map[string]chan struct{}),
		eventStreams:          make(map[string]chan *proto.SystemEvent),
		eventQueue:            NewEventQueue(eventQueueSize),
		offlinePeers:          make([]State, 0),
		notifier:              newNotifier(),
		mgmAddress:            mgmAddress,
		resolvedDomainsStates: map[domain.Domain]ResolvedDomainInfo{},
	}
}

func (d *Status) SetRelayMgr(manager *relayClient.Manager) {
	d.mux.Lock()
	defer d.mux.Unlock()
	d.relayMgr = manager
}

// ReplaceOfflinePeers replaces
func (d *Status) ReplaceOfflinePeers(replacement []State) {
	d.mux.Lock()
	defer d.mux.Unlock()
	d.offlinePeers = make([]State, len(replacement))
	copy(d.offlinePeers, replacement)

	// todo we should set to true in case if the list changed only
	d.peerListChangedForNotification = true
}

// AddPeer adds peer to Daemon status map
func (d *Status) AddPeer(peerPubKey string, fqdn string, ip string) error {
	d.mux.Lock()
	defer d.mux.Unlock()

	_, ok := d.peers[peerPubKey]
	if ok {
		return errors.New("peer already exist")
	}
	d.peers[peerPubKey] = State{
		PubKey:     peerPubKey,
		IP:         ip,
		ConnStatus: StatusDisconnected,
		FQDN:       fqdn,
		Mux:        new(sync.RWMutex),
	}
	d.peerListChangedForNotification = true
	return nil
}

// GetPeer adds peer to Daemon status map
func (d *Status) GetPeer(peerPubKey string) (State, error) {
	d.mux.Lock()
	defer d.mux.Unlock()

	state, ok := d.peers[peerPubKey]
	if !ok {
		return State{}, configurer.ErrPeerNotFound
	}
	return state, nil
}

// RemovePeer removes peer from Daemon status map
func (d *Status) RemovePeer(peerPubKey string) error {
	d.mux.Lock()
	defer d.mux.Unlock()

	_, ok := d.peers[peerPubKey]
	if !ok {
		return errors.New("no peer with to remove")
	}

	delete(d.peers, peerPubKey)
	d.peerListChangedForNotification = true
	return nil
}

// UpdatePeerState updates peer status
func (d *Status) UpdatePeerState(receivedState State) error {
	d.mux.Lock()
	defer d.mux.Unlock()

	peerState, ok := d.peers[receivedState.PubKey]
	if !ok {
		return errors.New("peer doesn't exist")
	}

	if receivedState.IP != "" {
		peerState.IP = receivedState.IP
	}

	skipNotification := shouldSkipNotify(receivedState.ConnStatus, peerState)

	if receivedState.ConnStatus != peerState.ConnStatus {
		peerState.ConnStatus = receivedState.ConnStatus
		peerState.ConnStatusUpdate = receivedState.ConnStatusUpdate
		peerState.Relayed = receivedState.Relayed
		peerState.LocalIceCandidateType = receivedState.LocalIceCandidateType
		peerState.RemoteIceCandidateType = receivedState.RemoteIceCandidateType
		peerState.LocalIceCandidateEndpoint = receivedState.LocalIceCandidateEndpoint
		peerState.RemoteIceCandidateEndpoint = receivedState.RemoteIceCandidateEndpoint
		peerState.RelayServerAddress = receivedState.RelayServerAddress
		peerState.RosenpassEnabled = receivedState.RosenpassEnabled
	}

	d.peers[receivedState.PubKey] = peerState

	if skipNotification {
		return nil
	}

	d.notifyPeerListChanged()
	return nil
}

func (d *Status) AddPeerStateRoute(peer string, route string) error {
	d.mux.Lock()
	defer d.mux.Unlock()

	peerState, ok := d.peers[peer]
	if !ok {
		return errors.New("peer doesn't exist")
	}

	peerState.AddRoute(route)
	d.peers[peer] = peerState

	// todo: consider to make sense of this notification or not
	d.notifyPeerListChanged()
	return nil
}

func (d *Status) RemovePeerStateRoute(peer string, route string) error {
	d.mux.Lock()
	defer d.mux.Unlock()

	peerState, ok := d.peers[peer]
	if !ok {
		return errors.New("peer doesn't exist")
	}

	peerState.DeleteRoute(route)
	d.peers[peer] = peerState

	// todo: consider to make sense of this notification or not
	d.notifyPeerListChanged()
	return nil
}

func (d *Status) UpdatePeerICEState(receivedState State) error {
	d.mux.Lock()
	defer d.mux.Unlock()

	peerState, ok := d.peers[receivedState.PubKey]
	if !ok {
		return errors.New("peer doesn't exist")
	}

	if receivedState.IP != "" {
		peerState.IP = receivedState.IP
	}

	skipNotification := shouldSkipNotify(receivedState.ConnStatus, peerState)

	peerState.ConnStatus = receivedState.ConnStatus
	peerState.ConnStatusUpdate = receivedState.ConnStatusUpdate
	peerState.Relayed = receivedState.Relayed
	peerState.LocalIceCandidateType = receivedState.LocalIceCandidateType
	peerState.RemoteIceCandidateType = receivedState.RemoteIceCandidateType
	peerState.LocalIceCandidateEndpoint = receivedState.LocalIceCandidateEndpoint
	peerState.RemoteIceCandidateEndpoint = receivedState.RemoteIceCandidateEndpoint
	peerState.RosenpassEnabled = receivedState.RosenpassEnabled

	d.peers[receivedState.PubKey] = peerState

	if skipNotification {
		return nil
	}

	d.notifyPeerStateChangeListeners(receivedState.PubKey)
	d.notifyPeerListChanged()
	return nil
}

func (d *Status) UpdatePeerRelayedState(receivedState State) error {
	d.mux.Lock()
	defer d.mux.Unlock()

	peerState, ok := d.peers[receivedState.PubKey]
	if !ok {
		return errors.New("peer doesn't exist")
	}

	skipNotification := shouldSkipNotify(receivedState.ConnStatus, peerState)

	peerState.ConnStatus = receivedState.ConnStatus
	peerState.ConnStatusUpdate = receivedState.ConnStatusUpdate
	peerState.Relayed = receivedState.Relayed
	peerState.RelayServerAddress = receivedState.RelayServerAddress
	peerState.RosenpassEnabled = receivedState.RosenpassEnabled

	d.peers[receivedState.PubKey] = peerState

	if skipNotification {
		return nil
	}

	d.notifyPeerStateChangeListeners(receivedState.PubKey)
	d.notifyPeerListChanged()
	return nil
}

func (d *Status) UpdatePeerRelayedStateToDisconnected(receivedState State) error {
	d.mux.Lock()
	defer d.mux.Unlock()

	peerState, ok := d.peers[receivedState.PubKey]
	if !ok {
		return errors.New("peer doesn't exist")
	}

	skipNotification := shouldSkipNotify(receivedState.ConnStatus, peerState)

	peerState.ConnStatus = receivedState.ConnStatus
	peerState.Relayed = receivedState.Relayed
	peerState.ConnStatusUpdate = receivedState.ConnStatusUpdate
	peerState.RelayServerAddress = ""

	d.peers[receivedState.PubKey] = peerState

	if skipNotification {
		return nil
	}

	d.notifyPeerStateChangeListeners(receivedState.PubKey)
	d.notifyPeerListChanged()
	return nil
}

func (d *Status) UpdatePeerICEStateToDisconnected(receivedState State) error {
	d.mux.Lock()
	defer d.mux.Unlock()

	peerState, ok := d.peers[receivedState.PubKey]
	if !ok {
		return errors.New("peer doesn't exist")
	}

	skipNotification := shouldSkipNotify(receivedState.ConnStatus, peerState)

	peerState.ConnStatus = receivedState.ConnStatus
	peerState.Relayed = receivedState.Relayed
	peerState.ConnStatusUpdate = receivedState.ConnStatusUpdate
	peerState.LocalIceCandidateType = receivedState.LocalIceCandidateType
	peerState.RemoteIceCandidateType = receivedState.RemoteIceCandidateType
	peerState.LocalIceCandidateEndpoint = receivedState.LocalIceCandidateEndpoint
	peerState.RemoteIceCandidateEndpoint = receivedState.RemoteIceCandidateEndpoint

	d.peers[receivedState.PubKey] = peerState

	if skipNotification {
		return nil
	}

	d.notifyPeerStateChangeListeners(receivedState.PubKey)
	d.notifyPeerListChanged()
	return nil
}

// UpdateWireGuardPeerState updates the WireGuard bits of the peer state
func (d *Status) UpdateWireGuardPeerState(pubKey string, wgStats configurer.WGStats) error {
	d.mux.Lock()
	defer d.mux.Unlock()

	peerState, ok := d.peers[pubKey]
	if !ok {
		return errors.New("peer doesn't exist")
	}

	peerState.LastWireguardHandshake = wgStats.LastHandshake
	peerState.BytesRx = wgStats.RxBytes
	peerState.BytesTx = wgStats.TxBytes

	d.peers[pubKey] = peerState

	return nil
}

func shouldSkipNotify(receivedConnStatus ConnStatus, curr State) bool {
	switch {
	case receivedConnStatus == StatusConnecting:
		return true
	case receivedConnStatus == StatusDisconnected && curr.ConnStatus == StatusConnecting:
		return true
	case receivedConnStatus == StatusDisconnected && curr.ConnStatus == StatusDisconnected:
		return curr.IP != ""
	default:
		return false
	}
}

// UpdatePeerFQDN update peer's state fqdn only
func (d *Status) UpdatePeerFQDN(peerPubKey, fqdn string) error {
	d.mux.Lock()
	defer d.mux.Unlock()

	peerState, ok := d.peers[peerPubKey]
	if !ok {
		return errors.New("peer doesn't exist")
	}

	peerState.FQDN = fqdn
	d.peers[peerPubKey] = peerState

	return nil
}

// FinishPeerListModifications this event invoke the notification
func (d *Status) FinishPeerListModifications() {
	d.mux.Lock()

	if !d.peerListChangedForNotification {
		d.mux.Unlock()
		return
	}
	d.peerListChangedForNotification = false
	d.mux.Unlock()

	d.notifyPeerListChanged()
}

// GetPeerStateChangeNotifier returns a change notifier channel for a peer
func (d *Status) GetPeerStateChangeNotifier(peer string) <-chan struct{} {
	d.mux.Lock()
	defer d.mux.Unlock()

	ch, found := d.changeNotify[peer]
	if found {
		return ch
	}

	ch = make(chan struct{})
	d.changeNotify[peer] = ch
	return ch
}

// GetLocalPeerState returns the local peer state
func (d *Status) GetLocalPeerState() LocalPeerState {
	d.mux.Lock()
	defer d.mux.Unlock()
	return d.localPeer.Clone()
}

// UpdateLocalPeerState updates local peer status
func (d *Status) UpdateLocalPeerState(localPeerState LocalPeerState) {
	d.mux.Lock()
	defer d.mux.Unlock()

	d.localPeer = localPeerState
	d.notifyAddressChanged()
}

// CleanLocalPeerState cleans local peer status
func (d *Status) CleanLocalPeerState() {
	d.mux.Lock()
	defer d.mux.Unlock()

	d.localPeer = LocalPeerState{}
	d.notifyAddressChanged()
}

// MarkManagementDisconnected sets ManagementState to disconnected
func (d *Status) MarkManagementDisconnected(err error) {
	d.mux.Lock()
	defer d.mux.Unlock()
	defer d.onConnectionChanged()

	d.managementState = false
	d.managementError = err
}

// MarkManagementConnected sets ManagementState to connected
func (d *Status) MarkManagementConnected() {
	d.mux.Lock()
	defer d.mux.Unlock()
	defer d.onConnectionChanged()

	d.managementState = true
	d.managementError = nil
}

// UpdateSignalAddress update the address of the signal server
func (d *Status) UpdateSignalAddress(signalURL string) {
	d.mux.Lock()
	defer d.mux.Unlock()
	d.signalAddress = signalURL
}

// UpdateManagementAddress update the address of the management server
func (d *Status) UpdateManagementAddress(mgmAddress string) {
	d.mux.Lock()
	defer d.mux.Unlock()
	d.mgmAddress = mgmAddress
}

// UpdateRosenpass update the Rosenpass configuration
func (d *Status) UpdateRosenpass(rosenpassEnabled, rosenpassPermissive bool) {
	d.mux.Lock()
	defer d.mux.Unlock()
	d.rosenpassPermissive = rosenpassPermissive
	d.rosenpassEnabled = rosenpassEnabled
}

// MarkSignalDisconnected sets SignalState to disconnected
func (d *Status) MarkSignalDisconnected(err error) {
	d.mux.Lock()
	defer d.mux.Unlock()
	defer d.onConnectionChanged()

	d.signalState = false
	d.signalError = err
}

// MarkSignalConnected sets SignalState to connected
func (d *Status) MarkSignalConnected() {
	d.mux.Lock()
	defer d.mux.Unlock()
	defer d.onConnectionChanged()

	d.signalState = true
	d.signalError = nil
}

func (d *Status) UpdateRelayStates(relayResults []relay.ProbeResult) {
	d.mux.Lock()
	defer d.mux.Unlock()
	d.relayStates = relayResults
}

func (d *Status) UpdateDNSStates(dnsStates []NSGroupState) {
	d.mux.Lock()
	defer d.mux.Unlock()
	d.nsGroupStates = dnsStates
}

func (d *Status) UpdateResolvedDomainsStates(originalDomain domain.Domain, resolvedDomain domain.Domain, prefixes []netip.Prefix) {
	d.mux.Lock()
	defer d.mux.Unlock()

	// Store both the original domain pattern and resolved domain
	d.resolvedDomainsStates[resolvedDomain] = ResolvedDomainInfo{
		Prefixes:     prefixes,
		ParentDomain: originalDomain,
	}
}

func (d *Status) DeleteResolvedDomainsStates(domain domain.Domain) {
	d.mux.Lock()
	defer d.mux.Unlock()

	// Remove all entries that have this domain as their parent
	for k, v := range d.resolvedDomainsStates {
		if v.ParentDomain == domain {
			delete(d.resolvedDomainsStates, k)
		}
	}
}

func (d *Status) GetRosenpassState() RosenpassState {
	d.mux.Lock()
	defer d.mux.Unlock()
	return RosenpassState{
		d.rosenpassEnabled,
		d.rosenpassPermissive,
	}
}

func (d *Status) GetManagementState() ManagementState {
	d.mux.Lock()
	defer d.mux.Unlock()
	return ManagementState{
		d.mgmAddress,
		d.managementState,
		d.managementError,
	}
}

func (d *Status) UpdateLatency(pubKey string, latency time.Duration) error {
	if latency <= 0 {
		return nil
	}

	d.mux.Lock()
	defer d.mux.Unlock()
	peerState, ok := d.peers[pubKey]
	if !ok {
		return errors.New("peer doesn't exist")
	}
	peerState.Latency = latency
	d.peers[pubKey] = peerState
	return nil
}

// IsLoginRequired determines if a peer's login has expired.
func (d *Status) IsLoginRequired() bool {
	d.mux.Lock()
	defer d.mux.Unlock()

	// if peer is connected to the management then login is not expired
	if d.managementState {
		return false
	}

	s, ok := gstatus.FromError(d.managementError)
	if ok && (s.Code() == codes.InvalidArgument || s.Code() == codes.PermissionDenied) {
		return true
	}
	return false
}

func (d *Status) GetSignalState() SignalState {
	d.mux.Lock()
	defer d.mux.Unlock()
	return SignalState{
		d.signalAddress,
		d.signalState,
		d.signalError,
	}
}

// GetRelayStates returns the stun/turn/permanent relay states
func (d *Status) GetRelayStates() []relay.ProbeResult {
	d.mux.Lock()
	defer d.mux.Unlock()
	if d.relayMgr == nil {
		return d.relayStates
	}

	// extend the list of stun, turn servers with relay address
	relayStates := slices.Clone(d.relayStates)

	// if the server connection is not established then we will use the general address
	// in case of connection we will use the instance specific address
	instanceAddr, err := d.relayMgr.RelayInstanceAddress()
	if err != nil {
		// TODO add their status
		for _, r := range d.relayMgr.ServerURLs() {
			relayStates = append(relayStates, relay.ProbeResult{
				URI: r,
				Err: err,
			})
		}
		return relayStates
	}

	relayState := relay.ProbeResult{
		URI: instanceAddr,
	}
	return append(relayStates, relayState)
}

func (d *Status) GetDNSStates() []NSGroupState {
	d.mux.Lock()
	defer d.mux.Unlock()

	// shallow copy is good enough, as slices fields are currently not updated
	return slices.Clone(d.nsGroupStates)
}

func (d *Status) GetResolvedDomainsStates() map[domain.Domain]ResolvedDomainInfo {
	d.mux.Lock()
	defer d.mux.Unlock()
	return maps.Clone(d.resolvedDomainsStates)
}

// GetFullStatus gets full status
func (d *Status) GetFullStatus() FullStatus {
	fullStatus := FullStatus{
		ManagementState: d.GetManagementState(),
		SignalState:     d.GetSignalState(),
		Relays:          d.GetRelayStates(),
		RosenpassState:  d.GetRosenpassState(),
		NSGroupStates:   d.GetDNSStates(),
	}

	d.mux.Lock()
	defer d.mux.Unlock()

	fullStatus.LocalPeerState = d.localPeer

	for _, status := range d.peers {
		fullStatus.Peers = append(fullStatus.Peers, status)
	}

	fullStatus.Peers = append(fullStatus.Peers, d.offlinePeers...)
	return fullStatus
}

// ClientStart will notify all listeners about the new service state
func (d *Status) ClientStart() {
	d.notifier.clientStart()
}

// ClientStop will notify all listeners about the new service state
func (d *Status) ClientStop() {
	d.notifier.clientStop()
}

// ClientTeardown will notify all listeners about the service is under teardown
func (d *Status) ClientTeardown() {
	d.notifier.clientTearDown()
}

// SetConnectionListener set a listener to the notifier
func (d *Status) SetConnectionListener(listener Listener) {
	d.notifier.setListener(listener)
}

// RemoveConnectionListener remove the listener from the notifier
func (d *Status) RemoveConnectionListener() {
	d.notifier.removeListener()
}

func (d *Status) onConnectionChanged() {
	d.notifier.updateServerStates(d.managementState, d.signalState)
}

// notifyPeerStateChangeListeners notifies route manager about the change in peer state
func (d *Status) notifyPeerStateChangeListeners(peerID string) {
	ch, found := d.changeNotify[peerID]
	if !found {
		return
	}

	close(ch)
	delete(d.changeNotify, peerID)
}

func (d *Status) notifyPeerListChanged() {
	d.notifier.peerListChanged(d.numOfPeers())
}

func (d *Status) notifyAddressChanged() {
	d.notifier.localAddressChanged(d.localPeer.FQDN, d.localPeer.IP)
}

func (d *Status) numOfPeers() int {
	return len(d.peers) + len(d.offlinePeers)
}

// PublishEvent adds an event to the queue and distributes it to all subscribers
func (d *Status) PublishEvent(
	severity proto.SystemEvent_Severity,
	category proto.SystemEvent_Category,
	msg string,
	userMsg string,
	metadata map[string]string,
) {
	event := &proto.SystemEvent{
		Id:          uuid.New().String(),
		Severity:    severity,
		Category:    category,
		Message:     msg,
		UserMessage: userMsg,
		Metadata:    metadata,
		Timestamp:   timestamppb.Now(),
	}

	d.eventMux.Lock()
	defer d.eventMux.Unlock()

	d.eventQueue.Add(event)

	for _, stream := range d.eventStreams {
		select {
		case stream <- event:
		default:
			log.Debugf("event stream buffer full, skipping event: %v", event)
		}
	}

	log.Debugf("event published: %v", event)
}

// SubscribeToEvents returns a new event subscription
func (d *Status) SubscribeToEvents() *EventSubscription {
	d.eventMux.Lock()
	defer d.eventMux.Unlock()

	id := uuid.New().String()
	stream := make(chan *proto.SystemEvent, 10)
	d.eventStreams[id] = stream

	return &EventSubscription{
		id:     id,
		events: stream,
	}
}

// UnsubscribeFromEvents removes an event subscription
func (d *Status) UnsubscribeFromEvents(sub *EventSubscription) {
	if sub == nil {
		return
	}

	d.eventMux.Lock()
	defer d.eventMux.Unlock()

	if stream, exists := d.eventStreams[sub.id]; exists {
		close(stream)
		delete(d.eventStreams, sub.id)
	}
}

// GetEventHistory returns all events in the queue
func (d *Status) GetEventHistory() []*proto.SystemEvent {
	return d.eventQueue.GetAll()
}

type EventQueue struct {
	maxSize int
	events  []*proto.SystemEvent
	mutex   sync.RWMutex
}

func NewEventQueue(size int) *EventQueue {
	return &EventQueue{
		maxSize: size,
		events:  make([]*proto.SystemEvent, 0, size),
	}
}

func (q *EventQueue) Add(event *proto.SystemEvent) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.events = append(q.events, event)

	if len(q.events) > q.maxSize {
		q.events = q.events[len(q.events)-q.maxSize:]
	}
}

func (q *EventQueue) GetAll() []*proto.SystemEvent {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	return slices.Clone(q.events)
}

type EventSubscription struct {
	id     string
	events chan *proto.SystemEvent
}

func (s *EventSubscription) Events() <-chan *proto.SystemEvent {
	return s.events
}
