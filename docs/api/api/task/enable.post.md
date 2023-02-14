# Enable task

**Route** : `POST` `/task/enable`<br />
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
// PostEnable
// Enable task.
//
// Line - 96
// Path - /app/controllers/task/controller.go
func (o *Controller) PostEnable(i iris.Context) interface{} {
}
```

### Request Params

* Annotation : `@Request(app/logics/task.EditStatus)`
* Struct : `<github.com/fuyibing/gmd/app/logics/task.EditStatus>`

  | Field | Type | Required | Condition | Description | Example |
  | ---- | ---- | :----: | ---- | ---- | ---- |
  | id | `int` | `Y` | gte=1 | 任务ID | 1 |

  *Example Code*: 

  ```json
  {
      "id": 1
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