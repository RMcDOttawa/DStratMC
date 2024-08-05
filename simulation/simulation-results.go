package simulation

type SimResults interface{}

type SimResultsInstance struct {
}

func NewSimResults() SimResults {
	results := &SimResultsInstance{}
	return results
}
