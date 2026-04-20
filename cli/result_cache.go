package cli

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/jreisinger/checkip/check"
)

const resultCacheVersion = 1

type resultCache struct {
	dir string
	now func() time.Time
}

type resultCacheRecord struct {
	Version            int             `json:"version"`
	Check              string          `json:"check"`
	IP                 string          `json:"ip"`
	StoredAt           time.Time       `json:"stored_at"`
	ExpiresAt          time.Time       `json:"expires_at"`
	Type               string          `json:"type"`
	MissingCredentials string          `json:"missing_credentials,omitempty"`
	IPAddrIsMalicious  bool            `json:"ip_addr_is_malicious"`
	IPAddrInfo         json.RawMessage `json:"ip_addr_info,omitempty"`
}

func newResultCache(dir string, now func() time.Time) (*resultCache, error) {
	if now == nil {
		now = time.Now
	}

	if dir == "" {
		var err error
		dir, err = check.CacheDir("results", "v1")
		if err != nil {
			return nil, err
		}
	} else if err := os.MkdirAll(dir, 0750); err != nil {
		return nil, err
	}

	return &resultCache{dir: dir, now: now}, nil
}

func (c *resultCache) load(definition check.Definition, ip string) (check.Check, bool) {
	recordPath := c.path(definition.Name, ip)
	buf, err := os.ReadFile(recordPath)
	if err != nil {
		return check.Check{}, false
	}

	var record resultCacheRecord
	if err := json.Unmarshal(buf, &record); err != nil {
		return check.Check{}, false
	}
	if record.Version != resultCacheVersion {
		return check.Check{}, false
	}
	if record.Check != definition.Name || record.IP != ip {
		return check.Check{}, false
	}
	if !record.ExpiresAt.After(c.now()) {
		return check.Check{}, false
	}
	if record.MissingCredentials != "" {
		return check.Check{}, false
	}

	resultType, ok := parseResultType(record.Type)
	if !ok {
		return check.Check{}, false
	}

	result := check.Check{
		Description:       definition.Name,
		Type:              resultType,
		IpAddrIsMalicious: record.IPAddrIsMalicious,
	}

	if len(record.IPAddrInfo) != 0 && string(record.IPAddrInfo) != "null" {
		if definition.NewInfo == nil {
			return check.Check{}, false
		}
		info := definition.NewInfo()
		if info == nil {
			return check.Check{}, false
		}
		if err := json.Unmarshal(record.IPAddrInfo, info); err != nil {
			return check.Check{}, false
		}
		result.IpAddrInfo = info
	}

	return result, true
}

func (c *resultCache) store(definition check.Definition, ip string, result check.Check) {
	record := resultCacheRecord{
		Version:           resultCacheVersion,
		Check:             definition.Name,
		IP:                ip,
		StoredAt:          c.now(),
		ExpiresAt:         c.now().Add(definition.PersistentTTL),
		Type:              result.Type.String(),
		IPAddrIsMalicious: result.IpAddrIsMalicious,
	}

	if result.IpAddrInfo != nil {
		buf, err := result.IpAddrInfo.Json()
		if err != nil {
			return
		}
		record.IPAddrInfo = buf
	}

	recordPath := c.path(definition.Name, ip)
	tmp, err := os.CreateTemp(c.dir, filepath.Base(recordPath)+".tmp-*")
	if err != nil {
		return
	}

	tmpPath := tmp.Name()
	removeTmp := true
	defer func() {
		tmp.Close()
		if removeTmp {
			os.Remove(tmpPath)
		}
	}()

	enc := json.NewEncoder(tmp)
	if err := enc.Encode(record); err != nil {
		return
	}
	if err := tmp.Close(); err != nil {
		return
	}
	if err := os.Rename(tmpPath, recordPath); err != nil {
		return
	}

	removeTmp = false
}

func (c *resultCache) path(definitionName, ip string) string {
	sum := sha256.Sum256([]byte(definitionName + "|" + ip))
	return filepath.Join(c.dir, hex.EncodeToString(sum[:])+".json")
}

func parseResultType(s string) (check.Type, bool) {
	switch s {
	case check.Info.String():
		return check.Info, true
	case check.IsMalicious.String():
		return check.IsMalicious, true
	case check.InfoAndIsMalicious.String():
		return check.InfoAndIsMalicious, true
	default:
		return 0, false
	}
}
