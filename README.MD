# 流程
 
 接受停车云推送的数据

    入场数据
        
        保留orderid picUrl 到db中

        上传数据到第三方云

    出场数据

        上传数据到第三方云

        上传photourl 到第三方云
        从db 中检索到 picurl ， 一起上传到第三方云


# DB 结构

orderid, picUrl

## 创建表

create table `PICInfo` (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  `orderid` varchar(255) NOT NULL,
  `picUrl` varchar(255) NOT NULL,
  PRIMARY KEY (`orderid`)
)

## 插入数据

insert into `PICInfo` (`orderid`, `picUrl`) values ('123', 'http://www.baidu.com')

## 查询数据

select * from `PICInfo` where orderid = '123'