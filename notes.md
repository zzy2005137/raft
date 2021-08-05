





数据上链与查询主要通过`service`包的`ServiceSetup`实现，通过添加其他方法实现想要功能

把`ServiceSetup`对象地址赋予web应用，使其也能够调用方法



操作数据库part

web直接读取couchDB



### 任务分解

+ 启动web
+ go操作couchDB
+ web直接读取couchDB

