package main

import (
	"fmt"
	"log"
	"time"
	"math"
	"errors"
	"strconv"
	"encoding/json"
	"net/http"
	"github.com/patrickmn/go-cache"
)

// cacheTTL is the duration of time to hold the key in cache
const cacheTTL = 60 * time.Second

// names for possible actions
const actionAdd = "add"
const actionSubtract = "subtract"
const actionMultiply = "multiply"
const actionDivide = "divide"

// ResponseData contains the fields to store response
type ResponseData struct {
	Action string
	X,Y,Answer float64
	Cached bool
}

var resp ResponseData

// ResponseData helper function to generate key to be used to store cache
func (r ResponseData) GenerateKey() string {
    return fmt.Sprintf("%v", r.X) + ":" + r.Action + ":" + fmt.Sprintf("%v", r.Y)
}

// Round result to 2 decimals
func RoundResult(result float64) (float64) {
	return math.Round(result*100)/100
}

// Add is our function that sums two numbers
func Add(x,y float64) (float64) {
	return RoundResult(x + y)
}

// Subtract is our function that deducts two numbers
func Subtract(x,y float64) (float64) {
	return RoundResult(x - y)
}

// Multiply helper function
func Multiply(x,y float64) (float64) {
	return RoundResult(x * y)
}

// Divide helper function
func Divide(x,y float64) (float64) {
	return RoundResult(x / y)
}

// -------------------- Cache related ----------------------

// Cache interface
type CacheItf interface {
	Set(key string, data interface{}, expiration time.Duration) error
	Get(key string) ([]byte, error)
}

var cacheDatabase CacheItf

type AppCache struct {
	client *cache.Cache
}

func InitCache() {
	cacheDatabase = &AppCache{
		client: cache.New(5*time.Minute, 10*time.Minute),
	}
}

func (r *AppCache) Get(key string) ([]byte, error) {
	res, exist := r.client.Get(key)
	if !exist {
		return nil, nil
	}

	resByte, ok := res.([]byte)
	if !ok {
		return nil, errors.New("Format is not arr of bytes")
	}

	return resByte, nil
}

func (r *AppCache) Set(key string, data interface{}, expiration time.Duration) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	r.client.Set(key, b, expiration)
	return nil
}

// -------------------- main ----------------------

func main () {
	InitCache()

	http.HandleFunc("/add",addHandler)
	http.HandleFunc("/subtract",substractHandler)
	http.HandleFunc("/multiply",multiplyHandler)
	http.HandleFunc("/divide",divideHandler)
	http.ListenAndServe(":8080",nil)
}

// -------------------- Handlers: ----------------------

// -------------------- addHandler ----------------------

func addHandler (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	stringX := r.URL.Query().Get("x")
	stringY := r.URL.Query().Get("y")

	x, err := strconv.ParseFloat(stringX, 64)
	if err != nil {
    	http.Error(w, "Missing/wrong format for key parameter X", 400)
        return
    }

    y, err := strconv.ParseFloat(stringY, 64)
    if err != nil{
    	http.Error(w, "Missing/wrong format for key parameter Y", 400)
        return
    }

	// load basic response data
	resp.X = x
	resp.Y = y
	resp.Action = actionAdd
	resp.Cached = false

	// generate key for this action
	key := resp.GenerateKey()

	// try getting cache
	b, err := cacheDatabase.Get(key)
	if err != nil {
		// error
		log.Fatal(err)
	}

	if b != nil {
		// cache exist
		err := json.Unmarshal(b, &resp)
		if err != nil {
			log.Fatal(err)
		}
		resp.Cached = true

		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}

		w.Write(jsonResp)
		return
	} else {
		resp.Answer = Add(resp.X,resp.Y)
	}

	// set cache
	err = cacheDatabase.Set(key, resp, cacheTTL)
	if err != nil {
		log.Fatal(err)
	}

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}

	w.Write(jsonResp)
}

// -------------------- substractHandler ----------------------

func substractHandler (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	stringX := r.URL.Query().Get("x")
	stringY := r.URL.Query().Get("y")

	x, err := strconv.ParseFloat(stringX, 64)
	if err != nil {
    	http.Error(w, "Missing/wrong format for key parameter X", 400)
        return
    }

    y, err := strconv.ParseFloat(stringY, 64)
    if err != nil {
    	http.Error(w, "Missing/wrong format for key parameter Y", 400)
        return
    }

	// load basic response data
	resp.X = x
	resp.Y = y
	resp.Action = actionSubtract
	resp.Cached = false

	// generate key for this action
	key := resp.GenerateKey()

	// try getting cache
	b, err := cacheDatabase.Get(key)
	if err != nil {
		// error
		log.Fatal(err)
	}

	if b != nil {
		// cache exist
		err := json.Unmarshal(b, &resp)
		if err != nil {
			log.Fatal(err)
		}
		resp.Cached = true

		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}

		w.Write(jsonResp)
		return
	} else {
		resp.Answer = Subtract(resp.X,resp.Y)
	}

	// set cache
	err = cacheDatabase.Set(key, resp, cacheTTL)
	if err != nil {
		log.Fatal(err)
	}

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}

	w.Write(jsonResp)
}

// -------------------- multiplyHandler ----------------------

func multiplyHandler (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	stringX := r.URL.Query().Get("x")
	stringY := r.URL.Query().Get("y")

	x, err := strconv.ParseFloat(stringX, 64)
	if err != nil {
    	http.Error(w, "Missing/wrong format for key parameter X", 400)
        return
    }

    y, err := strconv.ParseFloat(stringY, 64)
    if err != nil {
    	http.Error(w, "Missing/wrong format for key parameter Y", 400)
        return
    }

	// load basic response data
	resp.X = x
	resp.Y = y
	resp.Action = actionMultiply
	resp.Cached = false

	// generate key for this action
	key := resp.GenerateKey()

	// try getting cache
	b, err := cacheDatabase.Get(key)
	if err != nil {
		// error
		log.Fatal(err)
	}

	if b != nil {
		// cache exist
		err := json.Unmarshal(b, &resp)
		if err != nil {
			log.Fatal(err)
		}
		resp.Cached = true

		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}

		w.Write(jsonResp)
		return
	} else {
		resp.Answer = Multiply(resp.X,resp.Y)
	}

	// set cache
	err = cacheDatabase.Set(key, resp, cacheTTL)
	if err != nil {
		log.Fatal(err)
	}

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}

	w.Write(jsonResp)
}

// -------------------- divideHandler ----------------------

func divideHandler (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	stringX := r.URL.Query().Get("x")
	stringY := r.URL.Query().Get("y")

	x, err := strconv.ParseFloat(stringX, 64)
	if err != nil {
    	http.Error(w, "Missing/wrong format for key parameter X", 400)
        return
    }

    y, err := strconv.ParseFloat(stringY, 64)
    if err != nil || y == 0 {
    	http.Error(w, "Missing/wrong format for key parameter Y", 400)
        return
    }

	// load basic response data
	resp.X = x
	resp.Y = y
	resp.Action = actionDivide
	resp.Cached = false

	// generate key for this action
	key := resp.GenerateKey()

	// try getting cache
	b, err := cacheDatabase.Get(key)
	if err != nil {
		// error
		log.Fatal(err)
	}

	if b != nil {
		// cache exist
		err := json.Unmarshal(b, &resp)
		if err != nil {
			log.Fatal(err)
		}
		resp.Cached = true

		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}

		w.Write(jsonResp)
		return
	} else {
		resp.Answer = Divide(resp.X,resp.Y)
	}

	// set cache
	err = cacheDatabase.Set(key, resp, cacheTTL)
	if err != nil {
		log.Fatal(err)
	}

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}

	w.Write(jsonResp)
}