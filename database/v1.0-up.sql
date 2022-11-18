CREATE TABLE `test` (
    `id` CHAR(36),
    `foo` VARCHAR(128),
    `bar` VARCHAR(128),
    PRIMARY KEY (`id`)
);

CREATE TABLE `fakultas` (
    `id` CHAR(36),
    `nama` VARCHAR(32),
    PRIMARY KEY (`id`)
);

CREATE TABLE `admin` (
    `id` CHAR(36),
    `nama` VARCHAR(128),
    `email` VARCHAR(254),
    `password` CHAR(64),
    `fakultas_id` CHAR(36),
    PRIMARY KEY (`id`),
    CONSTRAINT `FK_Admin_Fakultas` FOREIGN KEY (`fakultas_id`) REFERENCES `fakultas` (`id`)
);

CREATE TABLE `peneliti` (
    `id` CHAR(36),
    `nidn` VARCHAR(64),
    `nama` VARCHAR(128),
    `jenis_kelamin` ENUM('Laki-Laki', 'Perempuan'),
    `scopus_author_id` VARCHAR(128),
    `gscholar_author_id` VARCHAR(128),
    `fakultas_id` CHAR(36),
    `diciptakan_oleh_id` CHAR(36),
    PRIMARY KEY (`id`),
    CONSTRAINT `FK_Peneliti_Fakultas` FOREIGN KEY (`fakultas_id`) REFERENCES `fakultas` (`id`),
    CONSTRAINT `FK_Peneliti_DiciptakanOleh` FOREIGN KEY (`diciptakan_oleh_id`) REFERENCES `admin` (`id`)
);

INSERT INTO `fakultas` (`id`, `nama`) VALUES
('27105b83-abdc-42d9-9442-f8076f0d86f3', 'FK'),
('32ae0d39-5e85-4323-8717-9f1ec8a2729d', 'FT'),
('d6c8d67d-853f-4246-b253-2a0929210b5a', 'FMIPA');
