package boomer

import (
	"fmt"
	"math"
	"reflect"
	"runtime"
	"strings"
	"time"
)

var CaseEntry func()interface{}

type taskEntry struct {
	index		reflect.Value
	num			reflect.Value
	times		reflect.Value
	log			reflect.Value
	config		reflect.Value
	param		reflect.Value
	isRunning	reflect.Value
}

type caseInfo struct {
	name  string
	first int64
	last  int64
}

type userInfo struct {
	index 			int64
	num   			int64
	entry 			interface{}
	taskEntry 		*taskEntry
}

type errMsg struct {
	error string
}

func task() {
	_ = <- control.createChannel
	userInfo := <- control.userChannel
	_ = <- control.createChannel
	runTask(userInfo)
	stopCase(info.caseFunc[userInfo.index])
	stopWorker()
}

func runTask(userInfo *userInfo) {
	defer func() {
		if err := recover(); err != nil {
			control.log.error("run task fail:%s", err)
		}
	}()

	taskEntry := userInfo.taskEntry
	if taskEntry == nil {
		control.log.error("can not load WorkerEntry.")
		return
	}

	caseFunc := info.caseFunc[userInfo.index]
	weightCount := info.weightCount
	caseConfig := info.caseConfig

	setEntryInterface(taskEntry.isRunning, &control.isRunning)
	setEntryInt(taskEntry.index, userInfo.index)
	setEntryInt(taskEntry.num, userInfo.num)
	setEntryInterface(taskEntry.config, Conf.BaseConfig)
	setEntryInterface(taskEntry.param, Param)
	setEntryInterface(taskEntry.log, &userLog{})
	control.startListChannel <- 0
	if !Param.RandomMode {
		select{
		case <- control.stopChannel:
		case <- control.startChannel:
		}
	}
	if checkEntryFunc(caseFunc, "OnStart") {
		errMsg := &errMsg{}
		callEntryFunc(caseFunc["OnStart"], errMsg)
		if errMsg.error != "" {
			control.log.error("Call OnStart fail:%s", errMsg.error)
		}
	} else {
		control.log.error("OnStart not exist or is nil")
	}

	if !Param.RandomMode {
		control.taskListChannel <- 0
		select{
			case <- control.stopChannel:
			case <- control.taskOverChannel:
		}
	}

	times := int64(0)
	for {
		if control.isRunning == false {
			break
		}
		if (times >= Param.RunTimes) && (Param.RunTimes > 0) && (Param.RandomMode == false) {
			break
		}
		if (Param.Duration > 0) && ((MilliNow()-startTime) >= Param.Duration) {
			break
		}
		if times < math.MaxInt64 {
			times++
		}

		setEntryInt(taskEntry.times, times)
		if !Param.RandomMode {
			for i := range caseConfig {
				name := caseConfig[i].name
				if checkEntryFunc(caseFunc, name) {
					errMsg := &errMsg{}
					callEntryFunc(caseFunc[name], errMsg)
					if errMsg.error != "" {
						control.log.error("Call %s fail:%s", name, errMsg.error)
					}
				} else {
					control.log.error("%s not exist or is nil", name)
				}
				control.taskListChannel <- 0
				select{
					case <- control.stopChannel:
						RandomWait(time.Duration(100)*time.Millisecond, time.Duration(1000)*time.Millisecond)
						break
					case <- control.taskOverChannel:
				}
				RandomWait(time.Duration(Param.MinWait)*time.Millisecond, time.Duration(Param.MaxWait)*time.Millisecond, times)
			}
		} else {
			name := randomGetEntryFunc(caseConfig, weightCount, times)
			if checkEntryFunc(caseFunc, name) {
				errMsg := &errMsg{}
				callEntryFunc(caseFunc[name], errMsg)
				if errMsg.error != "" {
					control.log.error("Call %s fail:%s", name, errMsg.error)
				}
			} else {
				control.log.error("%s not exist or is nil", name)
			}
			RandomWait(time.Duration(Param.MinWait)*time.Millisecond, time.Duration(Param.MaxWait)*time.Millisecond, times)
		}
	}
	stopCase(caseFunc)
	stopWorker()
}

func stopCase(caseFunc map[string]reflect.Value) {
	defer func() {
		if err := recover(); err != nil {
			control.log.error("StopTask fail:%s", err)
		}
	}()
	if checkEntryFunc(caseFunc,"OnStop") {
		errMsg := &errMsg{}
		callEntryFunc(caseFunc["OnStop"], errMsg)
		if errMsg.error != "" {
			control.log.error("Call %s fail:%s", "OnStop", errMsg.error)
		}
	}
	return
}

func stopWorker() {
	defer func() {
		if err := recover(); err != nil {
			control.log.error("StopWorker fail:%s", err)
		}
	}()
	control.stopListChannel <- 0
	runtime.Goexit()
}

func loadCaseWeight() int64 {
	caseConfig := Conf.BaseConfig.Cases
	weightCount := int64(0)
	for i := range caseConfig {
		fnType := caseConfig[i].Type
		enable := caseConfig[i].Enable
		weight := caseConfig[i].Weight
		if !enable || strings.ToLower(fnType) != "go" || weight <= 0 {
			continue
		}
		weightCount += weight
	}
	return weightCount
}

func loadCaseConfig() []caseInfo {
	caseConfig := Conf.BaseConfig.Cases
	weightCount := int64(0)
	var caseList []caseInfo
	for i := range caseConfig {
		fnType := caseConfig[i].Type
		enable := caseConfig[i].Enable
		weight := caseConfig[i].Weight
		fn := caseConfig[i].Fn
		if !enable || strings.ToLower(fnType) != "go" || weight <= 0 {
			continue
		}
		caseList = append(caseList, caseInfo{name: fn, first: weightCount+1, last: weightCount+weight})
		weightCount += weight
	}
	return caseList
}

func loadCaseFunc(entry interface{}) map[string]reflect.Value {
	funcList := make(map[string]reflect.Value)
	packageName := loadCasePackageName()
	funcValue := reflect.ValueOf(entry)
	funcType := reflect.TypeOf(entry)
	for i := 0; i < funcValue.NumMethod(); i++ {
		if (funcType.Method(i).Name == "OnStart") || (funcType.Method(i).Name == "OnStop") {
			funcList[funcType.Method(i).Name] = funcValue.Method(i)
			continue
		}
		funcList[fmt.Sprintf("%s.%s", packageName, funcType.Method(i).Name)] = funcValue.Method(i)
	}
	return funcList
}

func setEntryInt(entryValue reflect.Value, value int64) {
	entryValue.SetInt(value)
}

func setEntryInterface(entryValue reflect.Value, value interface{}) {
	entryValue.Set(reflect.ValueOf(value))
}

func checkEntryFunc(funcMap map[string]reflect.Value, name string) bool {
	if _, ok := funcMap[name]; ok {
		if !funcMap[name].IsNil() {
			return true
		}
	}
	return false
}

func loadWorkerEntry(entry interface{}) *taskEntry {
	fieldValue := reflect.ValueOf(entry)
	fieldType := reflect.TypeOf(entry)
	switch fieldType.Kind() {
	case reflect.Ptr:
		if fieldValue.IsNil() {
			return nil
		}
		if fieldValue.Elem().Type().Kind() != reflect.Struct {
			return nil
		}
		return searchReflect(fieldValue.Elem(), fieldType.Elem())
	case reflect.Struct:
		return searchReflect(fieldValue, fieldType)
	}
	return nil
}

func searchReflect(reflectValue reflect.Value, reflectType reflect.Type) *taskEntry {
	for i := 0; i < reflectType.NumField(); i++ {
		fieldValue := reflectValue.Field(i)
		fieldType := reflectType.Field(i)
		if !fieldValue.CanSet() {
			continue
		}
		switch fieldType.Type.Kind() {
		case reflect.Ptr:
			if (fieldType.Type.String() == "*boomer.WorkerEntry") {
				if fieldValue.IsNil() {
					fieldValue.Set(reflect.ValueOf(&WorkerEntry{}))
				}
				return searchWorkEntry(fieldValue.Elem())
			}
			if fieldValue.IsNil() {
				continue
			}
			if fieldValue.Elem().Type().Kind() != reflect.Struct {
				continue
			}
			searchReflect(fieldValue.Elem(), fieldType.Type.Elem())
		case reflect.Struct:
			if (fieldType.Type.String() == "boomer.WorkerEntry") {
				return searchWorkEntry(fieldValue)
			}
			searchReflect(fieldValue, fieldType.Type)
		}
		continue
	}
	return nil
}

func searchWorkEntry(reflectValue reflect.Value) *taskEntry {
	return &taskEntry{
		index: reflectValue.FieldByName("Index"),
		num: reflectValue.FieldByName("Num"),
		times: reflectValue.FieldByName("Times"),
		log: reflectValue.FieldByName("Log"),
		config: reflectValue.FieldByName("Config"),
		param: reflectValue.FieldByName("Param"),
		isRunning: reflectValue.FieldByName("IsRunning"),
	}
}

func loadCasePackageName() string {
	nameList := strings.Split(runtime.FuncForPC(reflect.ValueOf(CaseEntry).Pointer()).Name(), ".")
	packageList := strings.Split(nameList[0], "/")
	packageName := packageList[len(packageList)-1]
	return packageName
}

func callEntryFunc(reflectValue reflect.Value, errMsg *errMsg) {
	defer func() {
		if e := recover(); e != nil {
			errMsg.error = fmt.Sprint(e)
		}
	}()
	reflectValue.Call(nil)
	return
}

func randomGetEntryFunc(caseConfig []caseInfo, weightCount int64, times int64) string {
	var name string
	index := RandInt64(1, weightCount, times)
	for i := range caseConfig {
		if (index < caseConfig[i].first) || (index >caseConfig[i].last) {
			continue
		}
		name = caseConfig[i].name
		break
	}
	return name
}
