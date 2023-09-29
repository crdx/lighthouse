package prober

// https://github.com/liamg/furious

import (
	"log/slog"

	"crdx.org/lighthouse/service"
)

type Prober struct {
	log *slog.Logger
}

func New() *Prober {
	return &Prober{}
}

func (*Prober) Config() service.Config {
	return config{}
}

func (self *Prober) Init(args *service.Args) error {
	self.log = args.Logger
	return nil
}

func (*Prober) Run() error {
	// gomod get github.com/liamg/furious
	//
	// target := scan.NewTargetIterator("192.168.1.0/24")
	// scanner := scan.NewSynScanner(target, time.Second, 100)
	// lo.Must0(scanner.Start())
	// results := lo.Must(scanner.Scan(context.Background(), []int{80, 81}))
	// for _, result := range results {
	// 	// fmt.Printf("%#v\n", result)
	// 	if result.IsHostUp() {
	// 		scanner.OutputResult(result)
	// 	}
	// }
	return nil
}
