package template

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestVersionCompare 测试版本比较功能
// 测试 VersionCompare 函数在不同版本比较场景下的行为
func TestVersionCompare(t *testing.T) {
	// 测试版本大于指定版本
	assert.Equal(t, true, VersionCompare("v1.2.8", []string{"v1.2.7"}))
	// 测试版本等于指定版本
	assert.Equal(t, true, VersionCompare("v1.2.8", []string{"v1.2.8"}))
	// 测试版本等于多个指定版本中的一个
	assert.Equal(t, true, VersionCompare("v1.2.8", []string{"v1.2.5", "v1.2.8"}))
	// 测试版本小于指定版本
	assert.Equal(t, false, VersionCompare("v1.2.7", []string{"v1.2.8"}))
	// 测试版本等于指定版本（带前导零）
	assert.Equal(t, true, VersionCompare("v0.0.30", []string{"v0.0.30"}))
	// 测试版本大于等于指定版本
	assert.Equal(t, true, VersionCompare("v0.0.30", []string{">=v0.0.30"}))
	// 测试版本大于等于较小版本
	assert.Equal(t, true, VersionCompare("v0.0.30", []string{">=v0.0.29"}))
	// 测试版本不大于等于较大版本
	assert.Equal(t, false, VersionCompare("v0.0.30", []string{">=v0.1.1"}))
	// 测试版本小于等于指定版本
	assert.Equal(t, true, VersionCompare("v0.0.30", []string{"<=v0.1.1"}))
}
