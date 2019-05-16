package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"domliang.com/empty/lib/empty_apns"
)

func main() {
	for true {
		log.Print(`start request`)
		r, err := crawWithLogin()
		if err != nil {
			log.Print(err)
			empty_apns.SendNotify(`has error in request`)
		}
		if !r {
			empty_apns.SendNotify(`has AVAILABLE Seat`)
		} else {
			log.Print("no AVAILABLE Seat")
		}
		log.Print("sleep 15 sec")
		time.Sleep(time.Second * 15)
	}

}

func crawWithLogin() (avaliable bool, err error) {
	refererUrl := "https://www.naati.com.au/MyNaati"
	loginUrl := "https://www.naati.com.au/MyNaati/Account/LogOn"
	testPage := "https://www.naati.com.au/MyNaati/MyTests/AvailableTestSessions?CredentialRequestId=137220&CredentialApplicationId=161565"
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    160 * time.Second,
		DisableCompression: true,
	}
	var checkClient = &http.Client{}
	client := &http.Client{
		Transport: tr,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}}
	form := url.Values{}
	// form.Add("UserName", "soukeiiku@gmail.com")
	// form.Add("Password", "Mynaati5920!")
	form.Add("UserName", "yejunchen6381@gmail.com")
	form.Add("Password", "5201314yjcYJC!")
	req, err := http.NewRequest("POST", loginUrl, strings.NewReader(form.Encode()))
	if err != nil {
		return false, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", refererUrl)
	req.Header.Add("User-Agent", `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/54.0.2840.71 Safari/537.36`)
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	log.Print(resp.StatusCode)
	if resp.StatusCode == 302 {
		cookieJar, _ := cookiejar.New(nil)
		cookieJar.SetCookies(resp.Request.URL, resp.Cookies())
		checkClient = &http.Client{
			Transport: tr,
			Jar:       cookieJar,
		}

		req, err := http.NewRequest("GET", testPage, nil)
		if err != nil {
			log.Print(err)
			return false, err
		}
		req.Header.Add("User-Agent", `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/54.0.2840.71 Safari/537.36`)
		checkResp, err := checkClient.Do(req)
		if err != nil {
			return false, nil
		}
		checkBodyBytes, err := ioutil.ReadAll(checkResp.Body)
		if err != nil {
			return false, err
		}
		checkBodyString := string(checkBodyBytes)
		matched, err := regexp.MatchString(`(?ism)^.*?NO\sOTHER\sAVAILABLE\sTEST\sSESSIONS.*?$`, checkBodyString)
		if err != nil {
			return false, err
		}
		if !matched {
			return true, nil
		} else {
			return false, nil
		}
	} else {
		return false, nil
	}

}

func craw() (avaliable bool, err error) {
	timeUnix := makeTimestamp()
	apiUrl := "https://www.naati.com.au/MyNaati/credentialapplication/testsessions?ApplicationId=164187&Token=202780594&_=" + strconv.FormatInt(timeUnix, 10)
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    160 * time.Second,
		DisableCompression: true,
	}

	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return false, err
	}
	req.Header.Add("User-Agent", `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/54.0.2840.71 Safari/537.36`)

	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Print(err)
			return false, err
		}
		bodyString := string(bodyBytes)
		log.Print(bodyString)
		log.Print(len(bodyString))
		if len(bodyString) > 2 {
			return true, nil
		} else {
			return false, nil
		}
	} else {
		return false, nil
	}
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
