package models

import "time"

var Config struct {
	Base struct {
		Root   string `conf:"default:Professor Zoom"`
		Import string `conf:"default:import"`
	}
	Web struct {
		APIHost         string        `conf:"default:0.0.0.0:3000"`
		DebugHost       string        `conf:"default:0.0.0.0:4000"`
		ReadTimeout     time.Duration `conf:"default:5s"`
		WriteTimeout    time.Duration `conf:"default:5s"`
		ShutdownTimeout time.Duration `conf:"default:5s"`
	}
	DB struct {
		Driver     string `conf:"default:postgres"`
		Connection string `conf:"default:host=db port=5432 user=postgres dbname=postgres password=postgres sslmode=disable,noprint"`
		Post       string `conf:"default:CREATE SCHEMA IF NOT EXISTS \"{athlete}\"; set search_path='{athlete}'"`
	}
}
