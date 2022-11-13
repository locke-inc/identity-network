package peer

// Need to achieve Consensus: https://en.wikipedia.org/wiki/Consensus_(computer_science)

// "Some cryptocurrencies, such as Ripple, use a system of validating nodes to validate the ledger.
// This system used by Ripple, called Ripple Protocol Consensus Algorithm (RPCA), works in rounds:
// Step 1: every server compiles a list of valid candidate transactions; Step 2: each server
// amalgamates all candidates coming from its Unique Nodes List (UNL) and votes on their veracity;
// Step 3: transactions passing the minimum threshold are passed to the next round; Step 4: the
// final round requires 80% agreement[30]"

// Challenges
// 1. Peers go offline and lose ability to communicate
// 2. Peers get out of sync as a result
// 3. Adding new peers and getting them up to speed, process MUST be secure
// 4. Key dissemination
