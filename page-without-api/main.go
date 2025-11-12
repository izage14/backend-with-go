package main

import (
    "fmt"
    "html/template"
    "net/http"
)

var users = []string{"Alice", "Bob"}

func main() {
    http.HandleFunc("/", homeHandler)
    http.HandleFunc("/add", addHandler)

    fmt.Println("Server running at http://localhost:8080")
    http.ListenAndServe(":8080", nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
    tmpl := `
    <html>
      <body>
        <h1>User List</h1>
        <ul>
          {{range .}}
            <li>{{.}}</li>
          {{end}}
        </ul>
        <form action="/add" method="POST">
          <input name="name" placeholder="New user" />
          <button type="submit">Add</button>
        </form>
      </body>
    </html>`
    t := template.Must(template.New("page").Parse(tmpl))
    t.Execute(w, users)
}

func addHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        name := r.FormValue("name")
        if name != "" {
            users = append(users, name)
        }
    }
    http.Redirect(w, r, "/", http.StatusSeeOther) // reloads the entire page
}