package main

type iGenerator interface {
	Run(opt *option, meta *metaDataService) (err error)
}
