package gather

type Gatherer interface {
	Set(...interface{})
	DoGather()
	Report() (string, error)
}
