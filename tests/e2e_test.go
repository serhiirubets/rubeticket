package tests

import (
	"github.com/serhiirubets/rubeticket/app"
	"net/http/httptest"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Setenv("ENV", "test")
	env := SetupTestEnv()

	router, err := app.InitApp(env.Conf, env.Logger)
	if err != nil {
		env.Logger.Error("Failed to initialize app", "error", err.Error())
		os.Exit(1)
	}

	server := httptest.NewServer(router)
	defer server.Close()

	os.Exit(m.Run())
}
