[![GoDoc](https://godoc.org/github.com/go_cricket?status.svg)](https://godoc.org/github.com/go_cricket)
[![Go Report Card](https://goreportcard.com/badge/github.com/akashshinde/go_cricket)](https://goreportcard.com/report/github.com/akashshinde/go_cricket)

#Usage

###Import 
```
  import "github.com/akashshinde/go_cricket"
```

### Start cricket watcher
```$xslt
event := make(chan gocricket.ResponseEvent)
cricket := gocricket.NewCricketWatcher("IND",event)

```
This will start goroutine to check cricket score for India.

###ResponseEvent
```$xslt
type ResponseEvent struct {
    EventType int
    Response
}
```

you will receive cricket response with EventType

###EventType
```$xslt
	EVENT_NO_CHANGE           = 0
	EVENT_OUT                 = 1
	EVENT_MATCH_STATUS_CHANGE = 2
	EVENT_OVER_CHANGED        = 3
	EVENT_RUN_CHANGE          = 4
```
