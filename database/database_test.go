package database

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewPostgresConn(t *testing.T) {
	Convey("Given a valid database URL and schema", t, func() {
		databaseURL := "postgres://username:password@localhost:5432/mydb"
		databaseSchema := "public"

		Convey("When NewPostgresConn is called", func() {
			conn := NewPostgresConn(databaseURL, databaseSchema)

			Convey("The connection should not be nil", func() {
				So(conn, ShouldNotBeNil)
			})
		})
	})
}
