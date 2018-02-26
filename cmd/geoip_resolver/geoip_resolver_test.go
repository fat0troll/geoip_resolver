package main

import (
	"encoding/json"
	"github.com/fat0troll/geoip_resolver/lib/appcontext"
	"github.com/fat0troll/geoip_resolver/lib/cache"
	"github.com/fat0troll/geoip_resolver/lib/requester"
	"github.com/icrowley/fake"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
)

var (
	// Pointer to application context
	c *appcontext.Context
	// Server to start in separate routine
	server *http.Server
	// 100 ddresses to use in test
	IPs []string
)

func initializePackages() {
	cache.New(c)
	requester.New(c)
}

// Internal test functions (used for powering actual tests)
func prepareContext() *appcontext.Context {
	c = appcontext.New()
	c.Init()
	return c
}

func prepareConfig(t *testing.T) {
	configPath := "/tmp/geoip-resolver-test-config.yml"
	configIOW, err := os.Create(configPath)
	if err != nil {
		t.Errorf("Failed to open configuration destination file descriptor")
		t.Errorf(err.Error())
		t.FailNow()
	}
	defer configIOW.Close()
	currentDirectory, err := os.Getwd()
	if err != nil {
		t.Errorf("Failed to get current directory")
		t.Errorf(err.Error())
		t.FailNow()
	}
	configIOR, err := os.Open(currentDirectory + "/config.yml.dist")
	if err != nil {
		t.Errorf("Failed to read example configuration")
		t.Errorf(err.Error())
		t.FailNow()
	}
	defer configIOR.Close()
	_, err = io.Copy(configIOW, configIOR)
	if err != nil {
		t.Errorf("Failed to copy configuration data to temporary file")
		t.Errorf(err.Error())
		t.FailNow()
	}

	err = configIOW.Sync()
	if err != nil {
		t.Errorf("Failed to write example configuration")
		t.Errorf(err.Error())
		t.FailNow()
	}

	c.InitializeConfig(configPath)
}

func sendValidIPRequests(t *testing.T) {
	IPs = make([]string, 0)
	for i := 0; i < 100; i++ {
		IPs = append(IPs, fake.IPv4())
	}
	for i := range IPs {
		responseData := make(map[string]string)
		response, err := http.Get("http://" + c.Cfg.HTTPServer.Host + ":" + c.Cfg.HTTPServer.Port + "/?ip=" + IPs[i])
		if err != nil {
			t.Errorf("Failed to dial HTTP server!")
			t.FailNow()
		}

		err = json.NewDecoder(response.Body).Decode(&responseData)
		if err != nil {
			t.Errorf("Failed to decode response!")
			t.Errorf(err.Error())
			t.FailNow()
		}
	}
}

func startLocalServer(t *testing.T) {
	server = &http.Server{
		Addr:         c.Cfg.HTTPServer.Host + ":" + c.Cfg.HTTPServer.Port,
		Handler:      c.HTTPServerMux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			c.Log.Error(err.Error())
		}
	}()

	// Check if we have HTTP server successfully launched.
	serverStarted := make(chan bool, 1)
	go func() {
		c.Log.Debug("Checking if HTTP server is running...")
		checkIDX := 0
		var failed = false
		for {
			if checkIDX > 1 {
				c.Log.Error("HTTP server failed to start!")
				failed = true
				break
			}
			checkIDX++

			_, err := http.Get("http://" + c.Cfg.HTTPServer.Host + ":" + c.Cfg.HTTPServer.Port)
			if err != nil {
				c.Log.Error(err.Error())
				c.Log.Error("Above error indicates that HTTP server isn't started yet, skipping this iteration...")
				time.Sleep(time.Second * 5)
				continue
			}
			break
		}

		if failed {
			serverStarted <- false
		} else {
			serverStarted <- true
		}
	}()

	status := <-serverStarted
	if status {
		c.Log.Info("Server started on http://" + c.Cfg.HTTPServer.Host + ":" + c.Cfg.HTTPServer.Port)
	} else {
		t.Errorf("HTTP server failed to start!")
		t.FailNow()
	}

}

// Tests

func TestStartServer(t *testing.T) {
	prepareContext()
	prepareConfig(t)
	initializePackages()
	startLocalServer(t)
}

func TestValidIP(t *testing.T) {
	sendValidIPRequests(t)
}
