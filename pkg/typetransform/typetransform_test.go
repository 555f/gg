package typetransform

import (
	"strings"
	"testing"

	stdtypes "go/types"

	"github.com/555f/gg/pkg/types"
	"github.com/dave/jennifer/jen"
)

func TestParse(t *testing.T) {
	type args struct {
		valueID  jen.Code
		assignID jen.Code
		op       string
		t        any
	}
	tests := []struct {
		name          string
		args          args
		wantHasError  bool
		wantParseCode string
		wantParamID   string
	}{
		{
			name: "",
			args: args{
				valueID:  jen.Id("strValue"),
				assignID: jen.Id("strAssignValue"),
				op:       ":=",
				t:        &types.Basic{Kind: stdtypes.String},
			},
			wantHasError:  false,
			wantParseCode: "",
			wantParamID:   "strValue",
		},
		{
			name: "int",
			args: args{
				valueID:  jen.Id("intValue"),
				assignID: jen.Id("intAssignValue"),
				op:       ":=",
				t:        &types.Basic{Kind: stdtypes.Int},
			},
			wantHasError:  true,
			wantParseCode: "intAssignValue, err := gostrings.ParseInt[int](intValue, 10, 64)\nif err != nil {\n	return nil, err\n}",
			wantParamID:   "intAssignValue",
		},
		{
			name: "int8",
			args: args{
				valueID:  jen.Id("int8Value"),
				assignID: jen.Id("int8AssignValue"),
				op:       ":=",
				t:        &types.Basic{Kind: stdtypes.Int8},
			},
			wantHasError:  true,
			wantParseCode: "int8AssignValue, err := gostrings.ParseInt[int8](int8Value, 10, 8)\nif err != nil {\n	return nil, err\n}",
			wantParamID:   "int8AssignValue",
		},
		{
			name: "int16",
			args: args{
				valueID:  jen.Id("int16Value"),
				assignID: jen.Id("int16AssignValue"),
				op:       ":=",
				t:        &types.Basic{Kind: stdtypes.Int16},
			},
			wantHasError:  true,
			wantParseCode: "int16AssignValue, err := gostrings.ParseInt[int16](int16Value, 10, 16)\nif err != nil {\n	return nil, err\n}",
			wantParamID:   "int16AssignValue",
		},
		{
			name: "int32",
			args: args{
				valueID:  jen.Id("int32Value"),
				assignID: jen.Id("int32AssignValue"),
				op:       ":=",
				t:        &types.Basic{Kind: stdtypes.Int32},
			},
			wantHasError:  true,
			wantParseCode: "int32AssignValue, err := gostrings.ParseInt[int32](int32Value, 10, 32)\nif err != nil {\n	return nil, err\n}",
			wantParamID:   "int32AssignValue",
		},
		{
			name: "int64",
			args: args{
				valueID:  jen.Id("int64Value"),
				assignID: jen.Id("int64AssignValue"),
				op:       ":=",
				t:        &types.Basic{Kind: stdtypes.Int64},
			},
			wantHasError:  true,
			wantParseCode: "int64AssignValue, err := gostrings.ParseInt[int64](int64Value, 10, 64)\nif err != nil {\n	return nil, err\n}",
			wantParamID:   "int64AssignValue",
		},
		{
			name: "uint",
			args: args{
				valueID:  jen.Id("uintValue"),
				assignID: jen.Id("uintAssignValue"),
				op:       ":=",
				t:        &types.Basic{Kind: stdtypes.Uint},
			},
			wantHasError:  true,
			wantParseCode: "uintAssignValue, err := gostrings.ParseUint[uint](uintValue, 10, 64)\nif err != nil {\n	return nil, err\n}",
			wantParamID:   "uintAssignValue",
		},
		{
			name: "uint8",
			args: args{
				valueID:  jen.Id("uint8Value"),
				assignID: jen.Id("uint8AssignValue"),
				op:       ":=",
				t:        &types.Basic{Kind: stdtypes.Uint8},
			},
			wantHasError:  true,
			wantParseCode: "uint8AssignValue, err := gostrings.ParseUint[uint8](uint8Value, 10, 8)\nif err != nil {\n	return nil, err\n}",
			wantParamID:   "uint8AssignValue",
		},
		{
			name: "uint16",
			args: args{
				valueID:  jen.Id("uint16Value"),
				assignID: jen.Id("uint16AssignValue"),
				op:       ":=",
				t:        &types.Basic{Kind: stdtypes.Uint16},
			},
			wantHasError:  true,
			wantParseCode: "uint16AssignValue, err := gostrings.ParseUint[uint16](uint16Value, 10, 16)\nif err != nil {\n	return nil, err\n}",
			wantParamID:   "uint16AssignValue",
		},
		{
			name: "uint32",
			args: args{
				valueID:  jen.Id("uint32Value"),
				assignID: jen.Id("uint32AssignValue"),
				op:       ":=",
				t:        &types.Basic{Kind: stdtypes.Uint32},
			},
			wantHasError:  true,
			wantParseCode: "uint32AssignValue, err := gostrings.ParseUint[uint32](uint32Value, 10, 32)\nif err != nil {\n	return nil, err\n}",
			wantParamID:   "uint32AssignValue",
		},
		{
			name: "uint64",
			args: args{
				valueID:  jen.Id("uint64Value"),
				assignID: jen.Id("uint64AssignValue"),
				op:       ":=",
				t:        &types.Basic{Kind: stdtypes.Uint64},
			},
			wantHasError:  true,
			wantParseCode: "uint64AssignValue, err := gostrings.ParseUint[uint64](uint64Value, 10, 64)\nif err != nil {\n	return nil, err\n}",
			wantParamID:   "uint64AssignValue",
		},
		{
			name: "float32",
			args: args{
				valueID:  jen.Id("float32Value"),
				assignID: jen.Id("float32AssignValue"),
				op:       ":=",
				t:        &types.Basic{Kind: stdtypes.Float32},
			},
			wantHasError:  true,
			wantParseCode: "float32AssignValue, err := gostrings.ParseFloat[float32](float32Value, 10, 32)\nif err != nil {\n	return nil, err\n}",
			wantParamID:   "float32AssignValue",
		},
		{
			name: "time",
			args: args{
				valueID:  jen.Id("timeValue"),
				assignID: jen.Id("timeAssignValue"),
				op:       ":=",
				t: &types.Named{
					Name: "Time",
					Pkg: &types.PackageType{
						Name: "Time",
						Path: "time",
					},
				},
			},
			wantHasError:  true,
			wantParseCode: "timeAssignValue, err := time.Parse(time.RFC3339, timeValue)\nif err != nil {\n	return nil, err\n}",
			wantParamID:   "timeAssignValue",
		},
		{
			name: "duration",
			args: args{
				valueID:  jen.Id("durValue"),
				assignID: jen.Id("durAssignValue"),
				op:       ":=",
				t: &types.Named{
					Name: "Duration",
					Pkg: &types.PackageType{
						Name: "Duration",
						Path: "time",
					},
				},
			},
			wantHasError:  true,
			wantParseCode: "durAssignValue, err := time.ParseDuration(durValue)\nif err != nil {\n	return nil, err\n}",
			wantParamID:   "durAssignValue",
		},
		{
			name: "url",
			args: args{
				valueID:  jen.Id("urlValue"),
				assignID: jen.Id("urlAssignValue"),
				op:       ":=",
				t: &types.Named{
					Name: "URL",
					Pkg: &types.PackageType{
						Name: "URL",
						Path: "net/url",
					},
				},
			},
			wantHasError:  true,
			wantParseCode: "urlAssignValue, err := url.Parse(urlValue)\nif err != nil {\n	return nil, err\n}",
			wantParamID:   "urlAssignValue",
		},
		{
			name: "named string",
			args: args{
				valueID:  jen.Id("nsValue"),
				assignID: jen.Id("nsAssignValue"),
				op:       ":=",
				t: &types.Named{
					Name: "MyString",
					Pkg:  &types.PackageType{Name: "MyString", Path: "test/path"},
					Type: &types.Basic{Kind: stdtypes.String},
				},
			},
			wantHasError:  false,
			wantParseCode: "nsAssignValue := path.MyString(nsValue)",
			wantParamID:   "nsAssignValue",
		},
		{
			name: "named int",
			args: args{
				valueID:  jen.Id("niValue"),
				assignID: jen.Id("niAssignValue"),
				op:       ":=",
				t: &types.Named{
					Name: "MyInt",
					Pkg:  &types.PackageType{Name: "MyInt", Path: "test/path"},
					Type: &types.Basic{Kind: stdtypes.Int},
				},
			},
			wantHasError:  false,
			wantParseCode: "niAssignValue, err := gostrings.ParseInt[int](niValue, 10, 64)\nif err != nil {\n	return nil, err\n}",
			wantParamID:   "path.MyInt(niAssignValue)",
		},
		{
			name: "named int",
			args: args{
				valueID:  jen.Id("niValue"),
				assignID: jen.Id("niAssignValue"),
				op:       ":=",
				t: &types.Named{
					Name: "MyInt",
					Pkg: &types.PackageType{
						Name: "MyInt",
						Path: "test/path",
					},
					Type: &types.Basic{Kind: stdtypes.Int},
				},
			},
			wantHasError:  false,
			wantParseCode: "niAssignValue, err := gostrings.ParseInt[int](niValue, 10, 64)\nif err != nil {\n	return nil, err\n}",
			wantParamID:   "path.MyInt(niAssignValue)",
		},
		{
			name: "null Int",
			args: args{
				valueID:  jen.Id("niValue"),
				assignID: jen.Id("niAssignValue"),
				op:       ":=",
				t: &types.Named{
					Name: "Int",
					Pkg: &types.PackageType{
						Name: "Int",
						Path: "gopkg.in/guregu/null.v4",
					},
					Type: &types.Struct{},
				},
			},
			wantHasError:  false,
			wantParseCode: "niAssignValue, err := gostrings.ParseInt[int](niValue, 10, 64)\nif err != nil {\n	return nil, err\n}",
			wantParamID:   "nullv4.IntFrom(niAssignValue)",
		},
		{
			name: "null Float",
			args: args{
				valueID:  jen.Id("nfValue"),
				assignID: jen.Id("nfAssignValue"),
				op:       ":=",
				t: &types.Named{
					Name: "Float",
					Pkg: &types.PackageType{
						Name: "Float",
						Path: "gopkg.in/guregu/null.v4",
					},
					Type: &types.Struct{},
				},
			},
			wantHasError:  false,
			wantParseCode: "nfAssignValue, err := gostrings.ParseFloat[float64](nfValue, 10, 64)\nif err != nil {\n	return nil, err\n}",
			wantParamID:   "nullv4.FloatFrom(nfAssignValue)",
		},
		{
			name: "null Bool",
			args: args{
				valueID:  jen.Id("nbValue"),
				assignID: jen.Id("nbAssignValue"),
				op:       ":=",
				t: &types.Named{
					Name: "Bool",
					Pkg: &types.PackageType{
						Name: "Bool",
						Path: "gopkg.in/guregu/null.v4",
					},
					Type: &types.Struct{},
				},
			},
			wantHasError:  false,
			wantParseCode: "nbAssignValue, err := gostrings.ParseBool[bool](nbValue)\nif err != nil {\n	return nil, err\n}",
			wantParamID:   "nullv4.BoolFrom(nbAssignValue)",
		},
		{
			name: "null Time",
			args: args{
				valueID:  jen.Id("ntValue"),
				assignID: jen.Id("ntAssignValue"),
				op:       ":=",
				t: &types.Named{
					Name: "Time",
					Pkg: &types.PackageType{
						Name: "Time",
						Path: "gopkg.in/guregu/null.v4",
					},
					Type: &types.Struct{},
				},
			},
			wantHasError:  false,
			wantParseCode: "ntAssignValue, err := time.Parse(time.RFC3339, ntValue)\nif err != nil {\n	return nil, err\n}",
			wantParamID:   "nullv4.TimeFrom(ntAssignValue)",
		},
		{
			name: "slice int",
			args: args{
				valueID:  jen.Id("siValue"),
				assignID: jen.Id("siAssignValue"),
				op:       ":=",
				t: &types.Slice{
					Value: types.BasicTyp[stdtypes.Int],
				},
			},
			wantHasError:  true,
			wantParseCode: "siAssignValue, err := gostrings.SplitInt[int](siValue, \";\", 10, 64)\nif err != nil {\n	return nil, err\n}",
			wantParamID:   "siAssignValue",
		},
		{
			name: "slice string",
			args: args{
				valueID:  jen.Id("ssValue"),
				assignID: jen.Id("ssAssignValue"),
				op:       ":=",
				t: &types.Slice{
					Value: types.BasicTyp[stdtypes.String],
				},
			},
			wantHasError:  true,
			wantParseCode: "ssAssignValue, err := gostrings.Split(ssValue, \";\")\nif err != nil {\n	return nil, err\n}",
			wantParamID:   "ssAssignValue",
		},
		{
			name: "map string",
			args: args{
				valueID:  jen.Id("msValue"),
				assignID: jen.Id("msAssignValue"),
				op:       ":=",
				t: &types.Map{
					Value: types.BasicTyp[stdtypes.String],
				},
			},
			wantHasError:  true,
			wantParseCode: "msAssignValue, err := gostrings.SplitKeyValString[string](msValue, \",\", \"=\")\nif err != nil {\n	return nil, err\n}",
			wantParamID:   "msAssignValue",
		},
		{
			name: "map int",
			args: args{
				valueID:  jen.Id("miValue"),
				assignID: jen.Id("miAssignValue"),
				op:       ":=",
				t: &types.Map{
					Value: types.BasicTyp[stdtypes.Int],
				},
			},
			wantHasError:  true,
			wantParseCode: "miAssignValue, err := gostrings.SplitKeyValInt[int](miValue, \",\", \"=\", 10, 64)\nif err != nil {\n	return nil, err\n}",
			wantParamID:   "miAssignValue",
		},
		{
			name: "map int32",
			args: args{
				valueID:  jen.Id("miValue"),
				assignID: jen.Id("miAssignValue"),
				op:       ":=",
				t: &types.Map{
					Value: types.BasicTyp[stdtypes.Int32],
				},
			},
			wantHasError:  true,
			wantParseCode: "miAssignValue, err := gostrings.SplitKeyValInt[int32](miValue, \",\", \"=\", 10, 64)\nif err != nil {\n	return nil, err\n}",
			wantParamID:   "miAssignValue",
		},
		{
			name: "map uint",
			args: args{
				valueID:  jen.Id("muiValue"),
				assignID: jen.Id("muiAssignValue"),
				op:       ":=",
				t: &types.Map{
					Value: types.BasicTyp[stdtypes.Uint],
				},
			},
			wantHasError:  true,
			wantParseCode: "muiAssignValue, err := gostrings.SplitKeyValUint[uint](muiValue, \",\", \"=\", 10, 64)\nif err != nil {\n	return nil, err\n}",
			wantParamID:   "muiAssignValue",
		},
		{
			name: "map uint32",
			args: args{
				valueID:  jen.Id("muiValue"),
				assignID: jen.Id("muiAssignValue"),
				op:       ":=",
				t: &types.Map{
					Value: types.BasicTyp[stdtypes.Uint32],
				},
			},
			wantHasError:  true,
			wantParseCode: "muiAssignValue, err := gostrings.SplitKeyValUint[uint32](muiValue, \",\", \"=\", 10, 64)\nif err != nil {\n	return nil, err\n}",
			wantParamID:   "muiAssignValue",
		},
		{
			name: "map float32",
			args: args{
				valueID:  jen.Id("mfValue"),
				assignID: jen.Id("mfAssignValue"),
				op:       ":=",
				t: &types.Map{
					Value: types.BasicTyp[stdtypes.Float32],
				},
			},
			wantHasError:  true,
			wantParseCode: "mfAssignValue, err := gostrings.SplitKeyValFloat[float32](mfValue, \",\", \"=\", 64)\nif err != nil {\n	return nil, err\n}",
			wantParamID:   "mfAssignValue",
		},
		{
			name: "map float64",
			args: args{
				valueID:  jen.Id("mfValue"),
				assignID: jen.Id("mfAssignValue"),
				op:       ":=",
				t: &types.Map{
					Value: types.BasicTyp[stdtypes.Float64],
				},
			},
			wantHasError:  true,
			wantParseCode: "mfAssignValue, err := gostrings.SplitKeyValFloat[float64](mfValue, \",\", \"=\", 64)\nif err != nil {\n	return nil, err\n}",
			wantParamID:   "mfAssignValue",
		},
		{
			name: "google uuid",
			args: args{
				valueID:  jen.Id("guidValue"),
				assignID: jen.Id("guidAssignValue"),
				op:       ":=",
				t: &types.Named{
					Name: "UUID",
					Pkg: &types.PackageType{
						Name: "UUID",
						Path: "github.com/google/uuid",
					},
				},
			},
			wantHasError:  true,
			wantParseCode: "guidAssignValue, err := uuid.Parse(guidValue)\nif err != nil {\n	return nil, err\n}",
			wantParamID:   "guidAssignValue",
		},
		{
			name: "satori uuid",
			args: args{
				valueID:  jen.Id("guidValue"),
				assignID: jen.Id("guidAssignValue"),
				op:       ":=",
				t: &types.Named{
					Name: "UUID",
					Pkg: &types.PackageType{
						Name: "UUID",
						Path: "github.com/satori/go.uuid",
					},
				},
			},
			wantHasError:  true,
			wantParseCode: "guidAssignValue, err := gouuid.FromString(guidValue)\nif err != nil {\n	return nil, err\n}",
			wantParamID:   "guidAssignValue",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotParseCode, gotParamID, gotHasError := For(tt.args.t).
				SetValueID(tt.args.valueID).
				SetAssignID(tt.args.assignID).
				SetOp(tt.args.op).
				Parse()

			gotParamIDStr := jen.Add(gotParamID).GoString()

			if gotHasError != tt.wantHasError {
				t.Errorf("Parse() gotHasError = %v, want %v", gotHasError, tt.wantHasError)
			}

			gotParseCodeStr := strings.TrimLeft(jen.Add(gotParseCode).GoString(), "\n")
			if gotParseCodeStr != tt.wantParseCode {
				t.Errorf("Parse() gotParseCode = %v, want %v", gotParseCodeStr, tt.wantParseCode)
			}

			if gotParamIDStr != tt.wantParamID {
				t.Errorf("Parse() gotParamID = %v, want %v", gotParamIDStr, tt.wantParamID)
			}
		})
	}
}

func TestFormat(t *testing.T) {
	type args struct {
		valueID  jen.Code
		assignID jen.Code
		op       string
		t        any
	}
	tests := []struct {
		name          string
		args          args
		wantHasError  bool
		wantParseCode string
		wantParamID   string
	}{
		{
			name: "bool",
			args: args{
				valueID: jen.Lit("false"),
				op:      ":=",
				t:       types.BasicTyp[stdtypes.Bool],
			},
			wantHasError:  false,
			wantParseCode: "",
			wantParamID:   "strconv.FormatBool(\"false\")",
		},
		{
			name: "int",
			args: args{
				valueID: jen.Lit(10),
				op:      ":=",
				t:       types.BasicTyp[stdtypes.Int],
			},
			wantHasError:  false,
			wantParseCode: "",
			wantParamID:   "strconv.FormatInt(int64(10), 10)",
		},
		{
			name: "int8",
			args: args{
				valueID: jen.Lit(10),
				op:      ":=",
				t:       types.BasicTyp[stdtypes.Int8],
			},
			wantHasError:  false,
			wantParseCode: "",
			wantParamID:   "strconv.FormatInt(int64(10), 10)",
		},
		{
			name: "int16",
			args: args{
				valueID: jen.Lit(10),
				op:      ":=",
				t:       types.BasicTyp[stdtypes.Int16],
			},
			wantHasError:  false,
			wantParseCode: "",
			wantParamID:   "strconv.FormatInt(int64(10), 10)",
		},
		{
			name: "int32",
			args: args{
				valueID: jen.Lit(10),
				op:      ":=",
				t:       types.BasicTyp[stdtypes.Int32],
			},
			wantHasError:  false,
			wantParseCode: "",
			wantParamID:   "strconv.FormatInt(int64(10), 10)",
		},
		{
			name: "int64",
			args: args{
				valueID: jen.Lit(10),
				op:      ":=",
				t:       types.BasicTyp[stdtypes.Int64],
			},
			wantHasError:  false,
			wantParseCode: "",
			wantParamID:   "strconv.FormatInt(10, 10)",
		},
		{
			name: "time",
			args: args{
				valueID: jen.Id("date"),
				op:      ":=",
				t:       types.NamedTyp["time.Time"],
			},
			wantHasError:  false,
			wantParseCode: "",
			wantParamID:   "date.Format(time.RFC3339)",
		},
		{
			name: "null int",
			args: args{
				valueID: jen.Id("intVar"),
				op:      ":=",
				t: &types.Named{
					Name: "Int",
					Pkg: &types.PackageType{
						Name: "Int",
						Path: "gopkg.in/guregu/null.v4",
					},
				},
			},
			wantHasError:  false,
			wantParseCode: "",
			wantParamID:   "intVar.ValueOrZero()",
		},
		{
			name: "null string",
			args: args{
				valueID: jen.Id("strVar"),
				op:      ":=",
				t: &types.Named{
					Name: "Int",
					Pkg: &types.PackageType{
						Name: "String",
						Path: "gopkg.in/guregu/null.v4",
					},
				},
			},
			wantHasError:  false,
			wantParseCode: "",
			wantParamID:   "strVar.ValueOrZero()",
		},
		{
			name: "null float",
			args: args{
				valueID: jen.Id("floatVar"),
				op:      ":=",
				t: &types.Named{
					Name: "Int",
					Pkg: &types.PackageType{
						Name: "Float",
						Path: "gopkg.in/guregu/null.v4",
					},
				},
			},
			wantHasError:  false,
			wantParseCode: "",
			wantParamID:   "floatVar.ValueOrZero()",
		},
		{
			name: "null bool",
			args: args{
				valueID: jen.Id("boolVar"),
				op:      ":=",
				t: &types.Named{
					Name: "Int",
					Pkg: &types.PackageType{
						Name: "Bool",
						Path: "gopkg.in/guregu/null.v4",
					},
				},
			},
			wantHasError:  false,
			wantParseCode: "",
			wantParamID:   "boolVar.ValueOrZero()",
		},
		{
			name: "null time",
			args: args{
				valueID: jen.Id("timeVar"),
				op:      ":=",
				t: &types.Named{
					Name: "Int",
					Pkg: &types.PackageType{
						Name: "Time",
						Path: "gopkg.in/guregu/null.v4",
					},
				},
			},
			wantHasError:  false,
			wantParseCode: "",
			wantParamID:   "timeVar.ValueOrZero()",
		},
		{
			name: "slice int",
			args: args{
				valueID: jen.Id("intSlice"),
				op:      ":=",
				t: &types.Slice{
					Value: types.BasicTyp[stdtypes.Int],
				},
			},
			wantHasError:  false,
			wantParseCode: "",
			wantParamID:   "gostrings.JoinInt[int](intSlice, \",\", 10)",
		},
		{
			name: "slice int64",
			args: args{
				valueID: jen.Id("int64Slice"),
				op:      ":=",
				t: &types.Slice{
					Value: types.BasicTyp[stdtypes.Int64],
				},
			},
			wantHasError:  false,
			wantParseCode: "",
			wantParamID:   "gostrings.JoinInt[int64](int64Slice, \",\", 10)",
		},
		{
			name: "slice float32",
			args: args{
				valueID: jen.Id("float32Slice"),
				op:      ":=",
				t: &types.Slice{
					Value: types.BasicTyp[stdtypes.Float32],
				},
			},
			wantHasError:  false,
			wantParseCode: "",
			wantParamID:   "gostrings.JoinFloat[float32](float32Slice, \",\", int32(102), 2, 64)",
		},
		{
			name: "slice float64",
			args: args{
				valueID: jen.Id("float64Slice"),
				op:      ":=",
				t: &types.Slice{
					Value: types.BasicTyp[stdtypes.Float64],
				},
			},
			wantHasError:  false,
			wantParseCode: "",
			wantParamID:   "gostrings.JoinFloat[float64](float64Slice, \",\", int32(102), 2, 64)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotParseCode, gotParamID, gotHasError := For(tt.args.t).
				SetValueID(tt.args.valueID).
				SetAssignID(tt.args.assignID).
				SetOp(tt.args.op).
				Format()

			gotParamIDStr := jen.Add(gotParamID).GoString()

			if gotHasError != tt.wantHasError {
				t.Errorf("Parse() gotHasError = %v, want %v", gotHasError, tt.wantHasError)
			}

			gotParseCodeStr := strings.TrimLeft(jen.Add(gotParseCode).GoString(), "\n")
			if gotParseCodeStr != tt.wantParseCode {
				t.Errorf("Parse() gotParseCode = %v, want %v", gotParseCodeStr, tt.wantParseCode)
			}

			if gotParamIDStr != tt.wantParamID {
				t.Errorf("Parse() gotParamID = %v, want %v", gotParamIDStr, tt.wantParamID)
			}
		})
	}
}
