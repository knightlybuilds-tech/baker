package cli

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/knightlybuilds-tech/baker/internal/bakery"
	"github.com/knightlybuilds-tech/baker/internal/config"
	"github.com/spf13/cobra"
)

func bakeryCmd() *cobra.Command {
	var port int
	var noOpen bool
	cmd := &cobra.Command{
		Use:   "bakery",
		Args:  cobra.NoArgs,
		Short: "Browse this project's memory in a local, read-only Bakery view",
		Long: `Serve the Bakery workspace on 127.0.0.1 with this project's own memory/ notes
as cards, a list and a connected web. Read-only: nothing is written or uploaded,
and no account is needed. Review, branch and team screens are demo-only.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := projectRoot()
			if err != nil {
				return err
			}
			name := filepath.Base(root)
			if p, err := config.LoadProject(root); err == nil && p.Name != "" {
				name = p.Name
			}
			ln, err := listenLoopback(port)
			if err != nil {
				return err
			}
			url := "http://" + ln.Addr().String() + "/"
			srv := &http.Server{Handler: bakery.Handler(root, name), ReadHeaderTimeout: 5 * time.Second}

			fmt.Printf("Bakery is open for %s at %s\n", name, url)
			fmt.Println("Read-only view of memory/ on this machine. Press Ctrl+C to stop.")
			if !noOpen {
				openBrowser(url)
			}

			ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
			defer stop()
			go func() {
				<-ctx.Done()
				shutdown, cancel := context.WithTimeout(context.Background(), 2*time.Second)
				defer cancel()
				_ = srv.Shutdown(shutdown)
			}()
			if err := srv.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
				return err
			}
			fmt.Println("\nBakery closed.")
			return nil
		},
	}
	cmd.Flags().IntVar(&port, "port", 4410, "port on 127.0.0.1 (tries the next few if busy)")
	cmd.Flags().BoolVar(&noOpen, "no-open", false, "don't open a browser")
	return cmd
}

func listenLoopback(port int) (net.Listener, error) {
	var last error
	for p := port; p < port+10; p++ {
		ln, err := net.Listen("tcp", net.JoinHostPort("127.0.0.1", strconv.Itoa(p)))
		if err == nil {
			return ln, nil
		}
		last = err
	}
	return nil, fmt.Errorf("no free port from %d to %d: %w", port, port+9, last)
}

func openBrowser(url string) {
	var c *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		c = exec.Command("open", url)
	case "windows":
		c = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		c = exec.Command("xdg-open", url)
	}
	if err := c.Start(); err == nil {
		go func() { _ = c.Wait() }()
	}
}
