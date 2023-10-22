DROP TABLE IF EXISTS fundamentals;

CREATE EXTERNAL TABLE fundamentals (
    `AnalystRatings` struct<
        `Rating`: float,
        `TargetPrice`: float,
        `StrongBuy`: int,
        `Buy`: int,
        `Hold`: int,
        `Sell`: int,
        `StrongSell`: int
    >,
    `ESGScores` string,
    `Earnings` string,
    `Financials` string,
    `Highlights` string,
    `Holders` array<struct<
        `type`: string,
        `change`: int,
        `change_p`: float,
        `currentShares`: int,
        `date`: string,
        `name`: string,
        `totalAssets`: float,
        `totalShares`: float
    >>,
    `InsiderTransactions` string,
    `SharesStats` string,
    `SplitsDividends` string,
    `Technicals` string,
    `Valuation` string,
    `outstandingShares` string,
    `General` struct<
        `Code`: string,
        `Type`: string,
        `Name`: string,
        `Exchange`: string,
        `CurrencyCode`: string,
        `CurrencyName`: string,
        `CurrencySymbol`: string,
        `CountryName`: string,
        `CountryISO`: string,
        `OpenFigi`: string,
        `Isin`: string,
        `Lei`: string,
        `PrimaryTicker`: string,
        `Cusip`: string,
        `Cik`: string,
        `EmployerIDNumber`: string,
        `FiscalYearEnd`: string,
        `IPODate`: string,
        `InternationalDomestic`: string,
        `Sector`: string,
        `Industry`: string,
        `GicSector`: string,
        `GicGroup`: string,
        `GicIndustry`: string,
        `GicSubIndustry`: string,
        `HomeCategory`: string,
        `IsDelisted`: boolean,
        `Description`: string,
        `Address`: string,
        `AddressData`: struct<
            `Street`: string,
            `City`: string,
            `State`: string,
            `Country`: string,
            `Zip`: string
        >,
        `Listings`: array<struct<
            `Code`: string,
            `Exchange`: string,
            `Name`: string
        >>,
        `Officers`: array<struct<
            `Name`: string,
            `Title`: string,
            `YearBorn`: string
        >>,
        `Phone`: string,
        `WebURL`: string,
        `LogoURL`: string,
        `FullTimeEmployees`: int,
        `UpdatedAt`: string
    >
)
ROW FORMAT SERDE 'org.openx.data.jsonserde.JsonSerDe'
LOCATION 's3://bcap-stock-data/fundamentals/';


DROP TABLE IF EXISTS fundamentals_plain_text;

CREATE EXTERNAL TABLE fundamentals_plain_text (
    data string
)
ROW FORMAT SERDE 'org.apache.hadoop.hive.serde2.lazy.LazySimpleSerDe'
STORED AS INPUTFORMAT 'org.apache.hadoop.mapred.TextInputFormat'
OUTPUTFORMAT 'org.apache.hadoop.hive.ql.io.HiveIgnoreKeyTextOutputFormat'
LOCATION 's3://bcap-stock-data/fundamentals/';