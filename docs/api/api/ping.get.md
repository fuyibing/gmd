# Health check

**Route** : `GET` `/ping`<br />
**Deploy** : `http://gmd.{{domain}}:8101`<br />
**Content Type** : `application/json`

```go
// Controller
// Default.
//
// Struct - github.com/fuyibing/gmd/app/controllers.Controller
//
// Line - 16
// Path - /app/controllers/controller.go
type Controller struct {
}
```

```go
// GetPing
// Health check.
//
// Line - 31
// Path - /app/controllers/controller.go
func (o *Controller) GetPing(i iris.Context) interface{} {
}
```

### Request Params

### Response Params # 1

* Annotation : `@Response(app/logics/index.PingResponse)`
* Struct : `<github.com/fuyibing/gmd/app/logics/index.PingResponse>`

  | Field | Type | Description |
  | ---- | ---- | ---- |
  | cpu | `int` | CPU core count |
  | goroutines | `int` | Coroutine counts |
  | memory | `float64` | Used system memory, Unit: MB |
  | pid | `int` | Process ID |
  | start_time | `string` | Started time |

  *Example Code*: 

  ```json
  {
      "data": {
          "cpu": 8,
          "goroutines": 32,
          "memory": 16.57,
          "pid": 3721,
          "start_time": "2022-01-19 14:21:25"
      },
      "dataType": "OBJECT",
      "errno": 0,
      "error": ""
  }
  ```

----

* Updated: `2023-02-02 12:22`