package data
import (
    "appengine"
    "appengine/user"
    "appengine/datastore"
    "html/template"
    "net/http"
    "time"
)
type Greeting struct {
    Author string
    Content string
    Date time.Time
}
var guestbookTemplate = template.Must(template.New("book").Parse(guestbookTemplateHTML))
const guestbookTemplateHTML = `
<html>
  <body>
    {{range .}}
      {{with .Author}}
        <p><b>{{.}}</b> wrote:</p>
      {{else}}
        <p>An anonymous person wrote:</p>
      {{end}}
      <pre>{{.Content}}</pre>
    {{end}}
    <form action="/data/sign" method="post">
      <div><textarea name="content" rows="3" cols="60"></textarea></div>
      <div><input type="submit" value="Sign Guestbook"></div>
    </form>
<script type="text/javascript"><!--
google_ad_client = "ca-pub-6722923021345262";
/* sample */
google_ad_slot = "7366834398";
google_ad_width = 728;
google_ad_height = 90;
//-->
</script>
<script type="text/javascript"
src="http://pagead2.googlesyndication.com/pagead/show_ads.js">
</script>
  </body>
</html>
`
func init () {
    http.HandleFunc("/data", index)
    http.HandleFunc("/data/sign", sign)
}
func index (w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    q := datastore.NewQuery("Greeting").Order("-Date").Limit(10)
    greetings := make([]Greeting, 0, 10)
    if _, err := q.GetAll(c, &greetings); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    if err := guestbookTemplate.Execute(w, greetings); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}
func sign (w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    g := Greeting {
        Content: r.FormValue("content"),
        Date: time.Now(),
    }
    if u := user.Current(c); u != nil {
        g.Author = u.String()
    }
    _, err := datastore.Put(c, datastore.NewIncompleteKey(c, "Greeting", nil), &g)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, "/data", http.StatusFound)
}
