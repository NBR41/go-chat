
use platform;

DROP TABLE if EXISTS `users`;
CREATE TABLE `users` (
    `user_id` bigint(20) unsigned not null AUTO_INCREMENT,
    `email` varchar(256) COLLATE utf8_unicode_ci not null,
    `pseudo` varchar(64) COLLATE utf8_unicode_ci not null,
    `password` varchar(32) not null,
    `super_admin` TINYINT(1) unsigned NOT NULL DEFAULT '0',
    PRIMARY KEY (`user_id`),
    UNIQUE KEY `uk_email` (`email`(200)),
    UNIQUE KEY `uk_pseudo` (`pseudo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

DROP TABLE if EXISTS `roles`;
CREATE TABLE `roles` (
    `role_id`  bigint(20) unsigned not null AUTO_INCREMENT,
    `name` varchar(64) COLLATE utf8_unicode_ci not null,
    PRIMARY KEY (`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

DROP TABLE if EXISTS `rooms`;
CREATE TABLE `rooms` (
    `room_id` bigint(20) unsigned not null AUTO_INCREMENT,
    `name` varchar(64) COLLATE utf8_unicode_ci not null,
    `private` TINYINT(1) unsigned NOT NULL DEFAULT '0',
    PRIMARY KEY (`room_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

DROP TABLE if EXISTS `room_roles`;
CREATE TABLE `room_roles` (
    `id` bigint(20) unsigned not null AUTO_INCREMENT,
    `user_id` bigint(20) unsigned not null,
    `role_id` bigint(20) unsigned not null,
    PRIMARY KEY (`ID`),
    CONSTRAINT `fk_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`),
    CONSTRAINT `fk_role_id` FOREIGN KEY (`role_id`) REFERENCES `roles` (`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

insert into `roles` values (1, 'super'), (2, 'admin'), (3, 'user');
