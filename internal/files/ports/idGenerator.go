package ports

type IdGenerator interface {
	MakeId() string
}
