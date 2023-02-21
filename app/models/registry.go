// author: wsfuyibing <websearch@163.com>
// date: 2023-02-17

package models

type (
	// Registry
	//
	// topic name and tag pair.
	Registry struct {
		Id        int    `xorm:"id pk autoincr"`
		TopicName string `xorm:"topic_name"`
		TopicTag  string `xorm:"topic_tag"`
		FilterTag string `xorm:"filter_tag"`

		GmtCreated Timeline `xorm:"gmt_created"`
		GmtUpdated Timeline `xorm:"gmt_updated"`
	}
)
