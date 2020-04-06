package main

import (
	"fmt"
)

var (
	genMgr = &generatorMgr{
		generatorMap: make(map[string]iGenerator, 4),
	}
)

type generatorMgr struct {
	generatorMap map[string]iGenerator
}

func (g *generatorMgr) Register(name string, gen iGenerator) (err error) {
	_, ok := g.generatorMap[name]
	if ok {
		err = fmt.Errorf("generator %s exists\n", name)
		return
	}
	g.generatorMap[name] = gen
	return
}

func (g *generatorMgr) Run(opt *option) (err error) {
	for _, generator := range g.generatorMap {
		err = generator.Run(opt)
		if err != nil {
			return
		}
	}
	return
}
