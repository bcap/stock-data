DROP TABLE IF EXISTS exchanges;

CREATE EXTERNAL TABLE exchanges (
    `Code` string,
    `Country` string,
    `CountryISO2` string,
    `CountryISO3` string,
    `Name` string,
    `OperatingMIC` array<string>
)
ROW FORMAT SERDE 'org.openx.data.jsonserde.JsonSerDe'
LOCATION 's3://bcap-stock-data/exchanges/';


DROP TABLE IF EXISTS exchanges_plain_text;

CREATE EXTERNAL TABLE exchanges_plain_text (
	data string
)
ROW FORMAT SERDE 'org.apache.hadoop.hive.serde2.lazy.LazySimpleSerDe'
STORED AS INPUTFORMAT 'org.apache.hadoop.mapred.TextInputFormat'
OUTPUTFORMAT 'org.apache.hadoop.hive.ql.io.HiveIgnoreKeyTextOutputFormat'
LOCATION 's3://bcap-stock-data/exchanges/';