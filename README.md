# worddiff
用于展示两个字符串的差异，支持按照一定规则分割词语，比如英语的空格，或者其他特殊字符

# 示例
```go
    // default 以空格分割字符串来比较
    wd := Default()
    s1 := "Charged Attacks from bow-using characters have 100% increased CRIT Rate. Additionally, Charged Attacks from bow-using characters will unleash a shockwave when they hit opponents, dealing one instance of AoE DMG. Can occur once every 1s."
    s2 := "Charged Attacks from bow-wielding characters have 100% increased CRIT Rate. Additionally, Charged Attacks from bow-using characters will unleash a shockwave when they hit opponents, dealing one instance of AoE DMG. Can occur once every 1s."
    results := wd.Diff(s1, s2)
    fmt.Printf("%+v\n", results)
    // new不传参数则表示每个字符单独比较
    wd = New()
    results = wd.Diff(s1, s2)
    fmt.Printf("%+v\n", results)
    // 自己设置规则
    wd = New(SetSeparator([]string{"-", "."}), MergeContinuousSeparator(true))
    results = wd.Diff(s1, s2)
    fmt.Printf("%+v\n", results)
```

输出
```
[{Type:0 Text:Charged Attacks from } {Type:-1 Text:bow-using} {Type:1 Text:bow-wielding} {Type:0 Text: characters have 100% increased CRIT Rate. Additionally, Charged Attacks from bow-using characters will unleash a shockwave when they hit opponents, dealing one instance of AoE DMG. Can occur once every 1s.}]
[{Type:0 Text:Charged Attacks from bow-} {Type:-1 Text:us} {Type:1 Text:wield} {Type:0 Text:ing characters have 100% increased CRIT Rate. Additionally, Charged Attacks from bow-using characters will unleash a shockwave when they hit opponents, dealing one instance of AoE DMG. Can occur once every 1s.}]
[{Type:0 Text:Charged Attacks from bow-} {Type:-1 Text:using characters have 100% increased CRIT Rate} {Type:1 Text:wielding characters have 100% increased CRIT Rate} {Type:0 Text:. Additionally, Charged Attacks from bow-using characters will unleash a shockwave when they hit opponents, dealing one instance of AoE DMG. Can occur once every 1s.}]

```