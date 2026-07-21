package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

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

func main() {
	configFiles := []string{"./cmd/nats-remote-1/env/.env", "./cmd/nats-remote-1/env/.env.local"}
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
	
	_ = js

	// Subscribe to a subject
	_, err = nc.Subscribe("to.leaf1.hello", func(msg *nats.Msg) {
		fmt.Println("received:", string(msg.Data))
	})
	if err != nil {
		log.Fatal("to.leaf1.hello error: ", err)
	}

	select {} // block forever
}