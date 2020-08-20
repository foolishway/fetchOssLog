package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func main() {
	http.HandleFunc("/getOSSlog/", accessWapper(handleGetOSSlog))
	http.HandleFunc("/getHistograms/", accessWapper(handleGetHistograms))
	http.ListenAndServe(":9001", nil)
}
func handleGetOSSlog(w http.ResponseWriter, r *http.Request) {
	reqParams := resolveReqParam(r)
	logContent := getOSSlog(reqParams)
	w.Write([]byte(logContent))
}

func handleGetHistograms(w http.ResponseWriter, r *http.Request) {
	reqParams := resolveReqParam(r)
	logContent := getHistograms(reqParams)
	w.Write([]byte(logContent))
}
func accessWapper(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		//w.WriteHeader(http.StatusOK)
		w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
		handler(w, r)
	}
}
type ReqParams struct {
	Query string
	From string
	To string
	Page string
	LogStoreName string
}
func resolveReqParam(r *http.Request) *ReqParams  {
	var reqParams ReqParams
	qr := r.URL.Query()
	query := qr["query"][0]
	from := qr["from"][0]
	to := qr["to"][0]
	page := qr["page"][0]
	logStoreName := qr["logstorename"][0]
	reqParams = ReqParams{Query: query, From: from, To: to, Page: page, LogStoreName: logStoreName}
	return &reqParams
}
func getHistograms(reqParams *ReqParams) string {
	reqUrl := ""
	return doRequest(reqParams, reqUrl)
}
func getOSSlog(reqParams *ReqParams) string {
	reqUrl := ""
	return doRequest(reqParams, reqUrl)
}
func getReqBody(reqParams *ReqParams) string {
	reqBody := "ProjectName=online-log-xkt"
	reqBody += "&from=" + reqParams.From
	reqBody += "&query=" + reqParams.Query
	reqBody += "&to=" + reqParams.To
	reqBody += "&Page=" + reqParams.Page
	reqBody += "&Size=20"
	reqBody += "&Reverse=true"
	reqBody += "&secToken=" + getAccessAuth("token")
	reqBody += "&LogStoreName=" + reqParams.LogStoreName
	return reqBody
}
func setReqHeader(req *http.Request)  {
	req.Header.Set("cookie", getAccessAuth("cookie"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
}
func doRequest(reqParams *ReqParams, reqUrl string) string {
	reqBody := getReqBody(reqParams)
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, reqUrl,
		strings.NewReader(reqBody))
	if err != nil {
		panic(err)
	}
	setReqHeader(req)

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return string(body)
}

func getAccessAuth(flag string) string {
	var configFile *os.File
	var err error
	switch flag {
	case "cookie":
		configFile, err = os.Open("./cookie")
	case "token":
		configFile, err = os.Open("./token")
	}
	if configFile == nil {
		panic("Can not get cookie and token")
	}
	defer configFile.Close()
	if err != nil {
		panic("Get token or cookie error")
	}
	auth, err := ioutil.ReadAll(configFile)

	if err != nil {
		panic("Read token or cookie error")
	}
	return string(auth)
}