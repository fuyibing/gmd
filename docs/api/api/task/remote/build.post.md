# Build task remote relations on mq server

**Route** : `POST` `/task/remote/build`<br />
**Deploy** : `http://gmd.{{domain}}:8101`<br />
**Content Type** : `application/json`

```go
// Controller
// Task.
//
// Struct - github.com/fuyibing/gmd/app/controllers/task.Controller
//
// Line - 19
// Path - /app/controllers/task/controller.go
type Controller struct {
}
```

```go
// PostRemoteBuild
// Build task remote relations on mq server.
//
// Line - 105
// Path - /app/controllers/task/controller.go
func (o *Controller) PostRemoteBuild(i iris.Context) interface{} {
}
```

### Request Params

* Annotation : `@Request(app/logics/task.RemoteBuildRequest)`
* Struct : `<github.com/fuyibing/gmd/app/logics/task.RemoteBuildRequest>`

  | Field | Type | Required | Condition | Description | Example |
  | ---- | ---- | :----: | ---- | ---- | ---- |
  | id | `int` | `Y` | gte=1 | Task ID | 1 |

  *Example Code*: 

  ```json
  {
      "id": 1
  }
  ```

### Response Params # 1

* Annotation : `@Response(app/logics/task.RemoteBuildResponse)`
* Struct : `<github.com/fuyibing/gmd/app/logics/task.RemoteBuildResponse>`

  | Field | Type | Description |
  | ---- | ---- | ---- |
  | id | `int` | Task ID |
  | title | `string` | Task name |

  *Example Code*: 

  ```json
  {
      "data": {
          "id": 1,
          "title": "Task name"
      },
      "dataType": "OBJECT",
      "errno": 0,
      "error": ""
  }
  ```

----

* Updated: `2023-02-02 12:22`