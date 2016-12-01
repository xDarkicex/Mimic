package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/crypto/bcrypt"

	mgo "gopkg.in/mgo.v2"

	"github.com/xDarkicex/goMimic/db"
	"github.com/xDarkicex/goMimic/structs"
)

var (
	session *mgo.Session
)

func init() {
	err := db.Dial()
	if err != nil {
		log.Fatal(err)
	}
	session = db.Session()
}

func main() {
	// create self calling go routine
	go func() {
		interruptChannel := make(chan os.Signal, 0)
		// look for system interruptions
		signal.Notify(interruptChannel, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
		// lock lower code untel interruptChannel receives signal
		<-interruptChannel

		// Other cleanup tasks
		fmt.Println("Closing connection")
		fmt.Println("Saving session")
		// Accually Close DB session, maintain DATA integrity
		session.Close()
		os.Exit(0)
	}()

	server := &http.Server{
		Addr:    ":8080",
		Handler: &Mimic{},
	}
	fmt.Println("Listening on 127.0.0.1:8080")
	log.Fatal(server.ListenAndServe())
}

// Mimic our server.
type Mimic struct{}

func (mimic *Mimic) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" {
		fmt.Fprintln(writer, "URI is weird!")
		return
	}
	fmt.Println(request.URL.Path)
	method := request.FormValue("method")
	if method == "" {
		fmt.Fprintln(writer, "No Method!")
		return
	}
	switch method {
	case "login":
		arguments, err := getArguments(request, []string{"username", "password"})
		if err != nil {
			fmt.Fprintln(writer, err.Error())
			return
		}
		// Do something with logging in.
		_ = arguments

	case "register":
		arguments, err := getArguments(request, []string{"username", "password"})
		if err != nil {
			fmt.Fprintln(writer, err.Error())
			return
		}
		// Do something with register.
		arg := arguments
		username := arg[0]
		password := arg[1]
		hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			fmt.Println("Hashing Password incomplete")
		}

		session := db.Session()
		defer session.Close()
		c := session.DB("MimicDB").C("User")
		// Insert Datas
		err = c.Insert(&structs.MongoUser{
			Username: username,
			Password: string(hashedPass),
		})
		if err != nil {
			fmt.Printf("error  new user %s\n", err)
			return
		}
		fmt.Println(arg, username, password, hashedPass)
		return

	}
}

// Gets arguments from a form value using a slice of strings as the key.
func getArguments(request *http.Request, keys []string) (arguments []string, err error) {
	err = errors.New("")
	for _, e := range keys {
		value := request.FormValue(e)
		if value == "" {
			err = fmt.Errorf("Required parameter: %s\n%s", e, err.Error())
		} else {
			arguments = append(arguments, value)
		}
	}
	if err.Error() == "" {
		err = nil
	}
	return arguments, err
}
