package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
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

func (srv *HomeGateServer) times(w http.ResponseWriter, r *http.Request) {
	srv.Lock()
	defer srv.Unlock()

	params := mux.Vars(r)
	gateName := params["gatename"]
	g, ok := srv.gates[gateName]
	if ok {
		fmt.Fprintf(w, "Last Open for Gate : %v --> %v\n", gateName, g.lastOpen)
		fmt.Fprintf(w, "Last Close for Gate: %v --> %v\n", gateName, g.lastClose)
	} else {
		fmt.Fprintf(w, "Could not find a Gate: %v", gateName)
	}
}

func (srv *HomeGateServer) checkGateRequestParams(w http.ResponseWriter, gateName, userEmail, password *string) (*userGate, error) {
	if gateName == nil || userEmail == nil || password == nil {
		http.Error(w, "Bad Request / missing parameter !!!", http.StatusBadRequest)
		return nil, fmt.Errorf("handler: Open called, missing parmeters")
	}

	g, ok := srv.gates[*gateName]
	if !ok {
		http.Error(w, "Bad Request / gate does not exist  !!!", http.StatusBadRequest)
		return nil, fmt.Errorf("handler: unknown gate: %v", *gateName)
	} else {
		_, ok := g.userEmails[*userEmail]
		if !ok {
			http.Error(w, "Bad Request / user does not exist  !!!", http.StatusBadRequest)
			return nil, fmt.Errorf("handler: gate [ %v ], unknown user: %v", *gateName, userEmail)
		}

		return g, nil
	}
}

func (srv *HomeGateServer) open(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var openEvent GateEvent
	err := json.NewDecoder(r.Body).Decode(&openEvent)
	if err != nil {
		log.Printf("handler: failed to decode the openEvent message: %v", err.Error())
		http.Error(w, "Bad Request / data is not openEvent !!!", http.StatusBadRequest)
		return
	}

	srv.Lock()
	defer srv.Unlock()

	g, err := srv.checkGateRequestParams(w, openEvent.GateName, openEvent.Email, openEvent.Password)
	if err != nil {
		log.Printf("%v", err)
		return
	}

	if g.rcState == nil {
		http.Error(w, "Bad Request gate device haven't connected yet !!!", http.StatusBadRequest)
		return
	}

	select {
	case g.rcState <- Open:
		g.lastOpen = time.Now()
		fmt.Fprintf(w, "%v's gate requested to OPEN Acknowledged!\n", g.name)
	default:
		log.Printf("client does not read events...")
		close(g.rcState)
	}

}

func (srv *HomeGateServer) close(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var closeEvent CloseEvent
	err := json.NewDecoder(r.Body).Decode(&closeEvent)
	if err != nil {
		log.Printf("handler: failed to decode the closeEvent message: %v", err.Error())
		http.Error(w, "Bad Request / data is not closeEvent !!!", http.StatusBadRequest)
		return
	}

	srv.Lock()
	defer srv.Unlock()

	g, err := srv.checkGateRequestParams(w, closeEvent.GateName, closeEvent.Email, closeEvent.Password)
	if err != nil {
		log.Printf("%v", err)
		return
	}

	if g.rcState == nil {
		http.Error(w, "Bad Request gate device haven't connected yet !!!", http.StatusBadRequest)
		return
	}

	select {
	case g.rcState <- Close:
		g.lastClose = time.Now()
		fmt.Fprintf(w, "%v's gate requested to CLOSE Acknowledged!\n", g.name)
	default:
		log.Printf("client does not read events...")
		close(g.rcState)
	}

}

func (srv *HomeGateServer) learnOpen(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var learnEvent LearnEvent
	err := json.NewDecoder(r.Body).Decode(&learnEvent)
	if err != nil {
		log.Printf("handler: failed to decode the learnEvent message: %v", err.Error())
		http.Error(w, "Bad Request / data is not learnEvent !!!", http.StatusBadRequest)
		return
	}

	srv.Lock()
	defer srv.Unlock()

	g, err := srv.checkGateRequestParams(w, learnEvent.GateName, learnEvent.Email, learnEvent.Password)
	if err != nil {
		log.Printf("%v", err)
		return
	}

	if g.rcState == nil {
		http.Error(w, "Bad Request gate device haven't connected yet !!!", http.StatusBadRequest)
		return
	}

	select {
	case g.rcState <- LearnOpen:
		fmt.Fprintf(w, "%v's gate requested to LearnOpen Acknowledged!\n", g.name)
	default:
		log.Printf("client does not read events...")
		close(g.rcState)
	}

	fmt.Fprintf(w, "%v's gate requested to LEARN Open button -  Acknowledged!\n", g.name)
}

func (srv *HomeGateServer) learnClose(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var learnEvent LearnEvent
	err := json.NewDecoder(r.Body).Decode(&learnEvent)
	if err != nil {
		log.Printf("handler: failed to decode the learnEvent message: %v", err.Error())
		http.Error(w, "Bad Request / data is not learnEvent !!!", http.StatusBadRequest)
		return
	}

	srv.Lock()
	defer srv.Unlock()

	g, err := srv.checkGateRequestParams(w, learnEvent.GateName, learnEvent.Email, learnEvent.Password)
	if err != nil {
		log.Printf("%v", err)
		return
	}

	if g.rcState == nil {
		http.Error(w, "Bad Request gate device haven't connected yet !!!", http.StatusBadRequest)
		return
	}

	select {
	case g.rcState <- LearnClose:
		fmt.Fprintf(w, "%v's gate requested to LearnClose Acknowledged!\n", g.name)
	default:
		log.Printf("client does not read events...")
		close(g.rcState)
	}

	fmt.Fprintf(w, "%v's gate requested to LEARN Close button -  Acknowledged!\n", g.name)
}

func (srv *HomeGateServer) testOpen(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var testEvent TestEvent
	err := json.NewDecoder(r.Body).Decode(&testEvent)
	if err != nil {
		log.Printf("handler: failed to decode the testEvent message: %v", err.Error())
		http.Error(w, "Bad Request / data is not testEvent !!!", http.StatusBadRequest)
		return
	}

	srv.Lock()
	defer srv.Unlock()

	g, err := srv.checkGateRequestParams(w, testEvent.GateName, testEvent.Email, testEvent.Password)
	if err != nil {
		log.Printf("%v", err)
		return
	}

	if g.rcState == nil {
		http.Error(w, "Bad Request gate device haven't connected yet !!!", http.StatusBadRequest)
		return
	}

	select {
	case g.rcState <- TestOpen:
		fmt.Fprintf(w, "%v's gate requested to TestOpen Acknowledged!\n", g.name)
	default:
		log.Printf("client does not read events...")
		close(g.rcState)
	}
}

func (srv *HomeGateServer) testClose(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var testEvent TestEvent
	err := json.NewDecoder(r.Body).Decode(&testEvent)
	if err != nil {
		log.Printf("handler: failed to decode the testEvent message: %v", err.Error())
		http.Error(w, "Bad Request / data is not testEvent !!!", http.StatusBadRequest)
		return
	}

	srv.Lock()
	defer srv.Unlock()

	g, err := srv.checkGateRequestParams(w, testEvent.GateName, testEvent.Email, testEvent.Password)
	if err != nil {
		log.Printf("%v", err)
		return
	}

	if g.rcState == nil {
		http.Error(w, "Bad Request gate device haven't connected yet !!!", http.StatusBadRequest)
		return
	}

	select {
	case g.rcState <- TestClose:
		fmt.Fprintf(w, "%v's gate requested to TestClose Acknowledged!\n", g.name)
	default:
		log.Printf("client does not read events...")
		close(g.rcState)
	}
}

func (srv *HomeGateServer) setOpen(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var setEvent SetEvent
	err := json.NewDecoder(r.Body).Decode(&setEvent)
	if err != nil {
		log.Printf("handler: failed to decode the setEvent message: %v", err.Error())
		http.Error(w, "Bad Request / data is not testEvent !!!", http.StatusBadRequest)
		return
	}

	srv.Lock()
	defer srv.Unlock()

	g, err := srv.checkGateRequestParams(w, setEvent.GateName, setEvent.Email, setEvent.Password)
	if err != nil {
		log.Printf("%v", err)
		return
	}

	if g.rcState == nil {
		http.Error(w, "Bad Request gate device haven't connected yet !!!", http.StatusBadRequest)
		return
	}

	select {
	case g.rcState <- SetOpen:
		fmt.Fprintf(w, "%v's gate requested to SetOpen Acknowledged!\n", g.name)
	default:
		log.Printf("client does not read events...")
		close(g.rcState)
	}
}

func (srv *HomeGateServer) setClose(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var setEvent SetEvent
	err := json.NewDecoder(r.Body).Decode(&setEvent)
	if err != nil {
		log.Printf("handler: failed to decode the setEvent message: %v", err.Error())
		http.Error(w, "Bad Request / data is not testEvent !!!", http.StatusBadRequest)
		return
	}

	srv.Lock()
	defer srv.Unlock()

	g, err := srv.checkGateRequestParams(w, setEvent.GateName, setEvent.Email, setEvent.Password)
	if err != nil {
		log.Printf("%v", err)
		return
	}

	if g.rcState == nil {
		http.Error(w, "Bad Request gate device haven't connected yet !!!", http.StatusBadRequest)
		return
	}

	select {
	case g.rcState <- SetClose:
		fmt.Fprintf(w, "%v's gate requested to SetClose Acknowledged!\n", g.name)
	default:
		log.Printf("client does not read events...")
		close(g.rcState)
	}
}

func (srv *HomeGateServer) stream(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	// validate user
	mt, message, err := c.ReadMessage()
	if err != nil {
		http.Error(w, "Bad Request / data is not StatusEvent !!!", http.StatusBadRequest)
		return
	}
	var statusEvent StatusEvent
	err = json.Unmarshal([]byte(message), &statusEvent)
	if err != nil {
		http.Error(w, "Bad Request / data is not StatusEvent !!!", http.StatusBadRequest)
		return
	}

	srv.Lock()
	g, err := srv.checkGateRequestParams(w, statusEvent.GateName, statusEvent.Email, statusEvent.Password)
	if err != nil {
		http.Error(w, "Bad Request / invalid gate", http.StatusBadRequest)
		return
	}

	// create a new channel for this gate
	if g.rcState != nil {
		close(g.rcState)
	}

	g.rcState = make(chan KeyPressed, 1)
	srv.Unlock()

	for rcEvent := range g.rcState {
		err = c.WriteMessage(mt, []byte(fmt.Sprintf("%v", rcEvent)))
		if err != nil {
			break
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
		log.Fatalln("Error in password hashing.")
	}

	// allocate a unique gate identifier to the user
	gateID, _ := sf.NextID()
	myGate := Gate{Name: fmt.Sprintf("gate-%v", gateID), UserEmails: []string{user.Email}}

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

func (srv *HomeGateServer) adminIndex(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Role") != "admin" {
		w.Write([]byte("Not authorized."))
		return
	}
	w.Write([]byte("Welcome, Admin."))
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

	var userGate Gate
	connection.Where("name = 	?", authUser.MyGateName).First(&userGate)
	if userGate.Name == "" {
		var err Error
		err = SetError(err, "user gate does not exists in database")
		err.sendToClient(w, http.StatusBadRequest)
		return
	}

	userData := map[string]interface{}{
		"name":    authUser.Name,
		"gates":   authUser.Gates,
		"my_gate": authUser.MyGateName,
		"users":   userGate.UserEmails,
	}

	json.NewEncoder(w).Encode(userData)
}
