# Publish multiple

**Route** : `POST` `/topic/batch`<br />
**Deploy** : `http://gmd.{{domain}}:8101`<br />
**Content Type** : `application/json`

Each request can publish multiple messages, up to 100<br />
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
// PostBatch
// Publish multiple.
//
// Line - 30
// Path - /app/controllers/topic/controller.go
func (o *Controller) PostBatch(i iris.Context) interface{} {
}
```

### Request Params

* Annotation : `@Request(app/logics/topic.BatchRequest)`
* Struct : `<github.com/fuyibing/gmd/app/logics/topic.BatchRequest>`

  | Field | Type | Required | Condition | Description | Example |
  | ---- | ---- | :----: | ---- | ---- | ---- |
  | topic_name | `string` | `Y` | min=2,max=30 | Topic name |  |
  | topic_tag | `string` | `Y` | min=2,max=60 | Topic tag |  |
  | messages | `[]` `interface` |   |  | Message list, Accept json string or json object in list | * |

  *Example Code*: 

  ```json
  {
      "messages": [
          "*"
      ],
      "topic_name": "",
      "topic_tag": ""
  }
  ```

### Response Params # 1

* Annotation : `@Response(app/logics/topic.BatchResponse)`
* Struct : `<github.com/fuyibing/gmd/app/logics/topic.BatchResponse>`

  | Field | Type | Description |
  | ---- | ---- | ---- |
  | count | `int` | Message count |
  | hash | `string` | Message hash |
  | registry_id | `int` | Registry id |

  *Example Code*: 

  ```json
  {
      "data": {
          "count": 3,
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