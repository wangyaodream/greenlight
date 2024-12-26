## 使用migrate

本项目使用golang-migrate来进行数据库迁移
```shell
# 生成迁移文件
migrate create -ext sql -dir db/migrations -seq init_schema
# -seq 生成有序的迁移文件类似于0001、0002
# -ext 生成的文件后缀
```
同样可以使用`migrate`命令来让数据库迁移
```shell
# 迁移数据库
migrate -path=./migrations -database=$GREENLIGHT_DB_DSN up
```


### 结构解释

#### cmd/api/
该目录下存放的是项目的入口文件，主要是项目的启动文件以及路由的配置文件，所有的请求和响应由该目录的程序发起和处理。

#### internal/data/
该目录存放的是项目的数据层，主要是数据库的操作，包括数据库的连接、数据库的操作等。

#### internal/validator
项目验证器，用于方便验证一切数据的合法性
