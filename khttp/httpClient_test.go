package khttp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"testing"
)

func init() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			num := r.URL.Query().Get("num")
			numInt, _ := strconv.Atoi(num)
			fmt.Fprintf(w, "{\"message\": \"%d\"}", numInt+1)
		case "POST":
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				fmt.Fprintf(w, "{\"message\": \"请求参数错误\"}")
				return
			}
			var data map[string]interface{}
			err = json.Unmarshal(body, &data)
			if err != nil {
				fmt.Fprintf(w, "{\"message\": \"请求参数错误\"}")
				return
			}
			num, ok := data["num"].(float64)
			if !ok {
				fmt.Fprintf(w, "{\"message\": \"请求参数错误\"}")
				return
			}
			fmt.Fprintf(w, "{\"message\": \"%d\"}", int(num)+1)
		default:
			fmt.Fprintf(w, "{\"message\": \"This is not a GET or POST request\"}")
		}
	})
	http.ListenAndServe(":8080", nil)
}

func Test_httpClient_Get(t *testing.T) {
	type fields struct {
		client *http.Client
	}
	type args struct {
		url string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Result
	}{
		{
			name: "test_get",
			fields: fields{
				client: &http.Client{},
			},
			args: args{
				url: "http://localhost:8080/?num=1",
			},
			want: Result{
				bytes: []byte(`{"message": "2"}`),
				error: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &httpClient{
				client: tt.fields.client,
			}
			if got := c.Get(tt.args.url); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_httpClient_Post(t *testing.T) {
	type fields struct {
		client *http.Client
	}
	type args struct {
		url  string
		data interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Result
	}{
		{
			name: "test_post",
			fields: fields{
				client: &http.Client{},
			},
			args: args{
				url: "http://localhost:8080/",
				data: map[string]interface{}{
					"num": 1,
				},
			},
			want: Result{
				bytes: []byte(`{"message": "2"}`),
				error: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &httpClient{
				client: tt.fields.client,
			}
			if got := c.Post(tt.args.url, tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Post() = %v, want %v", got, tt.want)
			}
		})
	}
}
