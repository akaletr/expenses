package app

type expenses struct {
}

func New() App {
	return &expenses{}
}

func (expenses *expenses) Init() error {
	return nil
}

func (expenses *expenses) Run() error {
	return nil
}

func (expenses *expenses) Stop() error {
	return nil
}
