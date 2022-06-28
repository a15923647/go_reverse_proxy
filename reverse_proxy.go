package main
import (
  "flag"
  "fmt"
  "net/http"
  "net/http/httputil"
  "log"
  "encoding/json"
  "os"
)

const DONT_CHANGE string = ""

var incomingPort = flag.Int("i", 443, "incoming port")
var default_dstHost = flag.String("d", "", "default destination host")
var default_dstPort = flag.Int("p", 443, "default destination port")
var default_dstScheme = flag.String("dsch", "https", "default destination scheme")
var default_dstPath = flag.String("dpath", DONT_CHANGE, "default destination path, \"" + DONT_CHANGE + "\" without double quote means don't change")
var config_file_path = flag.String("cfg", "outRules.json", "path to config file")
var outRules = make(map[string]OutRule) //input domain[path] -> output host

type OutRule struct {
  DstHost string `json:"dstHost"`
  DstPort int `json:"dstPort"`
  DstScheme string `json:"dstScheme"`
  DstPath string `json:"dstPath"`
}

func show_outRule(out_rule OutRule) {
  fmt.Println("DstHost:", out_rule.DstHost)
  fmt.Println("DstPort:", out_rule.DstPort)
  fmt.Println("DstScheme:", out_rule.DstScheme)
  fmt.Println("DstPath:", out_rule.DstPath)
}

func NewOutRule() OutRule {
  out_rule := OutRule{}
  out_rule.DstHost = *default_dstHost
  out_rule.DstPort = *default_dstPort
  out_rule.DstScheme = *default_dstScheme
  out_rule.DstPath = *default_dstPath
  return out_rule
}

func check(err error) {
  if (err != nil) {
    log.Fatal(err.Error())
    panic(err)
  }
}

func outRules_get(host_p, path_p *string) OutRule {
  if r, found := outRules[*host_p + *path_p]; found {
    fmt.Println("use existing rule")
    return r
  } else if r, found := outRules[*host_p]; found {
    fmt.Println("use existing rule")
    return r
  }
  fmt.Println("use default rule")
  r := NewOutRule()
  return r
}

func main() {
  flag.Parse()
  if _, err := os.Stat(*config_file_path); !os.IsNotExist(err) {
    f, err := os.Open(*config_file_path)
    check(err)
    fi, err := f.Stat()
    check(err)
    buf := make([]byte, fi.Size())
    n1, err := f.Read(buf)
    check(err)
    fmt.Printf("Sucessfully load %d bytes from config file %s\n", n1, *config_file_path)
    err = json.Unmarshal(buf, &outRules)
    check(err)
    fmt.Println("Convert successfully,", outRules)
    for k, e := range outRules {
      fmt.Println(k)
      show_outRule(e)
    }
  }
  director := func(request *http.Request) {
    outRule := outRules_get(&request.URL.Host, &request.URL.Path)
    request.URL.Scheme = outRule.DstScheme
    request.URL.Host = fmt.Sprintf("%s:%d", outRule.DstHost, outRule.DstPort)
    if outRule.DstPath != DONT_CHANGE {
      request.URL.Path = outRule.DstPath
    }
  }
  reverseProxy := &httputil.ReverseProxy{Director: director}
  server := http.Server {
    Addr: fmt.Sprintf(":%d", *incomingPort),
    Handler: reverseProxy,
  }
  if err := server.ListenAndServe(); err != nil {
    log.Fatal(err.Error())
  }
}
