package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/n4djib/report-engine/pkg/config"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nkeys"
)

func ConnectToNATS(cfg ConfigVars) (*nats.Conn, nats.JetStreamContext, error) {
	// create key pair from seed
	kp, err := nkeys.FromSeed([]byte(cfg.NatsNkeySeed))
	if err != nil {
		// log.Fatal(err)
		return nil, nil, err
	}

	// extract public key
	pub, err := kp.PublicKey()
	if err != nil {
		// log.Fatal(err)
		return nil, nil, err
	}

	// provide signature function
	sigCB := func(nonce []byte) ([]byte, error) {
		return kp.Sign(nonce)
	}

	// connect to NATS
	nc, err := nats.Connect(
		cfg.NatsURL,
		nats.Nkey(pub, sigCB),
	)
	if err != nil {
		// log.Fatal(err)
		return nil, nil, err
	}

	js, err := nc.JetStream()
	if err != nil {
		// log.Fatal(err)
		return nil, nil, err
	}

	return nc, js, nil
}

func ensureStreams(js nats.JetStreamContext) error {
	// COMMANDS stream
	_, err := js.StreamInfo("COMMANDS")
	if err != nil {
		_, err = js.AddStream(&nats.StreamConfig{
			Name:     "COMMANDS",
			Subjects: []string{"to.*.*"},
			Storage:  nats.FileStorage,
		})
		if err != nil {
			return err
		}
		log.Println("Created stream: COMMANDS")
	}

	// AGGREGATED stream
	_, err = js.StreamInfo("AGGREGATED")
	if err != nil {
		_, err = js.AddStream(&nats.StreamConfig{
			Name:     "AGGREGATED",
			Subjects: []string{"from.*.*"},
			Storage:  nats.FileStorage,
		})
		if err != nil {
			return err
		}
		log.Println("Created stream: AGGREGATED")
	}

	return nil
}

func main() {
	configFiles := []string{"./cmd/nats-central/env/.env", "./cmd/nats-central/env/.env.local"}
	cfg := ConfigVars{}
	err := config.LoadConfigFromFiles(&cfg, configFiles)
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	nc, js, err := ConnectToNATS(cfg)
	if err != nil {
		log.Fatal("Failed to connect to NATS: ", err)
	}
	defer nc.Close()

	// Ensure streams exist
	if err := ensureStreams(js); err != nil {
		log.Fatal("stream setup failed:", err)
	}

	// Publish with timeout (JetStream)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
 

	// Publish a message
	_, err = js.PublishMsg(&nats.Msg{
		Subject: "to.leaf1.hello",
		Data:    []byte("Hello from central to leaf1"),
	}, nats.Context(ctx))
	if err != nil {
		log.Println("publish to.leaf1.hello error:", err)
	} else {
		log.Println("published: to.leaf1.hello")
	}

	_, err = js.PublishMsg(&nats.Msg{
		Subject: "to.leaf2.hello",
		Data:    []byte("Hello from central to leaf2"),
	}, nats.Context(ctx))
	if err != nil {
		log.Println("publish to.leaf2.hello error:", err)
	} else {
		log.Println("published: to.leaf2.hello")
	}

	// Subscribe to a subject
	_, err = js.Subscribe("from.leaf2.>", func(msg *nats.Msg) {
			fmt.Println("received:", string(msg.Data))

			// ACK is mandatory for reliability
			msg.Ack()
		},
		nats.Durable("central-leaf2"),
		nats.ManualAck(),
	)
	if err != nil {
		log.Fatal("from.leaf2.hello error:", err)
	}


	log.Println("Listening...")

	select {} // block forever
}