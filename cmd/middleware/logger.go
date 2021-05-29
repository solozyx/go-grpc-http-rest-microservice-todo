package middleware

import (
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func codeToLevel(code codes.Code) zapcore.Level {
	if code == codes.OK {
		return zap.DebugLevel
	}
	return grpc_zap.DefaultCodeToLevel(code)
}

// 生产环境无法打断点，只能看日志，微服务 小 轻 分布式部署
// 你调的服务，你不知道它在什么地方，只能靠日志排查问题
// 用户服务返回1个错误，如果不打日志，那其他人的错误就是我的错误，你要背锅
// 一般在生产环境，一个请求在多少毫秒内返回就ok，只要报错了都要具体去看，日志对于微服务非常重要

// 如果微服务某个系统挂了，打印错误日志、大屏爆红、
// 生产环境 info 级别的日志是不看的，info只是提示性的东西
// 每天几百笔交易、几千笔交易 都无所谓，所有业务都能正常运行，几十毫秒就处理完了
// debug info error fatal

func AddLogging(logger *zap.Logger, opts []grpc.ServerOption) []grpc.ServerOption {
	o := []grpc_zap.Option{
		grpc_zap.WithLevels(codeToLevel),
	}

	grpc_zap.ReplaceGrpcLoggerV2(logger)

	opts = append(opts, grpc_middleware.WithUnaryServerChain(
		grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
		grpc_zap.UnaryServerInterceptor(logger, o...),
	))

	opts = append(opts, grpc_middleware.WithStreamServerChain(
		grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
		grpc_zap.StreamServerInterceptor(logger, o...),
	))
	return opts
}
