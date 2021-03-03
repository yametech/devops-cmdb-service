package api

type IApiServer interface {
	Run(<-chan struct{}) error
	Stop() error
}

var _ IApiServer = BaseServer{}

type BaseServer struct {
}

func (b BaseServer) Run(i <-chan struct{}) error {
	for {
		select {
		case <-i:
			return nil
		}
	}
}

func (b BaseServer) Stop() error {
	panic("implement me")
}

var _ IApiServer = fakeServer{}

type fakeServer struct{}

func (f fakeServer) Run(i <-chan struct{}) error {
	panic("implement me")
}

func (f fakeServer) Stop() error {
	panic("implement me")
}
