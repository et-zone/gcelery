### github.com/et-zone/gcelery

### What about Gcelery
Gcelery, 功能和python的celery框架类似的异步任务框架，功能可参考celery。
+ 主要功能
     - 自定义worker路径，用于gcelery执行异步程序
     - client 集成到项目中
     - 支持tls 模式
     - 支持定时任务
     - 支持异步任务（用于持续执行相关业务）
     - 支持task任务，可以获取返回数据，（可以异步获取响应，需要自己起Goroutine）

#### 性能（8核16G 3.2GHz 服务器）
##### 极限测试
* Server支持10万级 op/s 
* Client支持30w op/min（gin是34.8w op/min,grpc是48-54w op/min），多client可满足
* 前Gcelery 可处理600w op/min,python celery也是说百万级 op/min  具体没测过  (说明，op是以最低响应时间计算的，暂不考虑业务本身的影响)
* 后面版本可能会提高client的处理op性能，根据时间而定

##### 传统测试（CPU跑满）
* 约1.4w op/s =84w op/min
### install
	`go get -u github.com/et-zone/gcelery`
 
