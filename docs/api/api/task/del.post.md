# Delete task

**Route** : `POST` `/task/del`<br />
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
// PostDel
// Delete task.
//
// Line - 33
// Path - /app/controllers/task/controller.go
func (o *Controller) PostDel(i iris.Context) interface{} {
}
```

### Request Params

### Response Params

----

* Updated: `2023-02-02 12:22`