package config_test

import (
	"context"
	"errors"

	. "github.com/dogmatiq/enginekit/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("type FuncVisitor", func() {
	entries := []TableEntry{
		Entry(
			"ApplicationConfig",
			&ApplicationConfig{
				ApplicationIdentity: MustNewIdentity("<app-name>", "<app-key>"),
			},
		),
		Entry(
			"AggregateConfig",
			&AggregateConfig{
				HandlerIdentity: MustNewIdentity("<aggregate-name>", "<aggregate-key>"),
			},
		),
		Entry(
			"ProcessConfig",
			&ProcessConfig{
				HandlerIdentity: MustNewIdentity("<process-name>", "<process-key>"),
			},
		),
		Entry(
			"IntegrationConfig",
			&IntegrationConfig{
				HandlerIdentity: MustNewIdentity("<integration-name>", "<integration-key>"),
			},
		),
		Entry(
			"ProjectionConfig",
			&ProjectionConfig{
				HandlerIdentity: MustNewIdentity("<projection-name>", "<projection-key>"),
			},
		),
	}

	When("the registered function is nil", func() {
		v := &FuncVisitor{}

		DescribeTable(
			"does not return an error",
			func(cfg Config) {
				err := cfg.Accept(context.Background(), v)
				Expect(err).ShouldNot(HaveOccurred())
			},
			entries...,
		)
	})

	When("the registered function is non-nil", func() {
		var arg Config
		v := &FuncVisitor{
			ApplicationConfig: func(_ context.Context, cfg *ApplicationConfig) error {
				arg = cfg
				return errors.New("<error>")
			},
			AggregateConfig: func(_ context.Context, cfg *AggregateConfig) error {
				arg = cfg
				return errors.New("<error>")
			},
			ProcessConfig: func(_ context.Context, cfg *ProcessConfig) error {
				arg = cfg
				return errors.New("<error>")
			},
			IntegrationConfig: func(_ context.Context, cfg *IntegrationConfig) error {
				arg = cfg
				return errors.New("<error>")
			},
			ProjectionConfig: func(_ context.Context, cfg *ProjectionConfig) error {
				arg = cfg
				return errors.New("<error>")
			},
		}

		BeforeEach(func() {
			arg = nil
		})

		DescribeTable(
			"does not return an error",
			func(cfg Config) {
				err := cfg.Accept(context.Background(), v)
				Expect(err).To(MatchError("<error>"))
				Expect(arg).To(Equal(cfg))
			},
			entries...,
		)
	})
})
