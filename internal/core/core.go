package core

// Core responsible for service business logic.
type Core struct {
	store     Store
	tokenizer Tokenizer
	uuid      UUID
	clock     Clock
}

// New build and returns one instance business logic core.
func New(store Store, tokenizer Tokenizer, uuidGenerator UUID, clock Clock) *Core {
	return &Core{
		store:     store,
		tokenizer: tokenizer,
		uuid:      uuidGenerator,
		clock:     clock,
	}
}
