package main

import (
  "encoding/json"
  "fmt"
  "io/ioutil"
  "net/http"
  "regexp"
  "strconv"
  "strings"
  "time"

  "github.com/labstack/echo"
  "github.com/labstack/echo/middleware"
)

type PeopleNowPresent struct {
  Value int `json:"value"`
}

type SpaceDescriptor struct {
  API      string `json:"api"`
  Space    string `json:"space"`
  Logo     string `json:"logo"`
  URL      string `json:"url"`
  Location struct {
    Address string  `json:"address"`
    Lon     float64 `json:"lon"`
    Lat     float64 `json:"lat"`
  } `json:"location"`
  State struct {
    Open       bool  `json:"open"`
    Lastchange int64 `json:"lastchange"`
  } `json:"state"`
  Contact struct {
    Email      string `json:"email"`
    Irc        string `json:"irc"`
    Ml         string `json:"ml"`
    Twitter    string `json:"twitter"`
    Facebook   string `json:"facebook"`
    Foursquare string `json:"foursquare"`
  } `json:"contact"`
  Sensors struct {
    PeopleNowPresent []PeopleNowPresent `json:"people_now_present"`
  } `json:"sensors"`
  IssueReportChannels []string `json:"issue_report_channels"`
  Projects            []string `json:"projects"`
  Cache               struct {
    Schedule string `json:"schedule"`
  } `json:"cache"`
}

type Status struct {
  Open             bool  `json:"open"`
  PeopleNowPresent int   `json:"people_now_present"`
  Lastchange       int64 `json:"lastchange"`
}

type HackerspaceEvents struct {
  Title string `json:"title"`
  Date  string `json:"date"`
  Begin string `json:"begin"`
  End   string `json:"end"`
}

type DiscourseApi struct {
  TopicList struct {
    Topics []struct {
      Title string `json:"title"`
    } `json:"topics"`
  } `json:"topic_list"`
}

func check(err error) {
  if err != nil {
    fmt.Println(err)
    // panic(err)
  }
}

// Read file from file system
func readFile(fileName string) []byte {
  file, err := ioutil.ReadFile(fileName)
  check(err)
  return file
}

// Read json file
func getJSONFile(fileName string) (jsonFile SpaceDescriptor) {
  tmp := readFile(fileName)
  err := json.Unmarshal(tmp, &jsonFile)
  check(err)
  return
}

// Get number of people in space
func peopleInSpace() int {
  req, err := http.NewRequest("GET", "https://lambdaspace.gr/hackers.txt", strings.NewReader(""))
  check(err)
  resp, err := http.DefaultClient.Do(req)
  check(err)
  defer resp.Body.Close()
  b, _ := ioutil.ReadAll(resp.Body)
  data, _ := strconv.Atoi(string(b))
  return data
}

// Update the current status of the space
func updateStatus(spaceStatus *bool, peoplePresent *int, lastChange *int64) {
  previewPeoplePresent := *peoplePresent
  *peoplePresent = peopleInSpace()
  *spaceStatus = *peoplePresent > 0
  if previewPeoplePresent != *peoplePresent {
    *lastChange = time.Now().Unix()
  }
}

func scheduledEvents() []HackerspaceEvents {
  var dat DiscourseApi
  ret := []HackerspaceEvents{}
  req, err := http.NewRequest("GET", "https://community.lambdaspace.gr/c/5/l/latest.json", strings.NewReader(""))
  check(err)
  resp, err := http.DefaultClient.Do(req)
  check(err)
  defer resp.Body.Close()
  b, _ := ioutil.ReadAll(resp.Body)
  err = json.Unmarshal(b, &dat)
  check(err)
  for _, element := range dat.TopicList.Topics {
    event := HackerspaceEvents{}
    fields := strings.Fields(element.Title)
    ryear, err := regexp.Compile(`^\d\d\/\d\d\/\d\d\d\d`)
    check(err)
    rhour, err := regexp.Compile(`^\d\d:\d\d`)
    check(err)
    if ryear.MatchString(fields[0]) {
      event.Date = fields[0]
      if rhour.MatchString(fields[1]) {
        event.Begin = fields[1]
        if fields[2] == "-" && rhour.MatchString(fields[3]) {
          event.End = fields[3]
          event.Title = strings.Join(fields[4:], " ")
        } else {
          event.End = ""
          event.Title = strings.Join(fields[2:], " ")
        }
        ret = append(ret, event)
      }
    }
  }
  return ret
}

func main() {
  e := echo.New()

  // Echo middlewares
  e.Use(middleware.Logger())
  e.Use(middleware.Recover())
  e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
    AllowMethods: []string{echo.GET},
  }))

  spaceDescriptor := getJSONFile("./LambdaSpaceAPI.json")
  peoplePresent := &spaceDescriptor.Sensors.PeopleNowPresent[0].Value
  spaceStatus := &spaceDescriptor.State.Open
  lastChange := &spaceDescriptor.State.Lastchange

  // Route compatible with SpaceApi spec
  e.GET("/api/v2.0/SpaceAPI", func(c echo.Context) error {
    updateStatus(spaceStatus, peoplePresent, lastChange)
    return c.JSON(http.StatusOK, spaceDescriptor)
  })

  // Route to serve space status
  e.GET("/api/v2.0/status", func(c echo.Context) error {
    updateStatus(spaceStatus, peoplePresent, lastChange)
    return c.JSON(http.StatusOK, Status{*spaceStatus, *peoplePresent, *lastChange})
  })

  // Route to serve events
  e.GET("/api/v2.0/events", func(c echo.Context) error {
    return c.JSON(http.StatusOK, scheduledEvents())
  })
  e.Logger.Fatal(e.Start(":1323"))
}
