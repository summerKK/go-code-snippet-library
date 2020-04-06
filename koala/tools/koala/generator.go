package main

type iGenerator interface {
	Run(opt *option) (err error)
}
