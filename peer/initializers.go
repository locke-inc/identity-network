package peer

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/libp2p/go-libp2p-core/peer"
	gorpc "github.com/libp2p/go-libp2p-gorpc"
)

func InitHandshake(p *Peer, dest peer.ID) {
	// Time the process to record in the drama
	start := time.Now()

	fmt.Println("Initiating handshake...")
	rpcClient := gorpc.NewClient(p.Host, HandshakeProtocolID)

	args1 := StartRelationshipArgs{
		CallingPeerID: p.Host.ID().String(),
	}
	var resp1 StartRelationshipResp
	err := rpcClient.Call(dest, "HandshakeService", "StartRelationship", args1, &resp1)
	if err != nil {
		log.Fatal(err)
	}

	if !resp1.ReadyForAuth {
		log.Fatal("They were not ready for auth, guess they have a commitment problem.")
	}

	// Input OTP
	stdReader := bufio.NewReader(os.Stdin)
	fmt.Print("> ")
	otp, err := stdReader.ReadString('\n')
	if err != nil {
		log.Println(err)
		return
	}

	args2 := AuthorizeRelationshipArgs{
		CallingPeerID: p.Host.ID().String(),
		Them:          p.Self.Person, // You is them to them!
		OTP:           otp,
	}
	var resp2 AuthorizeRelationshipResp
	err = rpcClient.Call(dest, "HandshakeService", "AuthorizeRelationship", args2, &resp2)
	if err != nil {
		log.Fatal(err)
	}

	// Check drama is valid
	// TODO
	if !resp2.TLD.isValid() {
		fmt.Println("Drama is NOT valid, abort")
		panic(errors.New("Drama is invalid"))
	}

	fmt.Println("Handshake was successful:", resp2)

	t := Transaction{
		Requester:   p.Host.ID().String(),
		RequestType: "handshake",
		Responder:   dest.String(),
		Result:      99, // 99 represents a successful OTP auth <---- this is a little cheeky; it's not 100 since we're never 100% sure of anything...
		Application: "handshake settled",
		ProcessTime: time.Since(start),
	}
	resp2.TLD.addBlock(t, resp2.SymKey)

	p.addPerson(&resp2.Them, &resp2.TLD, &resp2.SymKey)

	// Lastly, settle the relationship
	args3 := SettleRelationshipArgs{
		CallingPeerID: p.Host.ID().String(),
		Them:          p.Self.Person, // You is them to them!
		Drama:         resp2.TLD,
	}
	var resp3 SettleRelationshipResp
	err = rpcClient.Call(dest, "HandshakeService", "SettleRelationship", args3, &resp3)
	if err != nil {
		log.Fatal(err)
	}

	if !resp3.Success {
		panic(errors.New("Relationship was not settled for some reason..."))
	}

	fmt.Println("\n********** Relationship is settled! **********")
}
