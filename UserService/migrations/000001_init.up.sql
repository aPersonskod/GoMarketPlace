CREATE TABLE IF NOT EXISTS public."Users"
(
    "Id" uuid NOT NULL,
    "Name" text COLLATE pg_catalog."default" NOT NULL,
    "Email" text COLLATE pg_catalog."default" NOT NULL,
    "Password" text COLLATE pg_catalog."default" NOT NULL,
    "Wallet" integer NOT NULL,
    "Role" text COLLATE pg_catalog."default" NOT NULL DEFAULT 'user'::text,
    CONSTRAINT "PK_Users" PRIMARY KEY ("Id")
);

INSERT INTO public."Users" ("Id", "Name", "Email", "Password", "Wallet", "Role") 
VALUES ('49792511-261b-4edb-94a5-ecb8540e60ff', 'Петр Пятов', 'patochin@gmail.com', '12345', 10, 'user')
ON CONFLICT ("Id") DO NOTHING;

INSERT INTO public."Users" ("Id", "Name", "Email", "Password", "Wallet", "Role")
VALUES ('9e739988-9361-45de-8472-a5a8a5e73a0e', 'Andrey GoProfessional', 'the_greatest@gmail.com', '123321', 10, 'admin')
ON CONFLICT ("Id") DO NOTHING;
