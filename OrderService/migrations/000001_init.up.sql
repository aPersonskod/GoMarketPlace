CREATE TABLE IF NOT EXISTS public."Orders" (
    "Id" uuid NOT NULL,
    "CartId" uuid NOT NULL,
    "OrderedProductId" uuid NOT NULL,
    "Quantity" integer NOT NULL,
    CONSTRAINT "PK_Orders" PRIMARY KEY ("Id")
);

CREATE TABLE IF NOT EXISTS public."Places" (
    "Id" uuid NOT NULL,
    "Address" text NOT NULL,
    "WorkingTime" text NOT NULL,
    CONSTRAINT "PK_Places" PRIMARY KEY ("Id")
);

INSERT INTO public."Places" ("Id", "Address", "WorkingTime") 
VALUES ('98eac40c-77e6-44c8-8165-b9380b59a37b', '6-я Советская улица, 37', '09:00 - 21:00')
ON CONFLICT ("Id") DO NOTHING;

INSERT INTO public."Places" ("Id", "Address", "WorkingTime") 
VALUES ('f853bb36-6ad3-4d03-ad7e-9a3545d21429', 'Яхтенная ул., 3, корп. 2', '10:00 - 22:00')
ON CONFLICT ("Id") DO NOTHING;

CREATE TABLE IF NOT EXISTS public."ShoppingCarts" (
    "Id" uuid NOT NULL,
    "UserId" uuid NOT NULL,
    "PlaceId" uuid,
    "AmountToPay" integer DEFAULT 0 NOT NULL,
    "IsConfirmed" boolean DEFAULT false NOT NULL,
    "IsBought" boolean DEFAULT false NOT NULL,
    CONSTRAINT "PK_ShoppingCarts" PRIMARY KEY ("Id")
);