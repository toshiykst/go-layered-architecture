CREATE TABLE IF NOT EXISTS `users`
(
    `id`         VARCHAR(255) PRIMARY KEY NOT NULL,
    `name`       VARCHAR(255)             NOT NULL,
    `email`      VARCHAR(255)             NOT NULL,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

CREATE TABLE IF NOT EXISTS `groups`
(
    `id`         VARCHAR(255) PRIMARY KEY NOT NULL,
    `name`       VARCHAR(255)             NOT NULL,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

CREATE TABLE IF NOT EXISTS `group_users`
(
    `group_id`   VARCHAR(255) NOT NULL,
    `user_id`    VARCHAR(255) NOT NULL,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`group_id`, `user_id`),
    CONSTRAINT `fk_group_users_group_id` FOREIGN KEY (`group_id`) REFERENCES `groups` (`id`),
    CONSTRAINT `fk_group_users_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;
