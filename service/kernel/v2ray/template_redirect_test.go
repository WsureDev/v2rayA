package v2ray

import (
	"testing"

	"github.com/v2rayA/v2rayA/db/configure"
)

func TestAppendDokodemoRedirectAcceptsForwardedTCP(t *testing.T) {
	tmpl := new(Template)
	tmpl.AppendDokodemoTProxy(string(configure.TransparentRedirect), 52345, "transparent")

	if len(tmpl.Inbounds) != 1 {
		t.Fatalf("got %d inbounds, want 1", len(tmpl.Inbounds))
	}
	inbound := tmpl.Inbounds[0]
	if inbound.Listen != "0.0.0.0" {
		t.Fatalf("redirect listen = %q, want 0.0.0.0", inbound.Listen)
	}
	if inbound.Settings == nil || inbound.Settings.Network != "tcp" {
		t.Fatalf("redirect network = %#v, want tcp", inbound.Settings)
	}
	if inbound.StreamSettings == nil || inbound.StreamSettings.Sockopt == nil ||
		inbound.StreamSettings.Sockopt.Tproxy == nil ||
		*inbound.StreamSettings.Sockopt.Tproxy != "redirect" {
		t.Fatalf("redirect sockopt is incomplete: %#v", inbound.StreamSettings)
	}
	if !inbound.Settings.FollowRedirect {
		t.Fatal("redirect inbound must recover the original destination")
	}
}

func TestAppendDokodemoTproxyRemainsLoopbackDualProtocol(t *testing.T) {
	tmpl := new(Template)
	tmpl.AppendDokodemoTProxy(string(configure.TransparentTproxy), 52345, "transparent")

	inbound := tmpl.Inbounds[0]
	if inbound.Listen != "127.0.0.1" {
		t.Fatalf("tproxy listen = %q, want 127.0.0.1", inbound.Listen)
	}
	if inbound.Settings == nil || inbound.Settings.Network != "tcp,udp" {
		t.Fatalf("tproxy network = %#v, want tcp,udp", inbound.Settings)
	}
}
