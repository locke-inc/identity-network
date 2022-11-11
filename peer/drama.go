package peer

import "time"

type Drama struct {
	GenesisBlock Block
	Chain        []Block
	Difficulty   int
}

func CreateDrama(difficulty int) Drama {
	genesisBlock := Block{
		Hash:      "0",
		Timestamp: time.Now(),
	}
	return Drama{
		genesisBlock,
		[]Block{genesisBlock},
		difficulty,
	}
}

func InitDrama(difficulty int, t Transaction) Drama {
	genesisBlock := Block{
		Data:      t,
		Hash:      "0",
		Timestamp: time.Now(),
	}
	return Drama{
		genesisBlock,
		[]Block{genesisBlock},
		difficulty,
	}
}

func (b *Drama) addBlock(t Transaction) {
	lastBlock := b.Chain[len(b.Chain)-1]
	newBlock := Block{
		Data:         t,
		PreviousHash: lastBlock.Hash,
		Timestamp:    time.Now(),
	}
	newBlock.mine(b.Difficulty)
	b.Chain = append(b.Chain, newBlock)
}

func (b Drama) isValid() bool {
	for i := range b.Chain[1:] {
		previousBlock := b.Chain[i]
		currentBlock := b.Chain[i+1]
		if currentBlock.Hash != currentBlock.calculateHash() || currentBlock.PreviousHash != previousBlock.Hash {
			return false
		}
	}
	return true
}
