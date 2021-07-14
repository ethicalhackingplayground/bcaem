package bugcrowd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/tidwall/gjson"
)

type ScopeElement struct {
	Target string
}

type ProgramData struct {
	Url        string
	InScope    []ScopeElement
	OutOfScope []ScopeElement
}

var client = &http.Client{}

const (
	USER_AGENT          = "Mozilla/5.0 (X11; Linux x86_64; rv:82.0) Gecko/20100101 Firefox/82.0"
	BUGCROWD_LOGIN_PAGE = "https://bugcrowd.com/user/sign_in"
)

func Login(email string, password string) string {
	// Send GET to https://bugcrowd.com/user/sign_in
	// Get _crowdcontrol_session cookie
	// Get <meta name="csrf-token" content="Da...ktOQ==" />
	// Still under development

	req, err := http.NewRequest("GET", BUGCROWD_LOGIN_PAGE, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", USER_AGENT)
	client := &http.Client{
		// We don't need to follow redirects
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	crowdControlSession := ""
	csrfToken := ""
	for _, cookie := range resp.Header["Set-Cookie"] {
		if strings.HasPrefix(cookie, "_crowdcontrol_session") {
			crowdControlSession = strings.Split(strings.Split(cookie, ";")[0], "=")[1]
			break
		}
	}

	if crowdControlSession == "" {
		log.Fatal("Failed to get cookie. Something might have changed")
	}

	// Now we need to get the csrf-token...HTML parsing here we go
	body, _ := ioutil.ReadAll(resp.Body)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))

	if err != nil {
		log.Fatal("Failed to parse login response")
	}

	doc.Find("meta").Each(func(index int, s *goquery.Selection) {
		name, _ := s.Attr("name")
		if name == "csrf-token" {
			csrfToken, _ = s.Attr("content")
			//fmt.Println("TOKEN: ", url.QueryEscape(content))
		}
	})

	if csrfToken == "" {
		log.Fatal("Failed to get the CSRF token. Something might have changed")
	}

	// Now send the POST request
	req2, err := http.NewRequest("POST", BUGCROWD_LOGIN_PAGE, bytes.NewBuffer([]byte("utf8=%E2%9C%93&authenticity_token="+url.QueryEscape(csrfToken)+"&user%5Bredirect_to%5D=&user%5Bemail%5D="+url.QueryEscape(email)+"&user%5Bpassword%5D="+url.QueryEscape(password)+"&commit=Log+in")))
	if err != nil {
		log.Fatal(err)
	}

	req2.Header.Set("User-Agent", USER_AGENT)
	req2.Header.Set("Cookie", "_crowdcontrol_session="+crowdControlSession)
	resp2, err := client.Do(req2)
	if err != nil {
		panic(err)
	}
	defer resp2.Body.Close()

	sessionToken := ""
	for _, cookie := range resp2.Header["Set-Cookie"] {
		if strings.HasPrefix(cookie, "_crowdcontrol_session") {
			sessionToken = cookie
			break
		}
	}

	if resp2.StatusCode != 302 {
		log.Fatal("Login failed", resp2.StatusCode)
	}

	return sessionToken
}

func GetProgramHandles(sessionToken string, bbpOnly bool, pvtOnly bool) []string {
	allProgramsCount := 0
	currentProgramIndex := 0
	listEndpointURL := "https://bugcrowd.com/programs.json?"
	if pvtOnly {
		listEndpointURL = listEndpointURL + "accepted_invite[]=true&"
	}
	if bbpOnly {
		listEndpointURL = listEndpointURL + "vdp[]=false&"
	}
	listEndpointURL = listEndpointURL + "hidden[]=false&sort[]=invited-desc&sort[]=promoted-desc&offset[]="
	paths := []string{}

	for {
		req, err := http.NewRequest("GET", listEndpointURL+strconv.Itoa(currentProgramIndex), nil)
		if err != nil {
			log.Fatal(err)
		}

		req.Header.Set("Cookie", "_crowdcontrol_session="+sessionToken)
		req.Header.Set("User-Agent", USER_AGENT)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)

		if allProgramsCount == 0 {
			allProgramsCount = int(gjson.Get(string(body), "meta.totalHits").Int())
		}

		chunkData := gjson.Get(string(body), "programs.#.program_url")
		for i := 0; i < len(chunkData.Array()); i++ {
			paths = append(paths, chunkData.Array()[i].Str)
		}
		currentProgramIndex += 25

		if allProgramsCount <= currentProgramIndex {
			break
		}
	}

	return paths
}

func GetProgramScope(handle string, categories string, token string) (pData ProgramData) {
	pData.Url = "https://bugcrowd.com" + handle

	req, err := http.NewRequest("GET", pData.Url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Cookie", "_crowdcontrol_session="+token)
	req.Header.Set("User-Agent", USER_AGENT)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	// Times @arcwhite broke our HTML parsing: #1 and counting :D

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		fmt.Println("No url found")
		log.Fatal(err)
	}

	doc.Find(".react-component-researcher-target-groups").Each(func(index int, s *goquery.Selection) {
		json, ok := s.Attr("data-react-props")
		if !ok {
			log.Fatal("Error parsing HTML of ", pData.Url)
		}

		for _, scopeElement := range gjson.Get(string(json), "groups.#(in_scope==true).targets").Array() {
			var currentTarget struct {
				line       string
				categories []string
				target     string
				tags       []string
				name       string
			}
			for _, x := range scopeElement.Map()["target"].Map() {
				for _, y := range x.Array() {
					currentTarget.tags = append(currentTarget.tags, y.Map()["name"].Str)
					currentTarget.line = scopeElement.Map()["name"].Str
					if y.Map()["name"].Str == "Adobe Experience Manager" {
						pData.InScope = append(pData.InScope, ScopeElement{Target: currentTarget.line})
					}
					if y.Map()["name"].Str == "Website Testing" && strings.HasPrefix(currentTarget.line, "https") || strings.HasPrefix(currentTarget.line, "http") {

						req, _ := http.NewRequest("GET", currentTarget.line, nil)
						res, _ := client.Do(req)

						rbody, _ := ioutil.ReadAll(res.Body)
						if (strings.Contains(string(rbody), "/content/dam/")) ||
							(strings.Contains(string(rbody), "/libs/settings/")) ||
							(strings.Contains(string(rbody), "/libs/granite/")) ||
							(strings.Contains(string(rbody), "/etc/libs/")) {
							pData.InScope = append(pData.InScope, ScopeElement{Target: currentTarget.line})
						}

					}
				}
			}
		}
	})

	if len(pData.InScope) == 0 {
		return
	}

	return pData
}

func GetCategories(input string) []string {
	categories := map[string][]string{
		"url":      {"Website Testing"},
		"api":      {"API Testing"},
		"mobile":   {"Android", "iOS"},
		"android":  {"Android"},
		"apple":    {"iOS"},
		"other":    {"Other"},
		"hardware": {"Hardware Testing"},
		"all":      {},
	}

	selectedCategory, ok := categories[strings.ToLower(input)]
	if !ok {
		log.Fatal("Invalid category")
	}
	return selectedCategory
}

func GetAllProgramsScope(token string, bbpOnly bool, pvtOnly bool, categories string, concurrency int) (programs []ProgramData) {
	programHandles := GetProgramHandles(token, bbpOnly, pvtOnly)

	handles := make(chan string, concurrency)
	processGroup := new(sync.WaitGroup)
	processGroup.Add(concurrency)

	for i := 0; i < concurrency; i++ {
		go func() {
			for {
				handle := <-handles

				if handle == "" {
					break
				}

				programs = append(programs, GetProgramScope(handle, categories, token))
			}
			processGroup.Done()
		}()
	}

	for _, handle := range programHandles {
		handles <- handle
	}

	close(handles)
	processGroup.Wait()
	return programs
}

// PrintAllScope prints to stdout all scope elements of all targets
func PrintAllScope(token string, bbpOnly bool, pvtOnly bool, categories string, delimiter string, concurrency int) {
	programs := GetAllProgramsScope(token, bbpOnly, pvtOnly, categories, concurrency)
	for _, pData := range programs {
		for _, scopeElement := range pData.InScope {
			fmt.Println(scopeElement.Target)
		}
	}
}

/*
// ListPrograms prints a list of available programs
func ListPrograms(token string, bbpOnly bool, pvtOnly bool) {
	programPaths := GetProgramPagePaths(token, bbpOnly, pvtOnly)
	for _, path := range programPaths {
		fmt.Println("https://bugcrowd.com" + path)
	}
}*/
