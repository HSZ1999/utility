package log

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToLevel(t *testing.T) {

	cases := []struct {
		name string
		set  any
		want Level
	}{
		{"levevl->info", INFO, INFO},
		{"integer->debug", 0, DEBUG},
		{"string->debug", "debug", DEBUG},
		{"string->info", "info", INFO},
		{"string->warning", "warning", WARN},
		{"string->error", "error", ERROR},
		{"string->fatal", "fatal", FATAL},
		{"string->unknown", "unknown", defaultLevel},
		{"struct->defaultLevel", struct{}{}, defaultLevel},
		{"edgeLevel->defaultLevel", Level(-18), Level(-18)},
		{"true->error", true, defaultLevel},
		{"false->error", false, defaultLevel},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, ToLevel(c.set), c.want)
		})
	}

}

func TestLevelString(t *testing.T) {
	cases := []struct {
		Name   string
		Level  Level
		Expect string
	}{
		{
			"< debug",
			Level(-1),
			"[Level(-1)]",
		},
		{
			"> fatal",
			Level(100),
			"[Level(100)]",
		},
		{
			"debug",
			DEBUG,
			levels[0],
		},
		{
			"info",
			INFO,
			levels[1],
		},
		{
			"warning",
			WARN,
			levels[2],
		},
		{
			"error",
			ERROR,
			levels[3],
		},
		{
			"fatal",
			FATAL,
			levels[4],
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			require.Equal(t, c.Expect, c.Level.String())
		})
	}
}

func TestLevelLimit(t *testing.T) {

	cases := []struct {
		Name         string
		Level        Level
		Args         []any
		ExpectArgs   string
		Format       string
		ExpectFormat string
	}{
		{
			"< debug",
			Level(-1),
			[]any{"hello", "world"},
			"helloworld\n",
			"prefix: %s, %s",
			"prefix: hello, world\n",
		},
		{
			"debug",
			DEBUG,
			[]any{"-word"},
			"-word\n",
			"%s suffix",
			"-word suffix\n",
		},
		{
			"info",
			INFO,
			[]any{"%s %q", "string", "integer"},
			"%s %qstringinteger\n",
			"%s %s %q",
			"%s %q string \"integer\"\n",
		},
		{
			"warning",
			WARN,
			[]any{1, 12, "string"},
			"1 12string\n",
			"%d, %o, %s", // 12 == 0o14
			"1, 14, string\n",
		},
		{
			"error",
			ERROR,
			[]any{},
			"\n",
			"",
			"\n",
		},
		{
			"> error",
			Level(100),
			[]any{"beyond"},
			"beyond\n",
			"-- %s --",
			"-- beyond --",
		},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			checkOutput(t, c.Level, c.Args, c.ExpectArgs, c.Format, c.ExpectFormat)
		})
	}
}

func checkOutput(t *testing.T, level Level, args []any, expectArgs string, format string, expectFormat string) {
	recorder := new(bytes.Buffer)
	SetOutput(recorder)
	SetLevel(level)
	SetFlags(0)
	SetPrefix("")

	if level > DEBUG {
		Debug(args...)
		require.Equal(t, "", recorder.String())
		Debugf(format, args...)
		require.Equal(t, "", recorder.String())
	}

	if level > INFO {
		Debug(args...)
		require.Equal(t, "", recorder.String())
		Debugf(format, args...)
		require.Equal(t, "", recorder.String())

		Info(args...)
		require.Equal(t, "", recorder.String())
		Infof(format, args...)
		require.Equal(t, "", recorder.String())
	}

	if level > WARN {
		Debug(args...)
		require.Equal(t, "", recorder.String())
		Debugf(format, args...)
		require.Equal(t, "", recorder.String())

		Info(args...)
		require.Equal(t, "", recorder.String())
		Infof(format, args...)
		require.Equal(t, "", recorder.String())

		Warn(args...)
		require.Equal(t, "", recorder.String())
		Warnf(format, args...)
		require.Equal(t, "", recorder.String())
	}

	if level > ERROR {
		Debug(args...)
		require.Equal(t, "", recorder.String())
		Debugf(format, args...)
		require.Equal(t, "", recorder.String())

		Info(args...)
		require.Equal(t, "", recorder.String())
		Infof(format, args...)
		require.Equal(t, "", recorder.String())

		Warn(args...)
		require.Equal(t, "", recorder.String())
		Warnf(format, args...)
		require.Equal(t, "", recorder.String())

		Error(args...)
		require.Equal(t, "", recorder.String())
		Errorf(format, args...)
		require.Equal(t, "", recorder.String())
	}

	if level < DEBUG {
		recorder.Reset()
		Debug(args...)
		require.Equal(t, DEBUG.String()+expectArgs, recorder.String())

		recorder.Reset()
		Debugf(format, args...)
		require.Equal(t, DEBUG.String()+expectFormat, recorder.String())
	}

	if level < INFO {
		recorder.Reset()
		Info(args...)
		require.Equal(t, INFO.String()+expectArgs, recorder.String())

		recorder.Reset()
		Infof(format, args...)
		require.Equal(t, INFO.String()+expectFormat, recorder.String())
	}

	if level < WARN {
		recorder.Reset()
		Warn(args...)
		require.Equal(t, WARN.String()+expectArgs, recorder.String())

		recorder.Reset()
		Warnf(format, args...)
		require.Equal(t, WARN.String()+expectFormat, recorder.String())
	}

	if level < ERROR {
		recorder.Reset()
		Error(args...)
		require.Equal(t, ERROR.String()+expectArgs, recorder.String())

		recorder.Reset()
		Errorf(format, args...)
		require.Equal(t, ERROR.String()+expectFormat, recorder.String())
	}
}

func TestConfig(t *testing.T) {
	require.Equal(t, logger, DefaultLogger())
	newLog := new(defaultLogger)
	SetLogger(newLog)
	require.Equal(t, newLog, DefaultLogger())
}
