package sqlx

import "testing"

func TestEscapeLikeQuery(t *testing.T) {
	// 測試輸入包含 LIKE 特殊字元與反斜線
	in := `a%b_c\dd`
	out := EscapeLikeQuery(in)

	// 預期：\ 先翻倍，% 與 _ 加上跳脫
	want := `a\%b\_c\\dd`

	if out != want {
		t.Fatalf("EscapeLikeQuery mismatch:\nwant: %q\ngot:  %q", want, out)
	}
}

func TestBuildLikeQueryValue(t *testing.T) {
	value := `a%b_c\dd`

	// prefix match: value%
	gotStart := BuildLikeQueryValue(value, LikePosStart)
	wantStart := `a\%b\_c\\dd%`
	if gotStart != wantStart {
		t.Fatalf("LikePosStart mismatch:\nwant: %q\ngot:  %q", wantStart, gotStart)
	}

	// suffix match: %value
	gotEnd := BuildLikeQueryValue(value, LikePosEnd)
	wantEnd := `%a\%b\_c\\dd`
	if gotEnd != wantEnd {
		t.Fatalf("LikePosEnd mismatch:\nwant: %q\ngot:  %q", wantEnd, gotEnd)
	}

	// both match: %value%
	gotBoth := BuildLikeQueryValue(value, LikePosBoth)
	wantBoth := `%a\%b\_c\\dd%`
	if gotBoth != wantBoth {
		t.Fatalf("LikePosBoth mismatch:\nwant: %q\ngot:  %q", wantBoth, gotBoth)
	}

	// default: only escaped
	gotDefault := BuildLikeQueryValue(value, "unknown")
	wantDefault := `a\%b\_c\\dd`
	if gotDefault != wantDefault {
		t.Fatalf("default mismatch:\nwant: %q\ngot:  %q", wantDefault, gotDefault)
	}
}

func TestLikeEscapeClause(t *testing.T) {
	// 這裡只檢查輸出字串固定
	want := `ESCAPE '\'`
	if LikeEscapeClause() != want {
		t.Fatalf("LikeEscapeClause mismatch:\nwant: %q\ngot:  %q", want, LikeEscapeClause())
	}
}

func TestEscapeSQLString(t *testing.T) {
	in := `a\b'c"d`
	out := EscapeSQLString(in)

	// \ -> \\ , ' -> \' , " -> \"
	want := `a\\b\'c\"d`

	if out != want {
		t.Fatalf("EscapeSQLString mismatch:\nwant: %q\ngot:  %q", want, out)
	}
}

func TestFormatSQLForLog(t *testing.T) {
	// 測試空白壓縮 + 反斜線還原
	in := "SELECT   *   FROM   t   WHERE  a = 'x\\\\y' "
	out := FormatSQLForLog(in)

	// Fields 壓縮後：SELECT * FROM t WHERE a = 'x\\y'
	// UnescapeBackslash：'x\y'
	want := "SELECT * FROM t WHERE a = 'x\\y'"

	if out != want {
		t.Fatalf("FormatSQLForLog mismatch:\nwant: %q\ngot:  %q", want, out)
	}
}
