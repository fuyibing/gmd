# Publish one

**Route** : `POST` `/topic/publish`<br />
**Deploy** : `http://gmd.{{domain}}:8101`<br />
**Content Type** : `application/json`

Only 1 message can be published per request<br />
Asynchronous mode

```go
// Controller
// Topic.
//
// Struct - github.com/fuyibing/gmd/app/controllers/topic.Controller
//
// Line - 19
// Path - /app/controllers/topic/controller.go
type Controller struct {
}
```

```go
// PostPublish
// Publish one.
//
// Line - 42
// Path - /app/controllers/topic/controller.go
func (o *Controller) PostPublish(i iris.Context) interface{} {
}
```

### Request Params

* Annotation : `@Request(app/logics/topic.PublishRequest)`
* Struct : `<github.com/fuyibing/gmd/app/logics/topic.PublishRequest>`

  | Field | Type | Required | Condition | Description | Example |
  | ---- | ---- | :----: | ---- | ---- | ---- |
  | topic_name | `string` | `Y` | min=2,max=30 | Topic name |  |
  | topic_tag | `string` | `Y` | min=2,max=60 | Topic tag |  |
  | message | `interface` |   |  | Message content, Accept json string or json object | * |

  *Example Code*: 

  ```json
  {
      "message": "*",
      "topic_name": "",
      "topic_tag": ""
  }
  ```

### Response Params # 1

* Annotation : `@Response(app/logics/topic.PublishResponse)`
* Struct : `<github.com/fuyibing/gmd/app/logics/topic.PublishResponse>`

  | Field | Type | Description |
  | ---- | ---- | ---- |
  | hash | `string` | Message hash |
  | registry_id | `int` | Registry id |

  *Example Code*: 

  ```json
  {
      "data": {
          "hash": "C0837A1B5E264F19826F31457D51546D",
          "registry_id": 1
      },
      "dataType": "OBJECT",
      "errno": 0,
      "error": ""
  }
  ```

----

* Updated: `2023-02-02 12:22`