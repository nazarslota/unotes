package utils

import (
	"os"
	"os/exec"
	"testing"
	"time"
)

func TestGracefulShutdown(t *testing.T) {
	// Test that the shutdown channel is closed when the program is interrupted
	cmd := exec.Command(os.Args[0], "-test.run=TestGracefulShutdownHelper")
	err := cmd.Start()
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Second)
	if err := cmd.Process.Signal(os.Interrupt); err != nil {
		t.Fatal(err)
	}

	err = cmd.Wait()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}

func TestGracefulShutdownHelper(t *testing.T) {
	shutdown := GracefulShutdown()
	select {
	case <-shutdown:
		t.Error("GracefulShutdown returned immediately")
	default:
	}
	<-time.After(2 * time.Second)
	select {
	case <-shutdown:
		t.Error("GracefulShutdown returned after 2s")
	default:
	}
}

func TestBuildMongoURI(t *testing.T) {
	tests := []struct {
		name    string
		host    string
		port    string
		user    string
		pass    string
		want    string
		wanterr error
	}{
		{
			name:    "basic URI",
			host:    "localhost",
			port:    "27017",
			user:    "",
			pass:    "",
			want:    "mongodb://localhost:27017",
			wanterr: nil,
		},
		{
			name:    "URI with username and password",
			host:    "localhost",
			port:    "27017",
			user:    "user",
			pass:    "pass",
			want:    "mongodb://user:pass@localhost:27017",
			wanterr: nil,
		},
		{
			name:    "URI with username only",
			host:    "localhost",
			port:    "27017",
			user:    "user",
			pass:    "",
			want:    "mongodb://user@localhost:27017",
			wanterr: nil,
		},
		{
			name:    "URI with SRV scheme",
			host:    "localhost",
			port:    "",
			user:    "",
			pass:    "",
			want:    "mongodb+srv://localhost",
			wanterr: nil,
		},
		{
			name:    "URI with SRV scheme and username",
			host:    "localhost",
			port:    "",
			user:    "user",
			pass:    "",
			want:    "mongodb+srv://user@localhost",
			wanterr: nil,
		},
		{
			name:    "URI with SRV scheme and username and password",
			host:    "localhost",
			port:    "",
			user:    "user",
			pass:    "pass",
			want:    "mongodb+srv://user:pass@localhost",
			wanterr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BuildMongoURI(tt.host, tt.port, tt.user, tt.pass)
			if err != tt.wanterr {
				t.Errorf("BuildMongoURI() error = %v, wantErr %v", err, tt.wanterr)
				return
			}
			if got != tt.want {
				t.Errorf("BuildMongoURI() = %v, want %v", got, tt.want)
			}
		})
	}
}
