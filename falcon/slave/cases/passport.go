package cases

import "fmt"

func (c *CaseEntry) PassPortLogin(){
	resp, err :=c.http.Post(c.httpPostConfig,true)

	c.httpPost(resp, err, fmt.Sprintf("Http(%s)", c.httpPostConfig.Route))

}


func(c *CaseEntry)  CatFullinfo(){

	_,_ = c.http.Get(c.httpGetConfig, true)
}

// PATCH
//func(c *CaseEntry)  PATCHSignature(){
//	c.http.
//}