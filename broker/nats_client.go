package broker

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
)

type NatsClient struct {
	sync.Once
	sync.RWMutex
	// indicate if we're connected
	connected bool
	addrs     string
	conn      *nats.Conn
	subs      []*nats.Subscription
}

func NewBroker(addrs string) *NatsClient {
	subs := make([]*nats.Subscription, 0)
	return &NatsClient{addrs: addrs, subs: subs}
}

func (nc *NatsClient) Connect() error {
	nc.Lock()
	defer nc.Unlock()
	if nc.connected {
		return nil
	}
	log.Printf("connecting to nats hosts at : %s", nc.addrs)
	// Connect Options.
	opts := []nats.Option{nats.Name("NATS Client")}
	opts = setupConnOptions(opts)

	status := nats.CLOSED
	if nc.conn != nil {
		status = nc.conn.Status()
	}
	switch status {
	case nats.CONNECTED, nats.RECONNECTING, nats.CONNECTING:
		nc.connected = true
		return nil
	default: // DISCONNECTED or CLOSED or DRAINING
		c, err := nats.Connect(nc.addrs, opts...)
		if err != nil {
			return err
		}
		nc.conn = c
		nc.connected = true
		return nil
	}
}

func (nc *NatsClient) Disconnect() {
	if !nc.connected {
		return
	}
	nc.conn.Flush()
	nc.conn.Drain()
	for _, s := range nc.subs {
		defer log.Printf("unsubscribed from : %s", s.Subject)
		s.Unsubscribe()
	}
	nc.conn.Close()
}

func setupConnOptions(opts []nats.Option) []nats.Option {
	totalWait := 10 * time.Minute
	reconnectDelay := time.Second

	opts = append(opts, nats.ReconnectWait(reconnectDelay))
	opts = append(opts, nats.MaxReconnects(int(totalWait/reconnectDelay)))
	/*opts = append(opts, nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
		log.Printf("Disconnected due to: %s, will attempt reconnects for %.0fm", err, totalWait.Minutes())
	}))
	opts = append(opts, nats.ReconnectHandler(func(nc *nats.Conn) {
		log.Printf("Reconnected [%s]", nc.ConnectedUrl())
	}))
	opts = append(opts, nats.ClosedHandler(func(nc *nats.Conn) {
		log.Fatalf("Nats connection closed: %v", nc.LastError())
	}))*/
	return opts
}

func (nc *NatsClient) Send(subject string, request protoreflect.ProtoMessage, response protoreflect.ProtoMessage) error {
	if !nc.connected {
		return errors.New("nats_disconnected")
	}
	payload, err := proto.Marshal(request)
	if err != nil {
		return err
	}
	responsePayLoad, err := nc.conn.Request(subject, payload, time.Second)
	if err != nil {
		return err
	}
	if err := proto.Unmarshal(responsePayLoad.Data, response); err != nil {
		return err
	}
	return nil
}

func (nc *NatsClient) Subscribe(subject string, groupName string, h func([]byte) (protoreflect.ProtoMessage, error)) error {

	sub, err := nc.conn.QueueSubscribe(subject, groupName, func(m *nats.Msg) {
		go func() {
			response, err := h(m.Data)
			if err == nil {
				payload, err := proto.Marshal(response)
				if err == nil {
					if err = m.Respond(payload); err != nil {
						// Log the error, not much can be done here
						log.Printf("Eror in sending response %v", err)
					}
				} else {
					// Log the error, not much can be done here
					log.Printf("Eror in marshaling response %v", err)
				}
			} else {
				// Log the error, not much can be done here
				log.Printf("Eror in handling request %v", err)

			}
		}()
	})
	if err != nil {
		log.Printf("Unable to subscribe to subject %s - %v", subject, err)
		return err
	} else {
		nc.subs = append(nc.subs, sub)
	}
	return nil
}
