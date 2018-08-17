package Utilities

import "fmt"

type IConsulHelper interface {
	GetValue(key string) string
}

type ConsulHelper struct {
	//connection
	consulUrl string
	consulPath string
}

//PRIVATES
func (helper *ConsulHelper) connectToConsul(url string, consulPath string) (int, error){
	//connection == ???
	hostName := url + "/" + consulPath
	fmt.Printf("Consul host : %s \n", hostName)

	//TODO : Implement this
	return 1, nil
}

//PUBLICS
func NewConsulHelper(consulUrl string, consulPath string) (IConsulHelper, error){
	helper:=  &ConsulHelper{
		consulUrl:consulUrl,
		consulPath:consulPath,
	}

	res, err := helper.connectToConsul(consulUrl, consulPath)
	if err != nil {
		return  nil, err
	}

	//TODO : Connect
	res++

	return helper, nil
}

//IConsulHelper implementations
func (helper *ConsulHelper) GetValue(key string) string{
	//TODO : Implement this
	return "test"
}
