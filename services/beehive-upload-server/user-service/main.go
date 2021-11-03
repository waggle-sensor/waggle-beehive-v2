package main

// testing: docker run -d --name upload-server -p 8080:80  waggle/beehive-upload-server
// curl localhost:8080/user
import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
)

type NodeObj struct {
	ID string `json:"id"`
}

type APIResponse struct {
	Data interface{} `json:"data"`
}

type APIErrorResponse struct {
	Error string `json:"error"`
}

func GetUsers() (users []string, err error) {

	file, err := os.Open("/etc/passwd")

	if err != nil {
		err = fmt.Errorf("error opening file: %s", err.Error())
		return
	}

	defer file.Close()

	reader := bufio.NewReader(file)

	for {
		var line string
		line, err = reader.ReadString('\n')

		// skip all line starting with #
		if strings.HasPrefix(line, "#") {
			continue
		}

		// get the username and description
		lineSlice := strings.FieldsFunc(line, func(divide rune) bool {
			return divide == ':' // we divide at colon
		})

		if len(lineSlice) > 0 {
			user := lineSlice[0]
			if !strings.HasPrefix(user, "node-") {
				continue
			}
			users = append(users, lineSlice[0])
		}

		if err == io.EOF {
			err = nil
			break
		}
		if err != nil {
			err = fmt.Errorf("error reading file: %s", err.Error())
			return
		}

	}
	if users == nil {
		users = []string{}
	}
	return
}

func rootListener(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "alive")
}

func userListener(w http.ResponseWriter, req *http.Request) {

	fmt.Println("GET /user was called")
	users, err := GetUsers()
	if err != nil {
		fmt.Printf("returned error: %s\n", err.Error())

		ar := APIErrorResponse{}
		ar.Error = err.Error()
		response_json, _ := json.Marshal(ar)
		http.Error(w, string(response_json), http.StatusInternalServerError)
		return
	}

	ar := APIResponse{}
	ar.Data = users
	response_json, _ := json.Marshal(ar)
	// err := Sync()

	fmt.Fprint(w, string(response_json))

}

func CreateUser(username string) (err error) {
	// adduser -D -g "" "$username"
	// passwd -u "$username"
	// chown -R "$username:$username" "/home/$username"
	var out []byte
	out, err = exec.Command("adduser", "-D", "-g", "\"\"", username).CombinedOutput()
	if err != nil {
		err = fmt.Errorf("command failed: %s (%s)", err.Error(), string(out))
		return
	}

	out, err = exec.Command("passwd", "-u", username).CombinedOutput()
	if err != nil {
		err = fmt.Errorf("command failed: %s (%s)", err.Error(), string(out))
		return
	}

	out, err = exec.Command("chown", "-R", username+":"+username, "/home/"+username).CombinedOutput()
	if err != nil {
		err = fmt.Errorf("command failed: %s (%s)", err.Error(), string(out))
		return
	}

	return
}

func userCreateListener(w http.ResponseWriter, req *http.Request) {
	fmt.Println("POST /user/{user} was called")

	vars := mux.Vars(req)
	user := vars["user"]

	// user has to be lower-case, 0-9a-f and 16 characters long
	// node- prefix is optional

	user = strings.TrimPrefix(user, "node-")

	if len(user) != 16 {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "invalid user name, wrong length")
		return
	}

	re := regexp.MustCompile(`^[0-9a-f]{16}$`)

	if !re.MatchString(user) {

		//fmt.Fprint(w, "invalid user name")
		ar := APIErrorResponse{}
		ar.Error = "invalid user name"
		response_json, _ := json.Marshal(ar)
		http.Error(w, string(response_json), http.StatusInternalServerError)
		return
	}

	// Add prefix
	user = "node-" + user

	// check if user exists ?

	err := CreateUser(user)
	if err != nil {
		ar := APIErrorResponse{}
		ar.Error = fmt.Sprintf("Error creating user: %s", err.Error())
		response_json, _ := json.Marshal(ar)
		http.Error(w, string(response_json), http.StatusInternalServerError)
		return
	}

	ar := NodeObj{ID: user}
	response_json, _ := json.Marshal(ar)
	http.Error(w, string(response_json), http.StatusOK)
}

func main() {
	fmt.Println("starting...")

	//_ = Sync()
	r := mux.NewRouter()

	r.HandleFunc("/", rootListener).Methods("GET")
	r.HandleFunc("/user", userListener).Methods("GET")
	r.HandleFunc("/user/{user}", userCreateListener).Methods("POST")
	//http.HandleFunc("/user", userListener)
	//http.HandleFunc("/", rootListener)

	fmt.Println("listening on port 80...")
	//http.ListenAndServe(":80", nil)
	//http.Handle("/", r)
	http.ListenAndServe(":80", r)
}
