package config

import("time";sharedconfig "github.com/sirajul777/genieacs-platform/internal/shared/config";"github.com/sirajul777/genieacs-platform/internal/shared/logger")
type Config struct{Server ServerConfig `mapstructure:"server"`;Manager ManagerConfig `mapstructure:"manager"`;Logger logger.Config `mapstructure:"logger"`}
type ServerConfig struct{Address string `mapstructure:"address"`}
type ManagerConfig struct{Endpoint string `mapstructure:"endpoint"`;AgentID string `mapstructure:"agent_id"`;Token string `mapstructure:"token"`;Name string `mapstructure:"name"`;StateFile string `mapstructure:"state_file"`;PollInterval time.Duration `mapstructure:"poll_interval"`;HeartbeatInterval time.Duration `mapstructure:"heartbeat_interval"`}
func Load(file string)(Config,error){var cfg Config;err:=sharedconfig.NewLoader("agent",file).Load(&cfg,sharedconfig.Default{Key:"server.address",Value:":8081"},sharedconfig.Default{Key:"manager.endpoint",Value:"http://localhost:8080"},sharedconfig.Default{Key:"manager.state_file",Value:"agent-state.json"},sharedconfig.Default{Key:"manager.poll_interval",Value:"2s"},sharedconfig.Default{Key:"manager.heartbeat_interval",Value:"15s"},sharedconfig.Default{Key:"logger.environment",Value:"development"},sharedconfig.Default{Key:"logger.level",Value:"info"});return cfg,err}
