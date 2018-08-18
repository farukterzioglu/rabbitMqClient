package Utilities

import "fmt"
import consulapi "github.com/hashicorp/consul/api"

type IConsulHelper interface {
	GetValue(key string) string
	GetAllDirectory(prefix string) []string
}

type ConsulHelper struct {
	//connection
	consulUrl string
	consulPath string

	keyValueStore *consulapi.KV
	client *consulapi.Client
}

//PRIVATES
func (helper *ConsulHelper) connectToConsul(url string, consulPath string) error{
	//_ := url + "/" + consulPath

	config := consulapi.DefaultConfig()
	config.Address = "demo.consul.io"
	config.Scheme = "https"

	var err error
	helper.client, err = consulapi.NewClient(config)
	if err != nil {
		panic(err)
	}

	helper.keyValueStore = helper.client.KV()
	return nil
}

//PUBLICS
func NewConsulHelper(consulUrl string, consulPath string) (IConsulHelper, error){
	helper:=  &ConsulHelper{
		consulUrl:consulUrl,
		consulPath:consulPath,
	}

	err := helper.connectToConsul(consulUrl, consulPath)
	if err != nil {
		return  nil, err
	}

	return helper, nil
}

//IConsulHelper implementations
func (helper *ConsulHelper) GetValue(key string) string{
	kvp, _, err := helper.keyValueStore.Get(helper.consulPath + "/" + key, nil)
	if err != nil {
		fmt.Println(err)
	}

	if kvp == nil {
		panic(key + " is null")
	}

	return string(kvp.Value)
}


func (helper *ConsulHelper) GetAllDirectory(prefix string) []string{
	keys, _, err1 := helper.keyValueStore.Keys(prefix + "/", "/", nil)

	if err1 != nil {
		fmt.Println(err1)
	}

	for i:=0;  i< 1; i++  {
		fmt.Printf("%s : %s", string(keys[i]), string(keys[i]))
	}



	//TODO : Implement this
	return []string{"test" ,"test"}
}
