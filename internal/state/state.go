package state

type State struct {
	BaseURL       string
	CacheInterval int
}

func NewState(baseURL string, cacheInterval int) *State {
	return &State{baseURL, cacheInterval}
}
