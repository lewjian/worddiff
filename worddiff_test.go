package worddiff

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_split(t *testing.T) {
	type args struct {
		s                        string
		separators               map[string]struct{}
		mergeContinuousSeparator bool
	}
	tests := []struct {
		name string
		args args
		want []Word
	}{
		{
			name: "1",
			args: args{
				s: "hello world! nice to meet you",
				separators: map[string]struct{}{
					" ": {},
					"!": {},
				},
				mergeContinuousSeparator: false,
			},
			want: []Word{
				"hello", " ", "world", "!", " ", "nice", " ", "to", " ", "meet", " ", "you",
			},
		},
		{
			name: "2",
			args: args{
				s: "hello world!!!!   nice to meet you???？",
				separators: map[string]struct{}{
					" ": {},
					"!": {},
					"?": {},
				},
				mergeContinuousSeparator: true,
			},
			want: []Word{
				"hello", " ", "world", "!!!!", "   ", "nice", " ", "to", " ", "meet", " ", "you", "???", "？",
			},
		},
		{
			name: "3",
			args: args{
				s: "ashfiahuahfskhasfhaihfiahfi",
				separators: map[string]struct{}{
					" ": {},
					"!": {},
					"?": {},
				},
				mergeContinuousSeparator: true,
			},
			want: []Word{"ashfiahuahfskhasfhaihfiahfi"},
		},
		{
			name: "3",
			args: args{
				s: "ashfiahuahfskhasfhaihfiahfi ?!@$%",
				separators: map[string]struct{}{
					" ": {},
					"!": {},
					"?": {},
				},
				mergeContinuousSeparator: true,
			},
			want: []Word{"ashfiahuahfskhasfhaihfiahfi", " ", "?", "!", "@$%"},
		},
		{
			name: "4",
			args: args{
				s: "!!!!!!!!!!!!!",
				separators: map[string]struct{}{
					" ": {},
					"!": {},
					"?": {},
				},
				mergeContinuousSeparator: true,
			},
			want: []Word{"!!!!!!!!!!!!!"},
		},
		{
			name: "5",
			args: args{
				s: "The Wind and The Star Traveler",
				separators: map[string]struct{}{
					" ": {},
					"!": {},
					"?": {},
				},
				mergeContinuousSeparator: true,
			},
			want: []Word{"The", " ", "Wind", " ", "and", " ", "The", " ", "Star", " ", "Traveler"},
		},
		{
			name: "6",
			args: args{
				s:                        "我爱你中国",
				mergeContinuousSeparator: true,
			},
			want: []Word{"我", "爱", "你", "中", "国"},
		},
		{
			name: "6",
			args: args{
				s: "The  Hello??  !! come on！！<sad>",
				separators: map[string]struct{}{
					" ": {},
					"!": {},
					"?": {},
				},
				mergeContinuousSeparator: true,
			},
			want: []Word{"The", "  ", "Hello", "??", "  ", "!!", " ", "come", " ", "on！！<sad>"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := split(tt.args.s, tt.args.separators, tt.args.mergeContinuousSeparator); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("split() = %+v %d, want %+v %d", got, len(got), tt.want, len(tt.want))
			}
		})
	}
}

func Test_NewDefaultWd(t *testing.T) {
	wd := Default()
	s2 := "Charged Attacks from bow-wielding characters have 100% increased CRIT Rate. Additionally, Charged Attacks from bow-using characters will unleash a shockwave when they hit opponents, dealing one instance of AoE DMG. Can occur once every 1s."
	s1 := "Charged Attacks from bow-using characters have 100% increased CRIT Rate. Additionally, Charged Attacks from bow-using characters will unleash a shockwave when they hit opponents, dealing one instance of AoE DMG. Can occur once every 1s."
	results := wd.Diff(s1, s2)
	fmt.Printf("%+v\n", results)
	wd = New()
	results = wd.Diff(s1, s2)
	fmt.Printf("%+v\n", results)

}
