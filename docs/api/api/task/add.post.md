# Add new task

**Route** : `POST` `/task/add`<br />
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
// PostAdd
// Add new task.
//
// Line - 27
// Path - /app/controllers/task/controller.go
func (o *Controller) PostAdd(i iris.Context) interface{} {
}
```

### Request Params

* Annotation : `@Request(app/logics/task.AddRequest)`
* Struct : `<github.com/fuyibing/gmd/app/logics/task.AddRequest>`

  | Field | Type | Required | Condition | Description | Example |
  | ---- | ---- | :----: | ---- | ---- | ---- |
  | handler | `string` | `Y` | url | Callback address | https://example.com/orders/expired/remove |
  | title | `string` | `Y` | lte=80 | Task name | Example task |
  | topic_name | `string` | `Y` | gte=2,lte=30 | Topic name | orders |
  | topic_tag | `string` | `Y` | gte=2,lte=60 | Topic tag | created |
  | delay_seconds | `int` |   | gte=0,lte=86400 | Delay seconds, When this configuration is greater than 0, the message sent by the producer needs to wait for the specified seconds before consumption. <br />Unit: Second.<br />Default: 0 (not delay) | 0 |
  | remark | `string` |   |  | Description about task | Task remark |

  *Example Code*: 

  ```json
  {
      "delay_seconds": 0,
      "handler": "https://example.com/orders/expired/remove",
      "remark": "Task remark",
      "title": "Example task",
      "topic_name": "orders",
      "topic_tag": "created"
  }
  ```

### Response Params # 1

* Annotation : `@Response(app/logics/task.AddResponse)`
* Struct : `<github.com/fuyibing/gmd/app/logics/task.AddResponse>`

  | Field | Type | Description |
  | ---- | ---- | ---- |
  | delay_seconds | `int` | Delay seconds |
  | id | `int` | Task id |
  | title | `string` | Task name |
  | topic_name | `string` | Topic name |
  | topic_tag | `string` | Topic tag |

  *Example Code*: 

  ```json
  {
      "data": {
          "delay_seconds": 0,
          "id": 1,
          "title": "Example task",
          "topic_name": "orders",
          "topic_tag": "created"
      },
      "dataType": "OBJECT",
      "errno": 0,
      "error": ""
  }
  ```

----

* Updated: `2023-02-02 12:22`