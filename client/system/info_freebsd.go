//go:build freebsd

package system

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/netbirdio/netbird/client/system/detect_cloud"
	"github.com/netbirdio/netbird/client/system/detect_platform"
	"github.com/netbirdio/netbird/version"
)

// GetInfo retrieves and parses the system information
func GetInfo(ctx context.Context, staticInfo *StaticInfo) *Info {
	out := _getInfo()
	for strings.Contains(out, "broken pipe") {
		out = _getInfo()
		time.Sleep(500 * time.Millisecond)
	}
	osStr := strings.ReplaceAll(out, "\n", "")
	osStr = strings.ReplaceAll(osStr, "\r\n", "")
	osInfo := strings.Split(osStr, " ")

	osName, osVersion := readOsReleaseFile()

	systemHostname, _ := os.Hostname()

	info := &Info{
		GoOS:               runtime.GOOS,
		Kernel:             osInfo[0],
		Platform:           runtime.GOARCH,
		OS:                 osName,
		OSVersion:          osVersion,
		Hostname:           extractDeviceName(ctx, systemHostname),
		CPUs:               runtime.NumCPU(),
		WiretrusteeVersion: version.NetbirdVersion(),
		UIVersion:          extractUserAgent(ctx),
		KernelVersion:      osInfo[1],
	}
	if staticInfo != nil {
		info.SystemSerialNumber = staticInfo.SystemSerialNumber
		info.SystemProductName = staticInfo.SystemProductName
		info.SystemManufacturer = staticInfo.SystemManufacturer
		info.Environment = staticInfo.Environment
	}

	return info
}

func getStaticInfo(ctx context.Context) *StaticInfo {
	env := Environment{
		Cloud:    detect_cloud.Detect(ctx),
		Platform: detect_platform.Detect(ctx),
	}

	return &StaticInfo{
		Environment: env,
	}
}

func _getInfo() string {
	cmd := exec.Command("uname", "-sri")
	cmd.Stdin = strings.NewReader("some input")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Warnf("getInfo: %s", err)
	}

	return out.String()
}
