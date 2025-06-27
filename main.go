package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
	"os"
	"strings"
)

type Example struct {
	Summary string      `yaml:"summary"`
	Value   interface{} `yaml:"value"`
}

type ResponseContent struct {
	Examples map[string]Example `yaml:"examples"`
}

type Response struct {
	Description string                     `yaml:"description"`
	Content     map[string]ResponseContent `yaml:"content"`
}

type Operation struct {
	Responses map[string]Response `yaml:"responses"`
}

type PathItem struct {
	Get    *Operation `yaml:"get"`
	Post   *Operation `yaml:"post"`
	Put    *Operation `yaml:"put"`
	Delete *Operation `yaml:"delete"`
}

type Components struct {
	Responses map[string]Response `yaml:"responses"`
}

type OpenAPI struct {
	Paths      map[string]PathItem `yaml:"paths"`
	Components Components          `yaml:"components"`
}

type ResponseExample struct {
	StatusCode  int
	Payload     []byte
	ContentType string
}

var (
	responseExamples map[string]ResponseExample
	router           *mux.Router
)

func loadOpenAPI(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	var spec OpenAPI
	err = yaml.Unmarshal(data, &spec)
	if err != nil {
		log.Fatalf("Failed to parse YAML: %v", err)
	}

	responseExamples = make(map[string]ResponseExample)

	extractExamples := func(responses map[string]Response) {
		for statusCode, resp := range responses {
			code := mapNameToStatus(statusCode)
			for contentType, content := range resp.Content {
				for key, example := range content.Examples {
					payload, _ := json.MarshalIndent(example.Value, "", "  ")
					responseExamples[key] = ResponseExample{
						StatusCode:  code,
						Payload:     payload,
						ContentType: contentType,
					}
				}
			}
		}
	}

	extractExamples(spec.Components.Responses)

	router = mux.NewRouter()
	for path, item := range spec.Paths {
		for method, op := range map[string]*Operation{"GET": item.Get, "POST": item.Post, "PUT": item.Put, "DELETE": item.Delete} {
			if op != nil {
				extractExamples(op.Responses)
				r := router.HandleFunc(path, genericHandler).Methods(method)
				_ = r
				log.Printf("Registered %s %s", method, path)
			}
		}
	}

	log.Printf("Loaded %d response examples from OpenAPI spec", len(responseExamples))
}

func mapNameToStatus(name string) int {
	switch strings.ToLower(name) {
	case "200", "ok":
		return http.StatusOK
	case "201":
		return http.StatusCreated
	case "202", "inprocessing":
		return http.StatusAccepted
	case "400", "badrequest":
		return http.StatusBadRequest
	case "401", "unauthorized":
		return http.StatusUnauthorized
	case "403", "forbidden":
		return http.StatusForbidden
	case "404", "notfound":
		return http.StatusNotFound
	case "422", "rejected":
		return http.StatusUnprocessableEntity
	case "500", "internalerror", "internalservererror":
		return http.StatusInternalServerError
	default:
		return http.StatusOK
	}
}

func genericHandler(w http.ResponseWriter, r *http.Request) {
	var key string
	if r.Method == http.MethodGet {
		svars := mux.Vars(r)
		if val, ok := svars["contract-process-id"]; ok {
			key = val
		} else if val, ok := svars["contract-change-id"]; ok {
			key = val
		}
	} else {
		var bodyMap map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&bodyMap); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		if k, ok := bodyMap["_exampleKey"].(string); ok {
			key = k
		}
	}
	if strings.Contains(key, "received") {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	ex, ok := responseExamples[key]
	if !ok {
		if r.Method == http.MethodGet {
			errResp := map[string]interface{}{
				"headline": "Nicht gefunden",
				"status":   404,
				"detail":   "Ungültiger Vorgang '" + key + "'.",
			}
			respBytes, _ := json.Marshal(errResp)
			w.Header().Set("Content-Type", "application/problem+json")
			w.WriteHeader(http.StatusNotFound)
			_, err := w.Write(respBytes)
			if err != nil {
				http.Error(w, "unable to write response", http.StatusInternalServerError)
			}
		} else {
			http.Error(w, fmt.Sprintf("example '%s' not found", key), http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", ex.ContentType)
	w.WriteHeader(ex.StatusCode)
	_, err := w.Write(ex.Payload)
	if err != nil {
		http.Error(w, "unable to write response", http.StatusInternalServerError)
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <openapi.yaml>")
		os.Exit(1)
	}
	loadOpenAPI(os.Args[1])
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{
			"Agent-Id",
			"API-Version",
			"Authorization",
			"Consumer-Job-Id",
			"Content-Type",
			"User-Agent",
			"X-Request-ID",
		}),
	)(router)
	log.Println("Mock server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", corsHandler))
}
