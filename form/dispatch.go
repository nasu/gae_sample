package form

import (
    "fmt"
    "net/http"
    "html/template"
)

const guestbookForm = `
<html>
<body>
<form action="/form/res" method="post">
  <div><textarea name="content"></textarea></div>
  <div><input type="submit"></div>
</form>
</body>
</html>
`
const signTemplateHTML = `
<html>
  <body>
    <p>You wrote:</p>
    <pre>{{.}}</pre>
  </body>
</html>
`
var signTemplate = template.Must(template.New("form_res").Parse(signTemplateHTML))

func init() {
    http.HandleFunc("/form", index)
    http.HandleFunc("/form/res", form_res)
}

func form_res(w http.ResponseWriter, r *http.Request) {
    err := signTemplate.Execute(w, r.FormValue("content"))
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func index(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, guestbookForm)
}
