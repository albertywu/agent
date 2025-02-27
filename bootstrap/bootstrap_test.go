package bootstrap

import (
	"context"
	"testing"

	"github.com/buildkite/agent/v3/bootstrap/shell"
	"github.com/buildkite/agent/v3/redaction"
	"github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/opentracer"
)

var agentNameTests = []struct {
	agentName string
	expected  string
}{
	{"My Agent", "My-Agent"},
	{":docker: My Agent", "-docker--My-Agent"},
	{"My \"Agent\"", "My--Agent-"},
}

func TestDirForAgentName(t *testing.T) {
	t.Parallel()

	for _, test := range agentNameTests {
		assert.Equal(t, test.expected, dirForAgentName(test.agentName))
	}
}

var repositoryNameTests = []struct {
	repositoryName string
	expected       string
}{
	{"git@github.com:acme-inc/my-project.git", "git-github-com-acme-inc-my-project-git"},
	{"https://github.com/acme-inc/my-project.git", "https---github-com-acme-inc-my-project-git"},
}

func TestDirForRepository(t *testing.T) {
	t.Parallel()

	for _, test := range repositoryNameTests {
		assert.Equal(t, test.expected, dirForRepository(test.repositoryName))
	}
}

func TestGetValuesToRedact(t *testing.T) {
	t.Parallel()

	redactConfig := []string{
		"*_PASSWORD",
		"*_TOKEN",
	}
	environment := map[string]string{
		"BUILDKITE_PIPELINE": "unit-test",
		"DATABASE_USERNAME":  "AzureDiamond",
		"DATABASE_PASSWORD":  "hunter2",
	}

	valuesToRedact := redaction.GetValuesToRedact(shell.DiscardLogger, redactConfig, environment)

	assert.Equal(t, []string{"hunter2"}, valuesToRedact)
}

func TestGetValuesToRedactEmpty(t *testing.T) {
	t.Parallel()

	redactConfig := []string{}
	environment := map[string]string{
		"FOO":                "BAR",
		"BUILDKITE_PIPELINE": "unit-test",
	}

	valuesToRedact := redaction.GetValuesToRedact(shell.DiscardLogger, redactConfig, environment)

	var expected []string
	assert.Equal(t, expected, valuesToRedact)
	assert.Equal(t, 0, len(valuesToRedact))
}

func TestStartTracing(t *testing.T) {
	oriCtx := context.Background()
	var err error

	// When there's no Datadog tracing address, the tracer should be a no-op.
	cfg := Config{}
	b := New(Config{})
	b.shell, err = shell.NewWithContext(oriCtx)
	if err != nil {
		assert.FailNow(t, "Unexpected error while createing shell: %v", err)
	}
	span, ctx, stopper := b.startTracing(oriCtx)
	assert.IsType(t, opentracing.NoopTracer{}, opentracing.GlobalTracer())
	span.Finish()
	stopper()
	assert.Equal(t, span, opentracing.SpanFromContext(ctx))

	// With the Datadog tracing backend, the global tracer should be from Datadog.
	cfg = Config{
		TracingBackend: "datadog",
	}
	b = New(cfg)
	b.shell, err = shell.NewWithContext(oriCtx)
	if err != nil {
		assert.FailNow(t, "Unexpected error while createing shell: %v", err)
	}
	span, ctx, stopper = b.startTracing(oriCtx)
	assert.IsType(t, opentracer.New(), opentracing.GlobalTracer())
	span.Finish()
	stopper()
	assert.Equal(t, span, opentracing.SpanFromContext(ctx))
}
