package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Resolve cross-domain problems
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (srv *HomeGateServer) checkGateRequestParams(w http.ResponseWriter, r *http.Request) (*userGate, *GateEvent, error) {
	w.Header().Set("Content-Type", "application/json")

	dumpRequest(r)

	// get the email from the token
	userEmail := r.Header.Get("Email")
	if userEmail == "" {
		var err Error
		err = SetError(err, "Failed to get user email from token")
		err.sendToClient(w, http.StatusBadRequest)
		return nil, nil, fmt.Errorf(err.Message)
	}

	var evt GateEvent
	err := json.NewDecoder(r.Body).Decode(&evt)
	if err != nil {
		var err Error
		err = SetError(err, "checkGateRequestParams: failed to decode post message")
		err.sendToClient(w, http.StatusBadRequest)
		return nil, nil, fmt.Errorf(err.Message)
	}

	if evt.GateName == nil {
		return nil, nil, fmt.Errorf("checkGateRequestParams: missing gate name parmeter")
	}

	g, ok := srv.gates[*evt.GateName]
	if !ok {
		return nil, nil, fmt.Errorf("could not find gate: %v", *evt.GateName)
	}

	if _, userAllowed := g.userEmails[userEmail]; !userAllowed {
		return nil, nil, fmt.Errorf("user: %v, does not has access to gate: %v", userEmail, *evt.GateName)
	}

	return g, &evt, nil
}

func (srv *HomeGateServer) triggerGateCommand(w http.ResponseWriter, r *http.Request) (*GateEvent, error) {
	srv.Lock()
	defer srv.Unlock()

	g, evt, err := srv.checkGateRequestParams(w, r)
	if err != nil {
		var err2 Error
		err2 = SetError(err2, err.Error())
		err2.sendToClient(w, http.StatusBadRequest)
		return nil, fmt.Errorf(err2.Message)
	}

	if g.rcState == nil {
		var err Error
		err = SetError(err, "Bad Request gate device has not been connected yet !!!")
		err.sendToClient(w, http.StatusBadRequest)
		return nil, fmt.Errorf(err.Message)
	}

	var key KeyPressed
	if evt.IsOpen != nil {
		key = Open
	} else if evt.IsClose != nil {
		key = Close
	} else if evt.IsLearnOpen != nil {
		key = LearnOpen
	} else if evt.IsLearnClose != nil {
		key = LearnClose
	} else if evt.IsTestOpen != nil {
		key = TestOpen
	} else if evt.IsTestClose != nil {
		key = TestClose
	} else if evt.IsSetOpen != nil {
		key = SetOpen
	} else if evt.IsSetClose != nil {
		key = SetClose
	}

	err = sendCommandToGate(w, g, key)
	if err != nil {
		return nil, err
	}

	return evt, nil
}

func sendCommandToGate(w http.ResponseWriter, g *userGate, key KeyPressed) error {
	if g.rcState == nil {
		var err Error
		err = SetError(err, "Oh Man ! The gate device has not been connected yet, aborting request")
		err.sendToClient(w, http.StatusBadRequest)
		return fmt.Errorf(err.Message)
	}

	select {
	case g.rcState <- key:
		resp := map[string]interface{}{
			"gate_name": g.name,
			"is_error":  false,
			"message":   fmt.Sprintf("Yay, gate command: %v, Acknowledged!", key.String()),
		}

		json.NewEncoder(w).Encode(resp)
		return nil

	default:
		log.Printf("client does not read events...")
		close(g.rcState)
		g.rcState = nil
		resp := map[string]interface{}{
			"gate_name": g.name,
			"is_error":  true,
			"message":   fmt.Sprintf("Gate command: %v, ERROR!", key.String()),
		}

		json.NewEncoder(w).Encode(resp)
		return fmt.Errorf("client does not read events, closing stream")
	}
}

func (srv *HomeGateServer) command(w http.ResponseWriter, r *http.Request) {
	srv.triggerGateCommand(w, r)
}

func (srv *HomeGateServer) siri(w http.ResponseWriter, r *http.Request) {
	var evt SiriCommand
	err := json.NewDecoder(r.Body).Decode(&evt)
	if err != nil {
		var err Error
		err = SetError(err, "siri: failed to decode post message SiriCommand")
		err.sendToClient(w, http.StatusBadRequest)
		return
	}

	var cmd KeyPressed
	siriCommand := evt.OpenOrCloseCommand
	if strings.EqualFold(siriCommand, Open.String()) {
		cmd = Open
	} else if strings.EqualFold(siriCommand, Close.String()) {
		cmd = Close
	} else {
		var err Error
		err = SetError(err, fmt.Sprintf("siri: unknown gate command (must be either close or open): %v", siriCommand))
		err.sendToClient(w, http.StatusBadRequest)
		return
	}

	// check if the gate exists
	gateID := evt.GateID
	if gateID == "" {
		var err Error
		err = SetError(err, "siri: provided gate ID is empty")
		err.sendToClient(w, http.StatusBadRequest)
		return
	}

	// lock the server
	srv.Lock()
	defer srv.Unlock()

	g, ok := srv.gates[gateID]
	if !ok {
		var err Error
		err = SetError(err, fmt.Sprintf("siri: unknown gate: %v", gateID))
		err.sendToClient(w, http.StatusBadRequest)
		return
	}

	sendCommandToGate(w, g, cmd)
}

func (srv *HomeGateServer) stream(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Infof("stream:  upgrade error: %v", err)
		return
	}
	defer c.Close()

	// validate user
	mt, message, err := c.ReadMessage()
	if err != nil {
		var err Error
		err = SetError(err, "Bad Request / data is not a gate event !!!")
		err.sendToClient(w, http.StatusBadRequest)
		log.Infof("stream:  Bad Request / data is not a gate event: %v", err)
		return
	}

	var streamRequest DeviceStreamRequest
	err = json.Unmarshal([]byte(message), &streamRequest)
	if err != nil {
		var err Error
		err = SetError(err, "Bad Request post message, expected DeviceStreamRequest type !!!")
		err.sendToClient(w, http.StatusBadRequest)
		return
	}

	// lock the server
	srv.Lock()

	gateID := streamRequest.GateID
	g, ok := srv.gates[gateID]
	if !ok {
		var err Error
		err = SetError(err, fmt.Sprintf("stream: unknown gate: %v", gateID))
		err.sendToClient(w, http.StatusBadRequest)
		log.Infof("stream: unknown gate: %v", gateID)
		return
	}

	// set the pong handler
	c.SetPongHandler(func(string) error {
		log.Infof("stream: got a PONG from the gate: %v, [OK]", gateID)
		return nil
	})

	// create a new channel for this gate
	if g.rcState != nil {
		close(g.rcState)
		g.rcState = nil
	}

	g.rcState = make(chan KeyPressed, 1)

	// unlock the server now, hold a copy of the channel
	srv.Unlock()

	// send remote control events to the client and websocket pings
	pingTicker := time.NewTicker(10 * time.Second)
	defer pingTicker.Stop()
	for {
		select {
		case <-pingTicker.C:
			// time to send a ping to the client
			c.SetWriteDeadline(time.Now().Add(5 * time.Second))
			if err := c.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Infof("stream: sent a PING to gate: %v, [FAILED]", gateID)
				return
			}
		case rcEvent := <-g.rcState:
			log.Infof("stream: send a rc command: [%v] to gate: [%v]", rcEvent, gateID)
			c.SetWriteDeadline(time.Now().Add(5 * time.Second))
			err = c.WriteMessage(mt, []byte(fmt.Sprintf("%v", rcEvent)))
			if err != nil {
				log.Infof("stream: FAILED sending a rc command: [%v] to gate: [%v] [ERROR = %v]", rcEvent, gateID, err)
				return
			}
		}
	}
}

func (srv *HomeGateServer) signUp(w http.ResponseWriter, r *http.Request) {
	connection := GetDatabase()
	defer CloseDatabase(connection)

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		var err Error
		err = SetError(err, "Error in reading payload.")
		err.sendToClient(w, http.StatusBadRequest)
		return
	}

	gateName := user.MyGateName

	// check that this gate ID does not exists
	var dbGate Gate
	connection.Where("name = ?", gateName).First(&dbGate)

	if dbGate.Name != "" {
		var err Error
		err = SetError(err, "Provided user gate ID is already registered")
		err.sendToClient(w, http.StatusBadRequest)
		return
	}

	var dbuser User
	connection.Where("email = ?", user.Email).First(&dbuser)

	//check email is already registered or not
	if dbuser.Email != "" {
		var err Error
		err = SetError(err, "Email already in use")
		err.sendToClient(w, http.StatusBadRequest)
		return
	}

	user.Password, err = GeneratehashPassword(user.Password)
	if err != nil {
		var err Error
		err = SetError(err, "Error in password hashing")
		err.sendToClient(w, http.StatusInternalServerError)
		return
	}

	// allocate a unique gate identifier to the user
	//gateID, _ := sf.NextID()
	myGate := Gate{Name: fmt.Sprintf("gate-%v", gateName), UserEmails: []string{user.Email}}

	// insert gate details in database
	res := connection.Create(&myGate)
	if res.Error != nil {
		fmt.Println(res.Error)
	}

	// add the user create gate (currently a single one)
	user.MyGateName = myGate.Name
	user.Gates = []string{user.MyGateName}

	// insert user details in database
	res = connection.Create(&user)
	if res.Error != nil {
		fmt.Println(res.Error)
	}

	// update the in-memory state
	srv.Lock()

	g := &userGate{
		name:       myGate.Name,
		userEmails: make(map[string]bool),
		rcState:    nil,
	}

	// keep the users emails in a map for quick access
	for _, userEmail := range myGate.UserEmails {
		g.userEmails[userEmail] = true
	}

	// keep it in the global gate map
	srv.gates[myGate.Name] = g

	srv.Unlock()

	// return the response with a welcome message
	registerResponse := RegisterResponse{
		Name:    user.Name,
		Email:   user.Email,
		Message: fmt.Sprintf("%v registered OK", user.Name),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(registerResponse)
}

func (srv *HomeGateServer) signIn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	connection := GetDatabase()
	defer CloseDatabase(connection)

	var authDetails Authentication

	err := json.NewDecoder(r.Body).Decode(&authDetails)
	if err != nil {
		var err Error
		err = SetError(err, "Error in reading payload.")
		err.sendToClient(w, http.StatusBadRequest)
		return
	}

	var authUser User
	connection.Where("email = 	?", authDetails.Email).First(&authUser)

	if authUser.Email == "" {
		var err Error
		err = SetError(err, "Username or Password is incorrect")
		err.sendToClient(w, http.StatusBadRequest)
		return
	}

	check := CheckPasswordHash(authDetails.Password, authUser.Password)

	if !check {
		var err Error
		err = SetError(err, "Username or Password is incorrect")
		err.sendToClient(w, http.StatusBadRequest)
		return
	}

	// TODO: decide if we need a claim ?
	validToken, err := GenerateJWT(authUser.Email, "user")
	if err != nil {
		var err Error
		err = SetError(err, "Failed to generate token")
		err.sendToClient(w, http.StatusBadRequest)
		return
	}

	loginResponse := LoginResponse{
		ID:          authUser.ID,
		Message:     fmt.Sprintf("%v login OK", authUser.Name),
		UserName:    authUser.Name,
		Email:       authUser.Email,
		AccessToken: validToken,
		MyGateName:  authUser.MyGateName,
		Gates:       authUser.Gates,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(loginResponse)
}

func (srv *HomeGateServer) userIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// get the email from the token
	userEmail := r.Header.Get("Email")
	if userEmail == "" {
		var err Error
		err = SetError(err, "Failed to get user email from token")
		err.sendToClient(w, http.StatusBadRequest)
		return
	}

	connection := GetDatabase()
	defer CloseDatabase(connection)

	var authUser User
	connection.Where("email = 	?", userEmail).First(&authUser)

	if authUser.Email == "" {
		var err Error
		err = SetError(err, "user email does not exists in database")
		err.sendToClient(w, http.StatusBadRequest)
		return
	}

	var uGate Gate
	connection.Where("name = 	?", authUser.MyGateName).First(&uGate)
	if uGate.Name == "" {
		var err Error
		err = SetError(err, "user gate does not exists in database")
		err.sendToClient(w, http.StatusBadRequest)
		return
	}

	userData := map[string]interface{}{
		"name":    authUser.Name,
		"gates":   authUser.Gates,
		"my_gate": authUser.MyGateName,
		"users":   uGate.UserEmails,
	}

	json.NewEncoder(w).Encode(userData)
}
