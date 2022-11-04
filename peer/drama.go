package peer

import "time"

type Drama struct {
	genesisBlock Block
	chain        []Block
	difficulty   int
}

func CreateDrama(difficulty int) Drama {
	genesisBlock := Block{
		hash:      "0",
		timestamp: time.Now(),
	}
	return Drama{
		genesisBlock,
		[]Block{genesisBlock},
		difficulty,
	}
}

func (b *Drama) addBlock(from, to string, amount float64) {
	blockData := map[string]interface{}{
		"from":   from,
		"to":     to,
		"amount": amount,
	}
	lastBlock := b.chain[len(b.chain)-1]
	newBlock := Block{
		data:         blockData,
		previousHash: lastBlock.hash,
		timestamp:    time.Now(),
	}
	newBlock.mine(b.difficulty)
	b.chain = append(b.chain, newBlock)
}

func (b Drama) isValid() bool {
	for i := range b.chain[1:] {
		previousBlock := b.chain[i]
		currentBlock := b.chain[i+1]
		if currentBlock.hash != currentBlock.calculateHash() || currentBlock.previousHash != previousBlock.hash {
			return false
		}
	}
	return true
}
