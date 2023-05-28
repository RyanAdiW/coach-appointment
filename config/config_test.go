package config

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLoadConfig(t *testing.T) {
	Convey("Given a valid config file", t, func() {
		path := "."
		filename := "config.yaml"

		Convey("When LoadConfig is called", func() {
			config, err := LoadConfig(path, filename)

			Convey("The error should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("The database URL should be correct", func() {
				So(config.DatabaseURL, ShouldEqual, "postgres://postgres:mysecret@localhost:5432/fita-test?sslmode=disable")
			})

			Convey("The database schema should be correct", func() {
				So(config.DatabaseSchema, ShouldEqual, "public")
			})
		})
	})

	Convey("Given an invalid config file", t, func() {
		path := "."
		filename := "invalid_config.yaml"

		Convey("When LoadConfig is called", func() {
			_, err := LoadConfig(path, filename)

			Convey("The error should not be nil", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}
