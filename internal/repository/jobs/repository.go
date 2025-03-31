package jobs

type repository struct {
	db string // TODO: connect real DB here
}

func NewRepository(db string) *repository {
	return &repository{
		db: db,
	}
}
