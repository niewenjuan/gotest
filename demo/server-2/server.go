package main

import (
	"github.com/ServiceComb/go-chassis"
	"github.com/cairixian/gotest/basicUtil"
	p "github.com/cairixian/gotest/schema"
	"github.com/ServiceComb/go-chassis/core/server"
	_ "github.com/ServiceComb/go-chassis/config-center"
)

func main(){
	chassis.RegisterSchema("highway",&p.HighwayInvocation{},server.WithSchemaID("hellworld"))
	chassis.RegisterSchema("rest", &p.RestfulInvoke{},server.WithSchemaID("server2Restful"))
	basicUtil.Isok(chassis.Init())
	chassis.Run()
}
