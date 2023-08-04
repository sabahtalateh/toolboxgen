package discovery

import (
	"errors"
	"github.com/sabahtalateh/toolboxgen/internal/discovery"
	"github.com/sabahtalateh/toolboxgen/internal/mod"
	tutils "github.com/sabahtalateh/toolboxgen/tests"
	"os"
	"path/filepath"
	"testing"
)

func TestConvert(t *testing.T) {
	type testCase struct {
		name    string
		dir     string
		wantErr error
	}
	ttt := []testCase{
		{name: "not-mod-root", dir: "mod/nested", wantErr: mod.ErrNotModRoot},
		{name: "test", dir: "mod"},
	}
	for _, tt := range ttt {
		t.Run(tt.name, func(t *testing.T) {
			dir := tutils.Unwrap(os.Getwd())

			r, err := discovery.Discover(filepath.Join(dir, tt.dir))
			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("\ngot error:\n%s\nwant error:\n%s", err, tt.wantErr)
				return
			}
			if err != nil {
				panic(err)
			}

			println(r)

			// bb := tests.Unwrap(os.ReadFile(filepath.Join(dir, tt.dir, "want.yaml")))
			//
			// var want map[string]any
			// tests.Check(yaml3.Unmarshal(bb, &want))

			// got := convertFile(tt.dir, tt.dir)

			// if !reflect.DeepEqual(got, want) {
			// 	g := tests.Unwrap(yaml2.Marshal(ordered(got)))
			// 	w := tests.Unwrap(yaml2.Marshal(ordered(want)))
			//
			// 	t.Errorf("\ngot:\n\n%s\nwant\n\n%s", g, w)
			// }
		})
	}
}
