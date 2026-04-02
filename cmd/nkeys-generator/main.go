package main

import (
	"fmt"

	"github.com/nats-io/nkeys"
)

func main() {
	kp, err := nkeys.CreateUser()
	if err != nil {
		panic(err)
	}

	pub, _ := kp.PublicKey()
	seed, _ := kp.Seed()

	// // Save private key locally
	// err = os.WriteFile("leaf.nkey", seed, 0600)
	// if err != nil {
	// 	panic(err)
	// }

	fmt.Println("Private key (keep it secret):", string(seed))
	fmt.Println("Public key (send to server):", pub)
}
