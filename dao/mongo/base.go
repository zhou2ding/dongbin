package mongo

type Mongo struct {
}

func NewMongo() *Mongo {
	return &Mongo{}
}

func (m *Mongo) Start() error {
	return nil
}

func (m *Mongo) Stop() error {
	return nil
}
