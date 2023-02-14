// author: wsfuyibing <websearch@163.com>
// date: 2023-02-02

package conf

type (
	AccountAliyunmnsConfig struct {
		AccessId  string `yaml:"access-id" json:"access-id"`
		AccessKey string `yaml:"access-key" json:"access-key"`
		Endpoint  string `yaml:"endpoint" json:"endpoint"`
		Prefix    string `yaml:"prefix" json:"prefix"`
	}
)

func (o *AccountAliyunmnsConfig) init() *AccountAliyunmnsConfig {
	return o
}

func (o *AccountAliyunmnsConfig) initDefaults() {
}
