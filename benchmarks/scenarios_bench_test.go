// Copyright (c) 2016 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package benchmarks

import (
	"context"
	"log"
	"log/slog"
	"testing"

	"github.com/lthphuw/plogger"
	"github.com/lthphuw/plogger/benchmarks/internal/ztest"
	"go.uber.org/zap"
)

func BenchmarkWithoutFields(b *testing.B) {
	b.Logf("")
	b.Logf("Logging without any structured context.")

	b.Run("plogger", func(b *testing.B) {
		logger := newPLogger()
		entry := plogger.NewEntry().SetMsg(getMessage(0))
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(entry)
			}
		})
	})
	b.Run("Zap", func(b *testing.B) {
		logger := newZapLogger(zap.DebugLevel)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
	})
	b.Run("Zap.Check", func(b *testing.B) {
		logger := newZapLogger(zap.DebugLevel)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				if ce := logger.Check(zap.InfoLevel, getMessage(0)); ce != nil {
					ce.Write()
				}
			}
		})
	})
	b.Run("Zap.CheckSampled", func(b *testing.B) {
		logger := newSampledLogger(zap.DebugLevel)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				i++
				if ce := logger.Check(zap.InfoLevel, getMessage(i)); ce != nil {
					ce.Write()
				}
			}
		})
	})
	b.Run("Zap.Sugar", func(b *testing.B) {
		logger := newZapLogger(zap.DebugLevel).Sugar()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
	})
	b.Run("Zap.SugarFormatting", func(b *testing.B) {
		logger := newZapLogger(zap.DebugLevel).Sugar()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Infof("%v %v %v %s %v %v %v %v %v %s\n", fakeFmtArgs()...)
			}
		})
	})
	b.Run("apex/log", func(b *testing.B) {
		logger := newApexLog()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
	})
	b.Run("go-kit/kit/log", func(b *testing.B) {
		logger := newKitLog()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				if err := logger.Log(getMessage(0), getMessage(1)); err != nil {
					b.Fatal(err)
				}
			}
		})
	})
	b.Run("inconshreveable/log15", func(b *testing.B) {
		logger := newLog15()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
	})
	b.Run("sirupsen/logrus", func(b *testing.B) {
		logger := newLogrus()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
	})
	b.Run("stdlib.Println", func(b *testing.B) {
		logger := log.New(&ztest.Discarder{}, "", log.LstdFlags)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Println(getMessage(0))
			}
		})
	})
	b.Run("stdlib.Printf", func(b *testing.B) {
		logger := log.New(&ztest.Discarder{}, "", log.LstdFlags)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Printf("%v %v %v %s %v %v %v %v %v %s\n", fakeFmtArgs()...)
			}
		})
	})
	b.Run("rs/zerolog", func(b *testing.B) {
		logger := newZerolog()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info().Msg(getMessage(0))
			}
		})
	})
	b.Run("rs/zerolog.Formatting", func(b *testing.B) {
		logger := newZerolog()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info().Msgf("%v %v %v %s %v %v %v %v %v %s\n", fakeFmtArgs()...)
			}
		})
	})
	b.Run("rs/zerolog.Check", func(b *testing.B) {
		logger := newZerolog()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				if e := logger.Info(); e.Enabled() {
					e.Msg(getMessage(0))
				}
			}
		})
	})
	b.Run("slog", func(b *testing.B) {
		logger := newSlog()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0))
			}
		})
	})
	b.Run("slog.LogAttrs", func(b *testing.B) {
		logger := newSlog()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.LogAttrs(context.Background(), slog.LevelInfo, getMessage(0))
			}
		})
	})
}

func BenchmarkAddingFields(b *testing.B) {
	b.Logf("")
	b.Logf("Logging with additional context at each log site.")

	b.Run("plogger", func(b *testing.B) {
		logger := newPLogger()
		entry := plogger.NewEntry().SetMsg(getMessage(0)).SetFieldMap(ploggerFieldMap())
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(entry)
			}
		})
	})
	b.Run("Zap", func(b *testing.B) {
		logger := newZapLogger(zap.DebugLevel)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0), fakeFields()...)
			}
		})
	})
	b.Run("Zap.Check", func(b *testing.B) {
		logger := newZapLogger(zap.DebugLevel)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				if ce := logger.Check(zap.InfoLevel, getMessage(0)); ce != nil {
					ce.Write(fakeFields()...)
				}
			}
		})
	})
	b.Run("Zap.CheckSampled", func(b *testing.B) {
		logger := newSampledLogger(zap.DebugLevel)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				i++
				if ce := logger.Check(zap.InfoLevel, getMessage(i)); ce != nil {
					ce.Write(fakeFields()...)
				}
			}
		})
	})
	b.Run("Zap.Sugar", func(b *testing.B) {
		logger := newZapLogger(zap.DebugLevel).Sugar()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Infow(getMessage(0), fakeSugarFields()...)
			}
		})
	})
	b.Run("apex/log", func(b *testing.B) {
		logger := newApexLog()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.WithFields(fakeApexFields()).Info(getMessage(0))
			}
		})
	})
	b.Run("go-kit/kit/log", func(b *testing.B) {
		logger := newKitLog()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				if err := logger.Log(fakeSugarFields()...); err != nil {
					b.Fatal(err)
				}
			}
		})
	})
	b.Run("inconshreveable/log15", func(b *testing.B) {
		logger := newLog15()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0), fakeSugarFields()...)
			}
		})
	})
	b.Run("sirupsen/logrus", func(b *testing.B) {
		logger := newLogrus()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.WithFields(fakeLogrusFields()).Info(getMessage(0))
			}
		})
	})
	b.Run("rs/zerolog", func(b *testing.B) {
		logger := newZerolog()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				fakeZerologFields(logger.Info()).Msg(getMessage(0))
			}
		})
	})
	b.Run("rs/zerolog.Check", func(b *testing.B) {
		logger := newZerolog()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				if e := logger.Info(); e.Enabled() {
					fakeZerologFields(e).Msg(getMessage(0))
				}
			}
		})
	})
	b.Run("slog", func(b *testing.B) {
		logger := newSlog()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info(getMessage(0), fakeSlogArgs()...)
			}
		})
	})
	b.Run("slog.LogAttrs", func(b *testing.B) {
		logger := newSlog()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.LogAttrs(
					context.Background(),
					slog.LevelInfo,
					getMessage(0),
					fakeSlogFields()...)
			}
		})
	})
}
