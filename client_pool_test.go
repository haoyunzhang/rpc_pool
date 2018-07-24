package sdk_test

import (
	"net"
	"testing"
	"time"

	"github.com/maraino/go-mock"
	"github.com/smartystreets/goconvey/convey"
	"hfga.com.cn/tonghangyun/sdk"
)

func TestRPCClientPool(t *testing.T) {

	convey.Convey("初始化一个初始长度为0的缓存池", t, func() {
		stubConn := new(fakeConn)
		closeTime := stubConn.When("Close").Return(nil)
		stub := new(factoryFakeImpl)
		connTime := stub.When("Connect").Return(stubConn, nil)

		pool, err := sdk.NewRPCClientPool(0, 1, stub.Connect)
		convey.So(err, convey.ShouldBeNil)

		convey.Convey("如果没有实际请求链接", func() {
			connTime.Times(0)

			convey.Convey("默认的已链接的长度为0", func() {
				ok, err := stubConn.Verify()
				convey.So(ok, convey.ShouldBeTrue)
				convey.So(err, convey.ShouldBeNil)
				convey.So(pool.Len(), convey.ShouldEqual, 0)
			})
		})

		convey.Convey("如果请求获取链接", func() {
			connTime.Times(1)
			client, _ := pool.Get()

			convey.Convey("已链接的长度为0", func() {
				ok, err := stub.Verify()
				convey.So(ok, convey.ShouldBeTrue)
				convey.So(err, convey.ShouldBeNil)
				convey.So(client, convey.ShouldNotBeNil)
				convey.So(err, convey.ShouldBeNil)
				convey.So(pool.Len(), convey.ShouldEqual, 0)
			})

			convey.Convey("归还链接", func() {
				connTime.Times(1)
				closeTime.Times(0)
				client.Close()

				convey.Convey("已链接的长度为1", func() {
					ok, err := stub.Verify()
					convey.So(ok, convey.ShouldBeTrue)
					convey.So(err, convey.ShouldBeNil)
					ok, err = stubConn.Verify()
					convey.So(ok, convey.ShouldBeTrue)
					convey.So(err, convey.ShouldBeNil)

					convey.So(pool.Len(), convey.ShouldEqual, 1)
				})
			})

			convey.Convey("Close all the pool", func() {
				closeTime.Times(1)
				client.Close()
				pool.Close()

				convey.Convey("已链接的长度为0", func() {
					convey.So(pool.Len(), convey.ShouldEqual, 0)
					ok, err := stubConn.Verify()
					convey.So(ok, convey.ShouldBeTrue)
					convey.So(err, convey.ShouldBeNil)
				})
			})
		})
	})
}

type fakeConn struct {
	mock.Mock
}

func (obj *fakeConn) Close() error {
	ret := obj.Called()
	return ret.Error(0)
}

func (obj *fakeConn) Read(b []byte) (n int, err error) {
	return
}

func (obj *fakeConn) Write(b []byte) (n int, err error) {
	return
}

func (obj *fakeConn) LocalAddr() net.Addr {
	return nil
}

func (obj *fakeConn) RemoteAddr() net.Addr {
	return nil
}

func (obj *fakeConn) SetDeadline(t time.Time) error {
	return nil
}

func (obj *fakeConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (obj *fakeConn) SetWriteDeadline(t time.Time) error {
	return nil
}

type factoryFakeImpl struct {
	mock.Mock
}

func (obj *factoryFakeImpl) Connect() (net.Conn, error) {

	ret := obj.Called()
	return ret.Get(0).(*fakeConn), ret.Error(1)
}
