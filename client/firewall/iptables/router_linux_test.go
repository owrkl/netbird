//go:build !android

package iptables

import (
	"fmt"
	"net/netip"
	"os/exec"
	"testing"

	"github.com/coreos/go-iptables/iptables"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/netbirdio/netbird/client/firewall/test"
	"github.com/netbirdio/netbird/client/firewall/types"
	nbnet "github.com/netbirdio/netbird/util/net"
)

func isIptablesSupported() bool {
	_, err4 := exec.LookPath("iptables")
	return err4 == nil
}

func TestIptablesManager_RestoreOrCreateContainers(t *testing.T) {
	if !isIptablesSupported() {
		t.SkipNow()
	}

	iptablesClient, err := iptables.NewWithProtocol(iptables.ProtocolIPv4)
	require.NoError(t, err, "failed to init iptables client")

	manager, err := newRouter(iptablesClient, ifaceMock)
	require.NoError(t, err, "should return a valid iptables manager")
	require.NoError(t, manager.init(nil))

	defer func() {
		assert.NoError(t, manager.Reset(), "shouldn't return error")
	}()

	// Now 5 rules:
	// 1. established rule in forward chain
	// 2. jump rule to NAT chain
	// 3. jump rule to PRE chain
	// 4. static outbound masquerade rule
	// 5. static return masquerade rule
	require.Len(t, manager.rules, 5, "should have created rules map")

	exists, err := manager.iptablesClient.Exists(tableNat, chainPOSTROUTING, "-j", chainRTNAT)
	require.NoError(t, err, "should be able to query the iptables %s table and %s chain", tableNat, chainPOSTROUTING)
	require.True(t, exists, "postrouting jump rule should exist")

	exists, err = manager.iptablesClient.Exists(tableMangle, chainPREROUTING, "-j", chainRTPRE)
	require.NoError(t, err, "should be able to query the iptables %s table and %s chain", tableMangle, chainPREROUTING)
	require.True(t, exists, "prerouting jump rule should exist")

	pair := types.RouterPair{
		ID:          "abc",
		Source:      netip.MustParsePrefix("100.100.100.1/32"),
		Destination: netip.MustParsePrefix("100.100.100.0/24"),
		Masquerade:  true,
	}

	err = manager.AddNatRule(pair)
	require.NoError(t, err, "adding NAT rule should not return error")

	err = manager.Reset()
	require.NoError(t, err, "shouldn't return error")
}

func TestIptablesManager_AddNatRule(t *testing.T) {
	if !isIptablesSupported() {
		t.SkipNow()
	}

	for _, testCase := range test.InsertRuleTestCases {
		t.Run(testCase.Name, func(t *testing.T) {
			iptablesClient, err := iptables.NewWithProtocol(iptables.ProtocolIPv4)
			require.NoError(t, err, "failed to init iptables client")

			manager, err := newRouter(iptablesClient, ifaceMock)
			require.NoError(t, err, "shouldn't return error")
			require.NoError(t, manager.init(nil))

			defer func() {
				assert.NoError(t, manager.Reset(), "shouldn't return error")
			}()

			err = manager.AddNatRule(testCase.InputPair)
			require.NoError(t, err, "marking rule should be inserted")

			natRuleKey := types.GenRuleKey(types.NatFormat, testCase.InputPair)
			markingRule := []string{
				"-i", ifaceMock.Name(),
				"-m", "conntrack",
				"--ctstate", "NEW",
				"-s", testCase.InputPair.Source.String(),
				"-d", testCase.InputPair.Destination.String(),
				"-j", "MARK", "--set-mark",
				fmt.Sprintf("%#x", nbnet.PreroutingFwmarkMasquerade),
			}

			exists, err := iptablesClient.Exists(tableMangle, chainRTPRE, markingRule...)
			require.NoError(t, err, "should be able to query the iptables %s table and %s chain", tableMangle, chainRTPRE)
			if testCase.InputPair.Masquerade {
				require.True(t, exists, "marking rule should be created")
				foundRule, found := manager.rules[natRuleKey]
				require.True(t, found, "marking rule should exist in the map")
				require.Equal(t, markingRule, foundRule, "stored marking rule should match")
			} else {
				require.False(t, exists, "marking rule should not be created")
				_, found := manager.rules[natRuleKey]
				require.False(t, found, "marking rule should not exist in the map")
			}

			// Check inverse rule
			inversePair := types.GetInversePair(testCase.InputPair)
			inverseRuleKey := types.GenRuleKey(types.NatFormat, inversePair)
			inverseMarkingRule := []string{
				"!", "-i", ifaceMock.Name(),
				"-m", "conntrack",
				"--ctstate", "NEW",
				"-s", inversePair.Source.String(),
				"-d", inversePair.Destination.String(),
				"-j", "MARK", "--set-mark",
				fmt.Sprintf("%#x", nbnet.PreroutingFwmarkMasqueradeReturn),
			}

			exists, err = iptablesClient.Exists(tableMangle, chainRTPRE, inverseMarkingRule...)
			require.NoError(t, err, "should be able to query the iptables %s table and %s chain", tableMangle, chainRTPRE)
			if testCase.InputPair.Masquerade {
				require.True(t, exists, "inverse marking rule should be created")
				foundRule, found := manager.rules[inverseRuleKey]
				require.True(t, found, "inverse marking rule should exist in the map")
				require.Equal(t, inverseMarkingRule, foundRule, "stored inverse marking rule should match")
			} else {
				require.False(t, exists, "inverse marking rule should not be created")
				_, found := manager.rules[inverseRuleKey]
				require.False(t, found, "inverse marking rule should not exist in the map")
			}
		})
	}
}

func TestIptablesManager_RemoveNatRule(t *testing.T) {
	if !isIptablesSupported() {
		t.SkipNow()
	}

	for _, testCase := range test.RemoveRuleTestCases {
		t.Run(testCase.Name, func(t *testing.T) {
			iptablesClient, _ := iptables.NewWithProtocol(iptables.ProtocolIPv4)

			manager, err := newRouter(iptablesClient, ifaceMock)
			require.NoError(t, err, "shouldn't return error")
			require.NoError(t, manager.init(nil))
			defer func() {
				assert.NoError(t, manager.Reset(), "shouldn't return error")
			}()

			err = manager.AddNatRule(testCase.InputPair)
			require.NoError(t, err, "should add NAT rule without error")

			err = manager.RemoveNatRule(testCase.InputPair)
			require.NoError(t, err, "shouldn't return error")

			natRuleKey := types.GenRuleKey(types.NatFormat, testCase.InputPair)
			markingRule := []string{
				"-i", ifaceMock.Name(),
				"-m", "conntrack",
				"--ctstate", "NEW",
				"-s", testCase.InputPair.Source.String(),
				"-d", testCase.InputPair.Destination.String(),
				"-j", "MARK", "--set-mark",
				fmt.Sprintf("%#x", nbnet.PreroutingFwmarkMasquerade),
			}

			exists, err := iptablesClient.Exists(tableMangle, chainRTPRE, markingRule...)
			require.NoError(t, err, "should be able to query the iptables %s table and %s chain", tableMangle, chainRTPRE)
			require.False(t, exists, "marking rule should not exist")

			_, found := manager.rules[natRuleKey]
			require.False(t, found, "marking rule should not exist in the manager map")

			// Check inverse rule removal
			inversePair := types.GetInversePair(testCase.InputPair)
			inverseRuleKey := types.GenRuleKey(types.NatFormat, inversePair)
			inverseMarkingRule := []string{
				"!", "-i", ifaceMock.Name(),
				"-m", "conntrack",
				"--ctstate", "NEW",
				"-s", inversePair.Source.String(),
				"-d", inversePair.Destination.String(),
				"-j", "MARK", "--set-mark",
				fmt.Sprintf("%#x", nbnet.PreroutingFwmarkMasqueradeReturn),
			}

			exists, err = iptablesClient.Exists(tableMangle, chainRTPRE, inverseMarkingRule...)
			require.NoError(t, err, "should be able to query the iptables %s table and %s chain", tableMangle, chainRTPRE)
			require.False(t, exists, "inverse marking rule should not exist")

			_, found = manager.rules[inverseRuleKey]
			require.False(t, found, "inverse marking rule should not exist in the map")
		})
	}
}

func TestRouter_AddRouteFiltering(t *testing.T) {
	if !isIptablesSupported() {
		t.Skip("iptables not supported on this system")
	}

	iptablesClient, err := iptables.NewWithProtocol(iptables.ProtocolIPv4)
	require.NoError(t, err, "Failed to create iptables client")

	r, err := newRouter(iptablesClient, ifaceMock)
	require.NoError(t, err, "Failed to create router manager")
	require.NoError(t, r.init(nil))

	defer func() {
		err := r.Reset()
		require.NoError(t, err, "Failed to reset router")
	}()

	tests := []struct {
		name        string
		sources     []netip.Prefix
		destination netip.Prefix
		proto       types.Protocol
		sPort       *types.Port
		dPort       *types.Port
		direction   types.RuleDirection
		action      types.Action
		expectSet   bool
	}{
		{
			name:        "Basic TCP rule with single source",
			sources:     []netip.Prefix{netip.MustParsePrefix("192.168.1.0/24")},
			destination: netip.MustParsePrefix("10.0.0.0/24"),
			proto:       types.ProtocolTCP,
			sPort:       nil,
			dPort:       &types.Port{Values: []int{80}},
			direction:   types.RuleDirectionIN,
			action:      types.ActionAccept,
			expectSet:   false,
		},
		{
			name: "UDP rule with multiple sources",
			sources: []netip.Prefix{
				netip.MustParsePrefix("172.16.0.0/16"),
				netip.MustParsePrefix("192.168.0.0/16"),
			},
			destination: netip.MustParsePrefix("10.0.0.0/8"),
			proto:       types.ProtocolUDP,
			sPort:       &types.Port{Values: []int{1024, 2048}, IsRange: true},
			dPort:       nil,
			direction:   types.RuleDirectionOUT,
			action:      types.ActionDrop,
			expectSet:   true,
		},
		{
			name:        "All protocols rule",
			sources:     []netip.Prefix{netip.MustParsePrefix("10.0.0.0/8")},
			destination: netip.MustParsePrefix("0.0.0.0/0"),
			proto:       types.ProtocolALL,
			sPort:       nil,
			dPort:       nil,
			direction:   types.RuleDirectionIN,
			action:      types.ActionAccept,
			expectSet:   false,
		},
		{
			name:        "ICMP rule",
			sources:     []netip.Prefix{netip.MustParsePrefix("192.168.0.0/16")},
			destination: netip.MustParsePrefix("10.0.0.0/8"),
			proto:       types.ProtocolICMP,
			sPort:       nil,
			dPort:       nil,
			direction:   types.RuleDirectionIN,
			action:      types.ActionAccept,
			expectSet:   false,
		},
		{
			name:        "TCP rule with multiple source ports",
			sources:     []netip.Prefix{netip.MustParsePrefix("172.16.0.0/12")},
			destination: netip.MustParsePrefix("192.168.0.0/16"),
			proto:       types.ProtocolTCP,
			sPort:       &types.Port{Values: []int{80, 443, 8080}},
			dPort:       nil,
			direction:   types.RuleDirectionOUT,
			action:      types.ActionAccept,
			expectSet:   false,
		},
		{
			name:        "UDP rule with single IP and port range",
			sources:     []netip.Prefix{netip.MustParsePrefix("192.168.1.1/32")},
			destination: netip.MustParsePrefix("10.0.0.0/24"),
			proto:       types.ProtocolUDP,
			sPort:       nil,
			dPort:       &types.Port{Values: []int{5000, 5100}, IsRange: true},
			direction:   types.RuleDirectionIN,
			action:      types.ActionDrop,
			expectSet:   false,
		},
		{
			name:        "TCP rule with source and destination ports",
			sources:     []netip.Prefix{netip.MustParsePrefix("10.0.0.0/24")},
			destination: netip.MustParsePrefix("172.16.0.0/16"),
			proto:       types.ProtocolTCP,
			sPort:       &types.Port{Values: []int{1024, 65535}, IsRange: true},
			dPort:       &types.Port{Values: []int{22}},
			direction:   types.RuleDirectionOUT,
			action:      types.ActionAccept,
			expectSet:   false,
		},
		{
			name:        "Drop all incoming traffic",
			sources:     []netip.Prefix{netip.MustParsePrefix("0.0.0.0/0")},
			destination: netip.MustParsePrefix("192.168.0.0/24"),
			proto:       types.ProtocolALL,
			sPort:       nil,
			dPort:       nil,
			direction:   types.RuleDirectionIN,
			action:      types.ActionDrop,
			expectSet:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ruleKey, err := r.AddRouteFiltering(tt.sources, tt.destination, tt.proto, tt.sPort, tt.dPort, tt.action)
			require.NoError(t, err, "AddRouteFiltering failed")

			// Check if the rule is in the internal map
			rule, ok := r.rules[ruleKey.GetRuleID()]
			assert.True(t, ok, "Rule not found in internal map")

			// Log the internal rule
			t.Logf("Internal rule: %v", rule)

			// Check if the rule exists in iptables
			exists, err := iptablesClient.Exists(tableFilter, chainRTFWD, rule...)
			assert.NoError(t, err, "Failed to check rule existence")
			assert.True(t, exists, "Rule not found in iptables")

			// Verify rule content
			params := routeFilteringRuleParams{
				Sources:     tt.sources,
				Destination: tt.destination,
				Proto:       tt.proto,
				SPort:       tt.sPort,
				DPort:       tt.dPort,
				Action:      tt.action,
				SetName:     "",
			}

			expectedRule := genRouteFilteringRuleSpec(params)

			if tt.expectSet {
				setName := types.GenerateSetName(tt.sources)
				params.SetName = setName
				expectedRule = genRouteFilteringRuleSpec(params)

				// Check if the set was created
				_, exists := r.ipsetCounter.Get(setName)
				assert.True(t, exists, "IPSet not created")
			}

			assert.Equal(t, expectedRule, rule, "Rule content mismatch")

			// Clean up
			err = r.DeleteRouteRule(ruleKey)
			require.NoError(t, err, "Failed to delete rule")
		})
	}
}
