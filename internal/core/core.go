package core

// Core responsible for service business logic.
type Core struct {
	store     Store
	source    Source
	tokenizer Tokenizer
	uuid      UUID
	clock     Clock
}

// New build and returns one instance business logic core.
func New(store Store, source Source, tokenizer Tokenizer, uuidGenerator UUID, clock Clock) *Core {
	return &Core{
		store:     store,
		source:    source,
		tokenizer: tokenizer,
		uuid:      uuidGenerator,
		clock:     clock,
	}
}
