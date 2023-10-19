package types

import "testing"

func TestModule_ParseImportPath(t *testing.T) {
	type fields struct {
		ID       string
		Version  string
		Path     string
		Dir      string
		Indirect bool
		Main     bool
	}
	type args struct {
		s string
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantPkgPath string
		wantName    string
		wantErr     bool
	}{
		{
			name: "",
			fields: fields{
				ID:       "dev",
				Version:  "1.0",
				Path:     "gg",
				Dir:      "/home/vitaly/Documents/work/my/gg",
				Indirect: false,
				Main:     true,
			},
			args: args{
				s: "~/internal/auth/JWTContextKey",
			},
			wantPkgPath: "",
			wantName:    "",
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Module{
				ID:       tt.fields.ID,
				Version:  tt.fields.Version,
				Path:     tt.fields.Path,
				Dir:      tt.fields.Dir,
				Indirect: tt.fields.Indirect,
				Main:     tt.fields.Main,
			}
			gotPkgPath, gotName, err := m.ParseImportPath(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("Module.ParseImportPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotPkgPath != tt.wantPkgPath {
				t.Errorf("Module.ParseImportPath() gotPkgPath = %v, want %v", gotPkgPath, tt.wantPkgPath)
			}
			if gotName != tt.wantName {
				t.Errorf("Module.ParseImportPath() gotName = %v, want %v", gotName, tt.wantName)
			}
		})
	}
}
