package models

import (
	"github.com/mediocregopher/radix.v2/redis"
	"fmt"
	"errors"
	"goredisadmin/utils"
	"strconv"
)

func NewRedis(host string,port int,passwd string) (client *redis.Client,err error,conn,auth,ping bool)  {
	portstr:=strconv.Itoa(port)
	client=RedisMap[host+portstr]
	if client!=nil{
		result,err:=client.Ping()
		if result=="PONG"{
			return client,err,true,true,true
		}
	}

	client, err = redis.Dial("tcp", fmt.Sprintf("%v:%v",host,port))
	if err!=nil{
		utils.Logger.Printf("redis %v:%v 连接失败！！！",host,port)
		return client,err,conn,auth,ping
	}
	conn=true
	if passwd!=""{
		result,_:=client.Auth(passwd)
		if result!="OK"{
			utils.Logger.Printf("redis %v:%v 认证失败！！！",host,port)
			return client,err,conn,auth,ping
		}
	}
	auth=true
	result,err:=client.Ping()
	if result!="PONG"{
		utils.Logger.Printf("redis %v:%v ping失败！！！",host,port)
		if passwd==""{
			auth=false
		}
		return client,err,conn,auth,ping
	}
	ping=true
	RedisMap[host+portstr]=client
	return client,err,conn,auth,ping
}

var Redis,_,_,_,_=NewRedis(utils.Rc.Host,utils.Rc.Port,utils.Rc.Passwd)

var RedisMap=map[string]*redis.Client{}

func CheckredisResult(result string,err error) (error) {
	if err!=nil{
		return err
	}
	if result!="OK"{
		return errors.New(result)
	}
	return nil
}