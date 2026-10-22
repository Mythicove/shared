// Package mc 提供 Minecraft 通用领域类型与工具函数。
package mc

import (
	"crypto/md5"
	"fmt"
	"regexp"
)

// OfflineUUID 按 Mojang 规范生成离线玩家 UUID（version-3，基于 MD5("OfflinePlayer:"+name)）。
func OfflineUUID(name string) [16]byte {
	h := md5.Sum([]byte("OfflinePlayer:" + name))
	var u [16]byte
	copy(u[:], h[:])
	u[6] = (u[6] & 0x0f) | 0x30 // 置 version 3
	u[8] = (u[8] & 0x3f) | 0x80 // 置 RFC 4122 variant
	return u
}

// UUIDString 将 16 字节 UUID 格式化为 8-4-4-4-12 字符串。
func UUIDString(u [16]byte) string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:16])
}

// UUIDFromName 返回离线玩家名对应的 UUID 字符串。
func UUIDFromName(name string) string {
	return UUIDString(OfflineUUID(name))
}

var namePattern = regexp.MustCompile(`^[A-Za-z0-9_]{3,16}$`)

// ValidPlayerName 校验合法的 Java 版玩家名。
func ValidPlayerName(name string) bool {
	return namePattern.MatchString(name)
}
