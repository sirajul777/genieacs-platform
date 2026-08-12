module github.com/sirajul777/genieacs-platform

go 1.25

require (
	github.com/go-chi/chi/v5 v5.2.3
	github.com/spf13/cobra v1.9.1
	github.com/spf13/viper v1.20.1
	go.uber.org/zap v1.27.0
)

replace github.com/go-chi/chi/v5 => ./third_party/github.com/go-chi/chi/v5

replace github.com/spf13/cobra => ./third_party/github.com/spf13/cobra

replace github.com/spf13/viper => ./third_party/github.com/spf13/viper

replace go.uber.org/zap => ./third_party/go.uber.org/zap
