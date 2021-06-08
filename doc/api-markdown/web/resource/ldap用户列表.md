## 分组名称
> 资源管理

## 接口名称
> ldap用户列表

## 接口描述
> ldap用户列表

## 接口版本

> 1.0.0

##### HTTP请求方式

> GET

##### 接口路径
> /cmdb/web/ldap/user-list

##### 请求头参数
> | 参数       | 必填 | 描述            |
> | ---------- | :--- |  --------------- |
> | Content-Type |true|请求体编码类型,使用:application/json;charset=UTF-8|

##### 请求头样例
```
 {
    Content-Type:"application/json;charset=UTF-8"
    v:"1.7.4"
    sign:"38A4907D0AC57CBFB715BB5E69896B18"
    platform:"1"
    timestamp:"1561030834"
    token:"01fbec36c0b617cbaea00a89ccc08380"
}
```

##### Query参数
> | 参数       | 必填 | 描述            |
> | ---------- | :--- |  --------------- |



##### Query参数样例
```

```

##### 响应头参数
> | 参数       | 必填 | 描述            |
> | ---------- | :--- |  --------------- |

##### 响应体参数
> | 参数       | 必选 | 类型 | 说明            |
> | ---------- | :--- | :--- | --------------- |
> | code |true|Integer|状态码|
> | data |true|map|返回数据|
> | data.totalCount |true|int|总数|
> | data.list |true|list|用户列表|
> | data.list.uid |true|String|用户uid|
> | data.list.name |true|String|用户名称|
> | message |true|String| |


##### 响应体样例
```
{
    "data": {
        "totalCount": 1,
        "list": [
            {
                "uid": "xx",
                "name": "试试"
            }
        ]
    },
    "code": 200,
    "msg": ""
}
```
##### 错误码
> | 错误码      |错误描述|
> | :----------: | :---------------: |
> | 400 |通用错误提示,多用于toast弹窗|