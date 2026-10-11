package ezpass

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// CopyToClipboard copies text to the system clipboard.
// It uses OSC 52 terminal sequence as well as platform-native tools (wl-copy, xclip, pbcopy, clip.exe).
func CopyToClipboard(text string) error {
	// 1. Try platform clipboard commands
	commands := [][]string{
		{"wl-copy"},
		{"xclip", "-selection", "clipboard"},
		{"xsel", "--clipboard", "--input"},
		{"pbcopy"},
		{"clip.exe"},
	}

	toolSuccess := false
	for _, cmdArgs := range commands {
		path, err := exec.LookPath(cmdArgs[0])
		if err == nil {
			cmd := exec.Command(path, cmdArgs[1:]...)
			cmd.Stdin = strings.NewReader(text)
			if err := cmd.Run(); err == nil {
				toolSuccess = true
				break
			}
		}
	}

	// 2. Also emit OSC 52 escape sequence so terminal emulators / SSH sessions receive it
	oscErr := writeOSC52(text)
	if toolSuccess || oscErr == nil {
		return nil
	}

	return fmt.Errorf("could not copy to clipboard (install xclip, wl-copy, or use an OSC 52 compatible terminal)")
}

func writeOSC52(text string) error {
	encoded := base64.StdEncoding.EncodeToString([]byte(text))
	var seq string
	if os.Getenv("TMUX") != "" {
		// Passthrough sequence for tmux
		seq = fmt.Sprintf("\x1bPtmux;\x1b\x1b]52;c;%s\x07\x1b\\", encoded)
	} else {
		seq = fmt.Sprintf("\x1b]52;c;%s\x07", encoded)
	}

	// Prefer writing directly to /dev/tty if available, otherwise stderr
	var out io.Writer = os.Stderr
	tty, err := os.OpenFile("/dev/tty", os.O_WRONLY, 0)
	if err == nil {
		defer tty.Close()
		out = tty
	}

	_, err = fmt.Fprint(out, seq)
	return err
}
