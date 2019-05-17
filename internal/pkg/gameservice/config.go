package gameservice

import "2019_1_OPG_plus_2/internal/pkg/config"

var ServiceLocation = config.Game.GameServiceLocation

var Port = ":" + config.Game.GameServicePort
