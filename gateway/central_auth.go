package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	gorpc "github.com/libp2p/go-libp2p-gorpc"
	"github.com/locke-inc/identity-network/peer"
)

const (
	CentralAuthProtocolID = "/locke/central-auth"
	CentralAuthEndpoint   = "https://dqe8mcxdxf.execute-api.us-east-1.amazonaws.com/test3/"
)

type CentralAuthService struct {
	Peer *peer.Peer
}

type CentralAuthArgs struct {
	Personame    string `json:"personame"`
	PasswordHash string `json:"password_hash"`
}

type CentralAuthResp struct {
	JWT          string `json:"jwt"`
	SymmetricKey string `json:"symmetric_key"`
	Nonce        string `json:"nonce"`
	Salt         string `json:"salt"`
}

// A simple remote procedure that forwards auth to the Locke API
// This is here because peers probably shouldn't be allowed to
// query the Locke API directly -- everything should go through gateways
func (s *CentralAuthService) CentralAuth(ctx context.Context, args CentralAuthArgs, resp *CentralAuthResp) error {
	fmt.Println("Central authentication args:", args)
	b, _ := json.Marshal(args)
	body := bytes.NewBuffer(b)

	auth, err := http.Post(CentralAuthEndpoint+"authenticate/v2", "application/json", body)
	if err != nil {
		return err
	}

	defer auth.Body.Close()

	authResult, err := ioutil.ReadAll(auth.Body)
	if err != nil {
		return err
	}

	sb := string(authResult)
	log.Printf("Auth result is:", sb)

	err = json.Unmarshal(authResult, resp)
	if err != nil {
		return err
	}

	return nil
}

func (g *Gateway) listenForCentralAuth() {
	rpcHost := gorpc.NewServer(g.Peer.Host, CentralAuthProtocolID)

	svc := CentralAuthService{
		Peer: &g.Peer,
	}
	err := rpcHost.Register(&svc)
	if err != nil {
		panic(err)
	}

	fmt.Println("\nListening for central auth")
}
