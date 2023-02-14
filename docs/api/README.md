# GMD

**Host** : `http://gmd.{{domain}}:8101`<br />**Updated** : `2023-02-02 12:22`

> MQ Dispatcher by golang

### Table Contents

* *Default* <small>(1)</small>
  * [Health check](./api/ping.get.md) - <small>`GET`</small> <small>`/ping`</small>
* *Task* <small>(10)</small>
  * [Add new task](./api/task/add.post.md) - <small>`POST`</small> <small>`/task/add`</small>
  * [Delete task](./api/task/del.post.md) - <small>`POST`</small> <small>`/task/del`</small>
  * [Disable task](./api/task/disable.post.md) - <small>`POST`</small> <small>`/task/disable`</small>
  * [Edit task failed notification](./api/task/edit/failed.post.md) - <small>`POST`</small> <small>`/task/edit/failed`</small>
  * [Edit task subscriber](./api/task/edit/handler.post.md) - <small>`POST`</small> <small>`/task/edit/handler`</small>
  * [Edit task succeed notification](./api/task/edit/succeed.post.md) - <small>`POST`</small> <small>`/task/edit/succeed`</small>
  * [Edit task basic fields](./api/task/edit.post.md) - <small>`POST`</small> <small>`/task/edit`</small>
  * [Enable task](./api/task/enable.post.md) - <small>`POST`</small> <small>`/task/enable`</small>
  * [Build task remote relations on mq server](./api/task/remote/build.post.md) - <small>`POST`</small> <small>`/task/remote/build`</small>
  * [Destroy task remote relations of mq server](./api/task/remote/destroy.post.md) - <small>`POST`</small> <small>`/task/remote/destroy`</small>
* *Topic* <small>(2)</small>
  * [Publish multiple](./api/topic/batch.post.md) - <small>`POST`</small> <small>`/topic/batch`</small>
  * [Publish one](./api/topic/publish.post.md) - <small>`POST`</small> <small>`/topic/publish`</small>

----

* `/go.md` - match module name
* `/gdoc.json` - match application configurations
* `/app/controllers` - controller files location
* `/docs/api` - document storage location