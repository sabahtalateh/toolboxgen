package discovery

import (
	"errors"
	"github.com/sabahtalateh/toolboxgen/internal/discovery"
	"github.com/sabahtalateh/toolboxgen/internal/mod"
	"os"
	"path/filepath"
	"testing"

	"github.com/sabahtalateh/toolboxgen/tests"
)

func TestConvert(t *testing.T) {
	type testCase struct {
		name    string
		dir     string
		wantErr error
	}
	ttt := []testCase{
		{name: "nested", dir: "mod/nested", wantErr: mod.ErrNotModRoot},
	}
	for _, tt := range ttt {
		t.Run(tt.name, func(t *testing.T) {
			dir := tests.Unwrap(os.Getwd())

			r, err := discovery.Discover(filepath.Join(dir, tt.dir))
			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("\ngot error:\n%s\nwant error:\n%s", err, tt.wantErr)
				return
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
