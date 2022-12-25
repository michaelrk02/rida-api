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
    `h_index` INT,
    `is_remote` BOOLEAN,
    PRIMARY KEY (`id`),
    CONSTRAINT `FK_Peneliti_Fakultas` FOREIGN KEY (`fakultas_id`) REFERENCES `fakultas` (`id`),
    CONSTRAINT `FK_Peneliti_DiciptakanOleh` FOREIGN KEY (`diciptakan_oleh_id`) REFERENCES `admin` (`id`)
);

INSERT INTO `fakultas` (`id`, `nama`) VALUES
('27105b83-abdc-42d9-9442-f8076f0d86f3', 'FMIPA'),
('32ae0d39-5e85-4323-8717-9f1ec8a2729d', 'FT'),
('d6c8d67d-853f-4246-b253-2a0929210b5a', 'FKIP'),
('d79601f2-7a5b-4e8c-b6c3-ecbd212a5ea8', 'FIB'),
('c214f479-dbb7-47ce-b6b4-dc74e50d15d8', 'FISIP'),
('4f8ccd4f-4e1b-439c-bc35-1014d6d1ee44', 'FH'),
('db709f58-f936-4e67-aa50-a1c8af1b504d', 'FEB'),
('b62d129e-3517-4c87-806c-6b77fb3512ff', 'FATISDA'),
('3a093e35-4395-49ce-a505-85aa552b9b42', 'SV'),
('dd401155-4cb9-48fc-bbbc-767e1d7febd8', 'FP'),
('d6e24975-ca27-4388-8239-6eb4a930adb0', 'FKOR'),
('aec89f66-fc02-4e25-8ee0-b345d9cc04cf', 'FK'),
('50a5daaf-bdc4-470b-8a3e-b5783d65ce0f', 'FAPSI');
