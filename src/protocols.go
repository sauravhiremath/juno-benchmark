package src

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"
)

/*

RegisterModule initializes the given connection.

Reads the response and passes the ModuleID to moduleIDs channel

*/
func RegisterModule(c *net.Conn) {
	var i int = 0
	uuid := uuid.NewV4().String()
	iMsg, _ := json.Marshal(GetInitializeModuleMessage("Module-", uuid))
	fmt.Fprintf(*c, string(iMsg)+"\n")
	scanner := bufio.NewScanner(*c)
	for scanner.Scan() {
		// log.Println("Initial message " + strconv.Itoa(i) + scanner.Text())
		i++
		if i > 1 {
			break
		}
	}
	moduleIDs <- ("Module-" + uuid)
	// log.Println("[INFO] Finished registering module...")
}

/*

FunctionCall handles message transmission over given connection. Call function juno equivalent

Creates new request ID (UUID V4) for each transmission.

Calculates latency between request and response

*/
func FunctionCall(c *net.Conn, sub *sync.WaitGroup, moduleID string) {
	// Create unique requestID
	defer sub.Done()

	msg, _ := json.Marshal(GetFunctionCallMessage(moduleID, "test"))
	start := time.Now()
	fmt.Fprintf(*c, string(msg)+"\n")
	_, err := bufio.NewReader(*c).ReadString('\n')
	if err != nil {
		log.Panic("[ERROR] Message recieve failed: ", err)
		return
	}
	// log.Println(time.Since(start))
	stats <- time.Since(start)
	return
}

/*

DeclareFunction is a Go port for `Juno Declare Function`

Declares a function `test` for all the connected modules

*/
func DeclareFunction(c *net.Conn) {

	msg, _ := json.Marshal(GetDeclareFunctionMessage("test"))
	fmt.Fprintf(*c, string(msg)+"\n")
	_, err := bufio.NewReader(*c).ReadString('\n')
	if err != nil {
		log.Panic("[ERROR] Message recieve failed: ", err)
		return
	}
	// fmt.Println("[DEBUG] Function declared for moduleID")
}
