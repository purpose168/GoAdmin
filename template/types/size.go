package types

import "strconv"

// 超小屏幕 / 手机
// xs: 0

// 小屏幕 / 手机
// sm: 576px

// 中等屏幕 / 平板
// md: 768px

// 大屏幕 / 桌面
// lg: 992px

// 超大屏幕 / 宽屏桌面
// xl: 1200px

// S 是屏幕尺寸映射类型
type S map[string]string

// Size 创建屏幕尺寸映射
// 参数:
//   - sm: 小屏幕尺寸
//   - md: 中等屏幕尺寸
//   - lg: 大屏幕尺寸
//
// 返回: 屏幕尺寸映射
func Size(sm, md, lg int) S {
	var s = make(S)
	if sm > 0 && sm < 13 {
		s["sm"] = strconv.Itoa(sm)
	}
	if md > 0 && md < 13 {
		s["md"] = strconv.Itoa(md)
	}
	if lg > 0 && lg < 13 {
		s["lg"] = strconv.Itoa(lg)
	}
	return s
}

// LG 设置大屏幕尺寸
// 参数:
//   - lg: 大屏幕尺寸
//
// 返回: 更新后的屏幕尺寸映射
func (s S) LG(lg int) S {
	if lg > 0 && lg < 13 {
		s["lg"] = strconv.Itoa(lg)
	}
	return s
}

// XS 设置超小屏幕尺寸
// 参数:
//   - xs: 超小屏幕尺寸
//
// 返回: 更新后的屏幕尺寸映射
func (s S) XS(xs int) S {
	if xs > 0 && xs < 13 {
		s["xs"] = strconv.Itoa(xs)
	}
	return s
}

// XL 设置超大屏幕尺寸
// 参数:
//   - xl: 超大屏幕尺寸
//
// 返回: 更新后的屏幕尺寸映射
func (s S) XL(xl int) S {
	if xl > 0 && xl < 13 {
		s["xl"] = strconv.Itoa(xl)
	}
	return s
}

// SM 设置小屏幕尺寸
// 参数:
//   - sm: 小屏幕尺寸
//
// 返回: 更新后的屏幕尺寸映射
func (s S) SM(sm int) S {
	if sm > 0 && sm < 13 {
		s["sm"] = strconv.Itoa(sm)
	}
	return s
}

// MD 设置中等屏幕尺寸
// 参数:
//   - md: 中等屏幕尺寸
//
// 返回: 更新后的屏幕尺寸映射
func (s S) MD(md int) S {
	if md > 0 && md < 13 {
		s["md"] = strconv.Itoa(md)
	}
	return s
}

// SizeXS 创建超小屏幕尺寸映射
// 参数:
//   - xs: 超小屏幕尺寸
//
// 返回: 屏幕尺寸映射
func SizeXS(xs int) S {
	var s = make(S)
	if xs > 0 && xs < 13 {
		s["xs"] = strconv.Itoa(xs)
	}
	return s
}

// SizeXL 创建超大屏幕尺寸映射
// 参数:
//   - xl: 超大屏幕尺寸
//
// 返回: 屏幕尺寸映射
func SizeXL(xl int) S {
	var s = make(S)
	if xl > 0 && xl < 13 {
		s["xl"] = strconv.Itoa(xl)
	}
	return s
}

// SizeSM 创建小屏幕尺寸映射
// 参数:
//   - sm: 小屏幕尺寸
//
// 返回: 屏幕尺寸映射
func SizeSM(sm int) S {
	var s = make(S)
	if sm > 0 && sm < 13 {
		s["sm"] = strconv.Itoa(sm)
	}
	return s
}

// SizeMD 创建中等屏幕尺寸映射
// 参数:
//   - md: 中等屏幕尺寸
//
// 返回: 屏幕尺寸映射
func SizeMD(md int) S {
	var s = make(S)
	if md > 0 && md < 13 {
		s["md"] = strconv.Itoa(md)
	}
	return s
}

// SizeLG 创建大屏幕尺寸映射
// 参数:
//   - lg: 大屏幕尺寸
//
// 返回: 屏幕尺寸映射
func SizeLG(lg int) S {
	var s = make(S)
	if lg > 0 && lg < 13 {
		s["lg"] = strconv.Itoa(lg)
	}
	return s
}
