CREATE TABLE IF NOT EXISTS public."BuyReports" (
    "Id" uuid NOT NULL,
    "CartId" uuid NOT NULL,
    "SaleDate" timestamp with time zone NOT NULL,
    CONSTRAINT "PK_BuyReports" PRIMARY KEY ("Id")
);