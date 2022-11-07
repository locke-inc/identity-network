package peer_test

import (
	"reflect"
	"testing"

	"github.com/boltdb/bolt"
	"github.com/locke-inc/identity-network/peer"
)

func Test_initPeerStore(t *testing.T) {
	tests := []struct {
		name string
		want *bolt.DB
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := peer.InitPeerStore(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("initPeerStore() = %v, want %v", got, tt.want)
			}
		})
	}
}
