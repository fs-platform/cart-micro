package main

import (
	"github.com/fs-platform/cart-micro-service/domain/repository"
	service2 "github.com/fs-platform/cart-micro-service/domain/service"
	"github.com/fs-platform/cart-micro-service/handler"
	common "github.com/fs-platform/go-tool"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	ratelimit "github.com/micro/go-plugins/wrapper/ratelimiter/uber/v2"
	opentracing2 "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/opentracing/opentracing-go"
	"strconv"

	cart "github.com/fs-platform/cart-micro-service/proto/cart"
)

func main() {
	//配置中心
	consulConfig, err := common.GetConsulConfig("127.0.0.1", 8500, "/micro/config")
	if err != nil {
		log.Error(err)
	}
	mysqlInfo := common.GetMysqlFromConsul(consulConfig, "mysql")
	//注册中心
	consulRegistry := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"127.0.0.1:8500",
		}
	})
	//链路追踪
	t, io, err := common.NewTracer("go.micro.service.cart", "localhost:6831")
	if err != nil {
		log.Error(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	//数据库链接
	db, err := gorm.Open("mysql", mysqlInfo.User+":"+mysqlInfo.Pwd+
		"@tcp("+mysqlInfo.Host+":"+strconv.FormatInt(mysqlInfo.Port, 10)+")/"+mysqlInfo.Database+
		"?charset=utf8&parseTime=True&loc=Local")
	cartRepository := repository.NewCartRepository(db)
	cartRepository.InitTable()
	cartService := service2.NewCartDataService(cartRepository)

	// New Service
	service := micro.NewService(
		micro.Name("go.micro.service.cart"),
		micro.Version("latest"),
		//微服务Ip
		micro.Address("127.0.0.1:8400"),
		//添加服务发现,注册中心
		micro.Registry(consulRegistry),
		//添加链路追踪
		micro.WrapHandler(opentracing2.NewHandlerWrapper(opentracing.GlobalTracer())),
		//添加限流,每秒100个请求
		micro.WrapHandler(ratelimit.NewHandlerWrapper(100)),
	)
	// Register Handler
	cart.RegisterCartHandler(service.Server(), &handler.Cart{
		CartDataService: cartService,
	})
	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
