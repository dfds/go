package log

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"runtime"
	"strings"
)

var Logger *zap.Logger

func InitializeLogger(debug bool, logLevel string) {
	var logConf zap.Config
	if debug {
		logConf = zap.NewDevelopmentConfig()
	} else {
		logConf = zap.NewProductionConfig()
	}

	level, err := zapcore.ParseLevel(logLevel)
	if err != nil {
		fmt.Println(err)
		level = zapcore.InfoLevel
	}

	logConf.Level = zap.NewAtomicLevelAt(level)
	logConf.EncoderConfig.FunctionKey = "func"

	Logger, _ = logConf.Build(zap.AddStacktrace(zapcore.ErrorLevel), zap.AddCaller())
	Logger = Logger.WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return &FilterCore{
			Core: core,
			ignoreFields: map[string][]string{
				//"package": []string{
				//	"go.dfds.cloud/ssu-k8s/feats/operator/controller",
				//},
			},
			ignoreContainsFields: map[string][]string{
				"func": []string{
					"*NamespaceReconciler",
				},
			},
		}
	}))

	//Logger = zap.New(&FilterCore{
	//	Core:         Logger.Core(),
	//	ignoreFields: map[string][]string{},
	//}, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	Logger.Info(fmt.Sprintf("Logging enabled, log level set to %s", Logger.Level().String()))
}

type FilterCore struct {
	zapcore.Core
	ignoreFields         map[string][]string
	ignoreContainsFields map[string][]string
	withFields           []zapcore.Field
}

func (f *FilterCore) With(fields []zapcore.Field) zapcore.Core {
	newFields := append([]zapcore.Field{}, f.withFields...)
	newFields = append(newFields, fields...)
	return &FilterCore{
		Core:                 f.Core.With(fields),
		ignoreFields:         f.ignoreFields,
		ignoreContainsFields: f.ignoreContainsFields,
		withFields:           newFields,
	}
	//return f.Core.With(fields)
}

func (f *FilterCore) Check(entry zapcore.Entry, checkedEntry *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	//	return f.Core.Check(entry, checkedEntry)

	return checkedEntry.AddCore(entry, f)
}

func (f *FilterCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	allFields := append(f.withFields, fields...)
	allFields = append(allFields, PackageField(4))

	fmt.Println("Fields: ", allFields)

	for _, field := range allFields {
		if values, ok := f.ignoreFields[field.Key]; ok {
			var str string
			if field.Type == zapcore.StringType {
				str = field.String
			}

			for _, k := range values {
				if str == k {
					return nil
				}
			}
		}

		fmt.Println(field.Key)
		if values, ok := f.ignoreContainsFields[field.Key]; ok {
			var str string
			if field.Type == zapcore.StringType {
				str = field.String
			}

			for _, k := range values {
				if strings.Contains(str, k) {
					return nil
				}
			}
		}
	}
	return f.Core.Write(entry, allFields)
}

func PackageField(skip int) zap.Field {
	pc, _, _, ok := runtime.Caller(skip)
	if !ok {
		return zap.String("package", "unknown")
	}
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return zap.String("package", "unknown")
	}
	fullFuncName := fn.Name()

	parts := strings.Split(fullFuncName, "/")
	lastPart := strings.Split(parts[len(parts)-1], ".")
	packageName := strings.Replace(fullFuncName, fmt.Sprintf("/%s", parts[len(parts)-1]), "", -1)
	packageName = fmt.Sprintf("%s/%s", packageName, lastPart[0])

	return zap.String("package", packageName)
}
