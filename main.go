package main

import (
  "encoding/json"
  "log"
  "net/http"
  "strings"
  "io/ioutil"
  "flag"
  "path/filepath"
  "regexp"

  "github.com/julienschmidt/httprouter"
//  "github.com/k0kubun/pp"
)

var ConfigDirectory = flag.String("c", ".", "Configuration directory (default .)")
type Message struct {
  Status     string
  Messages    []string
}

func main() {
  StartServer()
}

func getLog(logFile string) string {
  logstr, err := ioutil.ReadFile(logFile)
  if err != nil {
    log.Fatal("Error opening config: ", err)
  }
  s := strings.Replace(string(logstr), "\n", "\t\t", -1)
  return strings.Replace(string(s), "kaikai", "\n", -1)
}

func channelList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  err := r.ParseForm()
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    return
  }
  fia, err := ioutil.ReadDir(".")
  reg := regexp.MustCompile(`\.log$`)
  l := make([]string, 0)
  for i := range fia {
    f := fia[i].Name()
    if (!reg.MatchString(f)) {
      continue
    }
    l = append(l, f)
  }
  w.WriteHeader(http.StatusOK)

  type JsonRes struct {
    Status     string
    Channnles    []string
  }
  data := JsonRes {
    Status: "ok",
    Channnles: l,
  }
  w.Header().Set("Content-Type", "application/json; charset=utf-8")
  t, err := json.Marshal(data)
  if err != nil {
    log.Println("Couldn't marshal hook response:", err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
  w.Write(t)
}

func groupLog(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  err := r.ParseForm()
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  str := getLog(filepath.Join(*ConfigDirectory, ps.ByName("group") + ".log"))
  w.WriteHeader(http.StatusOK)
  jsonResp(w, strings.Split(str, "\t\t"))
}

func channelLog(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  err := r.ParseForm()
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  str := getLog(filepath.Join(*ConfigDirectory, ps.ByName("channel") + ".log"))
  w.WriteHeader(http.StatusOK)
  jsonResp(w, strings.Split(str, "\t\t"))
}

func jsonResp(w http.ResponseWriter, msg []string) {
  w.Header().Set("Content-Type", "application/json; charset=utf-8")
  resp := Message {
    Status: "ok",
    Messages: msg[:len(msg)-1],
  }
  r, err := json.Marshal(resp)
  if err != nil {
    log.Println("Couldn't marshal hook response:", err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
  w.Write(r)
}

func StartServer() {
  router := httprouter.New()
  router.GET("/channel_list", channelList)
  router.GET("/channel/:channel", channelLog)
  router.GET("/group/:group", groupLog)

  log.Printf("Starting HTTP server on %d", 3002)
  log.Fatal(http.ListenAndServe(":3002", router))
}
