package boomer

import "time"

var msgChannel = make(chan *message)


func (r *slaveRunner) getMessage() {
	go func() {
		for {
			select {
			case msg := <- msgChannel:
				r.recvMessage(msg)
			case <-r.closeChan:
				return
			}
		}
	}()
}

func (r *slaveRunner) recvMessage(msg *message) {
	if msg.Type != "task_over" && msg.Type != "heartbeat" {
		control.log.info("Recv message type:%s", msg.Type)
	}
	switch msg.Type {
	case "create":
		if control.isRunning == false {
			go onCreate(msg)
			for {
				time.Sleep(time.Duration(1)*time.Millisecond)
				if control.isRunning == true {
					break
				}
			}
			go monitorStartChannel()
			go monitorStopChannel()
			go monitorTaskChannel()
		}
	case "hatch":
		go onHatch(msg)
	case "heartbeat":
		go onHeartbeat()
	case "hatch_over":
		go onHatchOver()
	case "task_over":
		go onTaskOver()
	case "monitor":
		go onMonitor()
	case "complete":
		go onComplete(msg)
	case "exit":
		go onExit()
	case "restart":
		go onRestart()
	case "config":
		go onConfig(msg)
	case "quit":
		if control.isRunning == true {
			onQuit()
		}
	case "stop":
		if control.isRunning == true {
			onStop()
		}
	}
	go r.onMessage(msg)
}
