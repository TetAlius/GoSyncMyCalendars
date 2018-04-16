package backend

type Worker interface {
	Synchronize(ID string) error
}
