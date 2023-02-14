# Edit task basic fields

**Route** : `POST` `/task/edit`<br />
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
// PostEdit
// Edit task basic fields.
//
// Line - 51
// Path - /app/controllers/task/controller.go
func (o *Controller) PostEdit(i iris.Context) interface{} {
}
```

### Request Params

* Annotation : `@Request(app/logics/task.EditRequest)`
* Struct : `<github.com/fuyibing/gmd/app/logics/task.EditRequest>`

  | Field | Type | Required | Condition | Description | Example |
  | ---- | ---- | :----: | ---- | ---- | ---- |
  | concurrency | `int32` | `Y` | gte=0 | Max concurrency, Max consuming message per consumer.<br />Default: 10.<br />Total: Nodes x Parallels * Concurrency.<br />Attention: If this value is set too large, the subscription service will be killed when there are too many messages in the queue (similar to DDOS) | 10 |
  | delay_seconds | `int` | `Y` | gte=0,lte=86400 | Delay seconds, When this configuration is greater than 0, the message sent by the producer needs to wait for the specified seconds before consumption. <br />Unit: Second.<br />Default: 0 (not delay) | 0 |
  | id | `int` | `Y` | gte=1 | Task id | 1 |
  | max_retry | `int` | `Y` | gte=0 | Max consume times, Max consume times if failed returned.<br />Default: 3. | 3 |
  | parallels | `int` | `Y` | gte=0,lte=5 | Max consumers, Start consumers count per node. <br />Default: 1 | 1 |
  | broadcasting | `int` |   |  | Broadcast enabled, When enabled, all consumers of each deployment node will consume.<br />0: Disabled<br />1: Enabled | 0 |
  | remark | `string` |   |  | Task remark | Description about task |
  | title | `string` |   |  | Task name | Example task |

  *Example Code*: 

  ```json
  {
      "broadcasting": 0,
      "concurrency": 10,
      "delay_seconds": 0,
      "id": 1,
      "max_retry": 3,
      "parallels": 1,
      "remark": "Description about task",
      "title": "Example task"
  }
  ```

### Response Params # 1

* Annotation : `@Response(app/logics/task.EditResponse)`
* Struct : `<github.com/fuyibing/gmd/app/logics/task.EditResponse>`

  | Field | Type | Description |
  | ---- | ---- | ---- |
  | affects | `int64` | Updated count |
  | id | `int` | Task id |
  | title | `string` | Task name |

  *Example Code*: 

  ```json
  {
      "data": {
          "affects": 1,
          "id": 1,
          "title": "Example task"
      },
      "dataType": "OBJECT",
      "errno": 0,
      "error": ""
  }
  ```

----

* Updated: `2023-02-02 12:22`