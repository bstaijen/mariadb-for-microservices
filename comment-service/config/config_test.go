package config_test

import (
	"os"
	"testing"

	"github.com/bstaijen/mariadb-for-microservices/comment-service/config"
)

func TestPort(t *testing.T) {
	os.Setenv("PORT", "1000")
	actual := config.LoadConfig().Port
	expected := 1000
	if expected != actual {
		t.Fatalf("Expected %v got %v", expected, actual)
	}
	os.Clearenv()
}

func TestPortEmpty(t *testing.T) {
	os.Clearenv()
	actual := config.LoadConfig().Port
	expected := 0
	if expected != actual {
		t.Fatalf("Expected %v got %v", expected, actual)
	}
}

func TestProfileServiceBaseurl(t *testing.T) {
	os.Setenv("PROFILE_SERVICE_URL", "/test")
	actual := config.LoadConfig().ProfileServiceBaseurl
	expected := "/test"
	if expected != actual {
		t.Fatalf("Expected %v got %v", expected, actual)
	}
	os.Clearenv()
}

func TestProfileServiceBaseurlEmpty(t *testing.T) {
	os.Clearenv()
	actual := config.LoadConfig().ProfileServiceBaseurl
	expected := ""
	if expected != actual {
		t.Fatalf("Expected %v got %v", expected, actual)
	}
}

func TestDBUsername(t *testing.T) {
	os.Setenv("DB_USERNAME", "user")
	actual := config.LoadConfig().DBUsername
	expected := "user"
	if expected != actual {
		t.Fatalf("Expected %s got %s", expected, actual)
	}
	os.Clearenv()
}

func TestDBUsernameEmpty(t *testing.T) {
	os.Clearenv()
	actual := config.LoadConfig().DBUsername
	expected := ""
	if expected != actual {
		t.Fatalf("Expected %s got %s", expected, actual)
	}
}

func TestDBPassword(t *testing.T) {
	os.Setenv("DB_PASSWORD", "pass")
	actual := config.LoadConfig().DBPassword
	expected := "pass"
	if expected != actual {
		t.Fatalf("Expected %s got %s", expected, actual)
	}
	os.Clearenv()
}

func TestDBPasswordEmpty(t *testing.T) {
	os.Clearenv()
	actual := config.LoadConfig().DBPassword
	expected := ""
	if expected != actual {
		t.Fatalf("Expected %s got %s", expected, actual)
	}
}

func TestDBHost(t *testing.T) {
	os.Setenv("DB_HOST", "localhost")
	actual := config.LoadConfig().DBHost
	expected := "localhost"
	if expected != actual {
		t.Fatalf("Expected %s got %s", expected, actual)
	}
	os.Clearenv()
}

func TestDBHostEmpty(t *testing.T) {
	os.Clearenv()
	actual := config.LoadConfig().DBHost
	expected := ""
	if expected != actual {
		t.Fatalf("Expected %s got %s", expected, actual)
	}
}

func TestDBPort(t *testing.T) {
	os.Setenv("DB_PORT", "3306")
	actual := config.LoadConfig().DBPort
	expected := 3306
	if expected != actual {
		t.Fatalf("Expected %v got %v", expected, actual)
	}
	os.Clearenv()
}

func TestDBPortEmpty(t *testing.T) {
	os.Clearenv()
	actual := config.LoadConfig().DBPort
	expected := 0
	if expected != actual {
		t.Fatalf("Expected %v got %v", expected, actual)
	}
}

func TestDB(t *testing.T) {
	os.Setenv("DB", "TestDatabase")
	actual := config.LoadConfig().Database
	expected := "TestDatabase"
	if expected != actual {
		t.Fatalf("Expected %s got %s", expected, actual)
	}
	os.Clearenv()
}

func TestDBEmpty(t *testing.T) {
	os.Clearenv()
	actual := config.LoadConfig().Database
	expected := ""
	if expected != actual {
		t.Fatalf("Expected %s got %s", expected, actual)
	}
}
