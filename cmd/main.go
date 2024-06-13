package main

import (
	"context"
	"time"

	logm "github.com/Jdemon/logrus-mask"
)

func main() {
	logm.NewLogger(&logm.Config{
		Level: "debug",
		Masking: logm.ConfigMasking{
			Enabled: true,
		},
	}, "logrus-mask")
	ctx := context.WithValue(context.Background(), logm.TraceID, "trace-id-value")

	logm.WithContext(ctx).WithField("name", "test example").WithField("password", "P@ssw0rd").Info("test")
	logm.Client(ctx, logm.ClientLogModel{
		SourceIP: "127.0.0.1",
	}, time.Now())
	logm.Legacy(ctx, logm.LegacyLogModel{
		StepName: "test step name",
	}, time.Now())
	logm.MessageQueue(ctx, logm.MessageQueueLogModel{
		MessageQueueTopic: "topic",
	}, time.Now())
}
