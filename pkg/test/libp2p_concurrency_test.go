//go:build unit || !integration

package test

import (
	"context"
	"errors"
	"fmt"
	"github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/phayes/freeport"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
	"time"
)

const testTopic = "foo"

func TestNodesDiscoverEachOtherOverPubSub(t *testing.T) {
	for i := 0; i < 100; i++ {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			runTest(t)
		})
	}
}

func runTest(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctx = network.WithDialPeerTimeout(ctx, 1*time.Second)

	host1 := newHost(t)
	defer host1.Close()

	host1PubSub := newPubSub(ctx, t, host1)
	host1Rec, host1Send, topic1Closer := startTopic(ctx, t, host1PubSub)
	defer topic1Closer()
	go func() {
		for {
			time.Sleep(100 * time.Millisecond)
			select {
			case <-ctx.Done():
				return
			default:
				host1Send(ctx, host1.ID().String())
			}
		}
	}()

	host2 := newHost(t)
	defer host2.Close()
	host2PubSub := newPubSub(ctx, t, host2)
	host2Rec, host2Send, topic2Closer := startTopic(ctx, t, host2PubSub)
	defer topic2Closer()

	go func() {
		for {
			time.Sleep(100 * time.Millisecond)
			select {
			case <-ctx.Done():
				return
			default:
				host2Send(ctx, host2.ID().String())
			}
		}
	}()

	host2.Peerstore().AddAddrs(host1.ID(), host1.Addrs(), peerstore.PermanentAddrTTL)
	require.NoError(t, host2.Connect(ctx, peer.AddrInfo{
		ID:    host1.ID(),
		Addrs: host1.Addrs(),
	}))

	require.Eventually(t, func() bool {
		for {
			select {
			case s := <-host1Rec:
				if host2.ID().String() == s {
					return true
				}
			default:
				return false
			}
		}
	}, 10*time.Second, 100*time.Millisecond)

	require.Eventually(t, func() bool {
		for {
			select {
			case s := <-host2Rec:
				if host1.ID().String() == s {
					return true
				}
			default:
				return false
			}
		}
	}, 10*time.Second, 100*time.Millisecond)
}

func newHost(t *testing.T) host.Host {
	port, err := freeport.GetFreePort()
	require.NoError(t, err)

	addrs := []string{
		fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port),
		fmt.Sprintf("/ip4/0.0.0.0/udp/%d/quic", port),
		fmt.Sprintf("/ip4/0.0.0.0/udp/%d/quic-v1", port),
		fmt.Sprintf("/ip6/::/tcp/%d", port),
		fmt.Sprintf("/ip6/::/udp/%d/quic", port),
		fmt.Sprintf("/ip6/::/udp/%d/quic-v1", port),
	}

	h, err := libp2p.New(libp2p.ListenAddrStrings(addrs...), libp2p.RandomIdentity)
	require.NoError(t, err)

	t.Log("started libp2p host", h.ID().String())

	return h
}

func newPubSub(ctx context.Context, t *testing.T, h host.Host) *pubsub.PubSub {
	pgParams := pubsub.NewPeerGaterParams(
		0.33, //nolint:gomnd
		pubsub.ScoreParameterDecay(2*time.Minute),
		pubsub.ScoreParameterDecay(10*time.Minute),
	)

	sub, err := pubsub.NewGossipSub(
		ctx,
		h,
		pubsub.WithPeerExchange(true),
		pubsub.WithPeerGater(pgParams),
		pubsub.WithRawTracer(tracer{t}),
	)
	require.NoError(t, err)

	return sub
}

func startTopic(ctx context.Context, t *testing.T, pubSub *pubsub.PubSub) (<-chan string, func(context.Context, string), func() error) {
	topic, err := pubSub.Join(testTopic)
	require.NoError(t, err)

	sub, err := topic.Subscribe()
	require.NoError(t, err)

	messages := make(chan string, 100)

	go func() {
		defer close(messages)
		for {
			msg, err := sub.Next(ctx)
			if errors.Is(err, pubsub.ErrSubscriptionCancelled) {
				return
			}
			require.NoError(t, err)

			messages <- string(msg.Data)
		}
	}()

	f := func(ctx context.Context, s string) {
		err := topic.Publish(ctx, []byte(s))
		if errors.Is(err, pubsub.ErrTopicClosed) {
			return
		}
		require.NoError(t, err)
	}

	return messages, f, func() error {
		sub.Cancel()
		return topic.Close()
	}
}

var _ pubsub.RawTracer = tracer{}

type tracer struct {
	t *testing.T
}

func (t tracer) AddPeer(p peer.ID, proto protocol.ID) {
	t.t.Log("AddPeer", p, proto)
}

func (t tracer) RemovePeer(p peer.ID) {
	t.t.Log("RemovePeer", p)
}

func (t tracer) Join(topic string) {
	t.t.Log("Join", topic)
}

func (t tracer) Leave(topic string) {
	t.t.Log("Leave", topic)
}

func (t tracer) Graft(p peer.ID, topic string) {
	t.t.Log("Graft", p, topic)
}

func (t tracer) Prune(p peer.ID, topic string) {
	t.t.Log("Prune", p, topic)
}

func (t tracer) ValidateMessage(msg *pubsub.Message) {
	t.t.Log("ValidateMessage", msg.ID)
}

func (t tracer) DeliverMessage(msg *pubsub.Message) {
	t.t.Log("DeliverMessage", msg.ID)
}

func (t tracer) RejectMessage(msg *pubsub.Message, reason string) {
	t.t.Log("RejectMessage", msg.ID, reason)
}

func (t tracer) DuplicateMessage(msg *pubsub.Message) {
	t.t.Log("DuplicateMessage", msg.ID)
}

func (t tracer) ThrottlePeer(p peer.ID) {
	t.t.Log("ThrottlePeer", p)
}

func (t tracer) RecvRPC(rpc *pubsub.RPC) {
	t.t.Log("RecvRPC", rpc)
}

func (t tracer) SendRPC(rpc *pubsub.RPC, p peer.ID) {
	t.t.Log("SendRPC", rpc, p)
}

func (t tracer) DropRPC(rpc *pubsub.RPC, p peer.ID) {
	t.t.Log("DropRPC", rpc, p)
}

func (t tracer) UndeliverableMessage(msg *pubsub.Message) {
	t.t.Log("UndeliverableMessage", msg.ID)
}
