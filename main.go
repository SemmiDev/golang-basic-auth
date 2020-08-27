package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func OutputJSON(w http.ResponseWriter, o interface{}) {
	res, err := json.Marshal(o)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
	w.Write([]byte("\n"))
}
func ActionPersonCareer(w http.ResponseWriter, r *http.Request) {
	if id := r.URL.Query().Get("id"); id != "" {
		OutputJSON(w, SelectPerson(id))
		return
	}

	OutputJSON(w, GetPerson())
}

func main() {
	mux := http.DefaultServeMux
	mux.HandleFunc("/personCareer", ActionPersonCareer)

	var handler http.Handler = mux
	handler = MiddlewareAuth(handler)
	handler = MiddlewareAllowOnlyGet(handler)

	server := new(http.Server)
	server.Addr = ":9000"
	server.Handler = handler

	log.Println(" ::: SERVER STARTED AT PORT 9000")
	server.ListenAndServe()

	/*
		curl -X GET -- user:password localhost:9000/personCareer
		curl -X GET -- user:password localhost:9000/personCareer?id=a001
	*/
}

// PERSON CAREER MODEL

var careers = []*Career{}
type Career struct {
	Id    string
	Name  string
	Age  uint8
	CareerLevel string
	Status bool
}

func GetPerson() []*Career {
	return careers
}

func SelectPerson(id string) *Career {
	for _, each := range careers {
		if each.Id == id {
			return each
		}
	}
	return nil
}

func init() {
	careers = append(careers, &Career{Id: "a001", Name: "Sam", Age: 19, CareerLevel: "Technical Architect", Status: true})
	careers = append(careers, &Career{Id: "a002", Name: "Ayatullah", Age: 19, CareerLevel: "CFO", Status: true})
	careers = append(careers, &Career{Id: "a003", Name: "Aditya", Age: 18, CareerLevel: "CTO", Status: true})
}

// MIDDLEWARE

const  USERNAME = "sammidev"
const  PASSWORD = "sammidev"

func MiddlewareAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		username, password, ok := request.BasicAuth()
		if !ok {
			writer.Write([]byte(`something went wrong`))
			return
		}
		isValid := (username == USERNAME) && (password == PASSWORD)
		if !isValid {
			writer.Write([]byte(`wrong username/password`))
			return
		}
		next.ServeHTTP(writer, request)

	})
}
func MiddlewareAllowOnlyGet(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != "GET" {
			writer.Write([]byte("ONLY GET IS ALLOWED BRO"))
			return
		}
		next.ServeHTTP(writer,request)
	})
}