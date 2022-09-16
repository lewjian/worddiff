package worddiff

import (
	"strings"
)

type (
	Word      string // 比较的最小单元，可以是字符串，字符，具体看如何分割字符串
	Operation int8
)

const (
	// DiffDelete 删除
	DiffDelete Operation = -1
	// DiffInsert 插入/新增
	DiffInsert Operation = 1
	// DiffEqual 相同
	DiffEqual Operation = 0
)

// Diff 代表一个diff结果
type Diff struct {
	Type Operation
	Text string
}

type wordDiff struct {
	diffType Operation
	word     Word
}

type WordDiff struct {
	// 字符分割，表示有这个字符则左右当做两个词来处理，当该不指定的时候则按照每个字符分割
	separator map[string]struct{}
	// 连续的分割附是否当做同一个词处理
	mergeContinuousSeparator bool
}

// New 初始化
func New(opts ...DiffOption) *WordDiff {
	wd := &WordDiff{
		separator:                make(map[string]struct{}),
		mergeContinuousSeparator: false,
	}
	for _, opt := range opts {
		opt(wd)
	}
	return wd
}

var defaultWd = New(MergeContinuousSeparator(true), SetSeparator([]string{" "}))

// Default 默认diff
// 分割字符为空格，合并连续空格
func Default() *WordDiff {
	return defaultWd
}

type DiffOption func(wd *WordDiff)

// SetSeparator 设置diff的分词符号
func SetSeparator(separators []string) DiffOption {
	return func(wd *WordDiff) {
		m := make(map[string]struct{}, len(separators))
		for i, _ := range separators {
			m[separators[i]] = struct{}{}
		}
		wd.separator = m
	}
}

// MergeContinuousSeparator 设置是否合并分隔符号
func MergeContinuousSeparator(merge bool) DiffOption {
	return func(wd *WordDiff) {
		wd.mergeContinuousSeparator = merge
	}
}

// Diff 执行比对
func (wd *WordDiff) Diff(oldStr, newStr string) []Diff {
	if oldStr == newStr {
		return []Diff{
			{
				Type: DiffEqual,
				Text: newStr,
			},
		}
	}
	oldWords := split(oldStr, wd.separator, wd.mergeContinuousSeparator)
	newWords := split(newStr, wd.separator, wd.mergeContinuousSeparator)
	return diff(oldWords, newWords)
}

// split 将字符串转为[]Word
func split(s string, separators map[string]struct{}, mergeContinuousSeparator bool) []Word {
	runeArr := []rune(s)
	rLen := len(runeArr)
	left := 0
	words := make([]Word, 0, rLen)
	for i := 0; i < rLen; i++ {
		item := string(runeArr[i])
		if separators == nil || len(separators) == 0 {
			words = append(words, Word(runeArr[left:i+1]))
			left = i + 1
		} else if _, ok := separators[item]; ok {
			// 是分隔符，left:i就是一个word
			if left != i {
				words = append(words, Word(runeArr[left:i]))
			}
			// 重置left
			left = i
			// 检查是否合并连续分割符
			for mergeContinuousSeparator && i+1 < rLen && runeArr[i+1] == runeArr[i] {
				i++
			}
			words = append(words, Word(runeArr[left:i+1]))
			left = i + 1
		}
	}
	if left <= rLen-1 {
		words = append(words, Word(runeArr[left:]))
	}
	return words
}

// shortestEditScript 生成最短的编辑脚本
func shortestEditScript(src, dst []Word) []wordDiff {
	n := len(src)
	m := len(dst)
	max := n + m
	trace := make([]map[int]int, 0, max)
	var x, y int

loop:
	for d := 0; d <= max; d++ {
		// 最多只有 d+1 个 k
		v := make(map[int]int, d+2)
		trace = append(trace, v)

		// 需要注意处理对角线
		if d == 0 {
			t := 0
			for len(src) > t && len(dst) > t && src[t] == dst[t] {
				t++
			}
			v[0] = t
			if t == len(src) && t == len(dst) {
				break loop
			}
			continue
		}

		lastV := trace[d-1]

		for k := -d; k <= d; k += 2 {
			// 向下
			if k == -d || (k != d && lastV[k-1] < lastV[k+1]) {
				x = lastV[k+1]
			} else { // 向右
				x = lastV[k-1] + 1
			}

			y = x - k

			// 处理对角线
			for x < n && y < m && src[x] == dst[y] {
				x, y = x+1, y+1
			}

			v[k] = x

			if x == n && y == m {
				break loop
			}
		}
	}

	// 反向回溯
	script := make([]Operation, 0, max)

	x = n
	y = m
	var k, prevK, prevX, prevY int

	for d := len(trace) - 1; d > 0; d-- {
		k = x - y
		lastV := trace[d-1]

		if k == -d || (k != d && lastV[k-1] < lastV[k+1]) {
			prevK = k + 1
		} else {
			prevK = k - 1
		}

		prevX = lastV[prevK]
		prevY = prevX - prevK

		for x > prevX && y > prevY {
			script = append(script, DiffEqual)
			x -= 1
			y -= 1
		}

		if x == prevX {
			script = append(script, DiffInsert)
		} else {
			script = append(script, DiffDelete)
		}

		x, y = prevX, prevY
	}

	if trace[0][0] != 0 {
		for i := 0; i < trace[0][0]; i++ {
			script = append(script, DiffEqual)
		}
	}
	diffs := make([]wordDiff, 0, len(script))
	srcIndex, dstIndex := 0, 0
	for i := len(script) - 1; i >= 0; i-- {
		op := script[i]
		switch op {
		case DiffInsert:
			diffs = append(diffs, wordDiff{
				diffType: DiffInsert,
				word:     dst[dstIndex],
			})
			dstIndex += 1
		case DiffEqual:
			diffs = append(diffs, wordDiff{
				diffType: DiffEqual,
				word:     src[srcIndex],
			})
			srcIndex += 1
			dstIndex += 1

		case DiffDelete:
			diffs = append(diffs, wordDiff{
				diffType: DiffDelete,
				word:     src[srcIndex],
			})
			srcIndex += 1
		}
	}

	return diffs
}

// diff 比对两个数据
func diff(src, dst []Word) []Diff {
	// 检查是否有公共前后部位
	// 公共前部
	prefixLen := commonPrefixLen(src, dst)
	prefix := src[:prefixLen]
	src = src[prefixLen:]
	dst = dst[prefixLen:]
	suffixLen := commonSuffixLength(src, dst)
	suffix := src[len(src)-suffixLen:]
	src = src[:len(src)-suffixLen]
	dst = dst[:len(dst)-suffixLen]
	// 计算
	diffs := compute(src, dst)
	if len(prefix) > 0 {
		head := merge(prefix, DiffEqual)
		diffs = append([]Diff{head}, diffs...)
	}
	if len(suffix) > 0 {
		tail := merge(suffix, DiffEqual)
		diffs = append(diffs, tail)
	}
	return diffs
}

func compute(text1, text2 []Word) []Diff {
	var longtext, shorttext []Word
	if len(text1) > len(text2) {
		longtext = text1
		shorttext = text2
	} else {
		longtext = text2
		shorttext = text1
	}
	var script []wordDiff
	if i := wordsIndex(longtext, shorttext); i != -1 {
		op := DiffInsert
		if len(text1) > len(text2) {
			op = DiffDelete
		}
		return []Diff{
			merge(longtext[:i], op),
			merge(shorttext, DiffEqual),
			merge(longtext[i+len(shorttext):], op),
		}
	} else if len(shorttext) == 1 {
		return []Diff{
			merge(text1, DiffDelete),
			merge(text2, DiffInsert),
		}
	} else {

		// 获取编辑距离
		script = shortestEditScript(text1, text2)
	}
	wordLen := len(script)
	diffs := make([]Diff, 0, wordLen)
	// 合并连续的相同Operation的字符串
	sb := strings.Builder{}
	for i := 0; i < wordLen; i++ {
		sb.WriteString(string(script[i].word))
		for wordLen > i+1 && script[i+1].diffType == script[i].diffType {
			sb.WriteString(string(script[i+1].word))
			i++
		}
		diffs = append(diffs, Diff{
			Type: script[i].diffType,
			Text: sb.String(),
		})
		sb.Reset()
	}
	return diffs
}

// wordsIndex is the equivalent of strings.Index for Word slices.
func wordsIndex(r1, r2 []Word) int {
	last := len(r1) - len(r2)
	for i := 0; i <= last; i++ {
		if wordsEqual(r1[i:i+len(r2)], r2) {
			return i
		}
	}
	return -1
}

func wordsEqual(r1, r2 []Word) bool {
	if len(r1) != len(r2) {
		return false
	}
	for i, c := range r1 {
		if c != r2[i] {
			return false
		}
	}
	return true
}

func merge(words []Word, op Operation) Diff {
	var sb strings.Builder
	for _, word := range words {
		sb.WriteString(string(word))
	}
	return Diff{
		Type: op,
		Text: sb.String(),
	}
}

// commonPrefixLen 公共前缀长度
func commonPrefixLen(text1, text2 []Word) int {
	n := 0
	for ; n < len(text1) && n < len(text2); n++ {
		if text1[n] != text2[n] {
			return n
		}
	}
	return n
}

// commonSuffixLength 公共尾部长度
func commonSuffixLength(text1, text2 []Word) int {
	i1 := len(text1)
	i2 := len(text2)
	for n := 0; ; n++ {
		i1--
		i2--
		if i1 < 0 || i2 < 0 || text1[i1] != text2[i2] {
			return n
		}
	}
}
