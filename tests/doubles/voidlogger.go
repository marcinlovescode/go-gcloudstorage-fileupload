package doubles

type VoidLogger struct{}

func (logger *VoidLogger) Debug(message interface{}, args ...interface{}) {
}
func (logger *VoidLogger) Info(message string, args ...interface{}) {
}
func (logger *VoidLogger) Warn(message string, args ...interface{}) {
}
func (logger *VoidLogger) Error(message interface{}, args ...interface{}) {
}
func (logger *VoidLogger) Fatal(message interface{}, args ...interface{}) {
}
