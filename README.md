## 使用migrate

本项目使用golang-migrate来进行数据库迁移
```shell
# 生成迁移文件
migrate create -ext sql -dir db/migrations -seq init_schema
# -seq 生成有序的迁移文件类似于0001、0002
# -ext 生成的文件后缀
```
