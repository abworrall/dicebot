package verbs

import(
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var(
	ChatHost = "localhost"
	ChatPort = "9898"
	ChatPath = "/run/textgen"
)

// https://github.com/oobabooga/text-generation-webui/blob/main/api-example.py#L20
// The horrid API takes an array where items have mixed types, so can't be represented cleanly in golang. Hack it up.
func NewJsonReq(prompts []string) []byte {
	str := "{\"data\": ["
	str += fmt.Sprintf("%q,", strings.Join(prompts, "\n"))
	str += " 100,"   // [200] max_new_tokens
	str += " true,"  // do_sample
	str += " 0.6,"   // [0.5] temperature
	str += " 0.9,"   // top_p
	str += " 1,"     // typical_p
	str += " 1.05,"  // repetition_penalty
	str += " 1.0,"   // encoder_repetition_penalty
	str += " 0,"     // top_k
	str += " 0,"     // min_length
	str += " 0,"     // no_repeat_ngram_size
	str += " 1,"     // num_beams
	str += " 0,"     // penalty_alpha
	str += " 1,"     // length_penalty
	str += " false," // early_stopping
	str += " -1"    // seed
	str += "]}"
	return []byte(str)
}

type Chatbot struct{
	Context []string
}

func (c *Chatbot)maybeInit() {
	if c.Context == nil {
		c.Context = []string{}
	}
}

func (c *Chatbot)Help() string { return "[reset], [context], [add a new prompt for the bot]" }

func (c *Chatbot)Process(vc VerbContext, args []string) string {
	c.maybeInit()

	if len(args) == 0 { return fmt.Sprintf("[Context: %v]", c.Context) }

	switch args[0] {
	case "context":
		return fmt.Sprintf("Chat context: %v\n", c.Context)
	case "reset":
		c.Context = c.Context[:0]
		return "Chat context reset"
	}

	c.AddContext(strings.Join(args, " "))
	return c.submit()
}


func (c *Chatbot)AddContext(line string) { c.Context = append(c.Context, line) }


func (c *Chatbot)submit() string {
	url := fmt.Sprintf("http://%s:%s%s", ChatHost, ChatPort, ChatPath)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(NewJsonReq(c.Context)))
	if err != nil {
		log.Printf("chat POST %s, err: %v\n", url, err)
		return "err"
	}
	defer resp.Body.Close()
	body,_ := ioutil.ReadAll(resp.Body)

	var jsonMap map[string]interface{}
	if err := json.Unmarshal([]byte(body), &jsonMap); err != nil { 
		log.Printf("chat, unmarshal err: %v\n", err)
		return "json err"
	} else if v := jsonMap["data"]; v == nil {
		log.Printf("chat, no 'data'\n")
		return "json data err"
	} else if vArray, ok := jsonMap["data"].([]interface{}); !ok {
		log.Printf("chat, empty 'data'\n")
		return "json data[] err"
	} else if vElem, ok := vArray[0].(string); !ok {
		log.Printf("chat, data[0] not string - %#v\n", vArray)
		return "json data[0] err"
	} else {
		// Annoyingly, prev prompts are at the beginning of vElem, so reset it.
		c.Context = c.Context[:0]
		c.AddContext(vElem)
		return vElem
	}

	return "??"
}

/*

REQUEST application/json POST BODY
{"data": ["prompts go here", 200, true, 0.5, 0.9, 1, 1.05, 1.0, 0, 0, 0, 1, 0, 1, false, -1]}


HTTP/1.1 200 OK
Content-Length: 1992
Content-Type: application/json
Date: Thu, 06 Apr 2023 18:36:17 GMT
Server: uvicorn

{"data":["response is here"],"is_generating":true,"duration":8.343884706497192,"average_duration":8.523337721824646}

*/
