# golang pkg

## 说明
- 下载到工程目录或软连接按需加载

## 模块
| 模块        | 备注          
| ----      | ----------------- 
| ctx       | gin统一响应接口        
| excel     | struct字段反射生成sheet
| log       | zeorlog接口封装 
| mqt       | 适配amqp、kafka、mqtt、stomp、nats，提供通用接口
| orm       | 基于gorm实现分表、分区、where分表等
| router    | [模块化加载路由](https://gitee.com/dontls/ginfast/blob/master/api/system/main.go)
| tree      | 通用Slice转换tree 