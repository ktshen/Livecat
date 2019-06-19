package main

import(
	"log"
	"net/http"
	//"redigogo"
	//"strings"
	// "encoding/json"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("./")))
	// http.HandleFunc("/testweb", Testweb)
	http.HandleFunc("/getreferrer", Getreferrer)
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func Getreferrer(w http.ResponseWriter, r *http.Request){
	 log.Println(r.Header)
	 return

}

// func Testweb(w http.ResponseWriter, r *http.Request) {
//
//
// 	// resmsg := map[string]string{
//  //       "mac" : r.FormValue("FBID"),
//  //       "fbid" : r.FormValue("MAC"),
//  //    }
// 	m := map[string]string{
//               "mac": r.FormValue("MAC"),
//               "fbid": r.FormValue("FBID"),
//        }
//
//     log.Println(r.FormValue("FBID"))
//
//        b, _ := json.Marshal(m)
//
//     resmsg := []byte(`{"mac":`+r.FormValue("MAC")+`,"fbid":`+r.FormValue("FBID")+`}`)
// 	w.Write([]byte("HI~~~"))
// }
