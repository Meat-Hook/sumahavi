package core

// Core responsible for service business logic.
type Core struct {
	store     Store
	index     InvertedIndex
	disk      Disk
	parser    Parser
	tokenizer Tokenizer
	uuid      UUID
	clock     Clock
}

// New build and returns one instance business logic core.
func New(store Store, index InvertedIndex, disk Disk, parser Parser, tokenizer Tokenizer, uuidGenerator UUID, clock Clock) *Core {
	return &Core{
		store:     store,
		index:     index,
		disk:      disk,
		parser:    parser,
		tokenizer: tokenizer,
		uuid:      uuidGenerator,
		clock:     clock,
	}
}
