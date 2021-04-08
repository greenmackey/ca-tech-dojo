CREATE TABLE IF NOT EXISTS `users` (
  `token` varchar(50) NOT NULL,
  `name` varchar(20) NOT NULL,
  `id` int NOT NULL AUTO_INCREMENT,
  PRIMARY KEY (`id`),
  UNIQUE KEY `token` (`token`)
);

CREATE TABLE IF NOT EXISTS `characters` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(20) NOT NULL,
  `likelihood` int UNSIGNED NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
);

CREATE TABLE IF NOT EXISTS `rel_user_character` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_token` varchar(50) NOT NULL,
  `character_id` int DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `character_id` (`character_id`),
  KEY `user_token` (`user_token`),
  CONSTRAINT `rel_user_character_ibfk_2` FOREIGN KEY (`character_id`) REFERENCES `characters` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `rel_user_character_ibfk_3` FOREIGN KEY (`user_token`) REFERENCES `users` (`token`) ON DELETE CASCADE ON UPDATE CASCADE
);