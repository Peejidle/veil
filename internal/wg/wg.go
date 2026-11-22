// Package wg wraps the WireGuard tools (wg-quick and wg) for veil.
package wg

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const sysNetDir = "/sys/class/net"

func IsUp(iface string) (bool, error) {
	_, err := os.Stat(filepath.Join(sysNetDir, iface))
	if err == nil {
		return true, nil
	}
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	return false, fmt.Errorf("checking interface %q: %w", iface, err)
}

// Up brings the named interface up with wg-quick.
func Up(ctx context.Context, iface string) error {
	_, err := run(ctx, "/usr/bin/wg-quick", "up", iface)
	return err
}

// Down takes the named interface down with wg-quick.
func Down(ctx context.Context, iface string) error {
	_, err := run(ctx, "/usr/bin/wg-quick", "down", iface)
	return err
}

// LatestHandshake returns when the interface last completed a handshake
// with its peer. A zero time means no handshake has happened yet.
func LatestHandshake(ctx context.Context, iface string) (time.Time, error) {
	out, err := run(ctx, "/usr/bin/wg", "show", iface, "latest-handshakes")
	if err != nil {
		return time.Time{}, err
	}
	fields := strings.Fields(out)
	if len(fields) < 2 {
		return time.Time{}, fmt.Errorf("unexpected wg output: %q", out)
	}
	secs, err := strconv.ParseInt(fields[1], 10, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("parsing handshake time %q: %w", fields[1], err)
	}
	if secs == 0 {
		return time.Time{}, nil
	}
	return time.Unix(secs, 0), nil
}

// run executes a command through sudo without prompting for a password
// and returns its output. On failure, the output goes into the error,
// because wg-quick explains what went wrong in its output.
func run(ctx context.Context, args ...string) (string, error) {
	sudoArgs := append([]string{"-n"}, args...)
	cmd := exec.CommandContext(ctx, "sudo", sudoArgs...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%s: %w: %s", strings.Join(args, " "), err, strings.TrimSpace(string(out)))
	}
	return string(out), nil
}
