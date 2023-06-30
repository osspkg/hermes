CREATE TABLE `hermes_addon_migrations`
(
    `id`         int unsigned                                                  NOT NULL AUTO_INCREMENT,
    `addon`      varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
    `data`       varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
    `created_at` timestamp                                                     NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `addon_data` (`addon`, `data`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;
