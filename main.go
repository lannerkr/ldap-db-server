package main

// version 1.0
// modified from ldapserver v1.2
// mongoDB connect
// write login history to monogDB/ldapDB/user_history
// version 1.1
// mongoDB connection customize
// version 1.2
// realm/status info db write add
// version 1.3
// update user ip to DB
// backup server rdcheck add
// backup server get framed-ip from DB add
// version 1.3.1
// update user history to seconday db (master server only)
// version 1.4
// beckend.go radius connection timeout(2sec) add
// server.go service is stopping when rdDown is true && secondary is false

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gopenguin/minimal-ldap-proxy/types"
	msg "github.com/lor00x/goldap/message"
)

var (
	serverAddr     string
	baseDn         string
	secondary      bool
	conf           string
	rdcheck        string
	rdDown         bool
	rdcuser        string
	rdcrealm       string
	FramedIPstring string
	svcStopped     bool = false
)

type Frontend struct {
	serverAddr string
	baseDn     string
	server     *Server
}

type LoginHistory struct {
	Realm     string    `bson:"realm"`
	UserName  string    `bson:"user_name"`
	LastLogin time.Time `bson:"last_login"`
	Enabled   string    `bson:"enabled"`
	FramedIP  string    `bson:"framedip"`
}

func init() {
	conf = os.Args[1]
	config(conf)
	rdcuser = configuration.RadHCuser
	rdcrealm = configuration.RadHCrealm
}

func main() {

	logpath := configuration.LogPath
	fpLog, err := os.OpenFile(logpath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer fpLog.Close()
	log.SetOutput(fpLog)
	Logger = log.New(fpLog, "", log.LstdFlags)

	serverAddr = configuration.LdapServer
	baseDn = configuration.LdapBaseDn
	secondary = configuration.Secondary
	frontend := &Frontend{
		serverAddr: serverAddr,
		baseDn:     baseDn,
		server:     NewServer(),
	}

	router := NewRouteMux()
	router.Bind(frontend.handleBind)
	router.Search(frontend.handleSearchUser).
		BaseDn(frontend.baseDn)
	router.Search(frontend.handleSearchGeneric)

	frontend.server.Handle(router)

	frontend.Serve()

	// time.Sleep(time.Second * 30)
	// Logger.Println("server is stopping")
	// frontend.server.Listener.Close()

	// time.Sleep(time.Second * 30)
	// Logger.Println("server is starting")
	// frontend.Serve()

	//if secondary {
	switch secondary {
	case false:
		go frontend.radiuscheck()
	case true:
		go secondarycheck()
	}

	//}

	// When CTRL+C, SIGINT and SIGTERM signal occurs
	// Then stop server gracefully
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	close(ch)

	frontend.Stop()
}

func (f *Frontend) handleBind(w ResponseWriter, m *Message) {
	r := m.GetBindRequest()

	res := NewBindResponse(LDAPResultSuccess)
	if string(r.Name()) == "" {
		w.Write(res)
		return
	}

}

func (f *Frontend) handleSearchUser(w ResponseWriter, m *Message) {
	r := m.GetSearchRequest()

	log.Printf("Searching on %s for %s with \n", r.BaseObject(), r.FilterString())

	oname := r.FilterString()
	log.Println("oname: ", oname)
	user, realm := userFromCn(oname)

	if user == "" {
		log.Printf("extract user: %v\n", user)
		res := NewSearchResultDoneResponse(LDAPResultNoSuchAttribute)
		w.Write(res)
		return
	}

	if user != "" {

		//go writedb(user, realm)

		switch rdDown {
		case false:
			// get ip from RADIUS
			FramedIPstring = beckend(user, realm, false)
		case true:
			// get ip from DB
			framedIP := getIPfromDB(user, realm)
			FramedIPstring = framedIP
		}

		go writedb(user, realm, FramedIPstring)

		if FramedIPstring == "Not Available" || FramedIPstring == "" {
			log.Printf("extract user: %v Framed-IP is not available\n", user)
			res := NewSearchResultDoneResponse(LDAPResultNoSuchAttribute)
			w.Write(res)
			return
		}

		var result types.Result
		result.Rdn = user
		result.Attributes = map[string][]string{
			"FRAMED-IP-ADDRESS": {FramedIPstring},
		}

		entry := NewSearchResultEntry(fmt.Sprintf("cn=%s", user))

		for key, value := range result.Attributes {
			var resAttDesc msg.AttributeDescription = msg.AttributeDescription(key)
			var attributeValues []msg.AttributeValue
			for _, v := range value {
				attributeValues = append(attributeValues, msg.AttributeValue(v))
			}
			entry.AddAttribute(resAttDesc, attributeValues...)
		}

		w.Write(entry)

		res := NewSearchResultDoneResponse(LDAPResultSuccess)
		w.Write(res)
	}

}

func (f *Frontend) handleSearchGeneric(w ResponseWriter, m *Message) {
	r := m.GetSearchRequest()
	log.Printf("Unhandled search request: %s\n", r.BaseObject())
	res := NewSearchResultDoneResponse(LDAPResultNoSuchObject)
	w.Write(res)
}

func (f *Frontend) Serve() {
	go func() {
		err := f.server.ListenAndServe(f.serverAddr)
		log.Println(err)
	}()
}

func (f *Frontend) Stop() {
	f.server.Stop()
}
