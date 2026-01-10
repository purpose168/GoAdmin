package types

import (
	"testing"
)

// TestFilterOperatorAddOrNot 测试 FilterOperator.AddOrNot() 方法
// 该方法用于判断是否需要添加操作符字段
func TestFilterOperatorAddOrNot(t *testing.T) {
	tests := []struct {
		name     string         // 测试用例名称
		operator FilterOperator // 筛选操作符
		expected bool           // 期望结果
	}{
		{
			name:     "FilterOperatorLike 不应添加操作符字段",
			operator: FilterOperatorLike,
			expected: false,
		},
		{
			name:     "FilterOperatorFree 不应添加操作符字段",
			operator: FilterOperatorFree,
			expected: false,
		},
		{
			name:     "空操作符不应添加操作符字段",
			operator: FilterOperator(""),
			expected: false,
		},
		{
			name:     "FilterOperatorGreater 应添加操作符字段",
			operator: FilterOperatorGreater,
			expected: true,
		},
		{
			name:     "FilterOperatorEqual 应添加操作符字段",
			operator: FilterOperatorEqual,
			expected: true,
		},
		{
			name:     "FilterOperatorNotEqual 应添加操作符字段",
			operator: FilterOperatorNotEqual,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.operator.AddOrNot()
			if result != tt.expected {
				t.Errorf("FilterOperator.AddOrNot() = %v, 期望 %v，操作符 %s", result, tt.expected, string(tt.operator))
			}
		})
	}
}

// TestFilterOperatorLabel 测试 FilterOperator.Label() 方法
// 该方法用于返回操作符的标签
func TestFilterOperatorLabel(t *testing.T) {
	tests := []struct {
		name     string         // 测试用例名称
		operator FilterOperator // 筛选操作符
		expected string         // 期望的标签
	}{
		{
			name:     "FilterOperatorLike 应返回空标签",
			operator: FilterOperatorLike,
			expected: "",
		},
		{
			name:     "FilterOperatorGreater 应返回 > 标签",
			operator: FilterOperatorGreater,
			expected: ">",
		},
		{
			name:     "FilterOperatorEqual 应返回 = 标签",
			operator: FilterOperatorEqual,
			expected: "=",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := string(tt.operator.Label())
			if result != tt.expected {
				t.Errorf("FilterOperator.Label() = %v, 期望 %v，操作符 %s", result, tt.expected, string(tt.operator))
			}
		})
	}
}
