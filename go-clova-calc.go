// clova extension for aws lambda
// build on windows
// > set GOOS=linux
// > set GOARCH=amd64
// > set CGO_ENABLED=0
// > go build -o calcurate go-clova-calc.go & build-lambda-zip.exe -o calcurate.zip calcurate

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type ClovaRequest struct {
	// Version string `json:"version"`
	// Session struct {
	// 	New               bool `json:"new"`
	// 	SessionAttributes struct {
	// 		X int `json:"x"`
	// 		Y int `json:"y"`
	// 		Operator string `json:"operator"`
	// 	} `json:"sessionAttributes"`
	// 	SessionID string `json:"sessionId"`
	// 	User      struct {
	// 		UserID      string `json:"userId"`
	// 		AccessToken string `json:"accessToken"`
	// 	} `json:"user"`
	// } `json:"session"`
	// Context struct {
	// 	System struct {
	// 		Application struct {
	// 			ApplicationID string `json:"applicationId"`
	// 		} `json:"application"`
	// 		User struct {
	// 			UserID      string `json:"userId"`
	// 			AccessToken string `json:"accessToken"`
	// 		} `json:"user"`
	// 		Device struct {
	// 			DeviceID string `json:"deviceId"`
	// 			Display  struct {
	// 				Size         string `json:"size"`
	// 				Orientation  string `json:"orientation"`
	// 				Dpi          int    `json:"dpi"`
	// 				ContentLayer struct {
	// 					Width  int `json:"width"`
	// 					Height int `json:"height"`
	// 				} `json:"contentLayer"`
	// 			} `json:"display"`
	// 		} `json:"device"`
	// 	} `json:"System"`
	// } `json:"context"`
	Request struct {
		Type   string `json:"type"`
		Intent struct {
			Name  string `json:"name"`
			Slots struct {
				X struct {
					Name  string `json:"name"`
					Value string `json:"value"`
				} `json:"x`
				Y struct {
					Name  string `json:"name"`
					Value string `json:"value"`
				} `json:"y"`
				Operator struct {
					Name  string `json:"name"`
					Value string `json:"value"`
				} `json:"operator"`
			} `json:"slots"`
		} `json:"intent"`
	} `json:"request"`
}

type MyResponse struct {
	StatusCode int `json:"statusCode"`
	Headers    struct {
		ContentType string `json:"Content-Type"`
	} `json:"headers"`
	Body struct {
		Version           string   `json:"version"`
		SessionAttributes struct{} `json:"sessionAttributes"`
		Response          struct {
			OutputSpeech struct {
				Type   string `json:"type"`
				Values struct {
					Type  string `json:"type"`
					Lang  string `json:"lang"`
					Value string `json:"value"`
				} `json:"values"`
			} `json:"outputSpeech"`
			Card             struct{}   `json:"card"`
			Directives       []struct{} `json:"directives"`
			Reprompt         struct{}   `json:"reprompt"`
			ShouldEndSession bool       `json:"shouldEndSession"`
		} `json:"response"`
	} `json:"body"`
}

func hello(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	reqJsonBytes := ([]byte)(request.Body)
	clovaRequest := new(ClovaRequest)

	var err error
	if err = json.Unmarshal(reqJsonBytes, clovaRequest); err != nil {
		log.Println("[ERROR0]", err)
	}

	response := MyResponse{}
	response.Body.Version = "1.0"
	response.Body.Response.ShouldEndSession = false
	response.Body.Response.OutputSpeech.Type = "SimpleSpeech"
	response.Body.Response.OutputSpeech.Values.Type = "PlainText"
	response.Body.Response.OutputSpeech.Values.Lang = "ja"

	var text string
	var result int
	switch clovaRequest.Request.Type {
	case "LaunchRequest":
		text = "たし算と引き算ができます。1たす1は？または、3ひく1は？のように話してみてください"

	case "IntentRequest":
		var x, y int
		if x, err = strconv.Atoi(clovaRequest.Request.Intent.Slots.X.Value); err != nil {
			log.Println("[ERROR1]", err)
			text = "すみません。理解できませんでした。"
			break
		}

		if y, err = strconv.Atoi(clovaRequest.Request.Intent.Slots.Y.Value); err != nil {
			log.Println("[ERROR2]", err)
			text = "すみません。理解できませんでした。"
			break
		}
		switch clovaRequest.Request.Intent.Slots.Operator.Value {
		case "たす", "プラス", "ぷらす", "足す", "足して", "たして":
			result = x + y
		case "まいなす", "マイナス", "ひいて", "引いて", "引く", "ひく":
			result = x - y
		case "かける", "かけて":
			result = x * y
		case "わる", "割る", "割って", "わって":
			result = x / y
		default:
			log.Println("[ERROR3] Intent request operator value parse error.")
			text = "すみません。理解できませんでした。"
			break
		}

		text = fmt.Sprintf("%d%s%dは、、、、%dです！！", x, clovaRequest.Request.Intent.Slots.Operator.Value, y, result)
	default:
		log.Println("[ERROR3]", "Intent request parse error.")
		text = "すみません。理解できませんでした。"
		break
	}

	response.Body.Response.OutputSpeech.Values.Value = text

	resJsonBytes, _ := json.Marshal(response.Body)
	return events.APIGatewayProxyResponse{
		Headers:    map[string]string{"Content-Type": "application/json;charset=UTF-8"},
		StatusCode: 200,
		Body:       string(resJsonBytes),
	}, nil
}

func main() {
	log.SetFlags(log.Lshortfile)
	lambda.Start(hello)
}
