DROP TABLE IF EXISTS fundamentals;

CREATE EXTERNAL TABLE fundamentals (
    `timestamp` timestamp,
    `analystratings` struct<
        `rating`: float,
        `targetprice`: float,
        `strongbuy`: int,
        `buy`: int,
        `hold`: int,
        `sell`: int,
        `strongsell`: int
    >,
    `esgscores` string,
    `earnings` string,
    `financials` string,
    `highlights` string,
    `holders` array<struct<
        `type`: string,
        `change`: int,
        `change_p`: float,
        `currentshares`: int,
        `date`: string,
        `name`: string,
        `totalassets`: float,
        `totalshares`: float
    >>,
    `insidertransactions` string,
    `sharesstats` string,
    `splitsdividends` string,
    `technicals` string,
    `valuation` string,
    `outstandingshares` string,
    `general` struct<
        `code`: string,
        `type`: string,
        `name`: string,
        `exchange`: string,
        `currencycode`: string,
        `currencyname`: string,
        `currencysymbol`: string,
        `countryname`: string,
        `countryiso`: string,
        `openfigi`: string,
        `isin`: string,
        `lei`: string,
        `primaryticker`: string,
        `cusip`: string,
        `cik`: string,
        `employeridnumber`: string,
        `fiscalyearend`: string,
        `ipodate`: string,
        `internationaldomestic`: string,
        `sector`: string,
        `industry`: string,
        `gicsector`: string,
        `gicgroup`: string,
        `gicindustry`: string,
        `gicsubindustry`: string,
        `homecategory`: string,
        `isdelisted`: boolean,
        `description`: string,
        `address`: string,
        `addressdata`: struct<
            `street`: string,
            `city`: string,
            `state`: string,
            `country`: string,
            `zip`: string
        >,
        `listings`: array<struct<
            `code`: string,
            `exchange`: string,
            `name`: string
        >>,
        `officers`: array<struct<
            `name`: string,
            `title`: string,
            `yearborn`: string
        >>,
        `phone`: string,
        `weburl`: string,
        `logourl`: string,
        `fulltimeemployees`: int,
        `updatedat`: string
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