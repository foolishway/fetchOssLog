package main

import (
	"io/ioutil"
	"net/http"
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
	reqUrl := "https://sls.console.aliyun.com/console/logstoreindex/getHistograms.json"
	return doRequest(reqParams, reqUrl)
}
func getOSSlog(reqParams *ReqParams) string {
	reqUrl := "https://sls.console.aliyun.com/console/logstoreindex/getLogs.json"
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
	reqBody += "&secToken=N2G0X6eEAX60FLWiGW3BH"
	reqBody += "&LogStoreName=" + reqParams.LogStoreName
	return reqBody
}
func setReqHeader(req *http.Request)  {
	req.Header.Set("cookie", `_ga=GA1.2.1491651477.1523590977; cnz=t31bEzM1XW0CAbIW+dyAAuvg; login_aliyunid_pk=1648557869747274; login_aliyunid_pks="BG+/SdGZcAj5BAz8paF1/S8m7ureLxikZ/aZUmsxuGSCf8="; aliyun_site=CN; aliyun_choice=CN; cna=aC+5F6hxsTICAW/I1R2PndAf; aliyun_lang=zh; _bl_uid=Fmkh5dLOsgp9en520p1UrIhm01qC; console_base_assets_version=3.15.2; pageSize=20; login_aliyunid_csrf=a232f6cb72b8418fad974a0f021009cc; login_aliyunid="liujiangtao @ tal-weilaichanpin"; login_aliyunid_ticket=Fg8*u*kzxTiisPwXXoqzVZ3RWfkWxThbav7gg6S4F7Efq1S1E2ml6JYlY4q9CyLstMknfiSc2GhOwNcWzj5bYLpKzKZ49O80KpzxYXWJ0WPzFXDzr7rhZ_Dua5Qyv2KMv85szYAdhP4$; login_aliyunid_sc=74u48x24xL7xCj1SQ9*cYL0T_GM6j755mQdVcUFqtlczQTRrnmoKTZaOB*_foDfr8SfFML2NJlw4LHYDaFQMytM2E9m9DF6EPqlFfzvNBasByQTUBTyIiKagGZz1Jondv85szYAdhP4$; JSESSIONID=I1666BD1-M1EIMG6ZZO352BXQRWAX1-AK2PDZDK-1C81; sls_console0=QNqJoyeHxNgfdgZ%2FPdGV1p6WFQT%2B1RDSF34o4q1qVfjv8Ae8XcXG4TAZroyqPOfnOsb9tjLzDXN5vcX2hdeUHZb7G2RmHtf6Td16vX%2FxtX9aJPKPGTB7R3yISHb%2BxRlrn64GVY5qatj7aBgJUPX4DTFOBgyPRjUcEOitQqu%2BcF398wYSU4jVMpal0ZKDbQxLifLpYKovi2k8WVVjDGMRddQ0iWhhWh37yeR6rphKyJq0f%2F93vm1FW1jMoM2WCg3v3Q4tVpKhAYBWxfP%2FrqYti5RKb1lpNa9sUj4MPQdBta2Ru%2BdgsaIVhZd9ZBvVIWcaT1IeLUA2sags%2F2G3%2BbOs%2FQ%3D%3D; FECS-XSRF-TOKEN=6cff7232-4275-46e1-b7ee-f86451f75435; FECS-UMID=%7B%22token%22%3A%22Y138508cf9023fe0a93bd52228358409e%22%2C%22timestamp%22%3A%2273911126565F50415048677B%22%7D; reverse=true; tfstk=clyNBQAKsOBaujDfwAM2CTMiV55OCta0WpuEs7Boh0OOWsXsXl500EnhAEnE74Bns; l=eBOfmXrRQZuU9AdWBO5ZKurza7791QOfhsPzaNbMiInca6ilZFsAoNQqS-ovodtjgtCX_F-ys7AKURFBrMzdgP4vV2cVbn1qnxvO.; isg=BKioDdymY3D3iE4K6N_jI0d0eZa60QzbWLlLWmLalCMWvU4nDeTua3a7tVVNjcSz`)
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