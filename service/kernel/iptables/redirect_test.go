package iptables

import (
	"strings"
	"testing"
)

func TestLegacyRedirectInstallsDNSHooksAfterGenericHooks(t *testing.T) {
	commands := (&legacyRedirect{}).GetSetupCommands().Cmds
	assertLater := func(generic, dns string) {
		t.Helper()
		genericAt := strings.Index(commands, generic)
		dnsAt := strings.Index(commands, dns)
		if genericAt < 0 || dnsAt < 0 {
			t.Fatalf("missing hook command: generic=%d dns=%d", genericAt, dnsAt)
		}
		if dnsAt <= genericAt {
			t.Fatalf("DNS hook must execute after generic hook because both use -I: generic=%d dns=%d", genericAt, dnsAt)
		}
	}

	assertLater(
		"iptables -w 2 -t nat -I PREROUTING -p tcp -j TP_PRE",
		"iptables -w 2 -t nat -I PREROUTING -p tcp --dport 53 -j DNS_REDIRECT",
	)
	assertLater(
		"iptables -w 2 -t nat -I OUTPUT -p tcp -j TP_OUT",
		"iptables -w 2 -t nat -I OUTPUT -p tcp --dport 53 -j DNS_REDIRECT",
	)
}

func TestLegacyRedirectDNSRulesDeclareProtocol(t *testing.T) {
	commands := (&legacyRedirect{}).GetSetupCommands().Cmds
	for _, want := range []string{
		"iptables -w 2 -t nat -A DNS_REDIRECT -p udp -j REDIRECT --to-port 52353",
		"iptables -w 2 -t nat -A DNS_REDIRECT -p tcp -j REDIRECT --to-port 52353",
	} {
		if !strings.Contains(commands, want) {
			t.Fatalf("missing protocol-qualified DNS redirect rule %q", want)
		}
	}
	for _, invalid := range []string{
		"iptables -w 2 -t nat -A DNS_REDIRECT -j REDIRECT --to-port 52353",
	} {
		if strings.Contains(commands, invalid) {
			t.Fatalf("found protocol-less DNS redirect rule %q", invalid)
		}
	}
}
