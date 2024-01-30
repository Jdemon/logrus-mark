package main

import logm "github.com/Jdemon/logrus-mask"

func main() {
	logm.NewLogger(&logm.Config{
		Level: "debug",
		Masking: logm.ConfigMasking{
			Enabled: true,
		},
	}, "logrus-mask")
	logm.WithField("name", "test example").WithField("password", "P@ssw0rd").Info("test")
}
