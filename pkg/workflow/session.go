package workflow

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/dymm/gorchestrator/pkg/messaging"
)

//Session about the assigned workflow
type Session struct {
	Key              uint64
	assignedWorkflow int
	CurrentStep      StepInformation
}

//StepInformation information about the current step
type StepInformation struct {
	Name        string
	Process     string
	Started     time.Time
	TimeoutTime time.Time
	Workitem    messaging.WorkItem
}

var sessionListMutex sync.Mutex
var sessionList map[uint64]*Session
var uniqueKey uint64

func init() {
	fmt.Println("Initialize package workflow")
	sessionList = make(map[uint64]*Session)
	uniqueKey = 0
}

//getSessionStoredInTheWorkItem return the workflow session from a work item
// return a session and a boolean which will be true if the session is found and false if not
func getSessionStoredInTheWorkItem(workItem messaging.WorkItem) (session *Session, found bool) {

	var sessionKeyStr string
	sessionKeyStr, found = workItem.GetValues()["sessionId"]
	if found == false {
		return
	}
	sessionKey, err := strconv.ParseUint(sessionKeyStr, 10, 64)
	if err == nil {
		session, err = getSession(sessionKey)
	}
	found = err == nil
	return
}

func createNewSession() *Session {
	sessionListMutex.Lock()
	defer sessionListMutex.Unlock()

	uniqueKey = uniqueKey + 1
	session := Session{Key: uniqueKey, assignedWorkflow: -1, CurrentStep: StepInformation{}}
	sessionList[uniqueKey] = &session

	fmt.Printf("Create session %d\n", session.Key)
	return &session
}

func getSession(key uint64) (*Session, error) {
	sessionListMutex.Lock()
	defer sessionListMutex.Unlock()

	session, found := sessionList[key]
	if found == false {
		return nil, fmt.Errorf("No session with key %d found", key)
	}
	return session, nil
}

// DeleteSession delete a session
func DeleteSession(session *Session) {
	sessionListMutex.Lock()
	defer sessionListMutex.Unlock()

	delete(sessionList, session.Key)
	fmt.Printf("Delete session %d\n", session.Key)
}

func setStepInformationInSession(session *Session, step Step, workitem messaging.WorkItem) {
	session.CurrentStep.Process = step.Process
	session.CurrentStep.Started = time.Now()
	session.CurrentStep.Workitem = workitem
	if step.Timeout > 0 {
		session.CurrentStep.TimeoutTime = session.CurrentStep.Started.Add(time.Duration(step.Timeout) * time.Second)
	}
}
