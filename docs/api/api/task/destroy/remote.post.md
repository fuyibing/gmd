# 销毁订阅任务

**路由** : `POST` `/task/destroy/remote`<br />
**部署** : `http://gmd.{{domain}}:8101`<br />
**格式** : `application/json`

```go
// Controller
// 订阅任务.
//
// 类型 - github.com/fuyibing/gmd/app/controllers/task.Controller
//
// 行号 - 17
// 路径 - /app/controllers/task/controller.go
type Controller struct {
}
```

```go
// PostDestroyRemote
// 销毁订阅任务.
//
// 行号 - 41
// 路径 - /app/controllers/task/controller.go
func (o *Controller) PostDestroyRemote(i iris.Context) interface{} {
}
```

### 【请求参数】

* 注解 : `@Request(app/logics/task.DestroyRemoteRequest)`
* 类型 : `<github.com/fuyibing/gmd/app/logics/task.DestroyRemoteRequest>`

  | 字段 | 类型 | 必需 | 条件 | 描述 | 示例 |
  | ---- | ---- | :----: | ---- | ---- | ---- |
  | id | `int` | `Y` | gte=1 | 任务ID | 1 |

  *示例代码*: 

  ```json
  {
      "id": 1
  }
  ```

### 【返回结果】 # 1

* 注解 : `@Response(app/logics/task.DestroyRemoteResponse)`
* 类型 : `<github.com/fuyibing/gmd/app/logics/task.DestroyRemoteResponse>`

  | 字段 | 类型 | 描述 |
  | ---- | ---- | ---- |
  | id | `int` | 任务ID |
  | title | `string` | 任务名称 |

  *示例代码*: 

  ```json
  {
      "data": {
          "id": 1,
          "title": "任务名称"
      },
      "dataType": "OBJECT",
      "errno": 0,
      "error": ""
  }
  ```

----

* 更新: `2023-02-02 21:06`