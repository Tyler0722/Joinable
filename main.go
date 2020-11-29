package main

import (
    "net/http"
    "net/url"
    "fmt"
    "bytes"
    "io/ioutil"
    "encoding/json"
    "html/template"
    "github.com/dgrijalva/jwt-go"
    "time"
    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"
)

/*
    (Interests) > (Birthday, gender) > (Avatar)
    1 > 2 > 3
*/

type User struct {
    Id uint32 `db:"id"`
    Email string `db:"email"`
    Username string `db:"username"`
    Password string `db:"password"`
    CreatedAt time.Time `db:"created_at"`
}

const clientId = "20440440045-pki18o0uqkvpimarec6prpr1k033oa9k.apps.googleusercontent.com"

var db *sqlx.DB

func googleHandler(res http.ResponseWriter, req *http.Request) {
    values := url.Values{}
    values.Set("scope", "profile email")
    values.Set("response_type", "code")
    values.Set("redirect_uri", "http://localhost:3000/callback")
    values.Set("client_id", "20440440045-pki18o0uqkvpimarec6prpr1k033oa9k.apps.googleusercontent.com")
    values.Set("access_type", "offline")
    url := fmt.Sprintf("%s?%s", "https://accounts.google.com/o/oauth2/v2/auth", values.Encode())
    http.Redirect(res, req, url, 301)
}

func keys(result map[string]interface{}) []string {
    keys := make([]string, 0, len(result))
    for k, v := range result {
        keys = append(keys, fmt.Sprintf("%v-%T", k, v))
    }
    return keys
}

// Layouts (base, landing)
func render(res http.ResponseWriter, data interface{}) error {
    temp, _ := template.ParseFiles("signup.html", "base.html")
    err := temp.ExecuteTemplate(res, "base", data)
    return err
}

func callbackHandler(res http.ResponseWriter, req *http.Request) {
    // result > (access_token, refresh_token, expires_in)
    q := req.URL.Query()
    if _, ok := q["error"]; ok {
        
    }
    bodyBytes := []byte(fmt.Sprintf(`
        {"client_id":"20440440045-pki18o0uqkvpimarec6prpr1k033oa9k.apps.googleusercontent.com",
        "client_secret":"2qw5mxGaiAvcDPCF8_zaEOsz","code":"%s","redirect_uri":"http://localhost:3000/callback","grant_type":"authorization_code"}
    `, q.Get("code")))
    resp, _ := http.Post("https://oauth2.googleapis.com/token", "application/json", bytes.NewBuffer(bodyBytes))
    body, _ := ioutil.ReadAll(resp.Body)
    var result map[string]interface{}
    if err := json.Unmarshal(body, &result); err != nil {
        
    }
    token, _, err := new(jwt.Parser).ParseUnverified(result["id_token"].(string), jwt.MapClaims{})
    if err != nil {
        
    }
    claims := token.Claims.(jwt.MapClaims)
    // expiresAt := time.Now().Add(time.Second * time.Duration(result["expires_in"].(float64)))

    userExists := false
    if db.Get(&userExists, "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 limit 1)", claims["email"]) != nil {

    } 

    if userExists {
        // Log user in
    } 
    // Redirect to signup page with some information already provided by Google
    http.Redirect(res, req, "/signup", 301)
    // db.MustExec(
    //     "INSERT INTO users (email, username, password) VALUES ($1, $2, 'stankhunt42')", 
    //     claims["email"], 
    //     claims["given_name"],
    // )
}

func signupHandler(w http.ResponseWriter, req *http.Request) {
    store, _ := NewPGStoreFromPool(db, []byte("stankhunt42"))
    session, _ := store.Get(req, "2323")
    session.Values["email"] = "example@example.com"
    session.Save(r, w)
    fmt.Printf("%v\n%T", session, session)
    render(w, nil)
}

func main() {
    instance, _ := sqlx.Connect("postgres", "dbname=d4msq5s6uclp49 user=uxcmqmiphaldvn password=c58fc494bdab90c327a8e13c9b49867e1a80a4bb5184acadb03578dc3ec05957 port=5432 host=ec2-54-146-4-66.compute-1.amazonaws.com")
    db = instance
    http.HandleFunc("/google", googleHandler)
    http.HandleFunc("/callback", callbackHandler)
    http.HandleFunc("/signup", signupHandler)
    http.ListenAndServe(":3000", nil)
}

