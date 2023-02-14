# Edit task succeed notification

**Route** : `POST` `/task/edit/succeed`<br />
**Deploy** : `http://gmd.{{domain}}:8101`<br />
**Content Type** : `application/json`

When the message consumption is successful, forward the delivery<br />
result to the successful callback

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
// PostEditSucceed
// Edit task succeed notification.
//
// Line - 87
// Path - /app/controllers/task/controller.go
func (o *Controller) PostEditSucceed(i iris.Context) interface{} {
}
```

### Request Params

* Annotation : `@Request(app/logics/task.EditSubscriber)`
* Struct : `<github.com/fuyibing/gmd/app/logics/task.EditSubscriber>`

  | Field | Type | Required | Condition | Description | Example |
  | ---- | ---- | :----: | ---- | ---- | ---- |
  | id | `int` | `Y` | gte=1 | Task id | 1 |
  | condition | `string` |   |  | Condition filter, Consume when the consumption content meets the filtering conditions, otherwise ignore the message. |  |
  | handler | `string` |   |  | Callback address, Where is the message delivered.<br />Protocol: http, https, tcp, rpc, ws, wss. | http://example.com/path/route?key=value |
  | ignore_codes | `string` |   |  | Ignore logic code, When the code returned by the business party is within the specified range, the consumption is considered successful.<br />Description: multiple codes are separated by commas | 1234,1234 |
  | method | `string` |   |  | Deliver method, Request method when delivering message. <br />Default: POST |  |
  | response_type | `int` |   |  | Response type, How to identify the return results of business parties.<br />0: https status code is 200.<br />1: Return json string and errno field value is zero string or integer. | 0 |
  | timeout | `int` |   |  | Timeout, If response not returned within specified seconds. | 10 |

  *Example Code*: 

  ```json
  {
      "condition": "",
      "handler": "http://example.com/path/route?key=value",
      "id": 1,
      "ignore_codes": "1234,1234",
      "method": "",
      "response_type": 0,
      "timeout": 10
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