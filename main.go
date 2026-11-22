package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/Peejidle/veil/internal/status"
	"github.com/Peejidle/veil/internal/wg"
)

const iface = "proton"

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "veil:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) != 1 {
		return errors.New("usage: veil on|off|toggle|status")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	switch args[0] {
	case "on":
		return wg.Up(ctx, iface)
	case "off":
		return wg.Down(ctx, iface)
	case "toggle":
		up, err := wg.IsUp(iface)
		if err != nil {
			return err
		}
		if up {
			return wg.Down(ctx, iface)
		}
		return wg.Up(ctx, iface)
	case "status":
		return printStatus(ctx)
	default:
		return fmt.Errorf("unknown command %q (want on|off|toggle|status)", args[0])
	}
}

func printStatus(ctx context.Context) error {
	r, err := status.Check(ctx, iface)
	if err != nil {
		return err
	}

	if r.Up {
		fmt.Println("connected")
	} else {
		fmt.Println("disconnected")
	}
	fmt.Println("public ip:", r.PublicIP)

	if r.Up {
		if r.Handshake.IsZero() {
			fmt.Println("handshake: never")
		} else {
			fmt.Printf("handshake: %s ago\n", time.Since(r.Handshake).Round(time.Second))
		}
	}
	return nil
}
