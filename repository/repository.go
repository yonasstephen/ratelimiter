package repository

// Repository interfaces the interaction with the underlying
// store where the rate limit data is persisted
type Repository interface {
	IncrementByKey(key string)
}
