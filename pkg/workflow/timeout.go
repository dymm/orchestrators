package workflow

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/dymm/gorchestrator/pkg/messaging"
)

//StartSessionTimeoutChecking start the timeout control loop
func StartSessionTimeoutChecking(timeout int, queue messaging.Queue, timeoutDestination string) {

	go timeoutLoop(timeout, queue, timeoutDestination)
}

func timeoutLoop(timeout int, queue messaging.Queue, timeoutDestination string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in timeoutLoop", r)
			os.Exit(-2)
		}
	}()

	for {
		workitemToSend := listAllTimeoutedSession(timeout, queue, timeoutDestination)
		for _, workitem := range workitemToSend {
			err := queue.Send(timeoutDestination, workitem)
			if err != nil {
				fmt.Println("Can't send a timeout message", err)
			}
		}
		time.Sleep(1 * time.Second)
	}

}

func listAllTimeoutedSession(timeout int, queue messaging.Queue, timeoutDestination string) []messaging.WorkItem {

	sessionListMutex.Lock()
	defer sessionListMutex.Unlock()

	sessionCount := len(sessionList)
	fmt.Printf("%d session in memory\n", len(sessionList))

	workItemToSend := make([]messaging.WorkItem, 0, sessionCount)
	now := time.Now()
	for key, session := range sessionList {
		if now.Sub(session.Started) >= time.Duration(timeout)*time.Second {

			session.Timeouted = true
			values := map[string]string{
				"sessionId": strconv.FormatUint(key, 10),
				"error":     `{"message":"Timeout"}`,
			}
			timeouted := messaging.NewWorkItem(values)
			workItemToSend = append(workItemToSend, timeouted)
		}
	}
	return workItemToSend
}
