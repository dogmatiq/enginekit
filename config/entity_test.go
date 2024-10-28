package config_test

import (
	"testing"

	. "github.com/dogmatiq/enginekit/config"
	"github.com/dogmatiq/enginekit/config/internal/configbuilder"
	"github.com/dogmatiq/enginekit/internal/test"
	"github.com/dogmatiq/enginekit/protobuf/identitypb"
	"github.com/dogmatiq/enginekit/protobuf/uuidpb"
)

func testEntity[
	T Entity,
	B configbuilder.EntityBuilder[T],
](
	t *testing.T,
	build func(func(B)) T,
) {
	t.Run("identity", func(t *testing.T) {
		t.Run("it returns the normalized identity", func(t *testing.T) {
			entity := build(
				func(b B) {
					b.Identity(
						func(b *configbuilder.IdentityBuilder) {
							b.Name("name")
							b.Key("19CB98D5-DD17-4DAF-AE00-1B413B7B899A")
						},
					)
				},
			)

			test.Expect(
				t,
				"unexpected identity",
				entity.Identity(),
				&identitypb.Identity{
					Name: "name",
					Key:  uuidpb.MustParse("19cb98d5-dd17-4daf-ae00-1b413b7b899a"),
				},
			)
		})

		t.Run("it panics if there is no identity", func(t *testing.T) {
			entity := build(
				func(b B) {},
			)

			test.ExpectPanic(
				t,
				`entity has no identity`,
				func() {
					entity.Identity()
				},
			)
		})

		t.Run("it panics if there are multiple identities", func(t *testing.T) {
			entity := build(
				func(b B) {
					b.Identity(
						func(b *configbuilder.IdentityBuilder) {
							b.Name("name1")
							b.Key("b3c0591b-4049-4f10-974f-05c99d2d6d83")
						},
					)
					b.Identity(
						func(b *configbuilder.IdentityBuilder) {
							b.Name("name2")
							b.Key("ee4089e4-7692-42ee-a4f4-450772eb39ad")
						},
					)
				},
			)

			test.ExpectPanic(
				t,
				`entity has 2 identities`,
				func() {
					entity.Identity()
				},
			)
		})
	})
}
