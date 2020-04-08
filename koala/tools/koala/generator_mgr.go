package main

import (
	"fmt"
	"github.com/emicklei/proto"
	"github.com/summerKK/go-code-snippet-library/koala/logger"
	"os"
	"path"
)

var (
	genMgr = &generatorMgr{
		generatorMap: make(map[string]iGenerator, 4),
		metaData:     &metaDataService{},
	}
)

type generatorMgr struct {
	generatorMap map[string]iGenerator
	metaData     *metaDataService
}

type metaDataService struct {
	Service      *proto.Service
	Message      []*proto.Message
	Rpc          []*proto.RPC
	Pkg          string
	ProtoFileDir string
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

func (g *generatorMgr) parseService(opt *option) (err error) {

	reader, err := os.Open(opt.Proto3FileName)
	if err != nil {
		logger.Logger.Infof("generator mgr generator open file %s failed:%v", opt.Proto3FileName, err)
		return
	}
	defer reader.Close()

	parser := proto.NewParser(reader)
	definition, _ := parser.Parse()

	proto.Walk(definition,
		proto.WithService(g.handleService),
		proto.WithMessage(g.handleMessage),
		proto.WithRPC(g.handleRpc),
		proto.WithPackage(g.handlePackage),
	)

	g.metaData.ProtoFileDir = path.Dir(opt.Proto3FileName)

	return
}

func (g *generatorMgr) handleService(s *proto.Service) {
	g.metaData.Service = s
}

func (g *generatorMgr) handleMessage(s *proto.Message) {
	g.metaData.Message = append(g.metaData.Message, s)
}

func (g *generatorMgr) handleRpc(s *proto.RPC) {
	g.metaData.Rpc = append(g.metaData.Rpc, s)
}

func (g *generatorMgr) handlePackage(s *proto.Package) {
	g.metaData.Pkg = s.Name
}

func (g *generatorMgr) Run(opt *option) (err error) {
	// 创建目录
	err = CreateDir(opt)
	if err != nil {
		logger.Logger.Infof("generator mgr [Run] create dir failed:%v", err)
		return
	}

	err = g.parseService(opt)
	if err != nil {
		logger.Logger.Infof("generator mgr [Run] parse meta data failed:%v", err)
		return
	}
	for _, generator := range g.generatorMap {
		err = generator.Run(opt, g.metaData)
		if err != nil {
			return
		}
	}

	return
}
